package service

import (
	"context"
	"errors"
	"time"

	"amlakcrm/backend/internal/domain"
	"amlakcrm/backend/internal/repository"
)

type BusinessService struct {
	store repository.Store
}

func NewBusinessService(store repository.Store) *BusinessService {
	return &BusinessService{store: store}
}

func (s *BusinessService) Create(ctx context.Context, owner domain.User, input domain.Business) (domain.Business, error) {
	if input.Name == "" || len(input.Phones) == 0 {
		return domain.Business{}, errors.New("نام کسب‌وکار و شماره تماس اصلی الزامی است")
	}
	input.OwnerUserID = owner.ID
	input.LicenseStatus = "trial"
	business, err := s.store.CreateBusiness(ctx, input)
	if err != nil {
		return domain.Business{}, err
	}
	_, err = s.store.CreateMember(ctx, domain.BusinessMember{
		BusinessID:        business.ID,
		UserID:            owner.ID,
		UserPhone:         owner.Phone,
		UserDisplayName:   owner.DisplayName,
		Role:              domain.RoleOwner,
		Permissions:       domain.DefaultPermissions(domain.RoleOwner),
		CommissionPercent: 100,
		Status:            domain.MemberActive,
	})
	return business, err
}

func (s *BusinessService) ListForUser(ctx context.Context, userID string) ([]domain.Business, error) {
	return s.store.ListBusinessesForUser(ctx, userID)
}

func (s *BusinessService) Dashboard(ctx context.Context, businessID string) (map[string]interface{}, error) {
	business, err := s.store.GetBusiness(ctx, businessID)
	if err != nil {
		return nil, err
	}
	members, _ := s.store.ListMembers(ctx, businessID)
	invites, _ := s.store.ListInvitations(ctx, businessID)
	properties, _ := s.store.ListPropertyFiles(ctx, businessID)
	pending := 0
	for _, invite := range invites {
		if invite.Status == domain.InvitationPending && time.Now().UTC().Before(invite.ExpiresAt) {
			pending++
		}
	}
	return map[string]interface{}{
		"business":           business,
		"consultantsCount":   len(members),
		"pendingInvitations": pending,
		"licenseStatus":      business.LicenseStatus,
		"propertiesCount":    len(properties),
		"customersCount":     0,
		"contractsCount":     0,
	}, nil
}

func (s *BusinessService) Members(ctx context.Context, businessID string) ([]domain.BusinessMember, error) {
	return s.store.ListMembers(ctx, businessID)
}

func (s *BusinessService) Leave(ctx context.Context, userID, businessID string) error {
	member, err := s.store.GetMemberByUser(ctx, businessID, userID)
	if err != nil {
		return errors.New("دسترسی به کسب‌وکار وجود ندارد")
	}
	if member.Role == domain.RoleOwner {
		return errors.New("مالک قبل از خروج باید مالکیت را منتقل کند یا کسب‌وکار را حذف کند")
	}
	member.Status = domain.MemberRemoved
	member.Permissions = nil
	_, err = s.store.UpdateMember(ctx, member)
	return err
}

func (s *BusinessService) UpdateMember(ctx context.Context, actorID, businessID, memberID string, role domain.Role, commission float64, status domain.MemberStatus) (domain.BusinessMember, error) {
	actor, err := s.store.GetMemberByUser(ctx, businessID, actorID)
	if err != nil {
		return domain.BusinessMember{}, errors.New("دسترسی به کسب‌وکار وجود ندارد")
	}
	if actor.Role != domain.RoleOwner && !domain.HasPermission(actor, domain.PermMembersManage) {
		return domain.BusinessMember{}, errors.New("شما دسترسی مدیریت مشاورین را ندارید")
	}
	members, _ := s.store.ListMembers(ctx, businessID)
	var target domain.BusinessMember
	for _, member := range members {
		if member.ID == memberID {
			target = member
			break
		}
	}
	if target.ID == "" {
		return domain.BusinessMember{}, repository.ErrNotFound
	}
	if role == domain.RoleManager && actor.Role != domain.RoleOwner {
		return domain.BusinessMember{}, errors.New("فقط مالک می‌تواند مدیر تعریف کند")
	}
	if target.Role == domain.RoleOwner && actor.Role != domain.RoleOwner {
		return domain.BusinessMember{}, errors.New("امکان تغییر مالک وجود ندارد")
	}
	if role != "" {
		target.Role = role
		target.Permissions = domain.DefaultPermissions(role)
	}
	if commission >= 0 {
		target.CommissionPercent = commission
	}
	if status != "" {
		target.Status = status
	}
	return s.store.UpdateMember(ctx, target)
}
