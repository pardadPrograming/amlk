package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"amlakcrm/backend/internal/domain"
	"amlakcrm/backend/internal/support"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoPlatformStore struct {
	Store
	db *mongo.Database
}

func NewMongoPlatformStore(primary Store, db *mongo.Database) *MongoPlatformStore {
	return &MongoPlatformStore{Store: primary, db: db}
}

func (s *MongoPlatformStore) EnsurePlatformIndexes(ctx context.Context) error {
	indexes := []struct {
		collection string
		models     []mongo.IndexModel
	}{
		{
			collection: "admin_accounts",
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "userId", Value: 1}}, Options: options.Index().SetUnique(true)},
			},
		},
		{
			collection: "cities",
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "normalizedName", Value: 1}}, Options: options.Index().SetUnique(true)},
			},
		},
		{
			collection: "system_areas",
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "cityId", Value: 1}, {Key: "normalizedName", Value: 1}}},
			},
		},
		{
			collection: "system_streets",
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "cityId", Value: 1}, {Key: "areaId", Value: 1}, {Key: "normalizedName", Value: 1}}},
			},
		},
		{
			collection: "system_neighborhoods",
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "cityId", Value: 1}, {Key: "areaId", Value: 1}, {Key: "streetId", Value: 1}, {Key: "normalizedName", Value: 1}}},
			},
		},
		{
			collection: "location_suggestions",
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "status", Value: 1}, {Key: "createdAt", Value: -1}}},
				{Keys: bson.D{{Key: "cityId", Value: 1}, {Key: "type", Value: 1}, {Key: "normalizedName", Value: 1}}},
			},
		},
		{
			collection: "platform_settings",
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "id", Value: 1}}, Options: options.Index().SetUnique(true)},
			},
		},
		{
			collection: "businesses",
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "id", Value: 1}}, Options: options.Index().SetUnique(true)},
			},
		},
		{
			collection: "business_members",
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "id", Value: 1}}, Options: options.Index().SetUnique(true)},
				{Keys: bson.D{{Key: "businessId", Value: 1}, {Key: "userId", Value: 1}}},
				{Keys: bson.D{{Key: "userId", Value: 1}, {Key: "status", Value: 1}}},
			},
		},
		{
			collection: "property_files",
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "id", Value: 1}}, Options: options.Index().SetUnique(true)},
				{Keys: bson.D{{Key: "businessId", Value: 1}, {Key: "ownerUserId", Value: 1}, {Key: "updatedAt", Value: -1}}},
				{Keys: bson.D{{Key: "businessId", Value: 1}, {Key: "ownerUserId", Value: 1}, {Key: "createdAt", Value: -1}}},
				{Keys: bson.D{{Key: "sharedFromFileId", Value: 1}}},
			},
		},
		{
			collection: "property_share_requests",
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "id", Value: 1}}, Options: options.Index().SetUnique(true)},
				{Keys: bson.D{{Key: "businessId", Value: 1}, {Key: "ownerUserId", Value: 1}, {Key: "createdAt", Value: -1}}},
				{Keys: bson.D{{Key: "businessId", Value: 1}, {Key: "requesterUserId", Value: 1}, {Key: "createdAt", Value: -1}}},
			},
		},
		{
			collection: "notifications",
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "id", Value: 1}}, Options: options.Index().SetUnique(true)},
				{Keys: bson.D{{Key: "userId", Value: 1}, {Key: "createdAt", Value: -1}}},
			},
		},
		{
			collection: "contacts",
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "id", Value: 1}}, Options: options.Index().SetUnique(true)},
				{Keys: bson.D{{Key: "businessId", Value: 1}, {Key: "updatedAt", Value: -1}}},
			},
		},
		{
			collection: "files",
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "id", Value: 1}}, Options: options.Index().SetUnique(true)},
				{Keys: bson.D{{Key: "uploaderId", Value: 1}, {Key: "purpose", Value: 1}, {Key: "targetId", Value: 1}}},
				{Keys: bson.D{{Key: "expiresAt", Value: 1}}},
			},
		},
		{
			collection: "channels",
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "id", Value: 1}}, Options: options.Index().SetUnique(true)},
				{Keys: bson.D{{Key: "type", Value: 1}, {Key: "ownerUserId", Value: 1}}},
				{Keys: bson.D{{Key: "type", Value: 1}, {Key: "businessId", Value: 1}}},
			},
		},
		{
			collection: "channel_vaults",
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "id", Value: 1}}, Options: options.Index().SetUnique(true)},
				{Keys: bson.D{{Key: "ownerUserId", Value: 1}, {Key: "isMain", Value: 1}}},
				{Keys: bson.D{{Key: "businessId", Value: 1}, {Key: "isMain", Value: 1}}},
				{Keys: bson.D{{Key: "channelId", Value: 1}}},
			},
		},
		{
			collection: "channel_members",
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "id", Value: 1}}, Options: options.Index().SetUnique(true)},
				{Keys: bson.D{{Key: "channelId", Value: 1}, {Key: "userId", Value: 1}}},
				{Keys: bson.D{{Key: "channelId", Value: 1}, {Key: "phone", Value: 1}}},
				{Keys: bson.D{{Key: "userId", Value: 1}, {Key: "status", Value: 1}}},
				{Keys: bson.D{{Key: "phone", Value: 1}, {Key: "status", Value: 1}}},
			},
		},
		{
			collection: "channel_messages",
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "id", Value: 1}}, Options: options.Index().SetUnique(true)},
				{Keys: bson.D{{Key: "channelId", Value: 1}, {Key: "createdAt", Value: -1}}},
			},
		},
		{
			collection: "channel_vault_files",
			models: []mongo.IndexModel{
				{Keys: bson.D{{Key: "id", Value: 1}}, Options: options.Index().SetUnique(true)},
				{Keys: bson.D{{Key: "channelId", Value: 1}, {Key: "createdAt", Value: -1}}},
			},
		},
	}
	for _, item := range indexes {
		if _, err := s.db.Collection(item.collection).Indexes().CreateMany(ctx, item.models); err != nil {
			return err
		}
	}
	return nil
}

func (s *MongoPlatformStore) CreateBusiness(ctx context.Context, business domain.Business) (domain.Business, error) {
	item, err := s.Store.CreateBusiness(ctx, business)
	if err != nil {
		return domain.Business{}, err
	}
	if err := s.saveBusiness(ctx, item); err != nil {
		return domain.Business{}, err
	}
	return item, nil
}

func (s *MongoPlatformStore) ListBusinessesForUser(ctx context.Context, userID string) ([]domain.Business, error) {
	members, err := findStoredData[domain.BusinessMember](ctx, s.db.Collection("business_members"), bson.M{
		"userId": userID,
		"status": bson.M{"$ne": domain.MemberRemoved},
	})
	if err != nil {
		return nil, err
	}
	if len(members) == 0 {
		return []domain.Business{}, nil
	}
	ids := make([]string, 0, len(members))
	for _, member := range members {
		if member.Status == domain.MemberActive {
			ids = append(ids, member.BusinessID)
		}
	}
	if len(ids) == 0 {
		return []domain.Business{}, nil
	}
	return findStoredData[domain.Business](ctx, s.db.Collection("businesses"), bson.M{"id": bson.M{"$in": ids}})
}

func (s *MongoPlatformStore) GetBusiness(ctx context.Context, id string) (domain.Business, error) {
	return findOneStoredData[domain.Business](ctx, s.db.Collection("businesses"), bson.M{"id": id})
}

func (s *MongoPlatformStore) UpdateBusiness(ctx context.Context, business domain.Business) (domain.Business, error) {
	item, err := s.Store.UpdateBusiness(ctx, business)
	if err != nil {
		existing, getErr := s.GetBusiness(ctx, business.ID)
		if getErr != nil {
			return domain.Business{}, getErr
		}
		item = business
		item.CreatedAt = existing.CreatedAt
		item.UpdatedAt = time.Now().UTC()
	}
	if err := s.saveBusiness(ctx, item); err != nil {
		return domain.Business{}, err
	}
	return item, nil
}

func (s *MongoPlatformStore) ListBusinesses(ctx context.Context) ([]domain.Business, error) {
	return findStoredData[domain.Business](ctx, s.db.Collection("businesses"), bson.M{})
}

func (s *MongoPlatformStore) CreateMember(ctx context.Context, member domain.BusinessMember) (domain.BusinessMember, error) {
	item, err := s.Store.CreateMember(ctx, member)
	if err != nil {
		return domain.BusinessMember{}, err
	}
	if err := s.saveMember(ctx, item); err != nil {
		return domain.BusinessMember{}, err
	}
	return item, nil
}

func (s *MongoPlatformStore) ListMembers(ctx context.Context, businessID string) ([]domain.BusinessMember, error) {
	return findStoredData[domain.BusinessMember](ctx, s.db.Collection("business_members"), bson.M{
		"businessId": businessID,
		"status":     bson.M{"$ne": domain.MemberRemoved},
	})
}

func (s *MongoPlatformStore) GetMemberByUser(ctx context.Context, businessID string, userID string) (domain.BusinessMember, error) {
	return findOneStoredData[domain.BusinessMember](ctx, s.db.Collection("business_members"), bson.M{
		"businessId": businessID,
		"userId":     userID,
		"status":     bson.M{"$ne": domain.MemberRemoved},
	})
}

func (s *MongoPlatformStore) UpdateMember(ctx context.Context, member domain.BusinessMember) (domain.BusinessMember, error) {
	item, err := s.Store.UpdateMember(ctx, member)
	if err != nil {
		existing, getErr := findOneStoredData[domain.BusinessMember](ctx, s.db.Collection("business_members"), bson.M{"id": member.ID})
		if getErr != nil {
			return domain.BusinessMember{}, getErr
		}
		item = member
		item.JoinedAt = existing.JoinedAt
		item.UpdatedAt = time.Now().UTC()
	}
	if err := s.saveMember(ctx, item); err != nil {
		return domain.BusinessMember{}, err
	}
	return item, nil
}

func (s *MongoPlatformStore) CreatePropertyFile(ctx context.Context, file domain.PropertyFile) (domain.PropertyFile, error) {
	item, err := s.Store.CreatePropertyFile(ctx, file)
	if err != nil {
		return domain.PropertyFile{}, err
	}
	if err := s.savePropertyFile(ctx, item); err != nil {
		return domain.PropertyFile{}, err
	}
	return item, nil
}

func (s *MongoPlatformStore) CreateNotification(ctx context.Context, notification domain.Notification) (domain.Notification, error) {
	now := time.Now().UTC()
	notification.ID = support.NewID()
	notification.CreatedAt = now
	if err := s.saveNotification(ctx, notification); err != nil {
		return domain.Notification{}, err
	}
	return notification, nil
}

func (s *MongoPlatformStore) CreateFile(ctx context.Context, file domain.FileObject) (domain.FileObject, error) {
	now := time.Now().UTC()
	if file.ID == "" {
		file.ID = support.NewID()
	}
	if file.CreatedAt.IsZero() {
		file.CreatedAt = now
	}
	file.UpdatedAt = now
	if err := s.saveFile(ctx, file); err != nil {
		return domain.FileObject{}, err
	}
	return file, nil
}

func (s *MongoPlatformStore) GetFile(ctx context.Context, fileID string) (domain.FileObject, error) {
	return findOneStoredData[domain.FileObject](ctx, s.db.Collection("files"), bson.M{"id": fileID})
}

func (s *MongoPlatformStore) UpdateFile(ctx context.Context, file domain.FileObject) (domain.FileObject, error) {
	existing, err := s.GetFile(ctx, file.ID)
	if err != nil {
		return domain.FileObject{}, err
	}
	if file.CreatedAt.IsZero() {
		file.CreatedAt = existing.CreatedAt
	}
	file.UpdatedAt = time.Now().UTC()
	if err := s.saveFile(ctx, file); err != nil {
		return domain.FileObject{}, err
	}
	return file, nil
}

func (s *MongoPlatformStore) ListExpiredFiles(ctx context.Context, before time.Time, limit int) ([]domain.FileObject, error) {
	if limit <= 0 {
		limit = 100
	}
	return findStoredDataWithOptions[domain.FileObject](ctx, s.db.Collection("files"), bson.M{
		"expiresAt": bson.M{"$ne": time.Time{}, "$lte": before},
	}, options.Find().SetSort(bson.D{{Key: "expiresAt", Value: 1}}).SetLimit(int64(limit)))
}

func (s *MongoPlatformStore) DeleteFile(ctx context.Context, fileID string) error {
	res, err := s.db.Collection("files").DeleteOne(ctx, bson.M{"id": fileID})
	if err != nil {
		return mongoErr(err)
	}
	if res.DeletedCount == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *MongoPlatformStore) ListNotifications(ctx context.Context, userID string, unreadOnly bool) ([]domain.Notification, error) {
	filter := bson.M{"userId": userID}
	if unreadOnly {
		filter["readAt"] = time.Time{}
	}
	return findStoredDataWithOptions[domain.Notification](ctx, s.db.Collection("notifications"), filter, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}))
}

func (s *MongoPlatformStore) MarkNotificationRead(ctx context.Context, userID string, notificationID string) error {
	existing, err := findOneStoredData[domain.Notification](ctx, s.db.Collection("notifications"), bson.M{"id": notificationID, "userId": userID})
	if err != nil {
		return err
	}
	existing.ReadAt = time.Now().UTC()
	return s.saveNotification(ctx, existing)
}

func (s *MongoPlatformStore) ListPropertyFiles(ctx context.Context, businessID string) ([]domain.PropertyFile, error) {
	return findStoredData[domain.PropertyFile](ctx, s.db.Collection("property_files"), bson.M{"businessId": businessID})
}

func (s *MongoPlatformStore) ListPropertyFilesForOwner(ctx context.Context, businessID string, ownerUserID string) ([]domain.PropertyFile, error) {
	return findStoredData[domain.PropertyFile](ctx, s.db.Collection("property_files"), bson.M{
		"businessId":  businessID,
		"ownerUserId": ownerUserID,
	})
}

func (s *MongoPlatformStore) ListLatestPropertyFiles(ctx context.Context, businessID string, ownerUserID string, vaultIDs []string, propertyType domain.PropertyFileType, limit int, offset int) ([]domain.PropertyFile, int, error) {
	if limit <= 0 {
		limit = 30
	}
	if offset < 0 {
		offset = 0
	}
	visibility := bson.A{bson.M{"ownerUserId": ownerUserID}}
	if len(vaultIDs) > 0 {
		visibility = append(visibility, bson.M{"vaultIds": bson.M{"$in": vaultIDs}})
	}
	filter := bson.M{
		"businessId": businessID,
		"$or":        visibility,
	}
	if propertyType != "" {
		filter["$and"] = bson.A{bson.M{"$or": bson.A{
			bson.M{"type": propertyType},
			bson.M{"types": propertyType},
		}}}
	}
	total, err := s.db.Collection("property_files").CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, mongoErr(err)
	}
	items, err := findStoredDataWithOptions[domain.PropertyFile](ctx, s.db.Collection("property_files"), filter, options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset)))
	if err != nil {
		return nil, 0, err
	}
	return items, int(total), nil
}

func (s *MongoPlatformStore) GetPropertyFile(ctx context.Context, businessID string, fileID string) (domain.PropertyFile, error) {
	return findOneStoredData[domain.PropertyFile](ctx, s.db.Collection("property_files"), bson.M{"id": fileID, "businessId": businessID})
}

func (s *MongoPlatformStore) UpdatePropertyFile(ctx context.Context, file domain.PropertyFile) (domain.PropertyFile, error) {
	item, err := s.Store.UpdatePropertyFile(ctx, file)
	if err != nil {
		existing, getErr := s.GetPropertyFile(ctx, file.BusinessID, file.ID)
		if getErr != nil {
			return domain.PropertyFile{}, getErr
		}
		item = file
		item.CreatedAt = existing.CreatedAt
		item.UpdatedAt = time.Now().UTC()
	}
	if err := s.savePropertyFile(ctx, item); err != nil {
		return domain.PropertyFile{}, err
	}
	return item, nil
}

func (s *MongoPlatformStore) CreatePropertyShareRequest(ctx context.Context, request domain.PropertyShareRequest) (domain.PropertyShareRequest, error) {
	now := time.Now().UTC()
	request.ID = support.NewID()
	request.CreatedAt = now
	request.UpdatedAt = now
	if err := s.savePropertyShareRequest(ctx, request); err != nil {
		return domain.PropertyShareRequest{}, err
	}
	return request, nil
}

func (s *MongoPlatformStore) GetPropertyShareRequest(ctx context.Context, businessID string, requestID string) (domain.PropertyShareRequest, error) {
	return findOneStoredData[domain.PropertyShareRequest](ctx, s.db.Collection("property_share_requests"), bson.M{"id": requestID, "businessId": businessID})
}

func (s *MongoPlatformStore) UpdatePropertyShareRequest(ctx context.Context, request domain.PropertyShareRequest) (domain.PropertyShareRequest, error) {
	if request.UpdatedAt.IsZero() {
		request.UpdatedAt = time.Now().UTC()
	}
	if err := s.savePropertyShareRequest(ctx, request); err != nil {
		return domain.PropertyShareRequest{}, err
	}
	return request, nil
}

func (s *MongoPlatformStore) ListPropertyShareRequestsForOwner(ctx context.Context, businessID string, ownerUserID string) ([]domain.PropertyShareRequest, error) {
	return findStoredDataWithOptions[domain.PropertyShareRequest](ctx, s.db.Collection("property_share_requests"), bson.M{"businessId": businessID, "ownerUserId": ownerUserID}, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}))
}

func (s *MongoPlatformStore) ListPropertyShareRequestsForRequester(ctx context.Context, businessID string, requesterUserID string) ([]domain.PropertyShareRequest, error) {
	return findStoredDataWithOptions[domain.PropertyShareRequest](ctx, s.db.Collection("property_share_requests"), bson.M{"businessId": businessID, "requesterUserId": requesterUserID}, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}))
}

func (s *MongoPlatformStore) CreateContact(ctx context.Context, contact domain.Contact) (domain.Contact, error) {
	item, err := s.Store.CreateContact(ctx, contact)
	if err != nil {
		return domain.Contact{}, err
	}
	if err := s.saveContact(ctx, item); err != nil {
		return domain.Contact{}, err
	}
	return item, nil
}

func (s *MongoPlatformStore) ListContacts(ctx context.Context, businessID string) ([]domain.Contact, error) {
	return findStoredData[domain.Contact](ctx, s.db.Collection("contacts"), bson.M{"businessId": businessID})
}

func (s *MongoPlatformStore) GetContact(ctx context.Context, businessID, contactID string) (domain.Contact, error) {
	return findOneStoredData[domain.Contact](ctx, s.db.Collection("contacts"), bson.M{"id": contactID, "businessId": businessID})
}

func (s *MongoPlatformStore) UpdateContact(ctx context.Context, contact domain.Contact) (domain.Contact, error) {
	item, err := s.Store.UpdateContact(ctx, contact)
	if err != nil {
		existing, getErr := s.GetContact(ctx, contact.BusinessID, contact.ID)
		if getErr != nil {
			return domain.Contact{}, getErr
		}
		item = contact
		item.CreatedAt = existing.CreatedAt
		item.UpdatedAt = time.Now().UTC()
	}
	if err := s.saveContact(ctx, item); err != nil {
		return domain.Contact{}, err
	}
	return item, nil
}

func (s *MongoPlatformStore) EnsureUserMainChannel(ctx context.Context, user domain.User) (domain.Channel, error) {
	existing, err := findOneStoredData[domain.Channel](ctx, s.db.Collection("channels"), bson.M{
		"type":        domain.ChannelTypeUserMain,
		"ownerUserId": user.ID,
	})
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, ErrNotFound) {
		return domain.Channel{}, err
	}
	channel, err := s.Store.EnsureUserMainChannel(ctx, user)
	if err != nil {
		return domain.Channel{}, err
	}
	if err := s.saveChannel(ctx, channel); err != nil {
		return domain.Channel{}, err
	}
	_, _ = s.UpsertChannelMember(ctx, domain.ChannelMember{
		ChannelID: channel.ID,
		UserID:    user.ID,
		Phone:     user.Phone,
		Status:    domain.ChannelMemberActive,
	})
	return channel, nil
}

func (s *MongoPlatformStore) EnsureBusinessMainChannel(ctx context.Context, business domain.Business) (domain.Channel, error) {
	existing, err := findOneStoredData[domain.Channel](ctx, s.db.Collection("channels"), bson.M{
		"type":       domain.ChannelTypeBusinessMain,
		"businessId": business.ID,
	})
	if err == nil {
		_ = s.ensureStoredBusinessChannelMembers(ctx, existing.ID, business.ID)
		return existing, nil
	}
	if !errors.Is(err, ErrNotFound) {
		return domain.Channel{}, err
	}
	channel, err := s.Store.EnsureBusinessMainChannel(ctx, business)
	if err != nil {
		return domain.Channel{}, err
	}
	if err := s.saveChannel(ctx, channel); err != nil {
		return domain.Channel{}, err
	}
	if err := s.ensureStoredBusinessChannelMembers(ctx, channel.ID, business.ID); err != nil {
		return domain.Channel{}, err
	}
	return channel, nil
}

func (s *MongoPlatformStore) EnsurePrivateChannel(ctx context.Context, actor domain.User, target domain.User) (domain.Channel, error) {
	if actor.ID == "" || target.ID == "" || actor.ID == target.ID {
		return domain.Channel{}, ErrNotFound
	}
	actorMemberships, err := findStoredData[domain.ChannelMember](ctx, s.db.Collection("channel_members"), bson.M{
		"userId": actor.ID,
		"status": domain.ChannelMemberActive,
	})
	if err != nil {
		return domain.Channel{}, err
	}
	ids := make([]string, 0, len(actorMemberships))
	for _, member := range actorMemberships {
		if member.ChannelID != "" {
			ids = append(ids, member.ChannelID)
		}
	}
	if len(ids) > 0 {
		channels, err := findStoredData[domain.Channel](ctx, s.db.Collection("channels"), bson.M{
			"id":   bson.M{"$in": ids},
			"type": domain.ChannelTypePrivate,
		})
		if err != nil {
			return domain.Channel{}, err
		}
		for _, channel := range channels {
			if _, err := findOneStoredData[domain.ChannelMember](ctx, s.db.Collection("channel_members"), bson.M{
				"channelId": channel.ID,
				"userId":    target.ID,
				"status":    domain.ChannelMemberActive,
			}); err == nil {
				return channel, nil
			} else if !errors.Is(err, ErrNotFound) {
				return domain.Channel{}, err
			}
		}
	}
	channel, err := s.createStoredChannel(ctx, domain.Channel{Type: domain.ChannelTypePrivate, Title: "private chat"})
	if err != nil {
		return domain.Channel{}, err
	}
	if _, err := s.UpsertChannelMember(ctx, domain.ChannelMember{ChannelID: channel.ID, UserID: actor.ID, Phone: actor.Phone, Status: domain.ChannelMemberActive}); err != nil {
		return domain.Channel{}, err
	}
	if _, err := s.UpsertChannelMember(ctx, domain.ChannelMember{ChannelID: channel.ID, UserID: target.ID, Phone: target.Phone, Status: domain.ChannelMemberActive}); err != nil {
		return domain.Channel{}, err
	}
	return channel, nil
}

func (s *MongoPlatformStore) EnsureUserMainVault(ctx context.Context, user domain.User) (domain.ChannelVault, error) {
	existing, err := findOneStoredData[domain.ChannelVault](ctx, s.db.Collection("channel_vaults"), bson.M{"ownerUserId": user.ID, "isMain": true})
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, ErrNotFound) {
		return domain.ChannelVault{}, err
	}
	channel, err := s.EnsureUserMainChannel(ctx, user)
	if err != nil {
		return domain.ChannelVault{}, err
	}
	return s.createStoredVault(ctx, user.ID, "", channel.ID, channel.Title, true)
}

func (s *MongoPlatformStore) EnsureBusinessMainVault(ctx context.Context, business domain.Business) (domain.ChannelVault, error) {
	existing, err := findOneStoredData[domain.ChannelVault](ctx, s.db.Collection("channel_vaults"), bson.M{"businessId": business.ID, "isMain": true})
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, ErrNotFound) {
		return domain.ChannelVault{}, err
	}
	channel, err := s.EnsureBusinessMainChannel(ctx, business)
	if err != nil {
		return domain.ChannelVault{}, err
	}
	return s.createStoredVault(ctx, "", business.ID, channel.ID, channel.Title, true)
}

func (s *MongoPlatformStore) CreateUserVault(ctx context.Context, user domain.User, title string) (domain.ChannelVault, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		title = "صندوقچه"
	}
	channel, err := s.createStoredChannel(ctx, domain.Channel{Type: domain.ChannelTypeUserVault, OwnerUserID: user.ID, Title: title})
	if err != nil {
		return domain.ChannelVault{}, err
	}
	vault, err := s.createStoredVault(ctx, user.ID, "", channel.ID, title, false)
	if err != nil {
		return domain.ChannelVault{}, err
	}
	channel.VaultID = vault.ID
	_ = s.saveChannel(ctx, channel)
	_, _ = s.UpsertChannelMember(ctx, domain.ChannelMember{ChannelID: channel.ID, UserID: user.ID, Phone: user.Phone, Status: domain.ChannelMemberActive})
	return vault, nil
}

func (s *MongoPlatformStore) CreateBusinessVault(ctx context.Context, business domain.Business, title string) (domain.ChannelVault, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		title = "صندوقچه"
	}
	channel, err := s.createStoredChannel(ctx, domain.Channel{Type: domain.ChannelTypeBusinessVault, BusinessID: business.ID, Title: title})
	if err != nil {
		return domain.ChannelVault{}, err
	}
	vault, err := s.createStoredVault(ctx, "", business.ID, channel.ID, title, false)
	if err != nil {
		return domain.ChannelVault{}, err
	}
	channel.VaultID = vault.ID
	_ = s.saveChannel(ctx, channel)
	_ = s.ensureStoredBusinessChannelMembers(ctx, channel.ID, business.ID)
	return vault, nil
}

func (s *MongoPlatformStore) GetChannelVault(ctx context.Context, vaultID string) (domain.ChannelVault, error) {
	return findOneStoredData[domain.ChannelVault](ctx, s.db.Collection("channel_vaults"), bson.M{"id": vaultID})
}

func (s *MongoPlatformStore) ListUserVaults(ctx context.Context, userID string) ([]domain.ChannelVault, error) {
	return findStoredData[domain.ChannelVault](ctx, s.db.Collection("channel_vaults"), bson.M{"ownerUserId": userID})
}

func (s *MongoPlatformStore) ListBusinessVaults(ctx context.Context, businessID string) ([]domain.ChannelVault, error) {
	return findStoredData[domain.ChannelVault](ctx, s.db.Collection("channel_vaults"), bson.M{"businessId": businessID})
}

func (s *MongoPlatformStore) createStoredChannel(ctx context.Context, channel domain.Channel) (domain.Channel, error) {
	now := time.Now().UTC()
	channel.ID = support.NewID()
	channel.CreatedAt = now
	channel.UpdatedAt = now
	if err := s.saveChannel(ctx, channel); err != nil {
		return domain.Channel{}, err
	}
	return channel, nil
}

func (s *MongoPlatformStore) createStoredVault(ctx context.Context, ownerUserID, businessID, channelID, title string, isMain bool) (domain.ChannelVault, error) {
	now := time.Now().UTC()
	vault := domain.ChannelVault{ID: support.NewID(), ChannelID: channelID, OwnerUserID: ownerUserID, BusinessID: businessID, Title: title, IsMain: isMain, CreatedAt: now, UpdatedAt: now}
	if err := s.saveChannelVault(ctx, vault); err != nil {
		return domain.ChannelVault{}, err
	}
	channel, err := s.GetChannel(ctx, channelID)
	if err == nil {
		channel.VaultID = vault.ID
		_ = s.saveChannel(ctx, channel)
	}
	return vault, nil
}

func (s *MongoPlatformStore) ensureStoredBusinessChannelMembers(ctx context.Context, channelID, businessID string) error {
	members, err := s.ListMembers(ctx, businessID)
	if err != nil {
		return err
	}
	for _, member := range members {
		if member.Status != domain.MemberActive {
			continue
		}
		_, err := s.UpsertChannelMember(ctx, domain.ChannelMember{
			ChannelID: channelID,
			UserID:    member.UserID,
			Phone:     member.UserPhone,
			Status:    domain.ChannelMemberActive,
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *MongoPlatformStore) GetChannel(ctx context.Context, channelID string) (domain.Channel, error) {
	return findOneStoredData[domain.Channel](ctx, s.db.Collection("channels"), bson.M{"id": channelID})
}

func (s *MongoPlatformStore) ListChannelsForUser(ctx context.Context, userID string, phones []string, businessIDs []string) ([]domain.Channel, error) {
	memberFilter := bson.M{
		"status": domain.ChannelMemberActive,
		"$or": []bson.M{
			{"userId": userID},
			{"phone": bson.M{"$in": phones}},
		},
	}
	members, err := findStoredData[domain.ChannelMember](ctx, s.db.Collection("channel_members"), memberFilter)
	if err != nil {
		return nil, err
	}
	ids := []string{}
	seen := map[string]struct{}{}
	for _, member := range members {
		if _, ok := seen[member.ChannelID]; ok {
			continue
		}
		seen[member.ChannelID] = struct{}{}
		ids = append(ids, member.ChannelID)
	}
	filter := bson.M{"id": bson.M{"$in": ids}}
	if len(businessIDs) > 0 {
		filter = bson.M{"$or": []bson.M{
			{"id": bson.M{"$in": ids}},
			{"type": domain.ChannelTypeBusinessMain, "businessId": bson.M{"$in": businessIDs}},
		}}
	}
	return findStoredData[domain.Channel](ctx, s.db.Collection("channels"), filter)
}

func (s *MongoPlatformStore) UpsertChannelMember(ctx context.Context, member domain.ChannelMember) (domain.ChannelMember, error) {
	now := time.Now().UTC()
	filter := bson.M{"channelId": member.ChannelID}
	if member.UserID != "" {
		filter["userId"] = member.UserID
	} else {
		filter["phone"] = member.Phone
	}
	existing, err := findOneStoredData[domain.ChannelMember](ctx, s.db.Collection("channel_members"), filter)
	if err == nil {
		member.ID = existing.ID
		member.CreatedAt = existing.CreatedAt
	} else if errors.Is(err, ErrNotFound) {
		member.ID = support.NewID()
		member.CreatedAt = now
	} else {
		return domain.ChannelMember{}, err
	}
	member.UpdatedAt = now
	if member.Status == "" {
		member.Status = domain.ChannelMemberActive
	}
	if member.Role == "" {
		member.Role = domain.ChannelMemberRoleMember
	}
	if err := s.saveChannelMember(ctx, member); err != nil {
		return domain.ChannelMember{}, err
	}
	return member, nil
}

func (s *MongoPlatformStore) GetChannelMember(ctx context.Context, channelID string, userID string) (domain.ChannelMember, error) {
	return findOneStoredData[domain.ChannelMember](ctx, s.db.Collection("channel_members"), bson.M{
		"channelId": channelID,
		"userId":    userID,
		"status":    domain.ChannelMemberActive,
	})
}

func (s *MongoPlatformStore) ListChannelMembers(ctx context.Context, channelID string) ([]domain.ChannelMember, error) {
	return findStoredData[domain.ChannelMember](ctx, s.db.Collection("channel_members"), bson.M{"channelId": channelID})
}

func (s *MongoPlatformStore) CreateChannelMessage(ctx context.Context, message domain.ChannelMessage) (domain.ChannelMessage, error) {
	now := time.Now().UTC()
	message.ID = support.NewID()
	message.CreatedAt = now
	message.UpdatedAt = now
	if err := s.saveChannelMessage(ctx, message); err != nil {
		return domain.ChannelMessage{}, err
	}
	return message, nil
}

func (s *MongoPlatformStore) GetChannelMessage(ctx context.Context, channelID string, messageID string) (domain.ChannelMessage, error) {
	return findOneStoredData[domain.ChannelMessage](ctx, s.db.Collection("channel_messages"), bson.M{
		"id":        messageID,
		"channelId": channelID,
	})
}

func (s *MongoPlatformStore) UpdateChannelMessage(ctx context.Context, message domain.ChannelMessage) (domain.ChannelMessage, error) {
	existing, err := s.GetChannelMessage(ctx, message.ChannelID, message.ID)
	if err != nil {
		return domain.ChannelMessage{}, err
	}
	message.CreatedAt = existing.CreatedAt
	message.UpdatedAt = time.Now().UTC()
	if err := s.saveChannelMessage(ctx, message); err != nil {
		return domain.ChannelMessage{}, err
	}
	return message, nil
}

func (s *MongoPlatformStore) DeleteChannelMessage(ctx context.Context, channelID string, messageID string) error {
	result, err := s.db.Collection("channel_messages").DeleteOne(ctx, bson.M{
		"id":        messageID,
		"channelId": channelID,
	})
	if err != nil {
		return mongoErr(err)
	}
	if result.DeletedCount == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *MongoPlatformStore) ListChannelMessages(ctx context.Context, channelID string, limit int, offset int) ([]domain.ChannelMessage, int, error) {
	filter := bson.M{"channelId": channelID}
	total, err := s.db.Collection("channel_messages").CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, mongoErr(err)
	}
	items, err := findStoredDataWithOptions[domain.ChannelMessage](ctx, s.db.Collection("channel_messages"), filter, options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset)))
	if err != nil {
		return nil, 0, err
	}
	return items, int(total), nil
}

func (s *MongoPlatformStore) MarkChannelMessagesSeen(ctx context.Context, channelID string, userID string, messageIDs []string) error {
	if userID == "" || len(messageIDs) == 0 {
		return nil
	}
	ids := make([]string, 0, len(messageIDs))
	for _, id := range messageIDs {
		if id != "" {
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		return nil
	}
	now := time.Now().UTC()
	_, err := s.db.Collection("channel_messages").UpdateMany(ctx, bson.M{
		"id":            bson.M{"$in": ids},
		"channelId":     channelID,
		"authorId":      bson.M{"$ne": userID},
		"seenBy.userId": bson.M{"$ne": userID},
	}, bson.M{
		"$push": bson.M{
			"seenBy":      domain.ChannelMessageSeen{UserID: userID, SeenAt: now},
			"data.seenBy": domain.ChannelMessageSeen{UserID: userID, SeenAt: now},
		},
		"$set": bson.M{
			"updatedAt":      now,
			"data.updatedAt": now,
		},
	})
	return err
}

func (s *MongoPlatformStore) CreateChannelVaultFile(ctx context.Context, file domain.ChannelVaultFile) (domain.ChannelVaultFile, error) {
	now := time.Now().UTC()
	file.ID = support.NewID()
	file.CreatedAt = now
	file.UpdatedAt = now
	if err := s.saveChannelVaultFile(ctx, file); err != nil {
		return domain.ChannelVaultFile{}, err
	}
	return file, nil
}

func (s *MongoPlatformStore) UpsertChannelVaultFile(ctx context.Context, file domain.ChannelVaultFile) (domain.ChannelVaultFile, error) {
	now := time.Now().UTC()
	filter := bson.M{"channelId": file.ChannelID, "vaultId": file.VaultID}
	if file.PropertyFileID != "" {
		filter["propertyFileId"] = file.PropertyFileID
	} else if file.FileID != "" {
		filter["fileId"] = file.FileID
	} else {
		return s.CreateChannelVaultFile(ctx, file)
	}
	existing, err := findOneStoredData[domain.ChannelVaultFile](ctx, s.db.Collection("channel_vault_files"), filter)
	if err == nil {
		file.ID = existing.ID
		file.CreatedAt = existing.CreatedAt
	} else if errors.Is(err, ErrNotFound) {
		file.ID = support.NewID()
		file.CreatedAt = now
	} else {
		return domain.ChannelVaultFile{}, err
	}
	file.UpdatedAt = now
	if err := s.saveChannelVaultFile(ctx, file); err != nil {
		return domain.ChannelVaultFile{}, err
	}
	return file, nil
}

func (s *MongoPlatformStore) DeletePropertyVaultReferences(ctx context.Context, propertyFileID string, keepVaultIDs []string) error {
	if strings.TrimSpace(propertyFileID) == "" {
		return nil
	}
	filter := bson.M{"propertyFileId": propertyFileID}
	keep := []string{}
	seen := map[string]struct{}{}
	for _, vaultID := range keepVaultIDs {
		vaultID = strings.TrimSpace(vaultID)
		if vaultID == "" {
			continue
		}
		if _, ok := seen[vaultID]; ok {
			continue
		}
		seen[vaultID] = struct{}{}
		keep = append(keep, vaultID)
	}
	if len(keep) > 0 {
		filter["vaultId"] = bson.M{"$nin": keep}
	}
	_, err := s.db.Collection("channel_vault_files").DeleteMany(ctx, filter)
	return mongoErr(err)
}

func (s *MongoPlatformStore) GetChannelVaultFile(ctx context.Context, channelID string, fileID string) (domain.ChannelVaultFile, error) {
	return findOneStoredData[domain.ChannelVaultFile](ctx, s.db.Collection("channel_vault_files"), bson.M{
		"id":        fileID,
		"channelId": channelID,
	})
}

func (s *MongoPlatformStore) ListChannelVaultFiles(ctx context.Context, channelID string, limit int, offset int) ([]domain.ChannelVaultFile, int, error) {
	filter := bson.M{"channelId": channelID}
	total, err := s.db.Collection("channel_vault_files").CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, mongoErr(err)
	}
	items, err := findStoredDataWithOptions[domain.ChannelVaultFile](ctx, s.db.Collection("channel_vault_files"), filter, options.Find().
		SetSort(bson.D{{Key: "createdAt", Value: -1}}).
		SetLimit(int64(limit)).
		SetSkip(int64(offset)))
	if err != nil {
		return nil, 0, err
	}
	return items, int(total), nil
}

func (s *MongoPlatformStore) saveBusiness(ctx context.Context, item domain.Business) error {
	return s.upsertStored(ctx, "businesses", bson.M{"id": item.ID}, bson.M{"id": item.ID, "data": item})
}

func (s *MongoPlatformStore) saveChannel(ctx context.Context, item domain.Channel) error {
	return s.upsertStored(ctx, "channels", bson.M{"id": item.ID}, bson.M{
		"id":          item.ID,
		"type":        item.Type,
		"ownerUserId": item.OwnerUserID,
		"businessId":  item.BusinessID,
		"vaultId":     item.VaultID,
		"updatedAt":   item.UpdatedAt,
		"data":        item,
	})
}

func (s *MongoPlatformStore) saveChannelVault(ctx context.Context, item domain.ChannelVault) error {
	return s.upsertStored(ctx, "channel_vaults", bson.M{"id": item.ID}, bson.M{
		"id":          item.ID,
		"channelId":   item.ChannelID,
		"ownerUserId": item.OwnerUserID,
		"businessId":  item.BusinessID,
		"isMain":      item.IsMain,
		"updatedAt":   item.UpdatedAt,
		"data":        item,
	})
}

func (s *MongoPlatformStore) saveChannelMember(ctx context.Context, item domain.ChannelMember) error {
	return s.upsertStored(ctx, "channel_members", bson.M{"id": item.ID}, bson.M{
		"id":          item.ID,
		"channelId":   item.ChannelID,
		"userId":      item.UserID,
		"phone":       item.Phone,
		"displayName": item.DisplayName,
		"role":        item.Role,
		"status":      item.Status,
		"isOnline":    item.IsOnline,
		"lastSeenAt":  item.LastSeenAt,
		"updatedAt":   item.UpdatedAt,
		"data":        item,
	})
}

func (s *MongoPlatformStore) saveChannelMessage(ctx context.Context, item domain.ChannelMessage) error {
	return s.upsertStored(ctx, "channel_messages", bson.M{"id": item.ID}, bson.M{
		"id":        item.ID,
		"channelId": item.ChannelID,
		"authorId":  item.AuthorID,
		"createdAt": item.CreatedAt,
		"updatedAt": item.UpdatedAt,
		"seenBy":    item.SeenBy,
		"data":      item,
	})
}

func (s *MongoPlatformStore) saveChannelVaultFile(ctx context.Context, item domain.ChannelVaultFile) error {
	return s.upsertStored(ctx, "channel_vault_files", bson.M{"id": item.ID}, bson.M{
		"id":                item.ID,
		"vaultId":           item.VaultID,
		"channelId":         item.ChannelID,
		"uploaderId":        item.UploaderID,
		"sourceType":        item.SourceType,
		"fileId":            item.FileID,
		"propertyFileId":    item.PropertyFileID,
		"propertyStatus":    item.PropertyStatus,
		"commissionPercent": item.CommissionPercent,
		"kind":              item.Kind,
		"createdAt":         item.CreatedAt,
		"data":              item,
	})
}

func (s *MongoPlatformStore) saveMember(ctx context.Context, item domain.BusinessMember) error {
	return s.upsertStored(ctx, "business_members", bson.M{"id": item.ID}, bson.M{
		"id":         item.ID,
		"businessId": item.BusinessID,
		"userId":     item.UserID,
		"status":     item.Status,
		"data":       item,
	})
}

func (s *MongoPlatformStore) savePropertyFile(ctx context.Context, item domain.PropertyFile) error {
	return s.upsertStored(ctx, "property_files", bson.M{"id": item.ID}, bson.M{
		"id":               item.ID,
		"businessId":       item.BusinessID,
		"ownerUserId":      item.OwnerUserID,
		"sharedFromFileId": item.SharedFromFileID,
		"updatedAt":        item.UpdatedAt,
		"data":             item,
	})
}

func (s *MongoPlatformStore) savePropertyShareRequest(ctx context.Context, item domain.PropertyShareRequest) error {
	return s.upsertStored(ctx, "property_share_requests", bson.M{"id": item.ID}, bson.M{
		"id":              item.ID,
		"businessId":      item.BusinessID,
		"propertyFileId":  item.PropertyFileID,
		"ownerUserId":     item.OwnerUserID,
		"requesterUserId": item.RequesterUserID,
		"status":          item.Status,
		"createdAt":       item.CreatedAt,
		"updatedAt":       item.UpdatedAt,
		"data":            item,
	})
}

func (s *MongoPlatformStore) saveFile(ctx context.Context, item domain.FileObject) error {
	return s.upsertStored(ctx, "files", bson.M{"id": item.ID}, bson.M{
		"id":          item.ID,
		"ownerId":     item.OwnerID,
		"uploaderId":  item.UploaderID,
		"purpose":     item.Purpose,
		"targetType":  item.TargetType,
		"targetId":    item.TargetID,
		"status":      item.Status,
		"key":         item.Key,
		"expiresAt":   item.ExpiresAt,
		"updatedAt":   item.UpdatedAt,
		"contentType": item.ContentType,
		"kind":        item.Kind,
		"data":        item,
	})
}

func (s *MongoPlatformStore) saveNotification(ctx context.Context, item domain.Notification) error {
	return s.upsertStored(ctx, "notifications", bson.M{"id": item.ID}, bson.M{
		"id":        item.ID,
		"userId":    item.UserID,
		"type":      item.Type,
		"readAt":    item.ReadAt,
		"createdAt": item.CreatedAt,
		"data":      item,
	})
}

func (s *MongoPlatformStore) saveContact(ctx context.Context, item domain.Contact) error {
	return s.upsertStored(ctx, "contacts", bson.M{"id": item.ID}, bson.M{
		"id":         item.ID,
		"businessId": item.BusinessID,
		"updatedAt":  item.UpdatedAt,
		"data":       item,
	})
}

func (s *MongoPlatformStore) upsertStored(ctx context.Context, collection string, filter bson.M, doc bson.M) error {
	_, err := s.db.Collection(collection).ReplaceOne(ctx, filter, doc, options.Replace().SetUpsert(true))
	return mongoErr(err)
}

func (s *MongoPlatformStore) SaveAdminAccount(ctx context.Context, account domain.AdminAccount) (domain.AdminAccount, error) {
	now := time.Now().UTC()
	existing, err := s.GetAdminAccountByUser(ctx, account.UserID)
	if err == nil {
		if account.ID == "" {
			account.ID = existing.ID
		}
		if account.CreatedAt.IsZero() {
			account.CreatedAt = existing.CreatedAt
		}
	} else if err != ErrNotFound {
		return domain.AdminAccount{}, err
	}
	if account.ID == "" {
		account.ID = support.NewID()
		account.CreatedAt = now
	}
	if account.Status == "" {
		account.Status = "active"
	}
	account.UpdatedAt = now
	_, err = s.db.Collection("admin_accounts").UpdateOne(
		ctx,
		bson.M{"userId": account.UserID},
		bson.M{"$set": account},
		options.Update().SetUpsert(true),
	)
	return account, mongoErr(err)
}

func (s *MongoPlatformStore) GetAdminAccountByUser(ctx context.Context, userID string) (domain.AdminAccount, error) {
	var account domain.AdminAccount
	err := s.db.Collection("admin_accounts").FindOne(ctx, bson.M{"userId": userID}).Decode(&account)
	return account, mongoErr(err)
}

func (s *MongoPlatformStore) ListAdminAccounts(ctx context.Context) ([]domain.AdminAccount, error) {
	var result []domain.AdminAccount
	err := findAll(ctx, s.db.Collection("admin_accounts"), bson.M{}, &result)
	return result, err
}

func (s *MongoPlatformStore) GetPlatformSettings(ctx context.Context) (domain.PlatformSettings, error) {
	var settings domain.PlatformSettings
	err := s.db.Collection("platform_settings").FindOne(ctx, bson.M{"id": "platform"}).Decode(&settings)
	if errors.Is(mongoErr(err), ErrNotFound) {
		return domain.PlatformSettings{ID: "platform"}, nil
	}
	return settings, mongoErr(err)
}

func (s *MongoPlatformStore) SavePlatformSettings(ctx context.Context, settings domain.PlatformSettings) (domain.PlatformSettings, error) {
	if settings.ID == "" {
		settings.ID = "platform"
	}
	settings.UpdatedAt = time.Now().UTC()
	_, err := s.db.Collection("platform_settings").UpdateOne(
		ctx,
		bson.M{"id": settings.ID},
		bson.M{"$set": settings},
		options.Update().SetUpsert(true),
	)
	return settings, mongoErr(err)
}

func (s *MongoPlatformStore) CreateCity(ctx context.Context, city domain.City) (domain.City, error) {
	now := time.Now().UTC()
	city.ID = support.NewID()
	city.CreatedAt = now
	city.UpdatedAt = now
	if city.Status == "" {
		city.Status = domain.SystemLocationActive
	}
	_, err := s.db.Collection("cities").InsertOne(ctx, city)
	return city, mongoErr(err)
}

func (s *MongoPlatformStore) ListCities(ctx context.Context) ([]domain.City, error) {
	var result []domain.City
	err := findAll(ctx, s.db.Collection("cities"), bson.M{"status": domain.SystemLocationActive}, &result)
	return result, err
}

func (s *MongoPlatformStore) GetCity(ctx context.Context, cityID string) (domain.City, error) {
	var city domain.City
	err := s.db.Collection("cities").FindOne(ctx, bson.M{"id": cityID}).Decode(&city)
	return city, mongoErr(err)
}

func (s *MongoPlatformStore) CreateSystemArea(ctx context.Context, area domain.SystemArea) (domain.SystemArea, error) {
	if _, err := s.GetCity(ctx, area.CityID); err != nil {
		return domain.SystemArea{}, err
	}
	now := time.Now().UTC()
	area.ID = support.NewID()
	area.CreatedAt = now
	area.UpdatedAt = now
	if area.Status == "" {
		area.Status = domain.SystemLocationActive
	}
	_, err := s.db.Collection("system_areas").InsertOne(ctx, area)
	return area, mongoErr(err)
}

func (s *MongoPlatformStore) ListSystemAreas(ctx context.Context, cityID string) ([]domain.SystemArea, error) {
	var result []domain.SystemArea
	err := findAll(ctx, s.db.Collection("system_areas"), bson.M{"cityId": cityID, "status": domain.SystemLocationActive}, &result)
	return result, err
}

func (s *MongoPlatformStore) GetSystemArea(ctx context.Context, areaID string) (domain.SystemArea, error) {
	var area domain.SystemArea
	err := s.db.Collection("system_areas").FindOne(ctx, bson.M{"id": areaID}).Decode(&area)
	return area, mongoErr(err)
}

func (s *MongoPlatformStore) CreateSystemStreet(ctx context.Context, street domain.SystemStreet) (domain.SystemStreet, error) {
	area, err := s.GetSystemArea(ctx, street.AreaID)
	if err != nil || area.CityID != street.CityID {
		return domain.SystemStreet{}, ErrNotFound
	}
	now := time.Now().UTC()
	street.ID = support.NewID()
	street.CreatedAt = now
	street.UpdatedAt = now
	if street.Status == "" {
		street.Status = domain.SystemLocationActive
	}
	_, err = s.db.Collection("system_streets").InsertOne(ctx, street)
	return street, mongoErr(err)
}

func (s *MongoPlatformStore) ListSystemStreets(ctx context.Context, cityID, areaID string) ([]domain.SystemStreet, error) {
	var result []domain.SystemStreet
	err := findAll(ctx, s.db.Collection("system_streets"), bson.M{"cityId": cityID, "areaId": areaID, "status": domain.SystemLocationActive}, &result)
	return result, err
}

func (s *MongoPlatformStore) GetSystemStreet(ctx context.Context, streetID string) (domain.SystemStreet, error) {
	var street domain.SystemStreet
	err := s.db.Collection("system_streets").FindOne(ctx, bson.M{"id": streetID}).Decode(&street)
	return street, mongoErr(err)
}

func (s *MongoPlatformStore) CreateSystemNeighborhood(ctx context.Context, neighborhood domain.SystemNeighborhood) (domain.SystemNeighborhood, error) {
	street, err := s.GetSystemStreet(ctx, neighborhood.StreetID)
	if err != nil || street.CityID != neighborhood.CityID || street.AreaID != neighborhood.AreaID {
		return domain.SystemNeighborhood{}, ErrNotFound
	}
	now := time.Now().UTC()
	neighborhood.ID = support.NewID()
	neighborhood.CreatedAt = now
	neighborhood.UpdatedAt = now
	if neighborhood.Status == "" {
		neighborhood.Status = domain.SystemLocationActive
	}
	_, err = s.db.Collection("system_neighborhoods").InsertOne(ctx, neighborhood)
	return neighborhood, mongoErr(err)
}

func (s *MongoPlatformStore) ListSystemNeighborhoods(ctx context.Context, cityID, areaID, streetID string) ([]domain.SystemNeighborhood, error) {
	var result []domain.SystemNeighborhood
	err := findAll(ctx, s.db.Collection("system_neighborhoods"), bson.M{
		"cityId": cityID, "areaId": areaID, "streetId": streetID, "status": domain.SystemLocationActive,
	}, &result)
	return result, err
}

func (s *MongoPlatformStore) GetSystemNeighborhood(ctx context.Context, neighborhoodID string) (domain.SystemNeighborhood, error) {
	var neighborhood domain.SystemNeighborhood
	err := s.db.Collection("system_neighborhoods").FindOne(ctx, bson.M{"id": neighborhoodID}).Decode(&neighborhood)
	return neighborhood, mongoErr(err)
}

func (s *MongoPlatformStore) CreateLocationSuggestion(ctx context.Context, suggestion domain.LocationSuggestion) (domain.LocationSuggestion, error) {
	suggestion.ID = support.NewID()
	suggestion.CreatedAt = time.Now().UTC()
	if suggestion.Status == "" {
		suggestion.Status = domain.LocationSuggestionPending
	}
	_, err := s.db.Collection("location_suggestions").InsertOne(ctx, suggestion)
	return suggestion, mongoErr(err)
}

func (s *MongoPlatformStore) ListLocationSuggestions(ctx context.Context, status domain.LocationSuggestionStatus) ([]domain.LocationSuggestion, error) {
	filter := bson.M{}
	if status != "" {
		filter["status"] = status
	}
	var result []domain.LocationSuggestion
	err := findAll(ctx, s.db.Collection("location_suggestions"), filter, &result)
	return result, err
}

func (s *MongoPlatformStore) GetLocationSuggestion(ctx context.Context, suggestionID string) (domain.LocationSuggestion, error) {
	var suggestion domain.LocationSuggestion
	err := s.db.Collection("location_suggestions").FindOne(ctx, bson.M{"id": suggestionID}).Decode(&suggestion)
	return suggestion, mongoErr(err)
}

func (s *MongoPlatformStore) UpdateLocationSuggestion(ctx context.Context, suggestion domain.LocationSuggestion) (domain.LocationSuggestion, error) {
	res, err := s.db.Collection("location_suggestions").ReplaceOne(ctx, bson.M{"id": suggestion.ID}, suggestion)
	if err != nil {
		return domain.LocationSuggestion{}, mongoErr(err)
	}
	if res.MatchedCount == 0 {
		return domain.LocationSuggestion{}, ErrNotFound
	}
	return suggestion, nil
}

func findAll[T any](ctx context.Context, collection *mongo.Collection, filter interface{}, result *[]T) error {
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return mongoErr(err)
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, result); err != nil {
		return mongoErr(err)
	}
	return nil
}

func findStoredData[T any](ctx context.Context, collection *mongo.Collection, filter interface{}) ([]T, error) {
	return findStoredDataWithOptions[T](ctx, collection, filter)
}

func findStoredDataWithOptions[T any](ctx context.Context, collection *mongo.Collection, filter interface{}, opts ...*options.FindOptions) ([]T, error) {
	cursor, err := collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, mongoErr(err)
	}
	defer cursor.Close(ctx)
	result := []T{}
	for cursor.Next(ctx) {
		var doc struct {
			Data T `bson:"data"`
		}
		if err := cursor.Decode(&doc); err != nil {
			return nil, mongoErr(err)
		}
		result = append(result, doc.Data)
	}
	if err := cursor.Err(); err != nil {
		return nil, mongoErr(err)
	}
	return result, nil
}

func findOneStoredData[T any](ctx context.Context, collection *mongo.Collection, filter interface{}) (T, error) {
	var doc struct {
		Data T `bson:"data"`
	}
	err := collection.FindOne(ctx, filter).Decode(&doc)
	if err != nil {
		var zero T
		return zero, mongoErr(err)
	}
	return doc.Data, nil
}

func mongoErr(err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, mongo.ErrNoDocuments) {
		return ErrNotFound
	}
	return err
}
