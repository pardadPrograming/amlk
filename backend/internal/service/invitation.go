package service

import (
	"context"
	"errors"
	"time"

	"amlakcrm/backend/internal/domain"
	"amlakcrm/backend/internal/repository"
)

type InvitationService struct {
	store repository.Store
}

func NewInvitationService(store repository.Store) *InvitationService {
	return &InvitationService{store: store}
}

func (s *InvitationService) Create(ctx context.Context, actorID, businessID, phone string, role domain.Role, commission float64) (domain.Invitation, error) {
	actor, err := s.store.GetMemberByUser(ctx, businessID, actorID)
	if err != nil {
		return domain.Invitation{}, errors.New("دسترسی به کسب‌وکار وجود ندارد")
	}
	if actor.Role != domain.RoleOwner && !domain.HasPermission(actor, domain.PermInvitationsManage) {
		return domain.Invitation{}, errors.New("شما دسترسی ارسال دعوت‌نامه را ندارید")
	}
	if role == domain.RoleOwner {
		return domain.Invitation{}, errors.New("دعوت با نقش مالک مجاز نیست")
	}
	if role == domain.RoleManager && actor.Role != domain.RoleOwner {
		return domain.Invitation{}, errors.New("فقط مالک می‌تواند مدیر دعوت کند")
	}
	normalized, err := NormalizePhone(phone)
	if err != nil {
		return domain.Invitation{}, err
	}
	business, err := s.store.GetBusiness(ctx, businessID)
	if err != nil {
		return domain.Invitation{}, err
	}
	invite := domain.Invitation{
		BusinessID:        businessID,
		BusinessName:      business.Name,
		InviterUserID:     actorID,
		InviteePhone:      normalized,
		Role:              role,
		CommissionPercent: commission,
		Status:            domain.InvitationPending,
		ExpiresAt:         time.Now().UTC().Add(7 * 24 * time.Hour),
	}
	return s.store.CreateInvitation(ctx, invite)
}

func (s *InvitationService) ListForBusiness(ctx context.Context, businessID string) ([]domain.Invitation, error) {
	return s.store.ListInvitations(ctx, businessID)
}

func (s *InvitationService) Inbox(ctx context.Context, phone string) ([]domain.Invitation, error) {
	return s.store.ListInvitationsForPhone(ctx, phone)
}

func (s *InvitationService) Accept(ctx context.Context, user domain.User, invitationID string) (domain.Invitation, error) {
	invite, err := s.store.GetInvitation(ctx, invitationID)
	if err != nil {
		return domain.Invitation{}, err
	}
	if invite.InviteePhone != user.Phone {
		return domain.Invitation{}, errors.New("این دعوت‌نامه برای شماره شما نیست")
	}
	if invite.Status != domain.InvitationPending || time.Now().UTC().After(invite.ExpiresAt) {
		return domain.Invitation{}, errors.New("دعوت‌نامه فعال نیست")
	}
	_, _ = s.store.CreateMember(ctx, domain.BusinessMember{
		BusinessID:        invite.BusinessID,
		UserID:            user.ID,
		UserPhone:         user.Phone,
		UserDisplayName:   user.DisplayName,
		Role:              invite.Role,
		Permissions:       domain.DefaultPermissions(invite.Role),
		CommissionPercent: invite.CommissionPercent,
		Status:            domain.MemberActive,
	})
	invite.Status = domain.InvitationAccepted
	invite.InviteeUserID = user.ID
	return s.store.UpdateInvitation(ctx, invite)
}

func (s *InvitationService) Reject(ctx context.Context, user domain.User, invitationID string) (domain.Invitation, error) {
	invite, err := s.store.GetInvitation(ctx, invitationID)
	if err != nil {
		return domain.Invitation{}, err
	}
	if invite.InviteePhone != user.Phone {
		return domain.Invitation{}, errors.New("این دعوت‌نامه برای شماره شما نیست")
	}
	if invite.Status != domain.InvitationPending {
		return domain.Invitation{}, errors.New("دعوت‌نامه فعال نیست")
	}
	invite.Status = domain.InvitationRejected
	invite.InviteeUserID = user.ID
	return s.store.UpdateInvitation(ctx, invite)
}
