package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"amlakcrm/backend/internal/domain"
	"amlakcrm/backend/internal/repository"
	"amlakcrm/backend/internal/searchengine"
	"amlakcrm/backend/internal/support"
)

const (
	maxPropertyMedia       = 20
	maxPropertyVideos      = 2
	maxImageBytes          = 500 * 1024
	maxVideoDurationSec    = 120
	maxVideoUploadBytes    = 200 * 1024 * 1024
	defaultObjectStorePath = "storage/objects"
)

type PropertyService struct {
	store     repository.Store
	objectDir string
}

func NewPropertyService(store repository.Store, objectDir string) *PropertyService {
	if objectDir == "" {
		objectDir = defaultObjectStorePath
	}
	return &PropertyService{store: store, objectDir: objectDir}
}

func (s *PropertyService) Create(ctx context.Context, userID, businessID string, input domain.PropertyFile) (domain.PropertyFile, error) {
	member, err := s.authorize(ctx, userID, businessID)
	if err != nil {
		return domain.PropertyFile{}, err
	}
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)
	input.InternalDescription = strings.TrimSpace(input.InternalDescription)
	if input.Title == "" {
		return domain.PropertyFile{}, errors.New("عنوان فایل الزامی است")
	}
	if len(normalizePropertyTypes(input.Type, input.Types)) == 0 {
		return domain.PropertyFile{}, errors.New("نوع فایل معتبر نیست")
	}
	input.Types = normalizePropertyTypes(input.Type, input.Types)
	input.Type = input.Types[0]
	if (input.HouseInfo.PropertyType == "apartment" || strings.EqualFold(input.HouseInfo.PropertyType, "آپارتمان")) && hasPropertyType(input.Types, domain.PropertyFilePartnership) {
		return domain.PropertyFile{}, errors.New("فایل آپارتمانی نمی‌تواند مشارکت باشد")
	}
	if err := validatePropertyNumbers(input); err != nil {
		return domain.PropertyFile{}, err
	}
	if err := normalizePropertyDistribution(&input, member.CommissionPercent); err != nil {
		return domain.PropertyFile{}, err
	}
	if len(input.Addresses) == 0 || len(input.Addresses) > 5 {
		return domain.PropertyFile{}, errors.New("برای هر فایل باید بین ۱ تا ۵ آدرس انتخاب شود")
	}
	addresses, err := s.resolveAddresses(ctx, businessID, input.Addresses)
	if err != nil {
		return domain.PropertyFile{}, err
	}
	input.BusinessID = businessID
	input.OwnerUserID = userID
	input.Addresses = addresses
	input.Media = nil
	created, err := s.store.CreatePropertyFile(ctx, input)
	if err != nil {
		return domain.PropertyFile{}, err
	}
	if err := s.syncPropertyVaults(ctx, userID, businessID, created); err != nil {
		return domain.PropertyFile{}, err
	}
	return created, nil
}

func (s *PropertyService) Update(ctx context.Context, userID, businessID, fileID string, input domain.PropertyFile) (domain.PropertyFile, error) {
	member, err := s.authorize(ctx, userID, businessID)
	if err != nil {
		return domain.PropertyFile{}, err
	}
	existing, err := s.store.GetPropertyFile(ctx, businessID, fileID)
	if err != nil {
		return domain.PropertyFile{}, err
	}
	if existing.OwnerUserID != userID {
		return domain.PropertyFile{}, errors.New("فقط ثبت‌کننده فایل می‌تواند آن را ویرایش کند")
	}
	input.ID = existing.ID
	input.BusinessID = businessID
	input.OwnerUserID = userID
	input.CreatedAt = existing.CreatedAt
	input.Media = existing.Media
	input.Title = strings.TrimSpace(input.Title)
	input.Description = strings.TrimSpace(input.Description)
	input.InternalDescription = strings.TrimSpace(input.InternalDescription)
	if input.Title == "" {
		return domain.PropertyFile{}, errors.New("عنوان فایل الزامی است")
	}
	if len(normalizePropertyTypes(input.Type, input.Types)) == 0 {
		return domain.PropertyFile{}, errors.New("نوع فایل معتبر نیست")
	}
	input.Types = normalizePropertyTypes(input.Type, input.Types)
	input.Type = input.Types[0]
	if err := validatePropertyNumbers(input); err != nil {
		return domain.PropertyFile{}, err
	}
	if err := normalizePropertyDistribution(&input, member.CommissionPercent); err != nil {
		return domain.PropertyFile{}, err
	}
	addresses, err := s.resolveAddresses(ctx, businessID, input.Addresses)
	if err != nil {
		return domain.PropertyFile{}, err
	}
	input.Addresses = addresses
	updated, err := s.store.UpdatePropertyFile(ctx, input)
	if err != nil {
		return domain.PropertyFile{}, err
	}
	if err := s.syncPropertyVaults(ctx, userID, businessID, updated); err != nil {
		return domain.PropertyFile{}, err
	}
	return updated, nil
}

func hasPropertyType(types []domain.PropertyFileType, target domain.PropertyFileType) bool {
	for _, item := range types {
		if item == target {
			return true
		}
	}
	return false
}

func (s *PropertyService) List(ctx context.Context, userID, businessID string) ([]domain.PropertyFile, error) {
	if _, err := s.authorize(ctx, userID, businessID); err != nil {
		return nil, err
	}
	items, err := s.store.ListPropertyFiles(ctx, businessID)
	if err != nil {
		return nil, err
	}
	result := make([]domain.PropertyFile, 0, len(items))
	for _, item := range items {
		if item.OwnerUserID == userID {
			result = append(result, item)
		}
	}
	return result, nil
}

func (s *PropertyService) Latest(ctx context.Context, userID, businessID string, propertyType domain.PropertyFileType, limit, offset int) ([]domain.PropertyFile, int, error) {
	if _, err := s.authorize(ctx, userID, businessID); err != nil {
		return nil, 0, err
	}
	if propertyType != "" &&
		propertyType != domain.PropertyFileSale &&
		propertyType != domain.PropertyFilePartnership &&
		propertyType != domain.PropertyFileRentLease {
		return nil, 0, errors.New("Ù†ÙˆØ¹ ÙÛŒÙ„ØªØ± ÙØ§ÛŒÙ„ Ù…Ø¹ØªØ¨Ø± Ù†ÛŒØ³Øª")
	}
	limit, offset = normalizePropertyPage(limit, offset)
	vaultIDs, err := s.accessibleVaultIDs(ctx, userID, businessID)
	if err != nil {
		return nil, 0, err
	}
	return s.store.ListLatestPropertyFiles(ctx, businessID, userID, vaultIDs, propertyType, limit, offset)
}

func (s *PropertyService) RequestShare(ctx context.Context, requester domain.User, businessID, fileID string, commissionPercent float64) (domain.PropertyShareRequest, error) {
	if _, err := s.authorize(ctx, requester.ID, businessID); err != nil {
		return domain.PropertyShareRequest{}, err
	}
	file, err := s.store.GetPropertyFile(ctx, businessID, fileID)
	if err != nil {
		return domain.PropertyShareRequest{}, err
	}
	if file.OwnerUserID == requester.ID {
		return domain.PropertyShareRequest{}, errors.New("امکان درخواست مشارکت روی فایل خودتان وجود ندارد")
	}
	if file.IsPartnershipCopy {
		return domain.PropertyShareRequest{}, errors.New("روی فایل مشارکتی نمی‌توان درخواست مشارکت جدید ثبت کرد")
	}
	if commissionPercent == 0 {
		commissionPercent = 25
	}
	if commissionPercent < 0 || commissionPercent > 100 {
		return domain.PropertyShareRequest{}, errors.New("درصد مشارکت معتبر نیست")
	}
	existing, _ := s.store.ListPropertyShareRequestsForRequester(ctx, businessID, requester.ID)
	for _, request := range existing {
		if request.PropertyFileID == fileID && (request.Status == domain.PropertySharePending || request.Status == domain.PropertyShareApproved) {
			return domain.PropertyShareRequest{}, errors.New("برای این فایل درخواست فعال دارید")
		}
	}
	request, err := s.store.CreatePropertyShareRequest(ctx, domain.PropertyShareRequest{
		BusinessID:        businessID,
		PropertyFileID:    file.ID,
		PropertyTitle:     file.Title,
		OwnerUserID:       file.OwnerUserID,
		RequesterUserID:   requester.ID,
		RequesterName:     requester.DisplayName,
		RequesterPhone:    requester.Phone,
		CommissionPercent: commissionPercent,
		Status:            domain.PropertySharePending,
	})
	if err != nil {
		return domain.PropertyShareRequest{}, err
	}
	file.SharingHistory = upsertSharingHistory(file.SharingHistory, request)
	_, _ = s.store.UpdatePropertyFile(ctx, file)
	_, _ = s.store.CreateNotification(ctx, domain.Notification{
		UserID:     file.OwnerUserID,
		Type:       "property_share_requested",
		Title:      "درخواست مشارکت فایل",
		Body:       fmt.Sprintf("%s برای فایل %s درخواست مشارکت با سهم %.0f%% ثبت کرد", requester.DisplayName, file.Title, commissionPercent),
		BusinessID: businessID,
		PropertyID: file.ID,
		RequestID:  request.ID,
	})
	return request, nil
}

func normalizePropertyPage(limit, offset int) (int, int) {
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

func (s *PropertyService) accessibleVaultIDs(ctx context.Context, userID, businessID string) ([]string, error) {
	user, err := s.store.GetUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	channels, err := s.store.ListChannelsForUser(ctx, userID, []string{user.Phone}, []string{businessID})
	if err != nil {
		return nil, err
	}
	seen := map[string]struct{}{}
	result := []string{}
	for _, channel := range channels {
		if channel.VaultID == "" {
			continue
		}
		if channel.BusinessID != "" && channel.BusinessID != businessID {
			continue
		}
		if _, ok := seen[channel.VaultID]; ok {
			continue
		}
		seen[channel.VaultID] = struct{}{}
		result = append(result, channel.VaultID)
	}
	if member, err := s.store.GetMemberByUser(ctx, businessID, userID); err == nil &&
		(member.Role == domain.RoleOwner || member.Role == domain.RoleManager || domain.HasPermission(member, domain.PermBusinessUpdate)) {
		if vaults, err := s.store.ListBusinessVaults(ctx, businessID); err == nil {
			for _, vault := range vaults {
				if vault.ID == "" {
					continue
				}
				if _, ok := seen[vault.ID]; ok {
					continue
				}
				seen[vault.ID] = struct{}{}
				result = append(result, vault.ID)
			}
		}
	}
	return result, nil
}

func (s *PropertyService) notifyMatchingContactRequests(ctx context.Context, businessID string, file domain.PropertyFile) {
	contacts, err := s.store.ListContacts(ctx, businessID)
	if err != nil {
		return
	}
	accessCache := map[string]map[string]struct{}{}
	for _, contact := range contacts {
		if contact.CreatedByID == "" || len(contact.Requests) == 0 {
			continue
		}
		if !s.fileVisibleToUser(ctx, contact.CreatedByID, businessID, file, accessCache) {
			continue
		}
		for _, request := range contact.Requests {
			if request.ID == "" || request.Status == "done" || request.Status == "finished" || request.Status == "suspended" {
				continue
			}
			match := searchengine.MatchProperty(request, file)
			if match.Tier == "" || match.Score == 0 {
				continue
			}
			_, _ = s.store.CreateNotification(ctx, domain.Notification{
				UserID:     contact.CreatedByID,
				Type:       "request_property_match",
				Title:      "فایل مناسب برای درخواست مشتری",
				Body:       fmt.Sprintf("فایل %s با درخواست %s از مشتری %s حدود %d%% مچ شد.", file.Title, request.Title, contact.DisplayName, match.Score),
				BusinessID: businessID,
				PropertyID: file.ID,
				RequestID:  request.ID,
			})
		}
	}
}

func (s *PropertyService) fileVisibleToUser(ctx context.Context, userID, businessID string, file domain.PropertyFile, accessCache map[string]map[string]struct{}) bool {
	if file.OwnerUserID == userID {
		return true
	}
	if len(file.VaultIDs) == 0 {
		return false
	}
	vaults, ok := accessCache[userID]
	if !ok {
		ids, err := s.accessibleVaultIDs(ctx, userID, businessID)
		if err != nil {
			return false
		}
		vaults = map[string]struct{}{}
		for _, id := range ids {
			vaults[id] = struct{}{}
		}
		accessCache[userID] = vaults
	}
	for _, vaultID := range file.VaultIDs {
		if _, ok := vaults[vaultID]; ok {
			return true
		}
	}
	return false
}

func (s *PropertyService) ShareRequests(ctx context.Context, userID, businessID, scope string) ([]domain.PropertyShareRequest, error) {
	if _, err := s.authorize(ctx, userID, businessID); err != nil {
		return nil, err
	}
	if scope == "requester" {
		return s.store.ListPropertyShareRequestsForRequester(ctx, businessID, userID)
	}
	return s.store.ListPropertyShareRequestsForOwner(ctx, businessID, userID)
}

func (s *PropertyService) DecideShareRequest(ctx context.Context, ownerID, businessID, requestID string, approve bool) (domain.PropertyShareRequest, error) {
	if _, err := s.authorize(ctx, ownerID, businessID); err != nil {
		return domain.PropertyShareRequest{}, err
	}
	request, err := s.store.GetPropertyShareRequest(ctx, businessID, requestID)
	if err != nil {
		return domain.PropertyShareRequest{}, err
	}
	if request.OwnerUserID != ownerID {
		return domain.PropertyShareRequest{}, errors.New("فقط مالک فایل می‌تواند درخواست مشارکت را تعیین تکلیف کند")
	}
	if request.Status != domain.PropertySharePending {
		return domain.PropertyShareRequest{}, errors.New("این درخواست قبلا تعیین تکلیف شده است")
	}
	now := time.Now().UTC()
	if approve {
		request.Status = domain.PropertyShareApproved
	} else {
		request.Status = domain.PropertyShareRejected
	}
	request.DecidedAt = now
	request.UpdatedAt = now
	request, err = s.store.UpdatePropertyShareRequest(ctx, request)
	if err != nil {
		return domain.PropertyShareRequest{}, err
	}
	file, err := s.store.GetPropertyFile(ctx, businessID, request.PropertyFileID)
	if err == nil {
		file.SharingHistory = upsertSharingHistory(file.SharingHistory, request)
		_, _ = s.store.UpdatePropertyFile(ctx, file)
	}
	notificationType := "property_share_rejected"
	title := "درخواست مشارکت رد شد"
	body := fmt.Sprintf("درخواست مشارکت شما برای فایل %s رد شد", request.PropertyTitle)
	if approve {
		notificationType = "property_share_approved"
		title = "درخواست مشارکت تایید شد"
		body = fmt.Sprintf("درخواست مشارکت شما برای فایل %s با سهم %.0f%% تایید شد", request.PropertyTitle, request.CommissionPercent)
	}
	_, _ = s.store.CreateNotification(ctx, domain.Notification{
		UserID:     request.RequesterUserID,
		Type:       notificationType,
		Title:      title,
		Body:       body,
		BusinessID: businessID,
		PropertyID: request.PropertyFileID,
		RequestID:  request.ID,
	})
	return request, nil
}

func (s *PropertyService) ReceiveSharedFile(ctx context.Context, requesterID, businessID, requestID string) (domain.PropertyFile, error) {
	if _, err := s.authorize(ctx, requesterID, businessID); err != nil {
		return domain.PropertyFile{}, err
	}
	request, err := s.store.GetPropertyShareRequest(ctx, businessID, requestID)
	if err != nil {
		return domain.PropertyFile{}, err
	}
	if request.RequesterUserID != requesterID {
		return domain.PropertyFile{}, errors.New("این درخواست متعلق به شما نیست")
	}
	if request.Status != domain.PropertyShareApproved && request.Status != domain.PropertyShareReceived {
		return domain.PropertyFile{}, errors.New("فقط درخواست تایید شده قابل دریافت است")
	}
	if request.SharedCopyFileID != "" {
		return s.store.GetPropertyFile(ctx, businessID, request.SharedCopyFileID)
	}
	original, err := s.store.GetPropertyFile(ctx, businessID, request.PropertyFileID)
	if err != nil {
		return domain.PropertyFile{}, err
	}
	copy := original
	copy.ID = ""
	copy.OwnerUserID = requesterID
	copy.SharedFromFileID = original.ID
	copy.SharedFromOwnerID = original.OwnerUserID
	copy.IsPartnershipCopy = true
	copy.PartnershipCommissionPercent = request.CommissionPercent
	copy.SharingHistory = nil
	copy.VaultIDs = nil
	copy.VaultPlacements = nil
	copy.CreatedAt = time.Time{}
	copy.UpdatedAt = time.Time{}
	created, err := s.store.CreatePropertyFile(ctx, copy)
	if err != nil {
		return domain.PropertyFile{}, err
	}
	request.Status = domain.PropertyShareReceived
	request.SharedCopyFileID = created.ID
	request.ReceivedAt = time.Now().UTC()
	request.UpdatedAt = request.ReceivedAt
	request, _ = s.store.UpdatePropertyShareRequest(ctx, request)
	original.SharingHistory = upsertSharingHistory(original.SharingHistory, request)
	_, _ = s.store.UpdatePropertyFile(ctx, original)
	return created, nil
}

func normalizePropertyTypes(primary domain.PropertyFileType, types []domain.PropertyFileType) []domain.PropertyFileType {
	seen := map[domain.PropertyFileType]bool{}
	result := []domain.PropertyFileType{}
	add := func(item domain.PropertyFileType) {
		if item != domain.PropertyFileSale && item != domain.PropertyFilePartnership && item != domain.PropertyFileRentLease {
			return
		}
		if seen[item] {
			return
		}
		seen[item] = true
		result = append(result, item)
	}
	for _, item := range types {
		add(item)
	}
	if len(result) == 0 {
		add(primary)
	}
	return result
}

func upsertSharingHistory(items []domain.PropertySharingHistory, request domain.PropertyShareRequest) []domain.PropertySharingHistory {
	entry := domain.PropertySharingHistory{
		RequestID:         request.ID,
		UserID:            request.RequesterUserID,
		UserName:          request.RequesterName,
		UserPhone:         request.RequesterPhone,
		CommissionPercent: request.CommissionPercent,
		Status:            request.Status,
		SharedCopyFileID:  request.SharedCopyFileID,
		CreatedAt:         request.CreatedAt,
		UpdatedAt:         request.UpdatedAt,
	}
	for i := range items {
		if items[i].RequestID == request.ID {
			items[i] = entry
			return items
		}
	}
	return append(items, entry)
}

func validatePropertyNumbers(input domain.PropertyFile) error {
	if input.SalePrice < 0 ||
		input.FinalPrice < 0 ||
		input.DepositPrice < 0 ||
		input.RentPrice < 0 ||
		input.MaxConvertibleDeposit < 0 ||
		input.HouseInfo.AreaM2 < 0 ||
		input.HouseInfo.Bedrooms < 0 ||
		input.HouseInfo.Floor < 0 ||
		input.HouseInfo.TotalFloors < 0 ||
		input.HouseInfo.AgeYears < 0 ||
		input.HouseInfo.TerraceCount < 0 ||
		input.HouseInfo.BackyardAreaM2 < 0 ||
		input.HouseInfo.GardenBuildingAreaM2 < 0 ||
		input.HouseInfo.GardenBuildingFloors < 0 ||
		input.HouseInfo.MasterServiceCount < 0 {
		return errors.New("مقادیر عددی نمی‌توانند منفی باشند")
	}
	return nil
}

func normalizePropertyDistribution(input *domain.PropertyFile, memberCommission float64) error {
	input.Status = normalizePropertyStatus(input.Status)
	placements := normalizePropertyVaultPlacements(input.VaultPlacements, input.VaultIDs)
	collaboratorTotal := 0.0
	vaultIDs := make([]string, 0, len(placements))
	for _, placement := range placements {
		if placement.CommissionPercent < 0 || placement.CommissionPercent > 100 {
			return errors.New("درصد کمیسیون صندوقچه باید بین ۰ تا ۱۰۰ باشد")
		}
		collaboratorTotal += placement.CommissionPercent
		vaultIDs = append(vaultIDs, placement.VaultID)
	}
	if memberCommission < 0 || memberCommission > 100 {
		return errors.New("درصد کمیسیون کاربر در املاک معتبر نیست")
	}
	if collaboratorTotal > memberCommission {
		return errors.New("جمع درصد صندوقچه‌ها بیشتر از سهم کمیسیون شما در این املاک است")
	}
	input.VaultPlacements = placements
	input.VaultIDs = vaultIDs
	input.BusinessCommissionPercent = 100 - memberCommission
	input.OwnerCommissionPercent = memberCommission - collaboratorTotal
	total := input.BusinessCommissionPercent + input.OwnerCommissionPercent + collaboratorTotal
	if total < 99.999 || total > 100.001 {
		return errors.New("جمع همه سهم‌های کمیسیون باید ۱۰۰ درصد باشد")
	}
	return nil
}

func normalizePropertyStatus(value string) string {
	value = strings.TrimSpace(value)
	switch value {
	case "active", "done", "suspended", "expired":
		return value
	default:
		return "active"
	}
}

func normalizePropertyVaultPlacements(items []domain.PropertyVaultPlacement, fallbackVaultIDs []string) []domain.PropertyVaultPlacement {
	seen := map[string]struct{}{}
	result := []domain.PropertyVaultPlacement{}
	if len(items) > 0 {
		for _, item := range items {
			vaultID := strings.TrimSpace(item.VaultID)
			if vaultID == "" {
				continue
			}
			if _, ok := seen[vaultID]; ok {
				continue
			}
			seen[vaultID] = struct{}{}
			result = append(result, domain.PropertyVaultPlacement{
				VaultID:           vaultID,
				CommissionPercent: item.CommissionPercent,
			})
		}
		return result
	}
	for _, vaultID := range fallbackVaultIDs {
		vaultID = strings.TrimSpace(vaultID)
		if vaultID == "" {
			continue
		}
		if _, ok := seen[vaultID]; ok {
			continue
		}
		seen[vaultID] = struct{}{}
		result = append(result, domain.PropertyVaultPlacement{VaultID: vaultID, CommissionPercent: 25})
	}
	return result
}

func (s *PropertyService) AddMedia(ctx context.Context, userID, businessID, fileID string, header *multipart.FileHeader) (domain.PropertyFile, error) {
	if _, err := s.authorize(ctx, userID, businessID); err != nil {
		return domain.PropertyFile{}, err
	}
	file, err := s.store.GetPropertyFile(ctx, businessID, fileID)
	if err != nil {
		return domain.PropertyFile{}, err
	}
	if len(file.Media) >= maxPropertyMedia {
		return domain.PropertyFile{}, errors.New("هر فایل حداکثر ۲۰ مدیا می‌تواند داشته باشد")
	}
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = detectContentType(header)
	}
	videoCount := 0
	for _, media := range file.Media {
		if media.Kind == "video" {
			videoCount++
		}
	}

	var media domain.PropertyMedia
	switch {
	case strings.HasPrefix(contentType, "image/"):
		media, err = s.processImage(ctx, businessID, fileID, header)
	case strings.HasPrefix(contentType, "video/"):
		if videoCount >= maxPropertyVideos {
			return domain.PropertyFile{}, errors.New("هر فایل حداکثر ۲ ویدئو می‌تواند داشته باشد")
		}
		if header.Size > maxVideoUploadBytes {
			return domain.PropertyFile{}, errors.New("حجم ویدئوی خام بیش از حد مجاز است")
		}
		media, err = s.processVideo(ctx, businessID, fileID, header)
	default:
		return domain.PropertyFile{}, errors.New("فقط عکس و ویدئو قابل آپلود است")
	}
	if err != nil {
		return domain.PropertyFile{}, err
	}
	file.Media = append(file.Media, media)
	return s.store.UpdatePropertyFile(ctx, file)
}

func (s *PropertyService) AddMediaFromUpload(ctx context.Context, userID, businessID, fileID, uploadedFileID string) (domain.PropertyFile, error) {
	if _, err := s.authorize(ctx, userID, businessID); err != nil {
		return domain.PropertyFile{}, err
	}
	file, err := s.store.GetPropertyFile(ctx, businessID, fileID)
	if err != nil {
		return domain.PropertyFile{}, err
	}
	if len(file.Media) >= maxPropertyMedia {
		return domain.PropertyFile{}, errors.New("هر فایل حداکثر ۲۰ مدیا می‌تواند داشته باشد")
	}
	videoCount := 0
	for _, media := range file.Media {
		if media.Kind == "video" {
			videoCount++
		}
	}
	object, err := s.store.GetFile(ctx, strings.TrimSpace(uploadedFileID))
	if err != nil {
		return domain.PropertyFile{}, errors.New("فایل آپلود شده پیدا نشد")
	}
	if object.UploaderID != "" && object.UploaderID != userID {
		return domain.PropertyFile{}, errors.New("فایل متعلق به کاربر دیگری است")
	}
	if object.Purpose != "" && object.Purpose != UploadPurposePropertyMedia {
		return domain.PropertyFile{}, errors.New("هدف فایل با درخواست همخوانی ندارد")
	}
	if object.TargetID != "" && object.TargetID != fileID {
		return domain.PropertyFile{}, errors.New("مقصد فایل با درخواست همخوانی ندارد")
	}
	if !object.ExpiresAt.IsZero() && time.Now().UTC().After(object.ExpiresAt) {
		return domain.PropertyFile{}, errors.New("مهلت استفاده از فایل تمام شده است")
	}
	if object.Kind != "image" && object.Kind != "video" {
		return domain.PropertyFile{}, errors.New("فقط عکس و ویدئو قابل آپلود است")
	}
	if object.Kind == "video" && videoCount >= maxPropertyVideos {
		return domain.PropertyFile{}, errors.New("هر فایل حداکثر ۲ ویدئو می‌تواند داشته باشد")
	}
	object.Purpose = UploadPurposePropertyMedia
	object.TargetType = "property"
	object.TargetID = fileID
	object.OwnerID = businessID
	object.Status = FileStatusAttached
	object.ExpiresAt = time.Time{}
	object, err = s.store.UpdateFile(ctx, object)
	if err != nil {
		return domain.PropertyFile{}, err
	}
	file.Media = append(file.Media, domain.PropertyMedia{
		ID:          support.NewID(),
		FileID:      object.ID,
		Kind:        object.Kind,
		URL:         object.URL,
		ContentType: object.ContentType,
		Size:        object.Size,
		CreatedAt:   time.Now().UTC(),
	})
	return s.store.UpdatePropertyFile(ctx, file)
}

func (s *PropertyService) authorize(ctx context.Context, userID, businessID string) (domain.BusinessMember, error) {
	member, err := s.store.GetMemberByUser(ctx, businessID, userID)
	if err != nil || member.Status != domain.MemberActive {
		return domain.BusinessMember{}, errors.New("دسترسی به کسب‌وکار وجود ندارد")
	}
	return member, nil
}

func (s *PropertyService) syncPropertyVaults(ctx context.Context, userID, businessID string, file domain.PropertyFile) error {
	placements := normalizePropertyVaultPlacements(file.VaultPlacements, file.VaultIDs)
	if len(placements) == 0 {
		return s.store.DeletePropertyVaultReferences(ctx, file.ID, nil)
	}
	vaults := make([]domain.ChannelVault, 0, len(placements))
	keepVaultIDs := make([]string, 0, len(placements))
	commissionByVault := map[string]float64{}
	for _, placement := range placements {
		vault, err := s.store.GetChannelVault(ctx, placement.VaultID)
		if err != nil {
			return err
		}
		if vault.OwnerUserID != "" && vault.OwnerUserID != userID {
			return errors.New("دسترسی به یکی از صندوقچه‌های انتخاب‌شده وجود ندارد")
		}
		if vault.BusinessID != "" && vault.BusinessID != businessID {
			return errors.New("صندوقچه انتخاب‌شده متعلق به این املاک نیست")
		}
		vaults = append(vaults, vault)
		keepVaultIDs = append(keepVaultIDs, vault.ID)
		commissionByVault[vault.ID] = placement.CommissionPercent
	}
	for _, vault := range vaults {
		_, err := s.store.UpsertChannelVaultFile(ctx, domain.ChannelVaultFile{
			VaultID:           vault.ID,
			ChannelID:         vault.ChannelID,
			UploaderID:        userID,
			Title:             file.Title,
			Note:              file.InternalDescription,
			SourceType:        "property_file",
			PropertyFileID:    file.ID,
			PropertyStatus:    file.Status,
			CommissionPercent: commissionByVault[vault.ID],
			Kind:              "property",
			URL:               "",
		})
		if err != nil {
			return err
		}
	}
	if err := s.store.DeletePropertyVaultReferences(ctx, file.ID, keepVaultIDs); err != nil {
		return err
	}
	return nil
}

func uniqueStrings(items []string) []string {
	seen := map[string]struct{}{}
	result := []string{}
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		result = append(result, item)
	}
	return result
}

func (s *PropertyService) resolveAddresses(ctx context.Context, businessID string, input []domain.PropertyAddress) ([]domain.PropertyAddress, error) {
	areas, _ := s.store.ListAreas(ctx, businessID)
	areaByID := map[string]domain.Area{}
	for _, area := range areas {
		areaByID[area.ID] = area
	}
	result := make([]domain.PropertyAddress, 0, len(input))
	for _, address := range input {
		area, ok := areaByID[address.AreaID]
		if !ok {
			return nil, errors.New("منطقه انتخاب‌شده معتبر نیست")
		}
		streets, _ := s.store.ListStreets(ctx, businessID, area.ID)
		var street domain.Street
		for _, item := range streets {
			if item.ID == address.StreetID {
				street = item
				break
			}
		}
		if street.ID == "" {
			return nil, errors.New("خیابان انتخاب‌شده معتبر نیست")
		}
		neighborhoods, _ := s.store.ListNeighborhoods(ctx, businessID, area.ID, street.ID)
		var neighborhood domain.Neighborhood
		for _, item := range neighborhoods {
			if item.ID == address.NeighborhoodID {
				neighborhood = item
				break
			}
		}
		if neighborhood.ID == "" {
			return nil, errors.New("محله انتخاب‌شده معتبر نیست")
		}
		result = append(result, domain.PropertyAddress{
			AreaID:             area.ID,
			AreaName:           area.Name,
			StreetID:           street.ID,
			StreetName:         street.Name,
			NeighborhoodID:     neighborhood.ID,
			NeighborhoodName:   neighborhood.Name,
			ManualExactAddress: strings.TrimSpace(address.ManualExactAddress),
		})
	}
	return result, nil
}

func (s *PropertyService) processImage(ctx context.Context, businessID, fileID string, header *multipart.FileHeader) (domain.PropertyMedia, error) {
	src, err := header.Open()
	if err != nil {
		return domain.PropertyMedia{}, err
	}
	defer src.Close()
	img, _, err := image.Decode(src)
	if err != nil {
		return domain.PropertyMedia{}, errors.New("تصویر معتبر نیست")
	}
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	current := img
	var encoded []byte
	for scale := 100; scale >= 45; scale -= 10 {
		if scale != 100 {
			current = resizeNearest(img, width*scale/100, height*scale/100)
		}
		for quality := 82; quality >= 45; quality -= 7 {
			var buf bytes.Buffer
			if err := jpeg.Encode(&buf, current, &jpeg.Options{Quality: quality}); err != nil {
				return domain.PropertyMedia{}, err
			}
			if buf.Len() <= maxImageBytes || quality == 45 {
				encoded = buf.Bytes()
				break
			}
		}
		if len(encoded) <= maxImageBytes {
			break
		}
	}
	if len(encoded) > maxImageBytes {
		return domain.PropertyMedia{}, errors.New("تصویر بعد از فشرده‌سازی هنوز بزرگ‌تر از ۵۰۰ کیلوبایت است")
	}
	key := fmt.Sprintf("properties/%s/%s/%s.jpg", businessID, fileID, support.NewID())
	if err := s.writeObject(key, encoded); err != nil {
		return domain.PropertyMedia{}, err
	}
	object, err := s.store.CreateFile(ctx, domain.FileObject{
		OwnerID:     businessID,
		Provider:    "s3-compatible",
		Bucket:      "amlak",
		Key:         key,
		URL:         "/objects/" + key,
		ContentType: "image/jpeg",
		Size:        int64(len(encoded)),
	})
	if err != nil {
		return domain.PropertyMedia{}, err
	}
	return domain.PropertyMedia{
		ID:          support.NewID(),
		FileID:      object.ID,
		Kind:        "image",
		URL:         object.URL,
		ContentType: "image/jpeg",
		Size:        int64(len(encoded)),
		Width:       current.Bounds().Dx(),
		Height:      current.Bounds().Dy(),
		CreatedAt:   time.Now().UTC(),
	}, nil
}

func (s *PropertyService) processVideo(ctx context.Context, businessID, fileID string, header *multipart.FileHeader) (domain.PropertyMedia, error) {
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return domain.PropertyMedia{}, errors.New("برای پردازش ویدئو، ffmpeg باید روی سرور نصب باشد")
	}
	if _, err := exec.LookPath("ffprobe"); err != nil {
		return domain.PropertyMedia{}, errors.New("برای بررسی مدت ویدئو، ffprobe باید روی سرور نصب باشد")
	}
	tempDir, err := os.MkdirTemp("", "amlak-video-*")
	if err != nil {
		return domain.PropertyMedia{}, err
	}
	defer os.RemoveAll(tempDir)

	inputPath := filepath.Join(tempDir, "input")
	outputPath := filepath.Join(tempDir, "output.mp4")
	if err := saveUpload(header, inputPath); err != nil {
		return domain.PropertyMedia{}, err
	}
	duration, err := probeDuration(ctx, inputPath)
	if err != nil {
		return domain.PropertyMedia{}, errors.New("مدت ویدئو قابل تشخیص نیست")
	}
	if duration > maxVideoDurationSec {
		return domain.PropertyMedia{}, errors.New("مدت ویدئو نباید بیشتر از ۲ دقیقه باشد")
	}
	cmd := exec.CommandContext(
		ctx,
		"ffmpeg", "-y", "-i", inputPath,
		"-vf", "scale=-2:480",
		"-c:v", "libx264", "-preset", "veryfast", "-crf", "28",
		"-c:a", "aac", "-b:a", "96k",
		"-movflags", "+faststart",
		outputPath,
	)
	if output, err := cmd.CombinedOutput(); err != nil {
		return domain.PropertyMedia{}, fmt.Errorf("پردازش ویدئو ناموفق بود: %s", strings.TrimSpace(string(output)))
	}
	body, err := os.ReadFile(outputPath)
	if err != nil {
		return domain.PropertyMedia{}, err
	}
	key := fmt.Sprintf("properties/%s/%s/%s.mp4", businessID, fileID, support.NewID())
	if err := s.writeObject(key, body); err != nil {
		return domain.PropertyMedia{}, err
	}
	object, err := s.store.CreateFile(ctx, domain.FileObject{
		OwnerID:     businessID,
		Provider:    "s3-compatible",
		Bucket:      "amlak",
		Key:         key,
		URL:         "/objects/" + key,
		ContentType: "video/mp4",
		Size:        int64(len(body)),
	})
	if err != nil {
		return domain.PropertyMedia{}, err
	}
	return domain.PropertyMedia{
		ID:          support.NewID(),
		FileID:      object.ID,
		Kind:        "video",
		URL:         object.URL,
		ContentType: "video/mp4",
		Size:        int64(len(body)),
		Height:      480,
		DurationSec: duration,
		CreatedAt:   time.Now().UTC(),
	}, nil
}

func (s *PropertyService) writeObject(key string, body []byte) error {
	path := filepath.Join(s.objectDir, filepath.FromSlash(key))
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, body, 0644)
}

func resizeNearest(src image.Image, width, height int) image.Image {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}
	dst := image.NewRGBA(image.Rect(0, 0, width, height))
	srcBounds := src.Bounds()
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			sx := srcBounds.Min.X + x*srcBounds.Dx()/width
			sy := srcBounds.Min.Y + y*srcBounds.Dy()/height
			dst.Set(x, y, src.At(sx, sy))
		}
	}
	return dst
}

func detectContentType(header *multipart.FileHeader) string {
	src, err := header.Open()
	if err != nil {
		return ""
	}
	defer src.Close()
	buffer := make([]byte, 512)
	n, _ := io.ReadFull(src, buffer)
	return http.DetectContentType(buffer[:n])
}

func saveUpload(header *multipart.FileHeader, path string) error {
	src, err := header.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := os.Create(path)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	return err
}

func probeDuration(ctx context.Context, path string) (int, error) {
	cmd := exec.CommandContext(ctx, "ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=nokey=1:noprint_wrappers=1", path)
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}
	value := strings.TrimSpace(string(output))
	seconds, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}
	return int(seconds + 0.5), nil
}
