package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"amlakcrm/backend/internal/config"
	"amlakcrm/backend/internal/domain"
	"amlakcrm/backend/internal/repository"
)

type PlatformService struct {
	store            repository.Store
	superAdminPhones map[string]struct{}
}

func NewPlatformService(store repository.Store, cfg config.Config) *PlatformService {
	phones := map[string]struct{}{}
	for _, raw := range strings.Split(cfg.SuperAdminPhones, ",") {
		phone, err := NormalizePhone(strings.TrimSpace(raw))
		if err == nil && phone != "" {
			phones[phone] = struct{}{}
		}
	}
	return &PlatformService{store: store, superAdminPhones: phones}
}

func (s *PlatformService) BootstrapAdmin(ctx context.Context, user domain.User) {
	if _, ok := s.superAdminPhones[user.Phone]; !ok {
		return
	}
	if _, err := s.store.GetAdminAccountByUser(ctx, user.ID); err == nil {
		return
	}
	_, _ = s.store.SaveAdminAccount(ctx, domain.AdminAccount{
		UserID: user.ID,
		Roles:  []domain.AdminRole{domain.AdminRoleSuperAdmin},
		Status: "active",
	})
}

func (s *PlatformService) AdminAccount(ctx context.Context, user domain.User) (domain.AdminAccount, error) {
	s.BootstrapAdmin(ctx, user)
	return s.store.GetAdminAccountByUser(ctx, user.ID)
}

func (s *PlatformService) Require(ctx context.Context, user domain.User, permission string) (domain.AdminAccount, error) {
	account, err := s.AdminAccount(ctx, user)
	if err != nil {
		return domain.AdminAccount{}, errors.New("platform admin access required")
	}
	if !domain.HasPlatformPermission(account, permission) {
		return domain.AdminAccount{}, errors.New("platform permission denied")
	}
	return account, nil
}

func (s *PlatformService) Me(ctx context.Context, user domain.User) (domain.AdminAccount, error) {
	return s.AdminAccount(ctx, user)
}

func (s *PlatformService) ListUsers(ctx context.Context, user domain.User) ([]domain.User, error) {
	if _, err := s.Require(ctx, user, domain.PermPlatformUsersRead); err != nil {
		return nil, err
	}
	return s.store.ListUsers(ctx)
}

func (s *PlatformService) ListBusinesses(ctx context.Context, user domain.User) ([]domain.Business, error) {
	if _, err := s.Require(ctx, user, domain.PermPlatformBusinessesRead); err != nil {
		return nil, err
	}
	return s.store.ListBusinesses(ctx)
}

func (s *PlatformService) ListAdminAccounts(ctx context.Context, user domain.User) ([]domain.AdminAccount, error) {
	if _, err := s.Require(ctx, user, domain.PermPlatformAdminsManage); err != nil {
		return nil, err
	}
	return s.store.ListAdminAccounts(ctx)
}

func (s *PlatformService) Settings(ctx context.Context, user domain.User) (domain.PlatformSettings, error) {
	if _, err := s.Require(ctx, user, domain.PermPlatformAdminsManage); err != nil {
		return domain.PlatformSettings{}, err
	}
	settings, err := s.store.GetPlatformSettings(ctx)
	if err != nil {
		return domain.PlatformSettings{}, err
	}
	settings.OTPAPIKeyMasked = maskSecret(settings.OTPAPIKey)
	settings.ServiceSMSMasked = maskSecret(settings.ServiceSMSAPIKey)
	return settings, nil
}

func (s *PlatformService) UpdateSettings(ctx context.Context, user domain.User, otpAPIKey, serviceSMSAPIKey string) (domain.PlatformSettings, error) {
	admin, err := s.Require(ctx, user, domain.PermPlatformAdminsManage)
	if err != nil {
		return domain.PlatformSettings{}, err
	}
	settings, err := s.store.GetPlatformSettings(ctx)
	if err != nil {
		return domain.PlatformSettings{}, err
	}
	otpAPIKey = strings.TrimSpace(otpAPIKey)
	serviceSMSAPIKey = strings.TrimSpace(serviceSMSAPIKey)
	if otpAPIKey != "" {
		settings.OTPAPIKey = otpAPIKey
	}
	if serviceSMSAPIKey != "" {
		settings.ServiceSMSAPIKey = serviceSMSAPIKey
	}
	settings.ID = "platform"
	settings.UpdatedByAdminID = admin.ID
	updated, err := s.store.SavePlatformSettings(ctx, settings)
	if err != nil {
		return domain.PlatformSettings{}, err
	}
	updated.OTPAPIKeyMasked = maskSecret(updated.OTPAPIKey)
	updated.ServiceSMSMasked = maskSecret(updated.ServiceSMSAPIKey)
	return updated, nil
}

func (s *PlatformService) SaveAdminAccount(ctx context.Context, actor domain.User, targetUserID string, roles []domain.AdminRole) (domain.AdminAccount, error) {
	admin, err := s.Require(ctx, actor, domain.PermPlatformAdminsManage)
	if err != nil {
		return domain.AdminAccount{}, err
	}
	if targetUserID == "" || len(roles) == 0 {
		return domain.AdminAccount{}, errors.New("user and at least one role are required")
	}
	account := domain.AdminAccount{
		UserID:           targetUserID,
		Roles:            roles,
		Status:           "active",
		CreatedByAdminID: admin.ID,
	}
	return s.store.SaveAdminAccount(ctx, account)
}

func (s *PlatformService) CreateCity(ctx context.Context, user domain.User, name string) (domain.City, error) {
	if _, err := s.Require(ctx, user, domain.PermPlatformLocationsManage); err != nil {
		return domain.City{}, err
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return domain.City{}, errors.New("city name is required")
	}
	return s.store.CreateCity(ctx, domain.City{Name: name, NormalizedName: NormalizePersianName(name)})
}

func (s *PlatformService) ListCities(ctx context.Context) ([]domain.City, error) {
	return s.store.ListCities(ctx)
}

func (s *PlatformService) CityLocations(ctx context.Context, cityID string) (map[string]interface{}, error) {
	if _, err := s.store.GetCity(ctx, cityID); err != nil {
		return nil, err
	}
	areas, err := s.store.ListSystemAreas(ctx, cityID)
	if err != nil {
		return nil, err
	}
	areaViews := make([]map[string]interface{}, 0, len(areas))
	for _, area := range areas {
		streets, _ := s.store.ListSystemStreets(ctx, cityID, area.ID)
		streetViews := make([]map[string]interface{}, 0, len(streets))
		for _, street := range streets {
			neighborhoods, _ := s.store.ListSystemNeighborhoods(ctx, cityID, area.ID, street.ID)
			streetViews = append(streetViews, map[string]interface{}{
				"id":            street.ID,
				"cityId":        street.CityID,
				"areaId":        street.AreaID,
				"name":          street.Name,
				"neighborhoods": neighborhoods,
			})
		}
		areaViews = append(areaViews, map[string]interface{}{
			"id":      area.ID,
			"cityId":  area.CityID,
			"name":    area.Name,
			"streets": streetViews,
		})
	}
	return map[string]interface{}{"cityId": cityID, "areas": areaViews}, nil
}

func (s *PlatformService) SearchCityLocations(ctx context.Context, cityID, query string) (map[string]interface{}, error) {
	query = NormalizePersianName(query)
	if query == "" {
		return s.CityLocations(ctx, cityID)
	}
	areas, _ := s.store.ListSystemAreas(ctx, cityID)
	areaMatches := []domain.SystemArea{}
	for _, area := range areas {
		if strings.Contains(area.NormalizedName, query) {
			areaMatches = append(areaMatches, area)
		}
	}
	return map[string]interface{}{"areas": areaMatches}, nil
}

func (s *PlatformService) CreateLocationSuggestion(ctx context.Context, user domain.User, input domain.LocationSuggestion) (domain.LocationSuggestion, error) {
	input.Name = strings.TrimSpace(input.Name)
	if input.CityID == "" || input.Name == "" || input.Type == "" {
		return domain.LocationSuggestion{}, errors.New("city, type and name are required")
	}
	if user.CityID == "" || user.CityID != input.CityID {
		return domain.LocationSuggestion{}, errors.New("suggesting locations for this city is not allowed")
	}
	if _, err := s.store.GetCity(ctx, input.CityID); err != nil {
		return domain.LocationSuggestion{}, errors.New("city not found")
	}
	input.SubmittedByUserID = user.ID
	input.NormalizedName = NormalizePersianName(input.Name)
	input.Status = domain.LocationSuggestionPending
	return s.store.CreateLocationSuggestion(ctx, input)
}

func (s *PlatformService) ListLocationSuggestions(ctx context.Context, user domain.User, status domain.LocationSuggestionStatus) ([]domain.LocationSuggestion, error) {
	if _, err := s.Require(ctx, user, domain.PermPlatformLocationsReview); err != nil {
		return nil, err
	}
	return s.store.ListLocationSuggestions(ctx, status)
}

func (s *PlatformService) ApproveLocationSuggestion(ctx context.Context, user domain.User, suggestionID string) (domain.LocationSuggestion, error) {
	admin, err := s.Require(ctx, user, domain.PermPlatformLocationsReview)
	if err != nil {
		return domain.LocationSuggestion{}, err
	}
	suggestion, err := s.store.GetLocationSuggestion(ctx, suggestionID)
	if err != nil {
		return domain.LocationSuggestion{}, err
	}
	if suggestion.Status != domain.LocationSuggestionPending {
		return domain.LocationSuggestion{}, errors.New("suggestion is not pending")
	}
	locationID, err := s.createSystemLocationFromSuggestion(ctx, suggestion)
	if err != nil {
		return domain.LocationSuggestion{}, err
	}
	now := time.Now().UTC()
	suggestion.Status = domain.LocationSuggestionApproved
	suggestion.ReviewedByAdminID = admin.ID
	suggestion.ApprovedSystemLocationID = locationID
	suggestion.ReviewedAt = now
	return s.store.UpdateLocationSuggestion(ctx, suggestion)
}

func (s *PlatformService) RejectLocationSuggestion(ctx context.Context, user domain.User, suggestionID, note string) (domain.LocationSuggestion, error) {
	admin, err := s.Require(ctx, user, domain.PermPlatformLocationsReview)
	if err != nil {
		return domain.LocationSuggestion{}, err
	}
	suggestion, err := s.store.GetLocationSuggestion(ctx, suggestionID)
	if err != nil {
		return domain.LocationSuggestion{}, err
	}
	if suggestion.Status != domain.LocationSuggestionPending {
		return domain.LocationSuggestion{}, errors.New("suggestion is not pending")
	}
	suggestion.Status = domain.LocationSuggestionRejected
	suggestion.ReviewedByAdminID = admin.ID
	suggestion.ReviewNote = strings.TrimSpace(note)
	suggestion.ReviewedAt = time.Now().UTC()
	return s.store.UpdateLocationSuggestion(ctx, suggestion)
}

func (s *PlatformService) createSystemLocationFromSuggestion(ctx context.Context, suggestion domain.LocationSuggestion) (string, error) {
	switch suggestion.Type {
	case domain.LocationSuggestionArea:
		area, err := s.store.CreateSystemArea(ctx, domain.SystemArea{
			CityID:         suggestion.CityID,
			Name:           suggestion.Name,
			NormalizedName: suggestion.NormalizedName,
		})
		return area.ID, err
	case domain.LocationSuggestionStreet:
		street, err := s.store.CreateSystemStreet(ctx, domain.SystemStreet{
			CityID:         suggestion.CityID,
			AreaID:         suggestion.ParentAreaID,
			Name:           suggestion.Name,
			NormalizedName: suggestion.NormalizedName,
		})
		return street.ID, err
	case domain.LocationSuggestionNeighborhood:
		neighborhood, err := s.store.CreateSystemNeighborhood(ctx, domain.SystemNeighborhood{
			CityID:         suggestion.CityID,
			AreaID:         suggestion.ParentAreaID,
			StreetID:       suggestion.ParentStreetID,
			Name:           suggestion.Name,
			NormalizedName: suggestion.NormalizedName,
		})
		return neighborhood.ID, err
	default:
		return "", errors.New("unsupported suggestion type")
	}
}

func NormalizePersianName(value string) string {
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, "ي", "ی")
	value = strings.ReplaceAll(value, "ك", "ک")
	return strings.Join(strings.Fields(value), " ")
}

func maskSecret(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if len(value) <= 8 {
		return "****"
	}
	return value[:4] + "****" + value[len(value)-4:]
}
