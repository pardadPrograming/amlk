package repository

import (
	"context"
	"time"

	"amlakcrm/backend/internal/domain"
)

type Store interface {
	SaveOTP(ctx context.Context, challenge domain.OTPChallenge) error
	GetOTP(ctx context.Context, phone string) (domain.OTPChallenge, error)
	GetLatestOTP(ctx context.Context) (domain.OTPChallenge, error)
	DeleteOTP(ctx context.Context, phone string) error

	UpsertUserByPhone(ctx context.Context, phone string) (domain.User, error)
	GetUser(ctx context.Context, id string) (domain.User, error)
	GetUserByPhone(ctx context.Context, phone string) (domain.User, error)
	UpdateUser(ctx context.Context, user domain.User) (domain.User, error)
	ListUsers(ctx context.Context) ([]domain.User, error)

	SaveSession(ctx context.Context, session domain.Session) error
	GetSessionByRefresh(ctx context.Context, refresh string) (domain.Session, error)
	ListSessions(ctx context.Context, userID string) ([]domain.Session, error)
	GetSession(ctx context.Context, sessionID string) (domain.Session, error)
	RevokeSession(ctx context.Context, refresh string) error
	RevokeSessionByID(ctx context.Context, userID string, sessionID string) error
	TouchSession(ctx context.Context, sessionID string, at time.Time) error

	CreateBusiness(ctx context.Context, business domain.Business) (domain.Business, error)
	ListBusinessesForUser(ctx context.Context, userID string) ([]domain.Business, error)
	GetBusiness(ctx context.Context, id string) (domain.Business, error)
	UpdateBusiness(ctx context.Context, business domain.Business) (domain.Business, error)
	ListBusinesses(ctx context.Context) ([]domain.Business, error)

	CreateMember(ctx context.Context, member domain.BusinessMember) (domain.BusinessMember, error)
	ListMembers(ctx context.Context, businessID string) ([]domain.BusinessMember, error)
	GetMemberByUser(ctx context.Context, businessID string, userID string) (domain.BusinessMember, error)
	UpdateMember(ctx context.Context, member domain.BusinessMember) (domain.BusinessMember, error)

	CreateInvitation(ctx context.Context, invitation domain.Invitation) (domain.Invitation, error)
	ListInvitations(ctx context.Context, businessID string) ([]domain.Invitation, error)
	ListInvitationsForPhone(ctx context.Context, phone string) ([]domain.Invitation, error)
	GetInvitation(ctx context.Context, id string) (domain.Invitation, error)
	UpdateInvitation(ctx context.Context, invitation domain.Invitation) (domain.Invitation, error)

	CreateFile(ctx context.Context, file domain.FileObject) (domain.FileObject, error)
	GetFile(ctx context.Context, fileID string) (domain.FileObject, error)
	UpdateFile(ctx context.Context, file domain.FileObject) (domain.FileObject, error)
	ListExpiredFiles(ctx context.Context, before time.Time, limit int) ([]domain.FileObject, error)
	DeleteFile(ctx context.Context, fileID string) error
	CreateNotification(ctx context.Context, notification domain.Notification) (domain.Notification, error)
	ListNotifications(ctx context.Context, userID string, unreadOnly bool) ([]domain.Notification, error)
	MarkNotificationRead(ctx context.Context, userID string, notificationID string) error

	EnsureUserMainChannel(ctx context.Context, user domain.User) (domain.Channel, error)
	EnsureBusinessMainChannel(ctx context.Context, business domain.Business) (domain.Channel, error)
	EnsurePrivateChannel(ctx context.Context, actor domain.User, target domain.User) (domain.Channel, error)
	EnsureUserMainVault(ctx context.Context, user domain.User) (domain.ChannelVault, error)
	EnsureBusinessMainVault(ctx context.Context, business domain.Business) (domain.ChannelVault, error)
	CreateUserVault(ctx context.Context, user domain.User, title string) (domain.ChannelVault, error)
	CreateBusinessVault(ctx context.Context, business domain.Business, title string) (domain.ChannelVault, error)
	GetChannelVault(ctx context.Context, vaultID string) (domain.ChannelVault, error)
	ListUserVaults(ctx context.Context, userID string) ([]domain.ChannelVault, error)
	ListBusinessVaults(ctx context.Context, businessID string) ([]domain.ChannelVault, error)
	GetChannel(ctx context.Context, channelID string) (domain.Channel, error)
	ListChannelsForUser(ctx context.Context, userID string, phones []string, businessIDs []string) ([]domain.Channel, error)
	UpsertChannelMember(ctx context.Context, member domain.ChannelMember) (domain.ChannelMember, error)
	GetChannelMember(ctx context.Context, channelID string, userID string) (domain.ChannelMember, error)
	ListChannelMembers(ctx context.Context, channelID string) ([]domain.ChannelMember, error)
	CreateChannelMessage(ctx context.Context, message domain.ChannelMessage) (domain.ChannelMessage, error)
	GetChannelMessage(ctx context.Context, channelID string, messageID string) (domain.ChannelMessage, error)
	UpdateChannelMessage(ctx context.Context, message domain.ChannelMessage) (domain.ChannelMessage, error)
	DeleteChannelMessage(ctx context.Context, channelID string, messageID string) error
	ListChannelMessages(ctx context.Context, channelID string, limit int, offset int) ([]domain.ChannelMessage, int, error)
	ChannelUnreadSummary(ctx context.Context, channelID string, userID string) (domain.ChannelUnreadSummary, error)
	MarkChannelMessagesSeen(ctx context.Context, channelID string, userID string, messageIDs []string) error
	CreateChannelVaultFile(ctx context.Context, file domain.ChannelVaultFile) (domain.ChannelVaultFile, error)
	UpsertChannelVaultFile(ctx context.Context, file domain.ChannelVaultFile) (domain.ChannelVaultFile, error)
	DeletePropertyVaultReferences(ctx context.Context, propertyFileID string, keepVaultIDs []string) error
	GetChannelVaultFile(ctx context.Context, channelID string, fileID string) (domain.ChannelVaultFile, error)
	ListChannelVaultFiles(ctx context.Context, channelID string, limit int, offset int) ([]domain.ChannelVaultFile, int, error)

	CreateArea(ctx context.Context, area domain.Area) (domain.Area, error)
	ListAreas(ctx context.Context, businessID string) ([]domain.Area, error)
	DeleteArea(ctx context.Context, businessID string, areaID string) error
	CreateStreet(ctx context.Context, street domain.Street) (domain.Street, error)
	ListStreets(ctx context.Context, businessID string, areaID string) ([]domain.Street, error)
	DeleteStreet(ctx context.Context, businessID string, areaID string, streetID string) error
	CreateNeighborhood(ctx context.Context, neighborhood domain.Neighborhood) (domain.Neighborhood, error)
	ListNeighborhoods(ctx context.Context, businessID string, areaID string, streetID string) ([]domain.Neighborhood, error)
	DeleteNeighborhood(ctx context.Context, businessID string, areaID string, streetID string, neighborhoodID string) error

	CreatePropertyFile(ctx context.Context, file domain.PropertyFile) (domain.PropertyFile, error)
	ListPropertyFiles(ctx context.Context, businessID string) ([]domain.PropertyFile, error)
	ListPropertyFilesForOwner(ctx context.Context, businessID string, ownerUserID string) ([]domain.PropertyFile, error)
	ListPropertyFilesForAccess(ctx context.Context, businessID string, ownerUserID string, vaultIDs []string) ([]domain.PropertyFile, error)
	ListLatestPropertyFiles(ctx context.Context, businessID string, ownerUserID string, vaultIDs []string, propertyType domain.PropertyFileType, limit int, offset int) ([]domain.PropertyFile, int, error)
	GetPropertyFile(ctx context.Context, businessID string, fileID string) (domain.PropertyFile, error)
	UpdatePropertyFile(ctx context.Context, file domain.PropertyFile) (domain.PropertyFile, error)
	CreatePropertyShareRequest(ctx context.Context, request domain.PropertyShareRequest) (domain.PropertyShareRequest, error)
	GetPropertyShareRequest(ctx context.Context, businessID string, requestID string) (domain.PropertyShareRequest, error)
	UpdatePropertyShareRequest(ctx context.Context, request domain.PropertyShareRequest) (domain.PropertyShareRequest, error)
	ListPropertyShareRequestsForOwner(ctx context.Context, businessID string, ownerUserID string) ([]domain.PropertyShareRequest, error)
	ListPropertyShareRequestsForRequester(ctx context.Context, businessID string, requesterUserID string) ([]domain.PropertyShareRequest, error)
	CreatePropertyOffer(ctx context.Context, offer domain.PropertyOffer) (domain.PropertyOffer, error)
	GetPropertyOffer(ctx context.Context, businessID string, offerID string) (domain.PropertyOffer, error)
	UpdatePropertyOffer(ctx context.Context, offer domain.PropertyOffer) (domain.PropertyOffer, error)
	ListPropertyOffersForUser(ctx context.Context, businessID string, userID string, scope string) ([]domain.PropertyOffer, error)

	CreateContact(ctx context.Context, contact domain.Contact) (domain.Contact, error)
	ListContacts(ctx context.Context, businessID string) ([]domain.Contact, error)
	GetContact(ctx context.Context, businessID, contactID string) (domain.Contact, error)
	UpdateContact(ctx context.Context, contact domain.Contact) (domain.Contact, error)

	SaveAdminAccount(ctx context.Context, account domain.AdminAccount) (domain.AdminAccount, error)
	GetAdminAccountByUser(ctx context.Context, userID string) (domain.AdminAccount, error)
	ListAdminAccounts(ctx context.Context) ([]domain.AdminAccount, error)
	GetPlatformSettings(ctx context.Context) (domain.PlatformSettings, error)
	SavePlatformSettings(ctx context.Context, settings domain.PlatformSettings) (domain.PlatformSettings, error)

	CreateCity(ctx context.Context, city domain.City) (domain.City, error)
	ListCities(ctx context.Context) ([]domain.City, error)
	GetCity(ctx context.Context, cityID string) (domain.City, error)
	CreateSystemArea(ctx context.Context, area domain.SystemArea) (domain.SystemArea, error)
	ListSystemAreas(ctx context.Context, cityID string) ([]domain.SystemArea, error)
	GetSystemArea(ctx context.Context, areaID string) (domain.SystemArea, error)
	CreateSystemStreet(ctx context.Context, street domain.SystemStreet) (domain.SystemStreet, error)
	ListSystemStreets(ctx context.Context, cityID, areaID string) ([]domain.SystemStreet, error)
	GetSystemStreet(ctx context.Context, streetID string) (domain.SystemStreet, error)
	CreateSystemNeighborhood(ctx context.Context, neighborhood domain.SystemNeighborhood) (domain.SystemNeighborhood, error)
	ListSystemNeighborhoods(ctx context.Context, cityID, areaID, streetID string) ([]domain.SystemNeighborhood, error)
	GetSystemNeighborhood(ctx context.Context, neighborhoodID string) (domain.SystemNeighborhood, error)
	CreateLocationSuggestion(ctx context.Context, suggestion domain.LocationSuggestion) (domain.LocationSuggestion, error)
	ListLocationSuggestions(ctx context.Context, status domain.LocationSuggestionStatus) ([]domain.LocationSuggestion, error)
	GetLocationSuggestion(ctx context.Context, suggestionID string) (domain.LocationSuggestion, error)
	UpdateLocationSuggestion(ctx context.Context, suggestion domain.LocationSuggestion) (domain.LocationSuggestion, error)
}
