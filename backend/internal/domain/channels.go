package domain

import "time"

type ChannelType string

const (
	ChannelTypePrivate       ChannelType = "private"
	ChannelTypeUserMain      ChannelType = "user_main"
	ChannelTypeBusinessMain  ChannelType = "business_main"
	ChannelTypeUserVault     ChannelType = "user_vault"
	ChannelTypeBusinessVault ChannelType = "business_vault"
)

type ChannelMemberStatus string
type ChannelMemberRole string

const (
	ChannelMemberActive  ChannelMemberStatus = "active"
	ChannelMemberInvited ChannelMemberStatus = "invited"

	ChannelMemberRoleMember ChannelMemberRole = "member"
	ChannelMemberRoleAdmin  ChannelMemberRole = "admin"
)

type Channel struct {
	ID          string      `json:"id"`
	Type        ChannelType `json:"type"`
	OwnerUserID string      `json:"ownerUserId,omitempty"`
	BusinessID  string      `json:"businessId,omitempty"`
	VaultID     string      `json:"vaultId,omitempty"`
	Title       string      `json:"title"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
}

type ChannelVault struct {
	ID          string    `json:"id"`
	ChannelID   string    `json:"channelId"`
	OwnerUserID string    `json:"ownerUserId,omitempty"`
	BusinessID  string    `json:"businessId,omitempty"`
	Title       string    `json:"title"`
	IsMain      bool      `json:"isMain"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ChannelMember struct {
	ID          string              `json:"id"`
	ChannelID   string              `json:"channelId"`
	UserID      string              `json:"userId,omitempty"`
	Phone       string              `json:"phone"`
	DisplayName string              `json:"displayName,omitempty"`
	Role        ChannelMemberRole   `json:"role"`
	Status      ChannelMemberStatus `json:"status"`
	IsOnline    bool                `json:"isOnline,omitempty"`
	LastSeenAt  time.Time           `json:"lastSeenAt,omitempty"`
	CreatedAt   time.Time           `json:"createdAt"`
	UpdatedAt   time.Time           `json:"updatedAt"`
}

type ChannelMedia struct {
	ID          string    `json:"id"`
	FileID      string    `json:"fileId"`
	Kind        string    `json:"kind"`
	URL         string    `json:"url"`
	ContentType string    `json:"contentType"`
	Size        int64     `json:"size"`
	CreatedAt   time.Time `json:"createdAt"`
	ExpiresAt   time.Time `json:"expiresAt,omitempty"`
}

type ChannelVaultFile struct {
	ID                string    `json:"id"`
	VaultID           string    `json:"vaultId,omitempty"`
	ChannelID         string    `json:"channelId"`
	UploaderID        string    `json:"uploaderId"`
	Title             string    `json:"title"`
	Note              string    `json:"note,omitempty"`
	SourceType        string    `json:"sourceType,omitempty"`
	FileID            string    `json:"fileId"`
	PropertyFileID    string    `json:"propertyFileId,omitempty"`
	PropertyStatus    string    `json:"propertyStatus,omitempty"`
	CommissionPercent float64   `json:"commissionPercent,omitempty"`
	Kind              string    `json:"kind"`
	URL               string    `json:"url"`
	ContentType       string    `json:"contentType"`
	Size              int64     `json:"size"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type ChannelMessage struct {
	ID             string                   `json:"id"`
	ChannelID      string                   `json:"channelId"`
	AuthorID       string                   `json:"authorId"`
	AuthorName     string                   `json:"authorName,omitempty"`
	Text           string                   `json:"text,omitempty"`
	Caption        string                   `json:"caption,omitempty"`
	Media          []ChannelMedia           `json:"media,omitempty"`
	ReplyToID      string                   `json:"replyToId,omitempty"`
	ReplyTo        *ChannelReplyPreview     `json:"replyTo,omitempty"`
	VaultFileRefID string                   `json:"vaultFileRefId,omitempty"`
	VaultFileRef   *ChannelVaultFilePreview `json:"vaultFileRef,omitempty"`
	SeenBy         []ChannelMessageSeen     `json:"seenBy,omitempty"`
	CreatedAt      time.Time                `json:"createdAt"`
	UpdatedAt      time.Time                `json:"updatedAt"`
}

type ChannelMessageSeen struct {
	UserID string    `json:"userId"`
	SeenAt time.Time `json:"seenAt"`
}

type ChannelReplyPreview struct {
	ID         string `json:"id"`
	AuthorID   string `json:"authorId"`
	AuthorName string `json:"authorName,omitempty"`
	Text       string `json:"text,omitempty"`
	Caption    string `json:"caption,omitempty"`
	MediaKind  string `json:"mediaKind,omitempty"`
	MediaCount int    `json:"mediaCount,omitempty"`
}

type ChannelVaultFilePreview struct {
	ID                string  `json:"id"`
	Title             string  `json:"title"`
	Kind              string  `json:"kind"`
	URL               string  `json:"url"`
	ContentType       string  `json:"contentType,omitempty"`
	Size              int64   `json:"size,omitempty"`
	PropertyStatus    string  `json:"propertyStatus,omitempty"`
	CommissionPercent float64 `json:"commissionPercent,omitempty"`
}
