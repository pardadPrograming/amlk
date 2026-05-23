package repository

import (
	"context"
	"log/slog"
	"time"

	"amlakcrm/backend/internal/cache"
	"amlakcrm/backend/internal/domain"
	"amlakcrm/backend/internal/events"
)

type CachedStoreOptions struct {
	TTL                 time.Duration
	Logger              *slog.Logger
	Publisher           events.Publisher
	UsersByIDCache      cache.Cache[domain.User]
	UserIDByPhoneCache  cache.Cache[string]
	SessionsByIDCache   cache.Cache[domain.Session]
	SessionByTokenCache cache.Cache[domain.Session]
	RefreshByIDCache    cache.Cache[string]
}

type CachedStore struct {
	primary        Store
	ttl            time.Duration
	logger         *slog.Logger
	publisher      events.Publisher
	usersByID      cache.Cache[domain.User]
	userIDByPhone  cache.Cache[string]
	sessionsByID   cache.Cache[domain.Session]
	sessionByToken cache.Cache[domain.Session]
	refreshByID    cache.Cache[string]
}

func NewCachedStore(primary Store, opts CachedStoreOptions) *CachedStore {
	ttl := opts.TTL
	if ttl <= 0 {
		ttl = time.Minute
	}
	usersByID := opts.UsersByIDCache
	if usersByID == nil {
		usersByID = cache.NewTTLCache[domain.User]()
	}
	userIDByPhone := opts.UserIDByPhoneCache
	if userIDByPhone == nil {
		userIDByPhone = cache.NewTTLCache[string]()
	}
	sessionsByID := opts.SessionsByIDCache
	if sessionsByID == nil {
		sessionsByID = cache.NewTTLCache[domain.Session]()
	}
	sessionByToken := opts.SessionByTokenCache
	if sessionByToken == nil {
		sessionByToken = cache.NewTTLCache[domain.Session]()
	}
	refreshByID := opts.RefreshByIDCache
	if refreshByID == nil {
		refreshByID = cache.NewTTLCache[string]()
	}
	return &CachedStore{
		primary:        primary,
		ttl:            ttl,
		logger:         opts.Logger,
		publisher:      opts.Publisher,
		usersByID:      usersByID,
		userIDByPhone:  userIDByPhone,
		sessionsByID:   sessionsByID,
		sessionByToken: sessionByToken,
		refreshByID:    refreshByID,
	}
}

func (s *CachedStore) SaveOTP(ctx context.Context, challenge domain.OTPChallenge) error {
	return s.primary.SaveOTP(ctx, challenge)
}

func (s *CachedStore) GetOTP(ctx context.Context, phone string) (domain.OTPChallenge, error) {
	return s.primary.GetOTP(ctx, phone)
}

func (s *CachedStore) GetLatestOTP(ctx context.Context) (domain.OTPChallenge, error) {
	return s.primary.GetLatestOTP(ctx)
}

func (s *CachedStore) DeleteOTP(ctx context.Context, phone string) error {
	return s.primary.DeleteOTP(ctx, phone)
}

func (s *CachedStore) UpsertUserByPhone(ctx context.Context, phone string) (domain.User, error) {
	user, err := s.primary.UpsertUserByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	s.cacheUser(ctx, user)
	s.publish(ctx, "user.upserted", user.ID, map[string]string{"phone": user.Phone})
	return user, nil
}

func (s *CachedStore) GetUser(ctx context.Context, id string) (domain.User, error) {
	if user, ok := s.usersByID.Get(ctx, userKey(id)); ok {
		return user, nil
	}
	user, err := s.primary.GetUser(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	s.cacheUser(ctx, user)
	return user, nil
}

func (s *CachedStore) GetUserByPhone(ctx context.Context, phone string) (domain.User, error) {
	if id, ok := s.userIDByPhone.Get(ctx, userPhoneKey(phone)); ok {
		if user, ok := s.usersByID.Get(ctx, userKey(id)); ok {
			return user, nil
		}
	}
	user, err := s.primary.GetUserByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	s.cacheUser(ctx, user)
	return user, nil
}

func (s *CachedStore) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {
	updated, err := s.primary.UpdateUser(ctx, user)
	if err != nil {
		return domain.User{}, err
	}
	s.cacheUser(ctx, updated)
	s.publish(ctx, "user.updated", updated.ID, map[string]string{"phone": updated.Phone})
	return updated, nil
}

func (s *CachedStore) ListUsers(ctx context.Context) ([]domain.User, error) {
	return s.primary.ListUsers(ctx)
}

func (s *CachedStore) SaveSession(ctx context.Context, session domain.Session) error {
	if err := s.primary.SaveSession(ctx, session); err != nil {
		return err
	}
	s.cacheSession(ctx, session)
	s.publish(ctx, "session.saved", session.ID, map[string]string{"userId": session.UserID})
	return nil
}

func (s *CachedStore) GetSessionByRefresh(ctx context.Context, refresh string) (domain.Session, error) {
	if session, ok := s.sessionByToken.Get(ctx, sessionRefreshKey(refresh)); ok {
		if session.RevokedAt.IsZero() && time.Now().UTC().Before(session.ExpiresAt) {
			return session, nil
		}
		s.evictSession(ctx, session)
	}
	session, err := s.primary.GetSessionByRefresh(ctx, refresh)
	if err != nil {
		return domain.Session{}, err
	}
	s.cacheSession(ctx, session)
	return session, nil
}

func (s *CachedStore) ListSessions(ctx context.Context, userID string) ([]domain.Session, error) {
	return s.primary.ListSessions(ctx, userID)
}

func (s *CachedStore) GetSession(ctx context.Context, sessionID string) (domain.Session, error) {
	if session, ok := s.sessionsByID.Get(ctx, sessionKey(sessionID)); ok {
		if session.RevokedAt.IsZero() && time.Now().UTC().Before(session.ExpiresAt) {
			return session, nil
		}
		s.evictSession(ctx, session)
	}
	session, err := s.primary.GetSession(ctx, sessionID)
	if err != nil {
		return domain.Session{}, err
	}
	s.cacheSession(ctx, session)
	return session, nil
}

func (s *CachedStore) RevokeSession(ctx context.Context, refresh string) error {
	session, _ := s.GetSessionByRefresh(ctx, refresh)
	if err := s.primary.RevokeSession(ctx, refresh); err != nil {
		return err
	}
	if session.ID != "" {
		s.evictSession(ctx, session)
		s.publish(ctx, "session.revoked", session.ID, map[string]string{"userId": session.UserID})
	}
	return nil
}

func (s *CachedStore) RevokeSessionByID(ctx context.Context, userID string, sessionID string) error {
	session, _ := s.GetSession(ctx, sessionID)
	if err := s.primary.RevokeSessionByID(ctx, userID, sessionID); err != nil {
		return err
	}
	if session.ID != "" {
		s.evictSession(ctx, session)
	}
	s.publish(ctx, "session.revoked", sessionID, map[string]string{"userId": userID})
	return nil
}

func (s *CachedStore) TouchSession(ctx context.Context, sessionID string, at time.Time) error {
	if err := s.primary.TouchSession(ctx, sessionID, at); err != nil {
		return err
	}
	session, err := s.primary.GetSession(ctx, sessionID)
	if err != nil {
		return nil
	}
	s.cacheSession(ctx, session)
	s.publish(ctx, "session.touched", session.ID, map[string]interface{}{"userId": session.UserID, "lastSeenAt": session.LastSeenAt})
	return nil
}

func (s *CachedStore) CreateBusiness(ctx context.Context, business domain.Business) (domain.Business, error) {
	return s.primary.CreateBusiness(ctx, business)
}

func (s *CachedStore) ListBusinessesForUser(ctx context.Context, userID string) ([]domain.Business, error) {
	return s.primary.ListBusinessesForUser(ctx, userID)
}

func (s *CachedStore) GetBusiness(ctx context.Context, id string) (domain.Business, error) {
	return s.primary.GetBusiness(ctx, id)
}

func (s *CachedStore) UpdateBusiness(ctx context.Context, business domain.Business) (domain.Business, error) {
	return s.primary.UpdateBusiness(ctx, business)
}

func (s *CachedStore) ListBusinesses(ctx context.Context) ([]domain.Business, error) {
	return s.primary.ListBusinesses(ctx)
}

func (s *CachedStore) CreateMember(ctx context.Context, member domain.BusinessMember) (domain.BusinessMember, error) {
	return s.primary.CreateMember(ctx, member)
}

func (s *CachedStore) ListMembers(ctx context.Context, businessID string) ([]domain.BusinessMember, error) {
	return s.primary.ListMembers(ctx, businessID)
}

func (s *CachedStore) GetMemberByUser(ctx context.Context, businessID string, userID string) (domain.BusinessMember, error) {
	return s.primary.GetMemberByUser(ctx, businessID, userID)
}

func (s *CachedStore) UpdateMember(ctx context.Context, member domain.BusinessMember) (domain.BusinessMember, error) {
	return s.primary.UpdateMember(ctx, member)
}

func (s *CachedStore) CreateInvitation(ctx context.Context, invitation domain.Invitation) (domain.Invitation, error) {
	return s.primary.CreateInvitation(ctx, invitation)
}

func (s *CachedStore) ListInvitations(ctx context.Context, businessID string) ([]domain.Invitation, error) {
	return s.primary.ListInvitations(ctx, businessID)
}

func (s *CachedStore) ListInvitationsForPhone(ctx context.Context, phone string) ([]domain.Invitation, error) {
	return s.primary.ListInvitationsForPhone(ctx, phone)
}

func (s *CachedStore) GetInvitation(ctx context.Context, id string) (domain.Invitation, error) {
	return s.primary.GetInvitation(ctx, id)
}

func (s *CachedStore) UpdateInvitation(ctx context.Context, invitation domain.Invitation) (domain.Invitation, error) {
	return s.primary.UpdateInvitation(ctx, invitation)
}

func (s *CachedStore) CreateFile(ctx context.Context, file domain.FileObject) (domain.FileObject, error) {
	return s.primary.CreateFile(ctx, file)
}

func (s *CachedStore) GetFile(ctx context.Context, fileID string) (domain.FileObject, error) {
	return s.primary.GetFile(ctx, fileID)
}

func (s *CachedStore) UpdateFile(ctx context.Context, file domain.FileObject) (domain.FileObject, error) {
	return s.primary.UpdateFile(ctx, file)
}

func (s *CachedStore) ListExpiredFiles(ctx context.Context, before time.Time, limit int) ([]domain.FileObject, error) {
	return s.primary.ListExpiredFiles(ctx, before, limit)
}

func (s *CachedStore) DeleteFile(ctx context.Context, fileID string) error {
	return s.primary.DeleteFile(ctx, fileID)
}

func (s *CachedStore) EnsureUserMainChannel(ctx context.Context, user domain.User) (domain.Channel, error) {
	return s.primary.EnsureUserMainChannel(ctx, user)
}

func (s *CachedStore) EnsureBusinessMainChannel(ctx context.Context, business domain.Business) (domain.Channel, error) {
	return s.primary.EnsureBusinessMainChannel(ctx, business)
}

func (s *CachedStore) EnsurePrivateChannel(ctx context.Context, actor domain.User, target domain.User) (domain.Channel, error) {
	return s.primary.EnsurePrivateChannel(ctx, actor, target)
}

func (s *CachedStore) EnsureUserMainVault(ctx context.Context, user domain.User) (domain.ChannelVault, error) {
	return s.primary.EnsureUserMainVault(ctx, user)
}

func (s *CachedStore) EnsureBusinessMainVault(ctx context.Context, business domain.Business) (domain.ChannelVault, error) {
	return s.primary.EnsureBusinessMainVault(ctx, business)
}

func (s *CachedStore) CreateUserVault(ctx context.Context, user domain.User, title string) (domain.ChannelVault, error) {
	return s.primary.CreateUserVault(ctx, user, title)
}

func (s *CachedStore) CreateBusinessVault(ctx context.Context, business domain.Business, title string) (domain.ChannelVault, error) {
	return s.primary.CreateBusinessVault(ctx, business, title)
}

func (s *CachedStore) GetChannelVault(ctx context.Context, vaultID string) (domain.ChannelVault, error) {
	return s.primary.GetChannelVault(ctx, vaultID)
}

func (s *CachedStore) ListUserVaults(ctx context.Context, userID string) ([]domain.ChannelVault, error) {
	return s.primary.ListUserVaults(ctx, userID)
}

func (s *CachedStore) ListBusinessVaults(ctx context.Context, businessID string) ([]domain.ChannelVault, error) {
	return s.primary.ListBusinessVaults(ctx, businessID)
}

func (s *CachedStore) GetChannel(ctx context.Context, channelID string) (domain.Channel, error) {
	return s.primary.GetChannel(ctx, channelID)
}

func (s *CachedStore) ListChannelsForUser(ctx context.Context, userID string, phones []string, businessIDs []string) ([]domain.Channel, error) {
	return s.primary.ListChannelsForUser(ctx, userID, phones, businessIDs)
}

func (s *CachedStore) UpsertChannelMember(ctx context.Context, member domain.ChannelMember) (domain.ChannelMember, error) {
	return s.primary.UpsertChannelMember(ctx, member)
}

func (s *CachedStore) GetChannelMember(ctx context.Context, channelID string, userID string) (domain.ChannelMember, error) {
	return s.primary.GetChannelMember(ctx, channelID, userID)
}

func (s *CachedStore) ListChannelMembers(ctx context.Context, channelID string) ([]domain.ChannelMember, error) {
	return s.primary.ListChannelMembers(ctx, channelID)
}

func (s *CachedStore) CreateChannelMessage(ctx context.Context, message domain.ChannelMessage) (domain.ChannelMessage, error) {
	return s.primary.CreateChannelMessage(ctx, message)
}

func (s *CachedStore) GetChannelMessage(ctx context.Context, channelID string, messageID string) (domain.ChannelMessage, error) {
	return s.primary.GetChannelMessage(ctx, channelID, messageID)
}

func (s *CachedStore) UpdateChannelMessage(ctx context.Context, message domain.ChannelMessage) (domain.ChannelMessage, error) {
	return s.primary.UpdateChannelMessage(ctx, message)
}

func (s *CachedStore) DeleteChannelMessage(ctx context.Context, channelID string, messageID string) error {
	return s.primary.DeleteChannelMessage(ctx, channelID, messageID)
}

func (s *CachedStore) ListChannelMessages(ctx context.Context, channelID string, limit int, offset int) ([]domain.ChannelMessage, int, error) {
	return s.primary.ListChannelMessages(ctx, channelID, limit, offset)
}

func (s *CachedStore) ChannelUnreadSummary(ctx context.Context, channelID string, userID string) (domain.ChannelUnreadSummary, error) {
	return s.primary.ChannelUnreadSummary(ctx, channelID, userID)
}

func (s *CachedStore) MarkChannelMessagesSeen(ctx context.Context, channelID string, userID string, messageIDs []string) error {
	return s.primary.MarkChannelMessagesSeen(ctx, channelID, userID, messageIDs)
}

func (s *CachedStore) CreateChannelVaultFile(ctx context.Context, file domain.ChannelVaultFile) (domain.ChannelVaultFile, error) {
	return s.primary.CreateChannelVaultFile(ctx, file)
}

func (s *CachedStore) UpsertChannelVaultFile(ctx context.Context, file domain.ChannelVaultFile) (domain.ChannelVaultFile, error) {
	return s.primary.UpsertChannelVaultFile(ctx, file)
}

func (s *CachedStore) DeletePropertyVaultReferences(ctx context.Context, propertyFileID string, keepVaultIDs []string) error {
	return s.primary.DeletePropertyVaultReferences(ctx, propertyFileID, keepVaultIDs)
}

func (s *CachedStore) GetChannelVaultFile(ctx context.Context, channelID string, fileID string) (domain.ChannelVaultFile, error) {
	return s.primary.GetChannelVaultFile(ctx, channelID, fileID)
}

func (s *CachedStore) ListChannelVaultFiles(ctx context.Context, channelID string, limit int, offset int) ([]domain.ChannelVaultFile, int, error) {
	return s.primary.ListChannelVaultFiles(ctx, channelID, limit, offset)
}

func (s *CachedStore) CreateArea(ctx context.Context, area domain.Area) (domain.Area, error) {
	return s.primary.CreateArea(ctx, area)
}

func (s *CachedStore) ListAreas(ctx context.Context, businessID string) ([]domain.Area, error) {
	return s.primary.ListAreas(ctx, businessID)
}

func (s *CachedStore) DeleteArea(ctx context.Context, businessID string, areaID string) error {
	return s.primary.DeleteArea(ctx, businessID, areaID)
}

func (s *CachedStore) CreateStreet(ctx context.Context, street domain.Street) (domain.Street, error) {
	return s.primary.CreateStreet(ctx, street)
}

func (s *CachedStore) ListStreets(ctx context.Context, businessID string, areaID string) ([]domain.Street, error) {
	return s.primary.ListStreets(ctx, businessID, areaID)
}

func (s *CachedStore) DeleteStreet(ctx context.Context, businessID string, areaID string, streetID string) error {
	return s.primary.DeleteStreet(ctx, businessID, areaID, streetID)
}

func (s *CachedStore) CreateNeighborhood(ctx context.Context, neighborhood domain.Neighborhood) (domain.Neighborhood, error) {
	return s.primary.CreateNeighborhood(ctx, neighborhood)
}

func (s *CachedStore) ListNeighborhoods(ctx context.Context, businessID string, areaID string, streetID string) ([]domain.Neighborhood, error) {
	return s.primary.ListNeighborhoods(ctx, businessID, areaID, streetID)
}

func (s *CachedStore) DeleteNeighborhood(ctx context.Context, businessID string, areaID string, streetID string, neighborhoodID string) error {
	return s.primary.DeleteNeighborhood(ctx, businessID, areaID, streetID, neighborhoodID)
}

func (s *CachedStore) CreatePropertyFile(ctx context.Context, file domain.PropertyFile) (domain.PropertyFile, error) {
	item, err := s.primary.CreatePropertyFile(ctx, file)
	if err == nil {
		s.publishSearchInvalidate(ctx, item.BusinessID, item.OwnerUserID)
		s.publishPropertyChanged(ctx, item.BusinessID, item.OwnerUserID, item.ID)
	}
	return item, err
}

func (s *CachedStore) CreateNotification(ctx context.Context, notification domain.Notification) (domain.Notification, error) {
	return s.primary.CreateNotification(ctx, notification)
}

func (s *CachedStore) ListNotifications(ctx context.Context, userID string, unreadOnly bool) ([]domain.Notification, error) {
	return s.primary.ListNotifications(ctx, userID, unreadOnly)
}

func (s *CachedStore) MarkNotificationRead(ctx context.Context, userID string, notificationID string) error {
	return s.primary.MarkNotificationRead(ctx, userID, notificationID)
}

func (s *CachedStore) ListPropertyFiles(ctx context.Context, businessID string) ([]domain.PropertyFile, error) {
	return s.primary.ListPropertyFiles(ctx, businessID)
}

func (s *CachedStore) ListPropertyFilesForOwner(ctx context.Context, businessID string, ownerUserID string) ([]domain.PropertyFile, error) {
	return s.primary.ListPropertyFilesForOwner(ctx, businessID, ownerUserID)
}

func (s *CachedStore) ListPropertyFilesForAccess(ctx context.Context, businessID string, ownerUserID string, vaultIDs []string) ([]domain.PropertyFile, error) {
	return s.primary.ListPropertyFilesForAccess(ctx, businessID, ownerUserID, vaultIDs)
}

func (s *CachedStore) ListLatestPropertyFiles(ctx context.Context, businessID string, ownerUserID string, vaultIDs []string, propertyType domain.PropertyFileType, limit int, offset int) ([]domain.PropertyFile, int, error) {
	return s.primary.ListLatestPropertyFiles(ctx, businessID, ownerUserID, vaultIDs, propertyType, limit, offset)
}

func (s *CachedStore) GetPropertyFile(ctx context.Context, businessID string, fileID string) (domain.PropertyFile, error) {
	return s.primary.GetPropertyFile(ctx, businessID, fileID)
}

func (s *CachedStore) UpdatePropertyFile(ctx context.Context, file domain.PropertyFile) (domain.PropertyFile, error) {
	item, err := s.primary.UpdatePropertyFile(ctx, file)
	if err == nil {
		s.publishSearchInvalidate(ctx, item.BusinessID, item.OwnerUserID)
		s.publishPropertyChanged(ctx, item.BusinessID, item.OwnerUserID, item.ID)
	}
	return item, err
}

func (s *CachedStore) CreatePropertyShareRequest(ctx context.Context, request domain.PropertyShareRequest) (domain.PropertyShareRequest, error) {
	return s.primary.CreatePropertyShareRequest(ctx, request)
}

func (s *CachedStore) GetPropertyShareRequest(ctx context.Context, businessID string, requestID string) (domain.PropertyShareRequest, error) {
	return s.primary.GetPropertyShareRequest(ctx, businessID, requestID)
}

func (s *CachedStore) UpdatePropertyShareRequest(ctx context.Context, request domain.PropertyShareRequest) (domain.PropertyShareRequest, error) {
	return s.primary.UpdatePropertyShareRequest(ctx, request)
}

func (s *CachedStore) ListPropertyShareRequestsForOwner(ctx context.Context, businessID string, ownerUserID string) ([]domain.PropertyShareRequest, error) {
	return s.primary.ListPropertyShareRequestsForOwner(ctx, businessID, ownerUserID)
}

func (s *CachedStore) ListPropertyShareRequestsForRequester(ctx context.Context, businessID string, requesterUserID string) ([]domain.PropertyShareRequest, error) {
	return s.primary.ListPropertyShareRequestsForRequester(ctx, businessID, requesterUserID)
}

func (s *CachedStore) CreateContact(ctx context.Context, contact domain.Contact) (domain.Contact, error) {
	item, err := s.primary.CreateContact(ctx, contact)
	if err == nil {
		s.publishSearchInvalidate(ctx, item.BusinessID, item.CreatedByID)
	}
	return item, err
}

func (s *CachedStore) ListContacts(ctx context.Context, businessID string) ([]domain.Contact, error) {
	return s.primary.ListContacts(ctx, businessID)
}

func (s *CachedStore) GetContact(ctx context.Context, businessID, contactID string) (domain.Contact, error) {
	return s.primary.GetContact(ctx, businessID, contactID)
}

func (s *CachedStore) UpdateContact(ctx context.Context, contact domain.Contact) (domain.Contact, error) {
	item, err := s.primary.UpdateContact(ctx, contact)
	if err == nil {
		s.publishSearchInvalidate(ctx, item.BusinessID, item.CreatedByID)
	}
	return item, err
}

func (s *CachedStore) publishSearchInvalidate(ctx context.Context, businessID, userID string) {
	events.SafePublish(ctx, s.logger, s.publisher, events.Event{
		Type:        events.SearchInvalidateEvent,
		AggregateID: businessID + ":" + userID,
		OccurredAt:  time.Now().UTC(),
		Payload: events.SearchInvalidatePayload{
			BusinessID: businessID,
			UserID:     userID,
		},
	})
}

func (s *CachedStore) publishPropertyChanged(ctx context.Context, businessID, userID, propertyID string) {
	events.SafePublish(ctx, s.logger, s.publisher, events.Event{
		Type:        events.PropertyChangedEvent,
		AggregateID: businessID + ":" + propertyID,
		OccurredAt:  time.Now().UTC(),
		Payload: events.PropertyChangedPayload{
			BusinessID: businessID,
			UserID:     userID,
			PropertyID: propertyID,
		},
	})
}

func (s *CachedStore) SaveAdminAccount(ctx context.Context, account domain.AdminAccount) (domain.AdminAccount, error) {
	return s.primary.SaveAdminAccount(ctx, account)
}

func (s *CachedStore) GetAdminAccountByUser(ctx context.Context, userID string) (domain.AdminAccount, error) {
	return s.primary.GetAdminAccountByUser(ctx, userID)
}

func (s *CachedStore) ListAdminAccounts(ctx context.Context) ([]domain.AdminAccount, error) {
	return s.primary.ListAdminAccounts(ctx)
}

func (s *CachedStore) GetPlatformSettings(ctx context.Context) (domain.PlatformSettings, error) {
	return s.primary.GetPlatformSettings(ctx)
}

func (s *CachedStore) SavePlatformSettings(ctx context.Context, settings domain.PlatformSettings) (domain.PlatformSettings, error) {
	return s.primary.SavePlatformSettings(ctx, settings)
}

func (s *CachedStore) CreateCity(ctx context.Context, city domain.City) (domain.City, error) {
	return s.primary.CreateCity(ctx, city)
}

func (s *CachedStore) ListCities(ctx context.Context) ([]domain.City, error) {
	return s.primary.ListCities(ctx)
}

func (s *CachedStore) GetCity(ctx context.Context, cityID string) (domain.City, error) {
	return s.primary.GetCity(ctx, cityID)
}

func (s *CachedStore) CreateSystemArea(ctx context.Context, area domain.SystemArea) (domain.SystemArea, error) {
	return s.primary.CreateSystemArea(ctx, area)
}

func (s *CachedStore) ListSystemAreas(ctx context.Context, cityID string) ([]domain.SystemArea, error) {
	return s.primary.ListSystemAreas(ctx, cityID)
}

func (s *CachedStore) GetSystemArea(ctx context.Context, areaID string) (domain.SystemArea, error) {
	return s.primary.GetSystemArea(ctx, areaID)
}

func (s *CachedStore) CreateSystemStreet(ctx context.Context, street domain.SystemStreet) (domain.SystemStreet, error) {
	return s.primary.CreateSystemStreet(ctx, street)
}

func (s *CachedStore) ListSystemStreets(ctx context.Context, cityID, areaID string) ([]domain.SystemStreet, error) {
	return s.primary.ListSystemStreets(ctx, cityID, areaID)
}

func (s *CachedStore) GetSystemStreet(ctx context.Context, streetID string) (domain.SystemStreet, error) {
	return s.primary.GetSystemStreet(ctx, streetID)
}

func (s *CachedStore) CreateSystemNeighborhood(ctx context.Context, neighborhood domain.SystemNeighborhood) (domain.SystemNeighborhood, error) {
	return s.primary.CreateSystemNeighborhood(ctx, neighborhood)
}

func (s *CachedStore) ListSystemNeighborhoods(ctx context.Context, cityID, areaID, streetID string) ([]domain.SystemNeighborhood, error) {
	return s.primary.ListSystemNeighborhoods(ctx, cityID, areaID, streetID)
}

func (s *CachedStore) GetSystemNeighborhood(ctx context.Context, neighborhoodID string) (domain.SystemNeighborhood, error) {
	return s.primary.GetSystemNeighborhood(ctx, neighborhoodID)
}

func (s *CachedStore) CreateLocationSuggestion(ctx context.Context, suggestion domain.LocationSuggestion) (domain.LocationSuggestion, error) {
	return s.primary.CreateLocationSuggestion(ctx, suggestion)
}

func (s *CachedStore) ListLocationSuggestions(ctx context.Context, status domain.LocationSuggestionStatus) ([]domain.LocationSuggestion, error) {
	return s.primary.ListLocationSuggestions(ctx, status)
}

func (s *CachedStore) GetLocationSuggestion(ctx context.Context, suggestionID string) (domain.LocationSuggestion, error) {
	return s.primary.GetLocationSuggestion(ctx, suggestionID)
}

func (s *CachedStore) UpdateLocationSuggestion(ctx context.Context, suggestion domain.LocationSuggestion) (domain.LocationSuggestion, error) {
	return s.primary.UpdateLocationSuggestion(ctx, suggestion)
}

func (s *CachedStore) cacheUser(ctx context.Context, user domain.User) {
	s.usersByID.Set(ctx, userKey(user.ID), user, s.ttl)
	s.userIDByPhone.Set(ctx, userPhoneKey(user.Phone), user.ID, s.ttl)
}

func (s *CachedStore) cacheSession(ctx context.Context, session domain.Session) {
	if !session.RevokedAt.IsZero() || time.Now().UTC().After(session.ExpiresAt) {
		s.evictSession(ctx, session)
		return
	}
	s.sessionsByID.Set(ctx, sessionKey(session.ID), session, s.ttl)
	s.sessionByToken.Set(ctx, sessionRefreshKey(session.RefreshToken), session, s.ttl)
	s.refreshByID.Set(ctx, sessionRefreshIDKey(session.ID), session.RefreshToken, s.ttl)
}

func (s *CachedStore) evictSession(ctx context.Context, session domain.Session) {
	s.sessionsByID.Delete(ctx, sessionKey(session.ID))
	if session.RefreshToken != "" {
		s.sessionByToken.Delete(ctx, sessionRefreshKey(session.RefreshToken))
	}
	if refresh, ok := s.refreshByID.Get(ctx, sessionRefreshIDKey(session.ID)); ok {
		s.sessionByToken.Delete(ctx, sessionRefreshKey(refresh))
	}
	s.refreshByID.Delete(ctx, sessionRefreshIDKey(session.ID))
}

func (s *CachedStore) publish(ctx context.Context, eventType, aggregateID string, payload interface{}) {
	events.SafePublish(ctx, s.logger, s.publisher, events.Event{
		Type:        eventType,
		AggregateID: aggregateID,
		OccurredAt:  time.Now().UTC(),
		Payload:     payload,
	})
}

func userKey(id string) string              { return "user:id:" + id }
func userPhoneKey(phone string) string      { return "user:phone:" + phone }
func sessionKey(id string) string           { return "session:id:" + id }
func sessionRefreshKey(token string) string { return "session:refresh:" + token }
func sessionRefreshIDKey(id string) string  { return "session:refresh-id:" + id }
