package service

import (
	"context"
	"errors"
	"strings"

	"amlakcrm/backend/internal/domain"
	"amlakcrm/backend/internal/repository"
)

type LocationService struct {
	store repository.Store
}

func NewLocationService(store repository.Store) *LocationService {
	return &LocationService{store: store}
}

func (s *LocationService) List(ctx context.Context, userID, businessID string) ([]map[string]interface{}, error) {
	if _, err := s.authorizeRead(ctx, userID, businessID); err != nil {
		return nil, err
	}
	cityID, err := s.userCityID(ctx, userID)
	if err != nil {
		return nil, err
	}
	areas, err := s.store.ListAreas(ctx, businessID)
	if err != nil {
		return nil, err
	}
	result := make([]map[string]interface{}, 0, len(areas))
	for _, area := range areas {
		if area.CityID != cityID {
			continue
		}
		streets, _ := s.store.ListStreets(ctx, businessID, area.ID)
		streetViews := make([]map[string]interface{}, 0, len(streets))
		for _, street := range streets {
			if street.CityID != cityID {
				continue
			}
			neighborhoods, _ := s.store.ListNeighborhoods(ctx, businessID, area.ID, street.ID)
			filteredNeighborhoods := make([]domain.Neighborhood, 0, len(neighborhoods))
			for _, neighborhood := range neighborhoods {
				if neighborhood.CityID == cityID {
					filteredNeighborhoods = append(filteredNeighborhoods, neighborhood)
				}
			}
			streetViews = append(streetViews, map[string]interface{}{
				"id":            street.ID,
				"businessId":    street.BusinessID,
				"cityId":        street.CityID,
				"areaId":        street.AreaID,
				"name":          street.Name,
				"createdAt":     street.CreatedAt,
				"updatedAt":     street.UpdatedAt,
				"neighborhoods": filteredNeighborhoods,
			})
		}
		result = append(result, map[string]interface{}{
			"id":         area.ID,
			"businessId": area.BusinessID,
			"cityId":     area.CityID,
			"name":       area.Name,
			"createdAt":  area.CreatedAt,
			"updatedAt":  area.UpdatedAt,
			"streets":    streetViews,
		})
	}
	return result, nil
}

func (s *LocationService) CreateArea(ctx context.Context, userID, businessID, name string) (domain.Area, error) {
	if _, err := s.authorizeManage(ctx, userID, businessID); err != nil {
		return domain.Area{}, err
	}
	cityID, err := s.userCityID(ctx, userID)
	if err != nil {
		return domain.Area{}, err
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return domain.Area{}, errors.New("نام منطقه الزامی است")
	}
	return s.store.CreateArea(ctx, domain.Area{BusinessID: businessID, CityID: cityID, Name: name})
}

func (s *LocationService) DeleteArea(ctx context.Context, userID, businessID, areaID string) error {
	if _, err := s.authorizeManage(ctx, userID, businessID); err != nil {
		return err
	}
	cityID, err := s.userCityID(ctx, userID)
	if err != nil {
		return err
	}
	if !s.areaInCity(ctx, businessID, areaID, cityID) {
		return repository.ErrNotFound
	}
	return s.store.DeleteArea(ctx, businessID, areaID)
}

func (s *LocationService) CreateStreet(ctx context.Context, userID, businessID, areaID, name string) (domain.Street, error) {
	if _, err := s.authorizeManage(ctx, userID, businessID); err != nil {
		return domain.Street{}, err
	}
	cityID, err := s.userCityID(ctx, userID)
	if err != nil {
		return domain.Street{}, err
	}
	if !s.areaInCity(ctx, businessID, areaID, cityID) {
		return domain.Street{}, repository.ErrNotFound
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return domain.Street{}, errors.New("نام خیابان الزامی است")
	}
	return s.store.CreateStreet(ctx, domain.Street{BusinessID: businessID, CityID: cityID, AreaID: areaID, Name: name})
}

func (s *LocationService) DeleteStreet(ctx context.Context, userID, businessID, areaID, streetID string) error {
	if _, err := s.authorizeManage(ctx, userID, businessID); err != nil {
		return err
	}
	cityID, err := s.userCityID(ctx, userID)
	if err != nil {
		return err
	}
	if !s.areaInCity(ctx, businessID, areaID, cityID) || !s.streetInCity(ctx, businessID, areaID, streetID, cityID) {
		return repository.ErrNotFound
	}
	return s.store.DeleteStreet(ctx, businessID, areaID, streetID)
}

func (s *LocationService) CreateNeighborhood(ctx context.Context, userID, businessID, areaID, streetID, name string) (domain.Neighborhood, error) {
	if _, err := s.authorizeManage(ctx, userID, businessID); err != nil {
		return domain.Neighborhood{}, err
	}
	cityID, err := s.userCityID(ctx, userID)
	if err != nil {
		return domain.Neighborhood{}, err
	}
	if !s.areaInCity(ctx, businessID, areaID, cityID) || !s.streetInCity(ctx, businessID, areaID, streetID, cityID) {
		return domain.Neighborhood{}, repository.ErrNotFound
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return domain.Neighborhood{}, errors.New("نام محله الزامی است")
	}
	return s.store.CreateNeighborhood(ctx, domain.Neighborhood{BusinessID: businessID, CityID: cityID, AreaID: areaID, StreetID: streetID, Name: name})
}

func (s *LocationService) DeleteNeighborhood(ctx context.Context, userID, businessID, areaID, streetID, neighborhoodID string) error {
	if _, err := s.authorizeManage(ctx, userID, businessID); err != nil {
		return err
	}
	cityID, err := s.userCityID(ctx, userID)
	if err != nil {
		return err
	}
	if !s.neighborhoodInCity(ctx, businessID, areaID, streetID, neighborhoodID, cityID) {
		return repository.ErrNotFound
	}
	return s.store.DeleteNeighborhood(ctx, businessID, areaID, streetID, neighborhoodID)
}

func (s *LocationService) userCityID(ctx context.Context, userID string) (string, error) {
	user, err := s.store.GetUser(ctx, userID)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(user.CityID) == "" {
		return "", errors.New("شهر پروفایل مشخص نیست")
	}
	return user.CityID, nil
}

func (s *LocationService) areaInCity(ctx context.Context, businessID, areaID, cityID string) bool {
	areas, err := s.store.ListAreas(ctx, businessID)
	if err != nil {
		return false
	}
	for _, area := range areas {
		if area.ID == areaID && area.CityID == cityID {
			return true
		}
	}
	return false
}

func (s *LocationService) streetInCity(ctx context.Context, businessID, areaID, streetID, cityID string) bool {
	streets, err := s.store.ListStreets(ctx, businessID, areaID)
	if err != nil {
		return false
	}
	for _, street := range streets {
		if street.ID == streetID && street.CityID == cityID {
			return true
		}
	}
	return false
}

func (s *LocationService) neighborhoodInCity(ctx context.Context, businessID, areaID, streetID, neighborhoodID, cityID string) bool {
	neighborhoods, err := s.store.ListNeighborhoods(ctx, businessID, areaID, streetID)
	if err != nil {
		return false
	}
	for _, neighborhood := range neighborhoods {
		if neighborhood.ID == neighborhoodID && neighborhood.CityID == cityID {
			return true
		}
	}
	return false
}

func (s *LocationService) authorizeRead(ctx context.Context, userID, businessID string) (domain.BusinessMember, error) {
	member, err := s.store.GetMemberByUser(ctx, businessID, userID)
	if err != nil || member.Status != domain.MemberActive {
		return domain.BusinessMember{}, errors.New("دسترسی به کسب و کار وجود ندارد")
	}
	return member, nil
}

func (s *LocationService) authorizeManage(ctx context.Context, userID, businessID string) (domain.BusinessMember, error) {
	member, err := s.authorizeRead(ctx, userID, businessID)
	if err != nil {
		return domain.BusinessMember{}, err
	}
	if member.Role == domain.RoleOwner || domain.HasPermission(member, domain.PermLocationsManage) {
		return member, nil
	}
	return domain.BusinessMember{}, errors.New("شما دسترسی مدیریت مناطق را ندارید")
}
