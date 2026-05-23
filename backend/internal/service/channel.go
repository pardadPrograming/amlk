package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"amlakcrm/backend/internal/domain"
	"amlakcrm/backend/internal/repository"
	"amlakcrm/backend/internal/support"
)

type ChannelService struct {
	store          repository.Store
	objectDir      string
	mediaOptimizer *MediaOptimizerClient
}

const (
	maxChannelImageBytes = 800 * 1024
	maxChannelVideoBytes = 50 * 1024 * 1024
	channelMediaTTL      = 180 * 24 * time.Hour
	maxVaultMembers      = 100
	onlinePresenceWindow = 2 * time.Minute
)

func NewChannelService(store repository.Store, objectDir string, mediaOptimizerURL ...string) *ChannelService {
	var optimizer *MediaOptimizerClient
	if len(mediaOptimizerURL) > 0 && strings.TrimSpace(mediaOptimizerURL[0]) != "" {
		optimizer = NewMediaOptimizerClient(mediaOptimizerURL[0])
	}
	return &ChannelService{store: store, objectDir: objectDir, mediaOptimizer: optimizer}
}

func (s *ChannelService) UserMain(ctx context.Context, user domain.User) (domain.Channel, error) {
	channel, err := s.store.EnsureUserMainChannel(ctx, user)
	if err == nil {
		_, _ = s.store.EnsureUserMainVault(ctx, user)
	}
	return channel, err
}

func (s *ChannelService) BusinessMain(ctx context.Context, userID, businessID string) (domain.Channel, error) {
	member, err := s.store.GetMemberByUser(ctx, businessID, userID)
	if err != nil || member.Status != domain.MemberActive {
		return domain.Channel{}, errors.New("دسترسی به کسب‌وکار وجود ندارد")
	}
	business, err := s.store.GetBusiness(ctx, businessID)
	if err != nil {
		return domain.Channel{}, err
	}
	channel, err := s.store.EnsureBusinessMainChannel(ctx, business)
	if err == nil {
		_, _ = s.store.EnsureBusinessMainVault(ctx, business)
	}
	return channel, err
}

func (s *ChannelService) List(ctx context.Context, user domain.User) ([]domain.Channel, error) {
	businesses, _ := s.store.ListBusinessesForUser(ctx, user.ID)
	businessIDs := make([]string, 0, len(businesses))
	for _, business := range businesses {
		businessIDs = append(businessIDs, business.ID)
		_, _ = s.store.EnsureBusinessMainChannel(ctx, business)
	}
	main, err := s.store.EnsureUserMainChannel(ctx, user)
	if err != nil {
		return nil, err
	}
	_, _ = s.store.EnsureUserMainVault(ctx, user)
	phones := []string{user.Phone}
	channels, err := s.store.ListChannelsForUser(ctx, user.ID, phones, businessIDs)
	if err != nil {
		return nil, err
	}
	seen := map[string]struct{}{main.ID: {}}
	result := []domain.Channel{main}
	for _, channel := range channels {
		if _, ok := seen[channel.ID]; ok {
			continue
		}
		seen[channel.ID] = struct{}{}
		if channel.Type == domain.ChannelTypePrivate {
			channel.Title = s.privateChannelTitle(ctx, user.ID, channel)
		}
		result = append(result, channel)
	}
	return result, nil
}

func (s *ChannelService) PrivateChat(ctx context.Context, actor domain.User, phone string) (domain.Channel, error) {
	normalized, err := NormalizePhone(phone)
	if err != nil {
		return domain.Channel{}, err
	}
	if normalized == actor.Phone {
		return domain.Channel{}, errors.New("cannot create private chat with yourself")
	}
	target, err := s.store.GetUserByPhone(ctx, normalized)
	if err != nil {
		return domain.Channel{}, errors.New("user was not found")
	}
	channel, err := s.store.EnsurePrivateChannel(ctx, actor, target)
	if err != nil {
		return domain.Channel{}, err
	}
	channel.Title = userChannelLabel(target)
	return channel, nil
}

func (s *ChannelService) UserVaults(ctx context.Context, user domain.User) ([]domain.ChannelVault, error) {
	_, _ = s.store.EnsureUserMainVault(ctx, user)
	return s.store.ListUserVaults(ctx, user.ID)
}

func (s *ChannelService) CreateUserVault(ctx context.Context, user domain.User, title string) (domain.ChannelVault, error) {
	return s.store.CreateUserVault(ctx, user, title)
}

func (s *ChannelService) BusinessVaults(ctx context.Context, userID, businessID string) ([]domain.ChannelVault, error) {
	member, err := s.store.GetMemberByUser(ctx, businessID, userID)
	if err != nil || member.Status != domain.MemberActive {
		return nil, errors.New("دسترسی به کسب‌وکار وجود ندارد")
	}
	business, err := s.store.GetBusiness(ctx, businessID)
	if err != nil {
		return nil, err
	}
	if s.canManageBusinessChannels(ctx, userID, businessID) {
		_, _ = s.store.EnsureBusinessMainVault(ctx, business)
		return s.store.ListBusinessVaults(ctx, businessID)
	}
	vaults, err := s.store.ListBusinessVaults(ctx, businessID)
	if err != nil {
		return nil, err
	}
	allowed := make([]domain.ChannelVault, 0, len(vaults))
	for _, vault := range vaults {
		if s.isChannelAdmin(ctx, userID, vault.ChannelID) {
			allowed = append(allowed, vault)
		}
	}
	if len(allowed) == 0 {
		return nil, errors.New("access to business vaults is restricted")
	}
	return allowed, nil
}

func (s *ChannelService) CreateBusinessVault(ctx context.Context, userID, businessID, title string) (domain.ChannelVault, error) {
	member, err := s.store.GetMemberByUser(ctx, businessID, userID)
	if err != nil || member.Status != domain.MemberActive {
		return domain.ChannelVault{}, errors.New("دسترسی به کسب‌وکار وجود ندارد")
	}
	if member.Role != domain.RoleOwner && member.Role != domain.RoleManager && !domain.HasPermission(member, domain.PermBusinessUpdate) {
		return domain.ChannelVault{}, errors.New("دسترسی ساخت صندوقچه املاک را ندارید")
	}
	if err := s.ensureBusinessVaultCapacity(ctx, businessID); err != nil {
		return domain.ChannelVault{}, err
	}
	business, err := s.store.GetBusiness(ctx, businessID)
	if err != nil {
		return domain.ChannelVault{}, err
	}
	return s.store.CreateBusinessVault(ctx, business, title)
}

func (s *ChannelService) InviteByPhone(ctx context.Context, actor domain.User, channelID, phone string) (domain.ChannelMember, error) {
	channel, err := s.authorizeRead(ctx, actor.ID, channelID)
	if err != nil {
		return domain.ChannelMember{}, err
	}
	if channel.Type == domain.ChannelTypeBusinessMain {
		return domain.ChannelMember{}, errors.New("کانال املاک اعضا را به صورت خودکار از اعضای املاک می‌گیرد")
	}
	normalized, err := NormalizePhone(phone)
	if err != nil {
		return domain.ChannelMember{}, err
	}
	if channel.Type == domain.ChannelTypePrivate {
		return domain.ChannelMember{}, errors.New("private chat members cannot be invited manually")
	}
	if channel.Type == domain.ChannelTypeUserMain && channel.OwnerUserID != actor.ID {
		return domain.ChannelMember{}, errors.New("only the owner can invite to this channel")
	}
	if channel.Type != domain.ChannelTypeUserMain &&
		channel.Type != domain.ChannelTypeUserVault &&
		channel.Type != domain.ChannelTypeBusinessVault {
		if _, err := s.authorizeWrite(ctx, actor.ID, channelID); err != nil {
			return domain.ChannelMember{}, err
		}
	}
	member := domain.ChannelMember{
		ChannelID: channelID,
		Phone:     normalized,
		Role:      domain.ChannelMemberRoleMember,
		Status:    domain.ChannelMemberActive,
	}
	if user, err := s.store.GetUserByPhone(ctx, normalized); err == nil {
		member.UserID = user.ID
	}
	if err := s.ensureVaultMemberCapacity(ctx, channel, member); err != nil {
		return domain.ChannelMember{}, err
	}
	return s.store.UpsertChannelMember(ctx, member)
}

func (s *ChannelService) SetMemberAdmin(ctx context.Context, actorID, channelID, memberID string, admin bool) (domain.ChannelMember, error) {
	if _, err := s.authorizeManageVault(ctx, actorID, channelID); err != nil {
		return domain.ChannelMember{}, err
	}
	members, err := s.store.ListChannelMembers(ctx, channelID)
	if err != nil {
		return domain.ChannelMember{}, err
	}
	for _, member := range members {
		if member.ID != memberID {
			continue
		}
		if member.Status != domain.ChannelMemberActive {
			return domain.ChannelMember{}, errors.New("عضو فعال نیست")
		}
		if admin {
			member.Role = domain.ChannelMemberRoleAdmin
		} else {
			member.Role = domain.ChannelMemberRoleMember
		}
		return s.store.UpsertChannelMember(ctx, member)
	}
	return domain.ChannelMember{}, errors.New("عضو کانال پیدا نشد")
}

func (s *ChannelService) Members(ctx context.Context, userID, channelID string) ([]domain.ChannelMember, error) {
	if _, err := s.authorizeRead(ctx, userID, channelID); err != nil {
		return nil, err
	}
	members, err := s.store.ListChannelMembers(ctx, channelID)
	if err != nil {
		return nil, err
	}
	return s.enrichChannelMembers(ctx, members), nil
}

func (s *ChannelService) Messages(ctx context.Context, userID, channelID string, limit, offset int) ([]domain.ChannelMessage, int, error) {
	if _, err := s.authorizeRead(ctx, userID, channelID); err != nil {
		return nil, 0, err
	}
	limit, offset = normalizeChannelPage(limit, offset)
	items, total, err := s.store.ListChannelMessages(ctx, channelID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	messageIDs := make([]string, 0, len(items))
	for _, item := range items {
		if item.AuthorID != userID {
			messageIDs = append(messageIDs, item.ID)
		}
	}
	if len(messageIDs) > 0 {
		_ = s.store.MarkChannelMessagesSeen(ctx, channelID, userID, messageIDs)
		for i := range items {
			if items[i].AuthorID != userID && !messageSeenBy(items[i], userID) {
				items[i].SeenBy = append(items[i].SeenBy, domain.ChannelMessageSeen{UserID: userID, SeenAt: time.Now().UTC()})
			}
		}
	}
	return items, total, nil
}

func (s *ChannelService) SendMessage(ctx context.Context, user domain.User, channelID string, input domain.ChannelMessage) (domain.ChannelMessage, error) {
	if _, err := s.authorizeWrite(ctx, user.ID, channelID); err != nil {
		return domain.ChannelMessage{}, err
	}
	input.ChannelID = channelID
	input.AuthorID = user.ID
	input.AuthorName = user.DisplayName
	input.Text = strings.TrimSpace(input.Text)
	input.Caption = strings.TrimSpace(input.Caption)
	media, err := s.claimChannelMedia(ctx, user.ID, channelID, input.Media)
	if err != nil {
		return domain.ChannelMessage{}, err
	}
	input.Media = media
	input.ReplyToID = strings.TrimSpace(input.ReplyToID)
	input.ReplyTo = nil
	input.VaultFileRefID = strings.TrimSpace(input.VaultFileRefID)
	input.VaultFileRef = nil
	if input.ReplyToID != "" {
		parent, err := s.store.GetChannelMessage(ctx, channelID, input.ReplyToID)
		if err != nil {
			return domain.ChannelMessage{}, errors.New("پیام ریپلای شده در این کانال پیدا نشد")
		}
		input.ReplyTo = channelReplyPreview(parent)
	}
	if input.VaultFileRefID != "" {
		file, err := s.store.GetChannelVaultFile(ctx, channelID, input.VaultFileRefID)
		if err != nil {
			return domain.ChannelMessage{}, errors.New("فایل رفرنس شده در صندوقچه این کانال پیدا نشد")
		}
		input.VaultFileRef = channelVaultFilePreview(file)
	}
	if input.Text == "" && input.Caption == "" && len(input.Media) == 0 && input.VaultFileRefID == "" {
		return domain.ChannelMessage{}, errors.New("متن یا مدیا الزامی است")
	}
	return s.store.CreateChannelMessage(ctx, input)
}

func (s *ChannelService) UpdateMessage(ctx context.Context, user domain.User, channelID, messageID string, input domain.ChannelMessage) (domain.ChannelMessage, error) {
	message, err := s.mutableMessage(ctx, user.ID, channelID, messageID)
	if err != nil {
		return domain.ChannelMessage{}, err
	}
	message.Text = strings.TrimSpace(input.Text)
	message.Caption = strings.TrimSpace(input.Caption)
	if message.Text == "" && message.Caption == "" && len(message.Media) == 0 && message.VaultFileRefID == "" {
		return domain.ChannelMessage{}, errors.New("Ù…ØªÙ† ÛŒØ§ Ú©Ù¾Ø´Ù† Ø¨Ø¹Ø¯ Ø§Ø² ÙˆÛŒØ±Ø§ÛŒØ´ Ù†Ø¨Ø§ÛŒØ¯ Ø®Ø§Ù„ÛŒ Ø¨Ø§Ø´Ø¯")
	}
	return s.store.UpdateChannelMessage(ctx, message)
}

func (s *ChannelService) DeleteMessage(ctx context.Context, user domain.User, channelID, messageID string) error {
	if _, err := s.mutableMessage(ctx, user.ID, channelID, messageID); err != nil {
		return err
	}
	return s.store.DeleteChannelMessage(ctx, channelID, messageID)
}

func (s *ChannelService) EnsureMessageMutable(ctx context.Context, actorID, channelID, messageID string) error {
	_, err := s.mutableMessage(ctx, actorID, channelID, messageID)
	return err
}

func (s *ChannelService) mutableMessage(ctx context.Context, actorID, channelID, messageID string) (domain.ChannelMessage, error) {
	channel, err := s.authorizeWrite(ctx, actorID, channelID)
	if err != nil {
		return domain.ChannelMessage{}, err
	}
	message, err := s.store.GetChannelMessage(ctx, channelID, messageID)
	if err != nil {
		return domain.ChannelMessage{}, err
	}
	if message.AuthorID != actorID {
		return domain.ChannelMessage{}, errors.New("only message author can modify the message")
	}
	if channel.Type == domain.ChannelTypePrivate && messageSeenByOther(message, actorID) {
		return domain.ChannelMessage{}, errors.New("message cannot be edited or deleted after another member has seen it")
	}
	return message, nil
}

func (s *ChannelService) VaultFiles(ctx context.Context, userID, channelID string, limit, offset int) ([]domain.ChannelVaultFile, int, error) {
	if _, err := s.authorizeRead(ctx, userID, channelID); err != nil {
		return nil, 0, err
	}
	limit, offset = normalizeChannelPage(limit, offset)
	return s.store.ListChannelVaultFiles(ctx, channelID, limit, offset)
}

func (s *ChannelService) VaultFile(ctx context.Context, userID, channelID, fileID string) (domain.ChannelVaultFile, error) {
	if _, err := s.authorizeRead(ctx, userID, channelID); err != nil {
		return domain.ChannelVaultFile{}, err
	}
	return s.store.GetChannelVaultFile(ctx, channelID, fileID)
}

func (s *ChannelService) SaveVaultFile(ctx context.Context, user domain.User, channelID string, header *multipart.FileHeader, title, note string) (domain.ChannelVaultFile, error) {
	if _, err := s.authorizeWrite(ctx, user.ID, channelID); err != nil {
		return domain.ChannelVaultFile{}, err
	}
	vault, err := s.vaultForChannel(ctx, channelID)
	if err != nil {
		return domain.ChannelVaultFile{}, err
	}
	media, err := s.saveObject(ctx, channelID, header)
	if err != nil {
		return domain.ChannelVaultFile{}, err
	}
	title = strings.TrimSpace(title)
	if title == "" {
		title = header.Filename
	}
	return s.store.CreateChannelVaultFile(ctx, domain.ChannelVaultFile{
		VaultID:     vault.ID,
		ChannelID:   channelID,
		UploaderID:  user.ID,
		Title:       title,
		Note:        strings.TrimSpace(note),
		SourceType:  "upload",
		FileID:      media.FileID,
		Kind:        media.Kind,
		URL:         media.URL,
		ContentType: media.ContentType,
		Size:        media.Size,
	})
}

func (s *ChannelService) SaveVaultFileFromUpload(ctx context.Context, user domain.User, channelID, fileID, title, note string) (domain.ChannelVaultFile, error) {
	if _, err := s.authorizeWrite(ctx, user.ID, channelID); err != nil {
		return domain.ChannelVaultFile{}, err
	}
	vault, err := s.vaultForChannel(ctx, channelID)
	if err != nil {
		return domain.ChannelVaultFile{}, err
	}
	object, err := s.claimUploadedFile(ctx, user.ID, fileID, UploadPurposeVaultFile, "channel", channelID, 0)
	if err != nil {
		return domain.ChannelVaultFile{}, err
	}
	title = strings.TrimSpace(title)
	if title == "" {
		title = "فایل صندوقچه"
	}
	return s.store.CreateChannelVaultFile(ctx, domain.ChannelVaultFile{
		VaultID:     vault.ID,
		ChannelID:   channelID,
		UploaderID:  user.ID,
		Title:       title,
		Note:        strings.TrimSpace(note),
		SourceType:  "upload",
		FileID:      object.ID,
		Kind:        object.Kind,
		URL:         object.URL,
		ContentType: object.ContentType,
		Size:        object.Size,
	})
}

func (s *ChannelService) vaultForChannel(ctx context.Context, channelID string) (domain.ChannelVault, error) {
	channel, err := s.store.GetChannel(ctx, channelID)
	if err != nil {
		return domain.ChannelVault{}, err
	}
	if channel.VaultID != "" {
		return s.store.GetChannelVault(ctx, channel.VaultID)
	}
	if channel.OwnerUserID != "" {
		user, err := s.store.GetUser(ctx, channel.OwnerUserID)
		if err != nil {
			return domain.ChannelVault{}, err
		}
		return s.store.EnsureUserMainVault(ctx, user)
	}
	if channel.BusinessID != "" {
		business, err := s.store.GetBusiness(ctx, channel.BusinessID)
		if err != nil {
			return domain.ChannelVault{}, err
		}
		return s.store.EnsureBusinessMainVault(ctx, business)
	}
	return domain.ChannelVault{}, errors.New("صندوقچه کانال پیدا نشد")
}

func (s *ChannelService) SaveMedia(ctx context.Context, user domain.User, channelID string, header *multipart.FileHeader) (domain.ChannelMedia, error) {
	if _, err := s.authorizeWrite(ctx, user.ID, channelID); err != nil {
		return domain.ChannelMedia{}, err
	}
	return s.saveOptimizedObject(ctx, channelID, header)
}

func (s *ChannelService) saveObject(ctx context.Context, channelID string, header *multipart.FileHeader) (domain.ChannelMedia, error) {
	if header == nil || header.Size <= 0 {
		return domain.ChannelMedia{}, errors.New("فایل معتبر نیست")
	}
	if header.Size > 50*1024*1024 {
		return domain.ChannelMedia{}, errors.New("حجم فایل کانال نباید بیشتر از ۵۰ مگابایت باشد")
	}
	contentType := detectChannelContentType(header)
	kind := "file"
	switch {
	case strings.HasPrefix(contentType, "image/"):
		kind = "image"
	case strings.HasPrefix(contentType, "video/"):
		kind = "video"
	case strings.HasPrefix(contentType, "audio/"):
		kind = "audio"
	}
	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".bin"
	}
	key := fmt.Sprintf("channels/%s/%s%s", channelID, support.NewID(), ext)
	if err := s.writeUpload(key, header); err != nil {
		return domain.ChannelMedia{}, err
	}
	object, err := s.store.CreateFile(ctx, domain.FileObject{
		OwnerID:     channelID,
		Provider:    "s3-compatible",
		Bucket:      "amlak",
		Key:         key,
		URL:         "/objects/" + key,
		ContentType: contentType,
		Size:        header.Size,
	})
	if err != nil {
		return domain.ChannelMedia{}, err
	}
	return domain.ChannelMedia{
		ID:          support.NewID(),
		FileID:      object.ID,
		Kind:        kind,
		URL:         object.URL,
		ContentType: contentType,
		Size:        header.Size,
		CreatedAt:   object.CreatedAt,
	}, nil
}

func (s *ChannelService) saveOptimizedObject(ctx context.Context, channelID string, header *multipart.FileHeader) (domain.ChannelMedia, error) {
	if header == nil || header.Size <= 0 {
		return domain.ChannelMedia{}, errors.New("invalid file")
	}
	contentType := detectChannelContentType(header)
	kind := "file"
	switch {
	case strings.HasPrefix(contentType, "image/"):
		kind = "image"
	case strings.HasPrefix(contentType, "video/"):
		kind = "video"
	case strings.HasPrefix(contentType, "audio/"):
		kind = "audio"
	}

	var (
		body []byte
		ext  string
		err  error
	)
	if kind == "image" || kind == "video" {
		optimized, optimizeErr := s.mediaOptimizer.Optimize(ctx, header, maxChannelImageBytes, maxChannelVideoBytes)
		if optimizeErr != nil {
			return domain.ChannelMedia{}, optimizeErr
		}
		body = optimized.Body
		contentType = optimized.ContentType
		if optimized.Kind != "" {
			kind = optimized.Kind
		}
		ext = optimized.Extension
	} else {
		if header.Size > maxChannelVideoBytes {
			return domain.ChannelMedia{}, errors.New("channel file size must be at most 50 MB")
		}
		body, err = readUpload(header)
		if err != nil {
			return domain.ChannelMedia{}, err
		}
		ext = filepath.Ext(header.Filename)
	}
	if kind == "image" && int64(len(body)) > maxChannelImageBytes {
		return domain.ChannelMedia{}, errors.New("channel images must be at most 800 KB after compression")
	}
	if kind == "video" && int64(len(body)) > maxChannelVideoBytes {
		return domain.ChannelMedia{}, errors.New("channel videos must be at most 50 MB after compression")
	}
	if ext == "" {
		ext = ".bin"
	}
	key := fmt.Sprintf("channels/%s/%s%s", channelID, support.NewID(), ext)
	if err := s.writeObject(key, body); err != nil {
		return domain.ChannelMedia{}, err
	}
	expiresAt := time.Now().UTC().Add(channelMediaTTL)
	object, err := s.store.CreateFile(ctx, domain.FileObject{
		OwnerID:     channelID,
		Provider:    "s3-compatible",
		Bucket:      "amlak",
		Key:         key,
		URL:         "/objects/" + key,
		ContentType: contentType,
		Size:        int64(len(body)),
		ExpiresAt:   expiresAt,
	})
	if err != nil {
		return domain.ChannelMedia{}, err
	}
	return domain.ChannelMedia{
		ID:          support.NewID(),
		FileID:      object.ID,
		Kind:        kind,
		URL:         object.URL,
		ContentType: contentType,
		Size:        int64(len(body)),
		CreatedAt:   object.CreatedAt,
		ExpiresAt:   object.ExpiresAt,
	}, nil
}

func (s *ChannelService) authorizeRead(ctx context.Context, userID, channelID string) (domain.Channel, error) {
	channel, err := s.store.GetChannel(ctx, channelID)
	if err != nil {
		return domain.Channel{}, err
	}
	if channel.OwnerUserID == userID {
		return channel, nil
	}
	if channel.BusinessID != "" {
		if member, err := s.store.GetMemberByUser(ctx, channel.BusinessID, userID); err == nil && member.Status == domain.MemberActive {
			if channel.Type != domain.ChannelTypeBusinessVault ||
				s.canManageBusinessChannels(ctx, userID, channel.BusinessID) ||
				s.isChannelAdmin(ctx, userID, channelID) {
				return channel, nil
			}
		}
	}
	if _, err := s.store.GetChannelMember(ctx, channelID, userID); err == nil {
		return channel, nil
	}
	if user, err := s.store.GetUser(ctx, userID); err == nil {
		members, _ := s.store.ListChannelMembers(ctx, channelID)
		for _, member := range members {
			if member.Status == domain.ChannelMemberActive && member.Phone != "" && member.Phone == user.Phone {
				return channel, nil
			}
		}
	}
	return domain.Channel{}, errors.New("دسترسی به کانال وجود ندارد")
}

func (s *ChannelService) authorizeWrite(ctx context.Context, userID, channelID string) (domain.Channel, error) {
	channel, err := s.authorizeRead(ctx, userID, channelID)
	if err != nil {
		return domain.Channel{}, err
	}
	if channel.Type == domain.ChannelTypeUserMain && channel.OwnerUserID != userID {
		return domain.Channel{}, errors.New("فقط صاحب کانال شخصی می‌تواند پیام ارسال کند")
	}
	if channel.Type == domain.ChannelTypeUserVault && channel.OwnerUserID != userID && !s.isChannelAdmin(ctx, userID, channelID) {
		return domain.Channel{}, errors.New("برای ارسال در این صندوقچه باید ادمین باشید")
	}
	if channel.Type == domain.ChannelTypeBusinessVault {
		if s.isChannelAdmin(ctx, userID, channelID) {
			return channel, nil
		}
		if s.canManageBusinessChannels(ctx, userID, channel.BusinessID) {
			return channel, nil
		}
		if member, err := s.store.GetMemberByUser(ctx, channel.BusinessID, userID); err == nil && member.Status == domain.MemberActive {
			if member.Role == domain.RoleOwner || member.Role == domain.RoleManager || domain.HasPermission(member, domain.PermBusinessUpdate) {
				return channel, nil
			}
		}
		return domain.Channel{}, errors.New("برای ارسال در این صندوقچه باید ادمین باشید")
	}
	return channel, nil
}

func (s *ChannelService) authorizeManageVault(ctx context.Context, userID, channelID string) (domain.Channel, error) {
	channel, err := s.authorizeRead(ctx, userID, channelID)
	if err != nil {
		return domain.Channel{}, err
	}
	if channel.OwnerUserID == userID {
		return channel, nil
	}
	if channel.BusinessID != "" {
		if member, err := s.store.GetMemberByUser(ctx, channel.BusinessID, userID); err == nil && member.Status == domain.MemberActive {
			if member.Role == domain.RoleOwner || member.Role == domain.RoleManager || domain.HasPermission(member, domain.PermBusinessUpdate) {
				return channel, nil
			}
		}
	}
	return domain.Channel{}, errors.New("دسترسی مدیریت صندوقچه را ندارید")
}

func (s *ChannelService) isChannelAdmin(ctx context.Context, userID, channelID string) bool {
	member, err := s.store.GetChannelMember(ctx, channelID, userID)
	return err == nil && member.Status == domain.ChannelMemberActive && member.Role == domain.ChannelMemberRoleAdmin
}

func (s *ChannelService) canManageBusinessChannels(ctx context.Context, userID, businessID string) bool {
	if strings.TrimSpace(businessID) == "" {
		return false
	}
	member, err := s.store.GetMemberByUser(ctx, businessID, userID)
	if err != nil || member.Status != domain.MemberActive {
		return false
	}
	return member.Role == domain.RoleOwner ||
		member.Role == domain.RoleManager ||
		domain.HasPermission(member, domain.PermBusinessUpdate) ||
		domain.HasPermission(member, domain.PermMembersManage)
}

func (s *ChannelService) ensureBusinessVaultCapacity(ctx context.Context, businessID string) error {
	members, err := s.store.ListMembers(ctx, businessID)
	if err != nil {
		return err
	}
	activeMembers := 0
	for _, member := range members {
		if member.Status == domain.MemberActive {
			activeMembers++
		}
	}
	if activeMembers > maxVaultMembers {
		return errors.New("each vault can have at most 100 members")
	}
	return nil
}

func (s *ChannelService) ensureVaultMemberCapacity(ctx context.Context, channel domain.Channel, next domain.ChannelMember) error {
	if channel.Type != domain.ChannelTypeUserVault && channel.Type != domain.ChannelTypeBusinessVault {
		return nil
	}
	members, err := s.store.ListChannelMembers(ctx, channel.ID)
	if err != nil {
		return err
	}
	activeCount := 0
	for _, member := range members {
		if member.Status != domain.ChannelMemberActive && member.Status != domain.ChannelMemberInvited {
			continue
		}
		sameUser := next.UserID != "" && member.UserID == next.UserID
		samePhone := next.Phone != "" && member.Phone == next.Phone
		if sameUser || samePhone {
			return nil
		}
		activeCount++
	}
	if activeCount >= maxVaultMembers {
		return errors.New("each vault can have at most 100 members")
	}
	return nil
}

func (s *ChannelService) enrichChannelMembers(ctx context.Context, members []domain.ChannelMember) []domain.ChannelMember {
	now := time.Now().UTC()
	enriched := make([]domain.ChannelMember, 0, len(members))
	for _, member := range members {
		if member.UserID == "" && member.Phone != "" {
			if user, err := s.store.GetUserByPhone(ctx, member.Phone); err == nil {
				member.UserID = user.ID
			}
		}
		if member.UserID != "" {
			if user, err := s.store.GetUser(ctx, member.UserID); err == nil {
				member.DisplayName = user.DisplayName
				if user.PrivacySettings.ShowActivityStatus {
					lastSeenAt, isOnline := s.userPresence(ctx, member.UserID, now)
					member.LastSeenAt = lastSeenAt
					member.IsOnline = isOnline
				}
			}
		}
		enriched = append(enriched, member)
	}
	return enriched
}

func (s *ChannelService) userPresence(ctx context.Context, userID string, now time.Time) (time.Time, bool) {
	sessions, err := s.store.ListSessions(ctx, userID)
	if err != nil {
		return time.Time{}, false
	}
	var lastSeenAt time.Time
	for _, session := range sessions {
		if session.LastSeenAt.After(lastSeenAt) {
			lastSeenAt = session.LastSeenAt
		}
	}
	return lastSeenAt, !lastSeenAt.IsZero() && now.Sub(lastSeenAt) <= onlinePresenceWindow
}

func (s *ChannelService) privateChannelTitle(ctx context.Context, userID string, channel domain.Channel) string {
	members, err := s.store.ListChannelMembers(ctx, channel.ID)
	if err != nil {
		return channel.Title
	}
	for _, member := range members {
		if member.UserID == "" || member.UserID == userID {
			continue
		}
		if user, err := s.store.GetUser(ctx, member.UserID); err == nil {
			return userChannelLabel(user)
		}
		if member.Phone != "" {
			return member.Phone
		}
	}
	for _, member := range members {
		if member.Phone != "" {
			return member.Phone
		}
	}
	return channel.Title
}

func userChannelLabel(user domain.User) string {
	if strings.TrimSpace(user.DisplayName) != "" {
		return strings.TrimSpace(user.DisplayName)
	}
	return user.Phone
}

func (s *ChannelService) StartMediaRetentionCleanup(ctx context.Context, interval time.Duration) {
	if interval <= 0 {
		interval = 24 * time.Hour
	}
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		s.cleanupExpiredChannelMedia()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				s.cleanupExpiredChannelMedia()
			}
		}
	}()
}

func (s *ChannelService) cleanupExpiredChannelMedia() {
	root := filepath.Join(s.objectDir, "channels")
	cutoff := time.Now().Add(-channelMediaTTL)
	_ = filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil || entry == nil || entry.IsDir() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return nil
		}
		if info.ModTime().Before(cutoff) {
			_ = os.Remove(path)
		}
		return nil
	})
}

func (s *ChannelService) writeUpload(key string, header *multipart.FileHeader) error {
	src, err := header.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	path := filepath.Join(s.objectDir, filepath.FromSlash(key))
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	return err
}

func (s *ChannelService) writeObject(key string, body []byte) error {
	path := filepath.Join(s.objectDir, filepath.FromSlash(key))
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, body, 0644)
}

func readUpload(header *multipart.FileHeader) ([]byte, error) {
	src, err := header.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()
	return io.ReadAll(src)
}

func detectChannelContentType(header *multipart.FileHeader) string {
	if header == nil {
		return ""
	}
	if contentType := strings.TrimSpace(header.Header.Get("Content-Type")); contentType != "" && contentType != "application/octet-stream" {
		return contentType
	}
	if ext := filepath.Ext(header.Filename); ext != "" {
		if contentType := mime.TypeByExtension(ext); contentType != "" {
			return contentType
		}
	}
	return header.Header.Get("Content-Type")
}

func normalizeChannelPage(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = 30
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}

func channelReplyPreview(message domain.ChannelMessage) *domain.ChannelReplyPreview {
	preview := &domain.ChannelReplyPreview{
		ID:         message.ID,
		AuthorID:   message.AuthorID,
		AuthorName: message.AuthorName,
		Text:       message.Text,
		Caption:    message.Caption,
		MediaCount: len(message.Media),
	}
	if len(message.Media) > 0 {
		preview.MediaKind = message.Media[0].Kind
	}
	return preview
}

func messageSeenBy(message domain.ChannelMessage, userID string) bool {
	for _, seen := range message.SeenBy {
		if seen.UserID == userID {
			return true
		}
	}
	return false
}

func messageSeenByOther(message domain.ChannelMessage, authorID string) bool {
	for _, seen := range message.SeenBy {
		if seen.UserID != "" && seen.UserID != authorID {
			return true
		}
	}
	return false
}

func channelVaultFilePreview(file domain.ChannelVaultFile) *domain.ChannelVaultFilePreview {
	return &domain.ChannelVaultFilePreview{
		ID:                file.ID,
		Title:             file.Title,
		Kind:              file.Kind,
		URL:               file.URL,
		ContentType:       file.ContentType,
		Size:              file.Size,
		PropertyStatus:    file.PropertyStatus,
		CommissionPercent: file.CommissionPercent,
	}
}

func (s *ChannelService) claimChannelMedia(ctx context.Context, userID, channelID string, items []domain.ChannelMedia) ([]domain.ChannelMedia, error) {
	items = normalizeChannelMedia(items)
	if len(items) == 0 {
		return nil, nil
	}
	result := make([]domain.ChannelMedia, 0, len(items))
	for _, item := range items {
		if item.FileID == "" {
			result = append(result, item)
			continue
		}
		object, err := s.claimUploadedFile(ctx, userID, item.FileID, UploadPurposeChannelMedia, "channel", channelID, channelMediaTTL)
		if err != nil {
			if item.URL != "" {
				result = append(result, item)
				continue
			}
			return nil, err
		}
		kind := object.Kind
		if kind == "" {
			kind = item.Kind
		}
		result = append(result, domain.ChannelMedia{
			ID:          firstNonEmptyString(item.ID, support.NewID()),
			FileID:      object.ID,
			Kind:        kind,
			URL:         object.URL,
			ContentType: object.ContentType,
			Size:        object.Size,
			CreatedAt:   object.CreatedAt,
			ExpiresAt:   object.ExpiresAt,
		})
	}
	return result, nil
}

func (s *ChannelService) claimUploadedFile(ctx context.Context, userID, fileID, purpose, targetType, targetID string, ttl time.Duration) (domain.FileObject, error) {
	fileID = strings.TrimSpace(fileID)
	if fileID == "" {
		return domain.FileObject{}, errors.New("شناسه فایل الزامی است")
	}
	object, err := s.store.GetFile(ctx, fileID)
	if err != nil {
		return domain.FileObject{}, errors.New("فایل آپلود شده پیدا نشد")
	}
	if object.UploaderID != "" && object.UploaderID != userID {
		return domain.FileObject{}, errors.New("فایل متعلق به کاربر دیگری است")
	}
	if object.Purpose != "" && object.Purpose != purpose {
		return domain.FileObject{}, errors.New("هدف فایل با درخواست همخوانی ندارد")
	}
	if object.TargetID != "" && object.TargetID != targetID {
		return domain.FileObject{}, errors.New("مقصد فایل با درخواست همخوانی ندارد")
	}
	if !object.ExpiresAt.IsZero() && time.Now().UTC().After(object.ExpiresAt) {
		return domain.FileObject{}, errors.New("مهلت استفاده از فایل تمام شده است")
	}
	object.Purpose = purpose
	object.TargetType = targetType
	object.TargetID = targetID
	object.OwnerID = targetID
	object.Status = FileStatusAttached
	if ttl > 0 {
		object.ExpiresAt = time.Now().UTC().Add(ttl)
	}
	return s.store.UpdateFile(ctx, object)
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func normalizeChannelMedia(items []domain.ChannelMedia) []domain.ChannelMedia {
	if len(items) == 0 {
		return nil
	}
	result := make([]domain.ChannelMedia, 0, len(items))
	for _, item := range items {
		item.ID = strings.TrimSpace(item.ID)
		item.FileID = strings.TrimSpace(item.FileID)
		item.Kind = strings.TrimSpace(item.Kind)
		item.URL = strings.TrimSpace(item.URL)
		item.ContentType = strings.TrimSpace(item.ContentType)
		if item.ID == "" || item.FileID == "" || item.URL == "" {
			continue
		}
		result = append(result, item)
	}
	return result
}
