package repository

import (
	"context"
	"errors"
	"sort"
	"strings"
	"sync"
	"time"

	"amlakcrm/backend/internal/domain"
	"amlakcrm/backend/internal/support"
)

var ErrNotFound = errors.New("not found")

type MemoryStore struct {
	mu                    sync.RWMutex
	otps                  map[string]domain.OTPChallenge
	latestOTP             string
	users                 map[string]domain.User
	usersPhone            map[string]string
	sessions              map[string]domain.Session
	adminAccounts         map[string]domain.AdminAccount
	platformSettings      domain.PlatformSettings
	businesses            map[string]domain.Business
	members               map[string]domain.BusinessMember
	invitations           map[string]domain.Invitation
	files                 map[string]domain.FileObject
	notifications         map[string]domain.Notification
	channels              map[string]domain.Channel
	channelVaults         map[string]domain.ChannelVault
	channelMembers        map[string]domain.ChannelMember
	channelMessages       map[string]domain.ChannelMessage
	channelVaultFiles     map[string]domain.ChannelVaultFile
	areas                 map[string]domain.Area
	streets               map[string]domain.Street
	neighborhoods         map[string]domain.Neighborhood
	propertyFiles         map[string]domain.PropertyFile
	propertyShareRequests map[string]domain.PropertyShareRequest
	contacts              map[string]domain.Contact
	cities                map[string]domain.City
	systemAreas           map[string]domain.SystemArea
	systemStreets         map[string]domain.SystemStreet
	systemNeighborhoods   map[string]domain.SystemNeighborhood
	locationSuggestions   map[string]domain.LocationSuggestion
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		otps:                  map[string]domain.OTPChallenge{},
		users:                 map[string]domain.User{},
		usersPhone:            map[string]string{},
		sessions:              map[string]domain.Session{},
		adminAccounts:         map[string]domain.AdminAccount{},
		businesses:            map[string]domain.Business{},
		members:               map[string]domain.BusinessMember{},
		invitations:           map[string]domain.Invitation{},
		files:                 map[string]domain.FileObject{},
		notifications:         map[string]domain.Notification{},
		channels:              map[string]domain.Channel{},
		channelVaults:         map[string]domain.ChannelVault{},
		channelMembers:        map[string]domain.ChannelMember{},
		channelMessages:       map[string]domain.ChannelMessage{},
		channelVaultFiles:     map[string]domain.ChannelVaultFile{},
		areas:                 map[string]domain.Area{},
		streets:               map[string]domain.Street{},
		neighborhoods:         map[string]domain.Neighborhood{},
		propertyFiles:         map[string]domain.PropertyFile{},
		propertyShareRequests: map[string]domain.PropertyShareRequest{},
		contacts:              map[string]domain.Contact{},
		cities:                map[string]domain.City{},
		systemAreas:           map[string]domain.SystemArea{},
		systemStreets:         map[string]domain.SystemStreet{},
		systemNeighborhoods:   map[string]domain.SystemNeighborhood{},
		locationSuggestions:   map[string]domain.LocationSuggestion{},
	}
}

func (s *MemoryStore) SaveOTP(ctx context.Context, challenge domain.OTPChallenge) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.otps[challenge.Phone] = challenge
	s.latestOTP = challenge.Phone
	return nil
}

func (s *MemoryStore) GetOTP(ctx context.Context, phone string) (domain.OTPChallenge, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.otps[phone]
	if !ok {
		return domain.OTPChallenge{}, ErrNotFound
	}
	return item, nil
}

func (s *MemoryStore) GetLatestOTP(ctx context.Context) (domain.OTPChallenge, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.latestOTP == "" {
		return domain.OTPChallenge{}, ErrNotFound
	}
	item, ok := s.otps[s.latestOTP]
	if !ok {
		return domain.OTPChallenge{}, ErrNotFound
	}
	return item, nil
}

func (s *MemoryStore) DeleteOTP(ctx context.Context, phone string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.otps, phone)
	return nil
}

func (s *MemoryStore) UpsertUserByPhone(ctx context.Context, phone string) (domain.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now().UTC()
	if id, ok := s.usersPhone[phone]; ok {
		return s.users[id], nil
	}
	user := domain.User{
		ID:        support.NewID(),
		Phone:     phone,
		CreatedAt: now,
		UpdatedAt: now,
		PrivacySettings: domain.PrivacySettings{
			ShowPhoneToTeam:    true,
			ShowActivityStatus: true,
			AllowInviteByPhone: true,
		},
	}
	s.users[user.ID] = user
	s.usersPhone[phone] = user.ID
	return user, nil
}

func (s *MemoryStore) GetUser(ctx context.Context, id string) (domain.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	user, ok := s.users[id]
	if !ok {
		return domain.User{}, ErrNotFound
	}
	return user, nil
}

func (s *MemoryStore) GetUserByPhone(ctx context.Context, phone string) (domain.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	id, ok := s.usersPhone[phone]
	if !ok {
		return domain.User{}, ErrNotFound
	}
	return s.users[id], nil
}

func (s *MemoryStore) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	user.UpdatedAt = time.Now().UTC()
	s.users[user.ID] = user
	s.usersPhone[user.Phone] = user.ID
	return user, nil
}

func (s *MemoryStore) ListUsers(ctx context.Context) ([]domain.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]domain.User, 0, len(s.users))
	for _, user := range s.users {
		result = append(result, user)
	}
	return result, nil
}

func (s *MemoryStore) SaveSession(ctx context.Context, session domain.Session) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sessions[session.RefreshToken] = session
	return nil
}

func (s *MemoryStore) GetSessionByRefresh(ctx context.Context, refresh string) (domain.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	session, ok := s.sessions[refresh]
	if !ok || !session.RevokedAt.IsZero() {
		return domain.Session{}, ErrNotFound
	}
	return session, nil
}

func (s *MemoryStore) ListSessions(ctx context.Context, userID string) ([]domain.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	now := time.Now().UTC()
	var result []domain.Session
	for _, session := range s.sessions {
		if session.UserID == userID && session.RevokedAt.IsZero() && now.Before(session.ExpiresAt) {
			result = append(result, session)
		}
	}
	return result, nil
}

func (s *MemoryStore) GetSession(ctx context.Context, sessionID string) (domain.Session, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, session := range s.sessions {
		if session.ID == sessionID && session.RevokedAt.IsZero() && time.Now().UTC().Before(session.ExpiresAt) {
			return session, nil
		}
	}
	return domain.Session{}, ErrNotFound
}

func (s *MemoryStore) RevokeSession(ctx context.Context, refresh string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	session, ok := s.sessions[refresh]
	if !ok {
		return ErrNotFound
	}
	session.RevokedAt = time.Now().UTC()
	s.sessions[refresh] = session
	return nil
}

func (s *MemoryStore) RevokeSessionByID(ctx context.Context, userID string, sessionID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for refresh, session := range s.sessions {
		if session.ID == sessionID && session.UserID == userID {
			session.RevokedAt = time.Now().UTC()
			s.sessions[refresh] = session
			return nil
		}
	}
	return ErrNotFound
}

func (s *MemoryStore) TouchSession(ctx context.Context, sessionID string, at time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for refresh, session := range s.sessions {
		if session.ID == sessionID && session.RevokedAt.IsZero() {
			session.LastSeenAt = at
			s.sessions[refresh] = session
			return nil
		}
	}
	return ErrNotFound
}

func (s *MemoryStore) CreateBusiness(ctx context.Context, business domain.Business) (domain.Business, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	business.ID = support.NewID()
	business.CreatedAt = time.Now().UTC()
	business.UpdatedAt = business.CreatedAt
	if business.LicenseStatus == "" {
		business.LicenseStatus = "trial"
	}
	s.businesses[business.ID] = business
	return business, nil
}

func (s *MemoryStore) ListBusinessesForUser(ctx context.Context, userID string) ([]domain.Business, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []domain.Business
	for _, member := range s.members {
		if member.UserID == userID && member.Status == domain.MemberActive {
			if business, ok := s.businesses[member.BusinessID]; ok {
				result = append(result, business)
			}
		}
	}
	return result, nil
}

func (s *MemoryStore) GetBusiness(ctx context.Context, id string) (domain.Business, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.businesses[id]
	if !ok {
		return domain.Business{}, ErrNotFound
	}
	return item, nil
}

func (s *MemoryStore) UpdateBusiness(ctx context.Context, business domain.Business) (domain.Business, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	business.UpdatedAt = time.Now().UTC()
	s.businesses[business.ID] = business
	return business, nil
}

func (s *MemoryStore) ListBusinesses(ctx context.Context) ([]domain.Business, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]domain.Business, 0, len(s.businesses))
	for _, business := range s.businesses {
		result = append(result, business)
	}
	return result, nil
}

func (s *MemoryStore) CreateMember(ctx context.Context, member domain.BusinessMember) (domain.BusinessMember, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	member.ID = support.NewID()
	member.JoinedAt = time.Now().UTC()
	member.UpdatedAt = member.JoinedAt
	s.members[member.ID] = member
	return member, nil
}

func (s *MemoryStore) ListMembers(ctx context.Context, businessID string) ([]domain.BusinessMember, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []domain.BusinessMember
	for _, member := range s.members {
		if member.BusinessID == businessID && member.Status != domain.MemberRemoved {
			result = append(result, member)
		}
	}
	return result, nil
}

func (s *MemoryStore) GetMemberByUser(ctx context.Context, businessID string, userID string) (domain.BusinessMember, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, member := range s.members {
		if member.BusinessID == businessID && member.UserID == userID && member.Status != domain.MemberRemoved {
			return member, nil
		}
	}
	return domain.BusinessMember{}, ErrNotFound
}

func (s *MemoryStore) UpdateMember(ctx context.Context, member domain.BusinessMember) (domain.BusinessMember, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	member.UpdatedAt = time.Now().UTC()
	s.members[member.ID] = member
	return member, nil
}

func (s *MemoryStore) CreateInvitation(ctx context.Context, invitation domain.Invitation) (domain.Invitation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	invitation.ID = support.NewID()
	invitation.CreatedAt = time.Now().UTC()
	invitation.UpdatedAt = invitation.CreatedAt
	s.invitations[invitation.ID] = invitation
	return invitation, nil
}

func (s *MemoryStore) ListInvitations(ctx context.Context, businessID string) ([]domain.Invitation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []domain.Invitation
	for _, invitation := range s.invitations {
		if invitation.BusinessID == businessID {
			result = append(result, invitation)
		}
	}
	return result, nil
}

func (s *MemoryStore) ListInvitationsForPhone(ctx context.Context, phone string) ([]domain.Invitation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []domain.Invitation
	for _, invitation := range s.invitations {
		if invitation.InviteePhone == phone {
			result = append(result, invitation)
		}
	}
	return result, nil
}

func (s *MemoryStore) GetInvitation(ctx context.Context, id string) (domain.Invitation, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.invitations[id]
	if !ok {
		return domain.Invitation{}, ErrNotFound
	}
	return item, nil
}

func (s *MemoryStore) UpdateInvitation(ctx context.Context, invitation domain.Invitation) (domain.Invitation, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	invitation.UpdatedAt = time.Now().UTC()
	s.invitations[invitation.ID] = invitation
	return invitation, nil
}

func (s *MemoryStore) CreateFile(ctx context.Context, file domain.FileObject) (domain.FileObject, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if file.ID == "" {
		file.ID = support.NewID()
	}
	now := time.Now().UTC()
	if file.CreatedAt.IsZero() {
		file.CreatedAt = now
	}
	file.UpdatedAt = now
	s.files[file.ID] = file
	return file, nil
}

func (s *MemoryStore) GetFile(ctx context.Context, fileID string) (domain.FileObject, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	file, ok := s.files[fileID]
	if !ok {
		return domain.FileObject{}, ErrNotFound
	}
	return file, nil
}

func (s *MemoryStore) UpdateFile(ctx context.Context, file domain.FileObject) (domain.FileObject, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, ok := s.files[file.ID]
	if !ok {
		return domain.FileObject{}, ErrNotFound
	}
	if file.CreatedAt.IsZero() {
		file.CreatedAt = existing.CreatedAt
	}
	file.UpdatedAt = time.Now().UTC()
	s.files[file.ID] = file
	return file, nil
}

func (s *MemoryStore) ListExpiredFiles(ctx context.Context, before time.Time, limit int) ([]domain.FileObject, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if limit <= 0 {
		limit = 100
	}
	result := []domain.FileObject{}
	for _, file := range s.files {
		if file.ExpiresAt.IsZero() || file.ExpiresAt.After(before) {
			continue
		}
		result = append(result, file)
		if len(result) >= limit {
			break
		}
	}
	return result, nil
}

func (s *MemoryStore) DeleteFile(ctx context.Context, fileID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.files[fileID]; !ok {
		return ErrNotFound
	}
	delete(s.files, fileID)
	return nil
}

func (s *MemoryStore) CreateNotification(ctx context.Context, notification domain.Notification) (domain.Notification, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	notification.ID = support.NewID()
	notification.CreatedAt = time.Now().UTC()
	s.notifications[notification.ID] = notification
	return notification, nil
}

func (s *MemoryStore) ListNotifications(ctx context.Context, userID string, unreadOnly bool) ([]domain.Notification, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := []domain.Notification{}
	for _, notification := range s.notifications {
		if notification.UserID != userID {
			continue
		}
		if unreadOnly && !notification.ReadAt.IsZero() {
			continue
		}
		result = append(result, notification)
	}
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})
	return result, nil
}

func (s *MemoryStore) MarkNotificationRead(ctx context.Context, userID string, notificationID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	notification, ok := s.notifications[notificationID]
	if !ok || notification.UserID != userID {
		return ErrNotFound
	}
	notification.ReadAt = time.Now().UTC()
	s.notifications[notificationID] = notification
	return nil
}

func (s *MemoryStore) EnsureUserMainChannel(ctx context.Context, user domain.User) (domain.Channel, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, channel := range s.channels {
		if channel.Type == domain.ChannelTypeUserMain && channel.OwnerUserID == user.ID {
			return channel, nil
		}
	}
	now := time.Now().UTC()
	channel := domain.Channel{
		ID:          support.NewID(),
		Type:        domain.ChannelTypeUserMain,
		OwnerUserID: user.ID,
		Title:       user.DisplayName,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if channel.Title == "" {
		channel.Title = user.Phone
	}
	s.channels[channel.ID] = channel
	member := domain.ChannelMember{
		ID:        support.NewID(),
		ChannelID: channel.ID,
		UserID:    user.ID,
		Phone:     user.Phone,
		Status:    domain.ChannelMemberActive,
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.channelMembers[member.ID] = member
	return channel, nil
}

func (s *MemoryStore) EnsureBusinessMainChannel(ctx context.Context, business domain.Business) (domain.Channel, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, channel := range s.channels {
		if channel.Type == domain.ChannelTypeBusinessMain && channel.BusinessID == business.ID {
			s.ensureBusinessChannelMembersLocked(channel.ID, business.ID)
			return channel, nil
		}
	}
	now := time.Now().UTC()
	channel := domain.Channel{
		ID:         support.NewID(),
		Type:       domain.ChannelTypeBusinessMain,
		BusinessID: business.ID,
		Title:      business.Name,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	s.channels[channel.ID] = channel
	s.ensureBusinessChannelMembersLocked(channel.ID, business.ID)
	return channel, nil
}

func (s *MemoryStore) EnsurePrivateChannel(ctx context.Context, actor domain.User, target domain.User) (domain.Channel, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if actor.ID == "" || target.ID == "" || actor.ID == target.ID {
		return domain.Channel{}, ErrNotFound
	}
	for _, channel := range s.channels {
		if channel.Type != domain.ChannelTypePrivate {
			continue
		}
		hasActor := false
		hasTarget := false
		for _, member := range s.channelMembers {
			if member.ChannelID != channel.ID || member.Status != domain.ChannelMemberActive {
				continue
			}
			hasActor = hasActor || member.UserID == actor.ID
			hasTarget = hasTarget || member.UserID == target.ID
		}
		if hasActor && hasTarget {
			return channel, nil
		}
	}
	now := time.Now().UTC()
	channel := domain.Channel{
		ID:        support.NewID(),
		Type:      domain.ChannelTypePrivate,
		Title:     "private chat",
		CreatedAt: now,
		UpdatedAt: now,
	}
	s.channels[channel.ID] = channel
	actorMemberID := support.NewID()
	s.channelMembers[actorMemberID] = domain.ChannelMember{
		ID:        actorMemberID,
		ChannelID: channel.ID,
		UserID:    actor.ID,
		Phone:     actor.Phone,
		Role:      domain.ChannelMemberRoleMember,
		Status:    domain.ChannelMemberActive,
		CreatedAt: now,
		UpdatedAt: now,
	}
	targetMemberID := support.NewID()
	s.channelMembers[targetMemberID] = domain.ChannelMember{
		ID:        targetMemberID,
		ChannelID: channel.ID,
		UserID:    target.ID,
		Phone:     target.Phone,
		Role:      domain.ChannelMemberRoleMember,
		Status:    domain.ChannelMemberActive,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return channel, nil
}

func (s *MemoryStore) EnsureUserMainVault(ctx context.Context, user domain.User) (domain.ChannelVault, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, vault := range s.channelVaults {
		if vault.OwnerUserID == user.ID && vault.IsMain {
			return vault, nil
		}
	}
	channel := s.ensureUserMainChannelLocked(user)
	return s.createVaultLocked(user.ID, "", channel.ID, channel.Title, true), nil
}

func (s *MemoryStore) EnsureBusinessMainVault(ctx context.Context, business domain.Business) (domain.ChannelVault, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, vault := range s.channelVaults {
		if vault.BusinessID == business.ID && vault.IsMain {
			return vault, nil
		}
	}
	channel := s.ensureBusinessMainChannelLocked(business)
	return s.createVaultLocked("", business.ID, channel.ID, channel.Title, true), nil
}

func (s *MemoryStore) CreateUserVault(ctx context.Context, user domain.User, title string) (domain.ChannelVault, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now().UTC()
	title = strings.TrimSpace(title)
	if title == "" {
		title = "صندوقچه"
	}
	channel := domain.Channel{
		ID:          support.NewID(),
		Type:        domain.ChannelTypeUserVault,
		OwnerUserID: user.ID,
		Title:       title,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	s.channels[channel.ID] = channel
	vault := s.createVaultLocked(user.ID, "", channel.ID, title, false)
	channel.VaultID = vault.ID
	s.channels[channel.ID] = channel
	member := domain.ChannelMember{ID: support.NewID(), ChannelID: channel.ID, UserID: user.ID, Phone: user.Phone, Status: domain.ChannelMemberActive, CreatedAt: now, UpdatedAt: now}
	s.channelMembers[member.ID] = member
	return vault, nil
}

func (s *MemoryStore) CreateBusinessVault(ctx context.Context, business domain.Business, title string) (domain.ChannelVault, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now().UTC()
	title = strings.TrimSpace(title)
	if title == "" {
		title = "صندوقچه"
	}
	channel := domain.Channel{
		ID:         support.NewID(),
		Type:       domain.ChannelTypeBusinessVault,
		BusinessID: business.ID,
		Title:      title,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	s.channels[channel.ID] = channel
	vault := s.createVaultLocked("", business.ID, channel.ID, title, false)
	channel.VaultID = vault.ID
	s.channels[channel.ID] = channel
	s.ensureBusinessChannelMembersLocked(channel.ID, business.ID)
	return vault, nil
}

func (s *MemoryStore) GetChannelVault(ctx context.Context, vaultID string) (domain.ChannelVault, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	vault, ok := s.channelVaults[vaultID]
	if !ok {
		return domain.ChannelVault{}, ErrNotFound
	}
	return vault, nil
}

func (s *MemoryStore) ListUserVaults(ctx context.Context, userID string) ([]domain.ChannelVault, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := []domain.ChannelVault{}
	for _, vault := range s.channelVaults {
		if vault.OwnerUserID == userID {
			result = append(result, vault)
		}
	}
	return result, nil
}

func (s *MemoryStore) ListBusinessVaults(ctx context.Context, businessID string) ([]domain.ChannelVault, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := []domain.ChannelVault{}
	for _, vault := range s.channelVaults {
		if vault.BusinessID == businessID {
			result = append(result, vault)
		}
	}
	return result, nil
}

func (s *MemoryStore) ensureUserMainChannelLocked(user domain.User) domain.Channel {
	for _, channel := range s.channels {
		if channel.Type == domain.ChannelTypeUserMain && channel.OwnerUserID == user.ID {
			return channel
		}
	}
	now := time.Now().UTC()
	title := user.DisplayName
	if title == "" {
		title = user.Phone
	}
	channel := domain.Channel{ID: support.NewID(), Type: domain.ChannelTypeUserMain, OwnerUserID: user.ID, Title: title, CreatedAt: now, UpdatedAt: now}
	s.channels[channel.ID] = channel
	member := domain.ChannelMember{ID: support.NewID(), ChannelID: channel.ID, UserID: user.ID, Phone: user.Phone, Status: domain.ChannelMemberActive, CreatedAt: now, UpdatedAt: now}
	s.channelMembers[member.ID] = member
	return channel
}

func (s *MemoryStore) ensureBusinessMainChannelLocked(business domain.Business) domain.Channel {
	for _, channel := range s.channels {
		if channel.Type == domain.ChannelTypeBusinessMain && channel.BusinessID == business.ID {
			s.ensureBusinessChannelMembersLocked(channel.ID, business.ID)
			return channel
		}
	}
	now := time.Now().UTC()
	channel := domain.Channel{ID: support.NewID(), Type: domain.ChannelTypeBusinessMain, BusinessID: business.ID, Title: business.Name, CreatedAt: now, UpdatedAt: now}
	s.channels[channel.ID] = channel
	s.ensureBusinessChannelMembersLocked(channel.ID, business.ID)
	return channel
}

func (s *MemoryStore) createVaultLocked(ownerUserID, businessID, channelID, title string, isMain bool) domain.ChannelVault {
	now := time.Now().UTC()
	vault := domain.ChannelVault{ID: support.NewID(), ChannelID: channelID, OwnerUserID: ownerUserID, BusinessID: businessID, Title: title, IsMain: isMain, CreatedAt: now, UpdatedAt: now}
	s.channelVaults[vault.ID] = vault
	channel := s.channels[channelID]
	channel.VaultID = vault.ID
	s.channels[channelID] = channel
	return vault
}

func (s *MemoryStore) ensureBusinessChannelMembersLocked(channelID, businessID string) {
	now := time.Now().UTC()
	for _, member := range s.members {
		if member.BusinessID != businessID || member.Status != domain.MemberActive {
			continue
		}
		exists := false
		for _, channelMember := range s.channelMembers {
			if channelMember.ChannelID == channelID && channelMember.UserID == member.UserID {
				exists = true
				break
			}
		}
		if exists {
			continue
		}
		id := support.NewID()
		s.channelMembers[id] = domain.ChannelMember{
			ID:        id,
			ChannelID: channelID,
			UserID:    member.UserID,
			Phone:     member.UserPhone,
			Status:    domain.ChannelMemberActive,
			CreatedAt: now,
			UpdatedAt: now,
		}
	}
}

func (s *MemoryStore) GetChannel(ctx context.Context, channelID string) (domain.Channel, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	channel, ok := s.channels[channelID]
	if !ok {
		return domain.Channel{}, ErrNotFound
	}
	return channel, nil
}

func (s *MemoryStore) ListChannelsForUser(ctx context.Context, userID string, phones []string, businessIDs []string) ([]domain.Channel, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	allowed := map[string]struct{}{}
	for _, member := range s.channelMembers {
		if member.Status != domain.ChannelMemberActive {
			continue
		}
		if member.UserID == userID {
			allowed[member.ChannelID] = struct{}{}
			continue
		}
		for _, phone := range phones {
			if phone != "" && member.Phone == phone {
				allowed[member.ChannelID] = struct{}{}
			}
		}
	}
	for _, businessID := range businessIDs {
		for _, channel := range s.channels {
			if channel.Type == domain.ChannelTypeBusinessMain && channel.BusinessID == businessID {
				allowed[channel.ID] = struct{}{}
			}
		}
	}
	result := []domain.Channel{}
	for id := range allowed {
		if channel, ok := s.channels[id]; ok {
			result = append(result, channel)
		}
	}
	return result, nil
}

func (s *MemoryStore) UpsertChannelMember(ctx context.Context, member domain.ChannelMember) (domain.ChannelMember, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now().UTC()
	if member.Role == "" {
		member.Role = domain.ChannelMemberRoleMember
	}
	if member.Status == "" {
		member.Status = domain.ChannelMemberActive
	}
	for id, existing := range s.channelMembers {
		sameUser := member.UserID != "" && existing.UserID == member.UserID
		samePhone := member.Phone != "" && existing.Phone == member.Phone
		if existing.ChannelID == member.ChannelID && (sameUser || samePhone) {
			member.ID = id
			member.CreatedAt = existing.CreatedAt
			member.UpdatedAt = now
			s.channelMembers[id] = member
			return member, nil
		}
	}
	member.ID = support.NewID()
	member.CreatedAt = now
	member.UpdatedAt = now
	s.channelMembers[member.ID] = member
	return member, nil
}

func (s *MemoryStore) GetChannelMember(ctx context.Context, channelID string, userID string) (domain.ChannelMember, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, member := range s.channelMembers {
		if member.ChannelID == channelID && member.UserID == userID && member.Status == domain.ChannelMemberActive {
			return member, nil
		}
	}
	return domain.ChannelMember{}, ErrNotFound
}

func (s *MemoryStore) ListChannelMembers(ctx context.Context, channelID string) ([]domain.ChannelMember, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := []domain.ChannelMember{}
	for _, member := range s.channelMembers {
		if member.ChannelID == channelID {
			result = append(result, member)
		}
	}
	return result, nil
}

func (s *MemoryStore) CreateChannelMessage(ctx context.Context, message domain.ChannelMessage) (domain.ChannelMessage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now().UTC()
	message.ID = support.NewID()
	message.CreatedAt = now
	message.UpdatedAt = now
	s.channelMessages[message.ID] = message
	return message, nil
}

func (s *MemoryStore) GetChannelMessage(ctx context.Context, channelID string, messageID string) (domain.ChannelMessage, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	message, ok := s.channelMessages[messageID]
	if !ok || message.ChannelID != channelID {
		return domain.ChannelMessage{}, ErrNotFound
	}
	return message, nil
}

func (s *MemoryStore) UpdateChannelMessage(ctx context.Context, message domain.ChannelMessage) (domain.ChannelMessage, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, ok := s.channelMessages[message.ID]
	if !ok || existing.ChannelID != message.ChannelID {
		return domain.ChannelMessage{}, ErrNotFound
	}
	message.CreatedAt = existing.CreatedAt
	message.UpdatedAt = time.Now().UTC()
	s.channelMessages[message.ID] = message
	return message, nil
}

func (s *MemoryStore) DeleteChannelMessage(ctx context.Context, channelID string, messageID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	message, ok := s.channelMessages[messageID]
	if !ok || message.ChannelID != channelID {
		return ErrNotFound
	}
	delete(s.channelMessages, messageID)
	return nil
}

func (s *MemoryStore) ListChannelMessages(ctx context.Context, channelID string, limit int, offset int) ([]domain.ChannelMessage, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	messages := []domain.ChannelMessage{}
	for _, message := range s.channelMessages {
		if message.ChannelID == channelID {
			messages = append(messages, message)
		}
	}
	sort.SliceStable(messages, func(i, j int) bool {
		return messages[i].CreatedAt.After(messages[j].CreatedAt)
	})
	total := len(messages)
	if offset >= total {
		return []domain.ChannelMessage{}, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return messages[offset:end], total, nil
}

func (s *MemoryStore) ChannelUnreadSummary(ctx context.Context, channelID string, userID string) (domain.ChannelUnreadSummary, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	messages := []domain.ChannelMessage{}
	for _, message := range s.channelMessages {
		if message.ChannelID == channelID {
			messages = append(messages, message)
		}
	}
	sort.SliceStable(messages, func(i, j int) bool {
		return messages[i].CreatedAt.After(messages[j].CreatedAt)
	})
	summary := domain.ChannelUnreadSummary{}
	for index, message := range messages {
		if message.AuthorID == userID || messageSeenByUser(message, userID) {
			continue
		}
		summary.UnreadCount++
		summary.FirstUnreadOffset = index
		summary.FirstUnreadMessageID = message.ID
	}
	return summary, nil
}

func (s *MemoryStore) MarkChannelMessagesSeen(ctx context.Context, channelID string, userID string, messageIDs []string) error {
	if userID == "" || len(messageIDs) == 0 {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	ids := map[string]struct{}{}
	for _, id := range messageIDs {
		if id != "" {
			ids[id] = struct{}{}
		}
	}
	if len(ids) == 0 {
		return nil
	}
	now := time.Now().UTC()
	for id := range ids {
		message, ok := s.channelMessages[id]
		if !ok || message.ChannelID != channelID || message.AuthorID == userID {
			continue
		}
		alreadySeen := false
		for _, seen := range message.SeenBy {
			if seen.UserID == userID {
				alreadySeen = true
				break
			}
		}
		if alreadySeen {
			continue
		}
		message.SeenBy = append(message.SeenBy, domain.ChannelMessageSeen{UserID: userID, SeenAt: now})
		message.UpdatedAt = now
		s.channelMessages[id] = message
	}
	return nil
}

func messageSeenByUser(message domain.ChannelMessage, userID string) bool {
	for _, seen := range message.SeenBy {
		if seen.UserID == userID {
			return true
		}
	}
	return false
}

func (s *MemoryStore) CreateChannelVaultFile(ctx context.Context, file domain.ChannelVaultFile) (domain.ChannelVaultFile, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now().UTC()
	file.ID = support.NewID()
	file.CreatedAt = now
	file.UpdatedAt = now
	s.channelVaultFiles[file.ID] = file
	return file, nil
}

func (s *MemoryStore) UpsertChannelVaultFile(ctx context.Context, file domain.ChannelVaultFile) (domain.ChannelVaultFile, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now().UTC()
	for id, existing := range s.channelVaultFiles {
		sameProperty := file.PropertyFileID != "" && existing.PropertyFileID == file.PropertyFileID
		sameFile := file.FileID != "" && existing.FileID == file.FileID
		if existing.ChannelID == file.ChannelID && existing.VaultID == file.VaultID && (sameProperty || sameFile) {
			file.ID = id
			file.CreatedAt = existing.CreatedAt
			file.UpdatedAt = now
			s.channelVaultFiles[id] = file
			return file, nil
		}
	}
	file.ID = support.NewID()
	file.CreatedAt = now
	file.UpdatedAt = now
	s.channelVaultFiles[file.ID] = file
	return file, nil
}

func (s *MemoryStore) DeletePropertyVaultReferences(ctx context.Context, propertyFileID string, keepVaultIDs []string) error {
	if strings.TrimSpace(propertyFileID) == "" {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	keep := map[string]struct{}{}
	for _, vaultID := range keepVaultIDs {
		vaultID = strings.TrimSpace(vaultID)
		if vaultID != "" {
			keep[vaultID] = struct{}{}
		}
	}
	for id, file := range s.channelVaultFiles {
		if file.PropertyFileID != propertyFileID {
			continue
		}
		if _, ok := keep[file.VaultID]; ok {
			continue
		}
		delete(s.channelVaultFiles, id)
	}
	return nil
}

func (s *MemoryStore) GetChannelVaultFile(ctx context.Context, channelID string, fileID string) (domain.ChannelVaultFile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	file, ok := s.channelVaultFiles[fileID]
	if !ok || file.ChannelID != channelID {
		return domain.ChannelVaultFile{}, ErrNotFound
	}
	return file, nil
}

func (s *MemoryStore) ListChannelVaultFiles(ctx context.Context, channelID string, limit int, offset int) ([]domain.ChannelVaultFile, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	files := []domain.ChannelVaultFile{}
	for _, file := range s.channelVaultFiles {
		if file.ChannelID == channelID {
			files = append(files, file)
		}
	}
	sort.SliceStable(files, func(i, j int) bool {
		return files[i].CreatedAt.After(files[j].CreatedAt)
	})
	total := len(files)
	if offset >= total {
		return []domain.ChannelVaultFile{}, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return files[offset:end], total, nil
}

func (s *MemoryStore) CreateArea(ctx context.Context, area domain.Area) (domain.Area, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	area.ID = support.NewID()
	area.CreatedAt = time.Now().UTC()
	area.UpdatedAt = area.CreatedAt
	s.areas[area.ID] = area
	return area, nil
}

func (s *MemoryStore) ListAreas(ctx context.Context, businessID string) ([]domain.Area, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []domain.Area
	for _, area := range s.areas {
		if area.BusinessID == businessID {
			result = append(result, area)
		}
	}
	return result, nil
}

func (s *MemoryStore) DeleteArea(ctx context.Context, businessID string, areaID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	area, ok := s.areas[areaID]
	if !ok || area.BusinessID != businessID {
		return ErrNotFound
	}
	delete(s.areas, areaID)
	for streetID, street := range s.streets {
		if street.BusinessID == businessID && street.AreaID == areaID {
			delete(s.streets, streetID)
		}
	}
	for neighborhoodID, neighborhood := range s.neighborhoods {
		if neighborhood.BusinessID == businessID && neighborhood.AreaID == areaID {
			delete(s.neighborhoods, neighborhoodID)
		}
	}
	return nil
}

func (s *MemoryStore) CreateStreet(ctx context.Context, street domain.Street) (domain.Street, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	area, ok := s.areas[street.AreaID]
	if !ok || area.BusinessID != street.BusinessID {
		return domain.Street{}, ErrNotFound
	}
	if area.CityID != "" && street.CityID != "" && area.CityID != street.CityID {
		return domain.Street{}, ErrNotFound
	}
	street.ID = support.NewID()
	street.CreatedAt = time.Now().UTC()
	street.UpdatedAt = street.CreatedAt
	s.streets[street.ID] = street
	return street, nil
}

func (s *MemoryStore) ListStreets(ctx context.Context, businessID string, areaID string) ([]domain.Street, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []domain.Street
	for _, street := range s.streets {
		if street.BusinessID == businessID && street.AreaID == areaID {
			result = append(result, street)
		}
	}
	return result, nil
}

func (s *MemoryStore) DeleteStreet(ctx context.Context, businessID string, areaID string, streetID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	street, ok := s.streets[streetID]
	if !ok || street.BusinessID != businessID || street.AreaID != areaID {
		return ErrNotFound
	}
	delete(s.streets, streetID)
	for neighborhoodID, neighborhood := range s.neighborhoods {
		if neighborhood.BusinessID == businessID && neighborhood.AreaID == areaID && neighborhood.StreetID == streetID {
			delete(s.neighborhoods, neighborhoodID)
		}
	}
	return nil
}

func (s *MemoryStore) CreateNeighborhood(ctx context.Context, neighborhood domain.Neighborhood) (domain.Neighborhood, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	street, ok := s.streets[neighborhood.StreetID]
	if !ok || street.BusinessID != neighborhood.BusinessID || street.AreaID != neighborhood.AreaID {
		return domain.Neighborhood{}, ErrNotFound
	}
	if street.CityID != "" && neighborhood.CityID != "" && street.CityID != neighborhood.CityID {
		return domain.Neighborhood{}, ErrNotFound
	}
	neighborhood.ID = support.NewID()
	neighborhood.CreatedAt = time.Now().UTC()
	neighborhood.UpdatedAt = neighborhood.CreatedAt
	s.neighborhoods[neighborhood.ID] = neighborhood
	return neighborhood, nil
}

func (s *MemoryStore) ListNeighborhoods(ctx context.Context, businessID string, areaID string, streetID string) ([]domain.Neighborhood, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []domain.Neighborhood
	for _, neighborhood := range s.neighborhoods {
		if neighborhood.BusinessID == businessID && neighborhood.AreaID == areaID && neighborhood.StreetID == streetID {
			result = append(result, neighborhood)
		}
	}
	return result, nil
}

func (s *MemoryStore) DeleteNeighborhood(ctx context.Context, businessID string, areaID string, streetID string, neighborhoodID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	neighborhood, ok := s.neighborhoods[neighborhoodID]
	if !ok || neighborhood.BusinessID != businessID || neighborhood.AreaID != areaID || neighborhood.StreetID != streetID {
		return ErrNotFound
	}
	delete(s.neighborhoods, neighborhoodID)
	return nil
}

func (s *MemoryStore) CreatePropertyFile(ctx context.Context, file domain.PropertyFile) (domain.PropertyFile, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	file.ID = support.NewID()
	file.CreatedAt = time.Now().UTC()
	file.UpdatedAt = file.CreatedAt
	s.propertyFiles[file.ID] = file
	return file, nil
}

func (s *MemoryStore) ListPropertyFiles(ctx context.Context, businessID string) ([]domain.PropertyFile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []domain.PropertyFile
	for _, file := range s.propertyFiles {
		if file.BusinessID == businessID {
			result = append(result, file)
		}
	}
	return result, nil
}

func (s *MemoryStore) ListPropertyFilesForOwner(ctx context.Context, businessID string, ownerUserID string) ([]domain.PropertyFile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []domain.PropertyFile
	for _, file := range s.propertyFiles {
		if file.BusinessID == businessID && file.OwnerUserID == ownerUserID {
			result = append(result, file)
		}
	}
	return result, nil
}

func (s *MemoryStore) ListPropertyFilesForAccess(ctx context.Context, businessID string, ownerUserID string, vaultIDs []string) ([]domain.PropertyFile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	accessibleVaults := map[string]struct{}{}
	for _, vaultID := range vaultIDs {
		if vaultID != "" {
			accessibleVaults[vaultID] = struct{}{}
		}
	}
	var result []domain.PropertyFile
	for _, file := range s.propertyFiles {
		if file.BusinessID != businessID {
			continue
		}
		if file.OwnerUserID == ownerUserID || propertyInVaults(file, accessibleVaults) {
			result = append(result, file)
		}
	}
	return result, nil
}

func (s *MemoryStore) ListLatestPropertyFiles(ctx context.Context, businessID string, ownerUserID string, vaultIDs []string, propertyType domain.PropertyFileType, limit int, offset int) ([]domain.PropertyFile, int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if limit <= 0 {
		limit = 30
	}
	if offset < 0 {
		offset = 0
	}
	result := []domain.PropertyFile{}
	accessibleVaults := map[string]struct{}{}
	for _, vaultID := range vaultIDs {
		if vaultID != "" {
			accessibleVaults[vaultID] = struct{}{}
		}
	}
	for _, file := range s.propertyFiles {
		if file.BusinessID != businessID {
			continue
		}
		visible := file.OwnerUserID == ownerUserID
		if !visible {
			for _, vaultID := range file.VaultIDs {
				if _, ok := accessibleVaults[vaultID]; ok {
					visible = true
					break
				}
			}
		}
		if !visible {
			continue
		}
		if propertyType != "" && !propertyFileHasType(file, propertyType) {
			continue
		}
		result = append(result, file)
	}
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].CreatedAt.After(result[j].CreatedAt)
	})
	total := len(result)
	if offset >= total {
		return []domain.PropertyFile{}, total, nil
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return result[offset:end], total, nil
}

func (s *MemoryStore) GetPropertyFile(ctx context.Context, businessID string, fileID string) (domain.PropertyFile, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	file, ok := s.propertyFiles[fileID]
	if !ok || file.BusinessID != businessID {
		return domain.PropertyFile{}, ErrNotFound
	}
	return file, nil
}

func (s *MemoryStore) UpdatePropertyFile(ctx context.Context, file domain.PropertyFile) (domain.PropertyFile, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.propertyFiles[file.ID]; !ok {
		return domain.PropertyFile{}, ErrNotFound
	}
	file.UpdatedAt = time.Now().UTC()
	s.propertyFiles[file.ID] = file
	return file, nil
}

func propertyFileHasType(file domain.PropertyFile, propertyType domain.PropertyFileType) bool {
	if file.Type == propertyType {
		return true
	}
	for _, item := range file.Types {
		if item == propertyType {
			return true
		}
	}
	return false
}

func propertyInVaults(file domain.PropertyFile, vaultIDs map[string]struct{}) bool {
	if len(vaultIDs) == 0 {
		return false
	}
	for _, vaultID := range file.VaultIDs {
		if _, ok := vaultIDs[vaultID]; ok {
			return true
		}
	}
	return false
}

func (s *MemoryStore) CreatePropertyShareRequest(ctx context.Context, request domain.PropertyShareRequest) (domain.PropertyShareRequest, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now().UTC()
	request.ID = support.NewID()
	request.CreatedAt = now
	request.UpdatedAt = now
	s.propertyShareRequests[request.ID] = request
	return request, nil
}

func (s *MemoryStore) GetPropertyShareRequest(ctx context.Context, businessID string, requestID string) (domain.PropertyShareRequest, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	request, ok := s.propertyShareRequests[requestID]
	if !ok || request.BusinessID != businessID {
		return domain.PropertyShareRequest{}, ErrNotFound
	}
	return request, nil
}

func (s *MemoryStore) UpdatePropertyShareRequest(ctx context.Context, request domain.PropertyShareRequest) (domain.PropertyShareRequest, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.propertyShareRequests[request.ID]; !ok {
		return domain.PropertyShareRequest{}, ErrNotFound
	}
	request.UpdatedAt = time.Now().UTC()
	s.propertyShareRequests[request.ID] = request
	return request, nil
}

func (s *MemoryStore) ListPropertyShareRequestsForOwner(ctx context.Context, businessID string, ownerUserID string) ([]domain.PropertyShareRequest, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := []domain.PropertyShareRequest{}
	for _, request := range s.propertyShareRequests {
		if request.BusinessID == businessID && request.OwnerUserID == ownerUserID {
			result = append(result, request)
		}
	}
	sort.SliceStable(result, func(i, j int) bool { return result[i].CreatedAt.After(result[j].CreatedAt) })
	return result, nil
}

func (s *MemoryStore) ListPropertyShareRequestsForRequester(ctx context.Context, businessID string, requesterUserID string) ([]domain.PropertyShareRequest, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := []domain.PropertyShareRequest{}
	for _, request := range s.propertyShareRequests {
		if request.BusinessID == businessID && request.RequesterUserID == requesterUserID {
			result = append(result, request)
		}
	}
	sort.SliceStable(result, func(i, j int) bool { return result[i].CreatedAt.After(result[j].CreatedAt) })
	return result, nil
}

func (s *MemoryStore) CreateContact(ctx context.Context, contact domain.Contact) (domain.Contact, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	contact.ID = support.NewID()
	contact.CreatedAt = time.Now().UTC()
	contact.UpdatedAt = contact.CreatedAt
	for i := range contact.Requests {
		if contact.Requests[i].ID == "" {
			contact.Requests[i].ID = support.NewID()
		}
		if contact.Requests[i].CreatedAt.IsZero() {
			contact.Requests[i].CreatedAt = contact.CreatedAt
		}
	}
	s.contacts[contact.ID] = contact
	return contact, nil
}

func (s *MemoryStore) ListContacts(ctx context.Context, businessID string) ([]domain.Contact, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []domain.Contact
	for _, contact := range s.contacts {
		if contact.BusinessID == businessID {
			result = append(result, contact)
		}
	}
	return result, nil
}

func (s *MemoryStore) GetContact(ctx context.Context, businessID, contactID string) (domain.Contact, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	contact, ok := s.contacts[contactID]
	if !ok || contact.BusinessID != businessID {
		return domain.Contact{}, ErrNotFound
	}
	return contact, nil
}

func (s *MemoryStore) UpdateContact(ctx context.Context, contact domain.Contact) (domain.Contact, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	existing, ok := s.contacts[contact.ID]
	if !ok || existing.BusinessID != contact.BusinessID {
		return domain.Contact{}, ErrNotFound
	}
	contact.CreatedAt = existing.CreatedAt
	contact.UpdatedAt = time.Now().UTC()
	for i := range contact.Requests {
		if contact.Requests[i].ID == "" {
			contact.Requests[i].ID = support.NewID()
		}
		if contact.Requests[i].CreatedAt.IsZero() {
			contact.Requests[i].CreatedAt = contact.UpdatedAt
		}
	}
	s.contacts[contact.ID] = contact
	return contact, nil
}

func (s *MemoryStore) SaveAdminAccount(ctx context.Context, account domain.AdminAccount) (domain.AdminAccount, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now().UTC()
	if existing, ok := s.adminAccounts[account.UserID]; ok {
		if account.ID == "" {
			account.ID = existing.ID
		}
		if account.CreatedAt.IsZero() {
			account.CreatedAt = existing.CreatedAt
		}
	}
	if account.ID == "" {
		account.ID = support.NewID()
		account.CreatedAt = now
	}
	if account.Status == "" {
		account.Status = "active"
	}
	account.UpdatedAt = now
	s.adminAccounts[account.UserID] = account
	return account, nil
}

func (s *MemoryStore) GetAdminAccountByUser(ctx context.Context, userID string) (domain.AdminAccount, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	account, ok := s.adminAccounts[userID]
	if !ok {
		return domain.AdminAccount{}, ErrNotFound
	}
	return account, nil
}

func (s *MemoryStore) ListAdminAccounts(ctx context.Context) ([]domain.AdminAccount, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]domain.AdminAccount, 0, len(s.adminAccounts))
	for _, account := range s.adminAccounts {
		result = append(result, account)
	}
	return result, nil
}

func (s *MemoryStore) GetPlatformSettings(ctx context.Context) (domain.PlatformSettings, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.platformSettings.ID == "" {
		return domain.PlatformSettings{ID: "platform"}, nil
	}
	return s.platformSettings, nil
}

func (s *MemoryStore) SavePlatformSettings(ctx context.Context, settings domain.PlatformSettings) (domain.PlatformSettings, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if settings.ID == "" {
		settings.ID = "platform"
	}
	settings.UpdatedAt = time.Now().UTC()
	s.platformSettings = settings
	return settings, nil
}

func (s *MemoryStore) CreateCity(ctx context.Context, city domain.City) (domain.City, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	city.ID = support.NewID()
	city.CreatedAt = time.Now().UTC()
	city.UpdatedAt = city.CreatedAt
	if city.Status == "" {
		city.Status = domain.SystemLocationActive
	}
	s.cities[city.ID] = city
	return city, nil
}

func (s *MemoryStore) ListCities(ctx context.Context) ([]domain.City, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]domain.City, 0, len(s.cities))
	for _, city := range s.cities {
		if city.Status == domain.SystemLocationActive {
			result = append(result, city)
		}
	}
	return result, nil
}

func (s *MemoryStore) GetCity(ctx context.Context, cityID string) (domain.City, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	city, ok := s.cities[cityID]
	if !ok {
		return domain.City{}, ErrNotFound
	}
	return city, nil
}

func (s *MemoryStore) CreateSystemArea(ctx context.Context, area domain.SystemArea) (domain.SystemArea, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.cities[area.CityID]; !ok {
		return domain.SystemArea{}, ErrNotFound
	}
	area.ID = support.NewID()
	area.CreatedAt = time.Now().UTC()
	area.UpdatedAt = area.CreatedAt
	if area.Status == "" {
		area.Status = domain.SystemLocationActive
	}
	s.systemAreas[area.ID] = area
	return area, nil
}

func (s *MemoryStore) ListSystemAreas(ctx context.Context, cityID string) ([]domain.SystemArea, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []domain.SystemArea
	for _, area := range s.systemAreas {
		if area.CityID == cityID && area.Status == domain.SystemLocationActive {
			result = append(result, area)
		}
	}
	return result, nil
}

func (s *MemoryStore) GetSystemArea(ctx context.Context, areaID string) (domain.SystemArea, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	area, ok := s.systemAreas[areaID]
	if !ok {
		return domain.SystemArea{}, ErrNotFound
	}
	return area, nil
}

func (s *MemoryStore) CreateSystemStreet(ctx context.Context, street domain.SystemStreet) (domain.SystemStreet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	area, ok := s.systemAreas[street.AreaID]
	if !ok || area.CityID != street.CityID {
		return domain.SystemStreet{}, ErrNotFound
	}
	street.ID = support.NewID()
	street.CreatedAt = time.Now().UTC()
	street.UpdatedAt = street.CreatedAt
	if street.Status == "" {
		street.Status = domain.SystemLocationActive
	}
	s.systemStreets[street.ID] = street
	return street, nil
}

func (s *MemoryStore) ListSystemStreets(ctx context.Context, cityID, areaID string) ([]domain.SystemStreet, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []domain.SystemStreet
	for _, street := range s.systemStreets {
		if street.CityID == cityID && street.AreaID == areaID && street.Status == domain.SystemLocationActive {
			result = append(result, street)
		}
	}
	return result, nil
}

func (s *MemoryStore) GetSystemStreet(ctx context.Context, streetID string) (domain.SystemStreet, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	street, ok := s.systemStreets[streetID]
	if !ok {
		return domain.SystemStreet{}, ErrNotFound
	}
	return street, nil
}

func (s *MemoryStore) CreateSystemNeighborhood(ctx context.Context, neighborhood domain.SystemNeighborhood) (domain.SystemNeighborhood, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	street, ok := s.systemStreets[neighborhood.StreetID]
	if !ok || street.CityID != neighborhood.CityID || street.AreaID != neighborhood.AreaID {
		return domain.SystemNeighborhood{}, ErrNotFound
	}
	neighborhood.ID = support.NewID()
	neighborhood.CreatedAt = time.Now().UTC()
	neighborhood.UpdatedAt = neighborhood.CreatedAt
	if neighborhood.Status == "" {
		neighborhood.Status = domain.SystemLocationActive
	}
	s.systemNeighborhoods[neighborhood.ID] = neighborhood
	return neighborhood, nil
}

func (s *MemoryStore) ListSystemNeighborhoods(ctx context.Context, cityID, areaID, streetID string) ([]domain.SystemNeighborhood, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []domain.SystemNeighborhood
	for _, neighborhood := range s.systemNeighborhoods {
		if neighborhood.CityID == cityID && neighborhood.AreaID == areaID && neighborhood.StreetID == streetID && neighborhood.Status == domain.SystemLocationActive {
			result = append(result, neighborhood)
		}
	}
	return result, nil
}

func (s *MemoryStore) GetSystemNeighborhood(ctx context.Context, neighborhoodID string) (domain.SystemNeighborhood, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	neighborhood, ok := s.systemNeighborhoods[neighborhoodID]
	if !ok {
		return domain.SystemNeighborhood{}, ErrNotFound
	}
	return neighborhood, nil
}

func (s *MemoryStore) CreateLocationSuggestion(ctx context.Context, suggestion domain.LocationSuggestion) (domain.LocationSuggestion, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	suggestion.ID = support.NewID()
	suggestion.CreatedAt = time.Now().UTC()
	if suggestion.Status == "" {
		suggestion.Status = domain.LocationSuggestionPending
	}
	s.locationSuggestions[suggestion.ID] = suggestion
	return suggestion, nil
}

func (s *MemoryStore) ListLocationSuggestions(ctx context.Context, status domain.LocationSuggestionStatus) ([]domain.LocationSuggestion, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var result []domain.LocationSuggestion
	for _, suggestion := range s.locationSuggestions {
		if status == "" || suggestion.Status == status {
			result = append(result, suggestion)
		}
	}
	return result, nil
}

func (s *MemoryStore) GetLocationSuggestion(ctx context.Context, suggestionID string) (domain.LocationSuggestion, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	suggestion, ok := s.locationSuggestions[suggestionID]
	if !ok {
		return domain.LocationSuggestion{}, ErrNotFound
	}
	return suggestion, nil
}

func (s *MemoryStore) UpdateLocationSuggestion(ctx context.Context, suggestion domain.LocationSuggestion) (domain.LocationSuggestion, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.locationSuggestions[suggestion.ID]; !ok {
		return domain.LocationSuggestion{}, ErrNotFound
	}
	s.locationSuggestions[suggestion.ID] = suggestion
	return suggestion, nil
}
