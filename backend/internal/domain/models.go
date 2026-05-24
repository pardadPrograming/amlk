package domain

import "time"

type Role string

const (
	RoleOwner      Role = "owner"
	RoleManager    Role = "manager"
	RoleConsultant Role = "consultant"
)

type InvitationStatus string

const (
	InvitationPending   InvitationStatus = "pending"
	InvitationAccepted  InvitationStatus = "accepted"
	InvitationRejected  InvitationStatus = "rejected"
	InvitationExpired   InvitationStatus = "expired"
	InvitationCancelled InvitationStatus = "cancelled"
)

type MemberStatus string

const (
	MemberActive   MemberStatus = "active"
	MemberDisabled MemberStatus = "disabled"
	MemberRemoved  MemberStatus = "removed"
)

type User struct {
	ID              string          `json:"id"`
	Phone           string          `json:"phone"`
	FirstName       string          `json:"firstName"`
	LastName        string          `json:"lastName"`
	DisplayName     string          `json:"displayName"`
	CityID          string          `json:"cityId,omitempty"`
	PrivacySettings PrivacySettings `json:"privacySettings"`
	CreatedAt       time.Time       `json:"createdAt"`
	UpdatedAt       time.Time       `json:"updatedAt"`
}

type PrivacySettings struct {
	ShowPhoneToTeam    bool `json:"showPhoneToTeam"`
	ShowActivityStatus bool `json:"showActivityStatus"`
	AllowInviteByPhone bool `json:"allowInviteByPhone"`
}

type Session struct {
	ID           string    `json:"id"`
	UserID       string    `json:"userId"`
	RefreshToken string    `json:"-"`
	UserAgent    string    `json:"userAgent"`
	DeviceName   string    `json:"deviceName"`
	DeviceType   string    `json:"deviceType"`
	Browser      string    `json:"browser"`
	OS           string    `json:"os"`
	IP           string    `json:"ip"`
	LastSeenAt   time.Time `json:"lastSeenAt"`
	ExpiresAt    time.Time `json:"expiresAt"`
	RevokedAt    time.Time `json:"revokedAt,omitempty"`
	CreatedAt    time.Time `json:"createdAt"`
}

type OTPChallenge struct {
	Phone       string    `json:"phone"`
	Code        string    `json:"-"`
	ExpiresAt   time.Time `json:"expiresAt"`
	LastSentAt  time.Time `json:"lastSentAt"`
	Attempts    int       `json:"attempts"`
	SendCount   int       `json:"sendCount"`
	Development string    `json:"developmentCode,omitempty"`
}

type Business struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Phones        []string  `json:"phones"`
	Address       string    `json:"address"`
	WorkingHours  string    `json:"workingHours"`
	LicenseNumber string    `json:"licenseNumber"`
	LicenseStatus string    `json:"licenseStatus"`
	LogoFileID    string    `json:"logoFileId,omitempty"`
	OwnerUserID   string    `json:"ownerUserId"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type BusinessMember struct {
	ID                string       `json:"id"`
	BusinessID        string       `json:"businessId"`
	UserID            string       `json:"userId"`
	UserPhone         string       `json:"userPhone"`
	UserDisplayName   string       `json:"userDisplayName"`
	Role              Role         `json:"role"`
	Permissions       []string     `json:"permissions"`
	CommissionPercent float64      `json:"commissionPercent"`
	Status            MemberStatus `json:"status"`
	JoinedAt          time.Time    `json:"joinedAt"`
	UpdatedAt         time.Time    `json:"updatedAt"`
}

type Invitation struct {
	ID                string           `json:"id"`
	BusinessID        string           `json:"businessId"`
	BusinessName      string           `json:"businessName"`
	InviterUserID     string           `json:"inviterUserId"`
	InviteePhone      string           `json:"inviteePhone"`
	InviteeUserID     string           `json:"inviteeUserId,omitempty"`
	Role              Role             `json:"role"`
	CommissionPercent float64          `json:"commissionPercent"`
	Status            InvitationStatus `json:"status"`
	ExpiresAt         time.Time        `json:"expiresAt"`
	CreatedAt         time.Time        `json:"createdAt"`
	UpdatedAt         time.Time        `json:"updatedAt"`
}

type FileObject struct {
	ID          string    `json:"id"`
	OwnerID     string    `json:"ownerId"`
	UploaderID  string    `json:"uploaderId,omitempty"`
	Purpose     string    `json:"purpose,omitempty"`
	TargetType  string    `json:"targetType,omitempty"`
	TargetID    string    `json:"targetId,omitempty"`
	Status      string    `json:"status,omitempty"`
	Provider    string    `json:"provider"`
	Bucket      string    `json:"bucket"`
	Key         string    `json:"key"`
	PreloadKeys []string  `json:"preloadKeys,omitempty"`
	URL         string    `json:"url"`
	Kind        string    `json:"kind,omitempty"`
	ContentType string    `json:"contentType"`
	Size        int64     `json:"size"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	ExpiresAt   time.Time `json:"expiresAt,omitempty"`
}

type Notification struct {
	ID         string    `json:"id"`
	UserID     string    `json:"userId"`
	DedupKey   string    `json:"dedupKey,omitempty"`
	Type       string    `json:"type"`
	Title      string    `json:"title"`
	Body       string    `json:"body"`
	BusinessID string    `json:"businessId,omitempty"`
	PropertyID string    `json:"propertyId,omitempty"`
	RequestID  string    `json:"requestId,omitempty"`
	ReadAt     time.Time `json:"readAt,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
}

type Area struct {
	ID         string    `json:"id"`
	BusinessID string    `json:"businessId"`
	CityID     string    `json:"cityId"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type Street struct {
	ID         string    `json:"id"`
	BusinessID string    `json:"businessId"`
	CityID     string    `json:"cityId"`
	AreaID     string    `json:"areaId"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type Neighborhood struct {
	ID         string    `json:"id"`
	BusinessID string    `json:"businessId"`
	CityID     string    `json:"cityId"`
	AreaID     string    `json:"areaId"`
	StreetID   string    `json:"streetId"`
	Name       string    `json:"name"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type PropertyFileType string

const (
	PropertyFileSale        PropertyFileType = "sale"
	PropertyFilePartnership PropertyFileType = "partnership"
	PropertyFileRentLease   PropertyFileType = "rent_lease"
)

type PropertyFile struct {
	ID                           string                   `json:"id"`
	BusinessID                   string                   `json:"businessId"`
	OwnerUserID                  string                   `json:"ownerUserId"`
	SharedFromFileID             string                   `json:"sharedFromFileId,omitempty"`
	SharedFromOwnerID            string                   `json:"sharedFromOwnerId,omitempty"`
	IsPartnershipCopy            bool                     `json:"isPartnershipCopy,omitempty"`
	PartnershipCommissionPercent float64                  `json:"partnershipCommissionPercent,omitempty"`
	Status                       string                   `json:"status"`
	Type                         PropertyFileType         `json:"type"`
	Types                        []PropertyFileType       `json:"types"`
	Title                        string                   `json:"title"`
	Description                  string                   `json:"description"`
	InternalDescription          string                   `json:"internalDescription,omitempty"`
	SalePrice                    int64                    `json:"salePrice,omitempty"`
	FinalPrice                   int64                    `json:"finalPrice,omitempty"`
	DepositPrice                 int64                    `json:"depositPrice,omitempty"`
	RentPrice                    int64                    `json:"rentPrice,omitempty"`
	Convertible                  bool                     `json:"convertible,omitempty"`
	MaxConvertibleDeposit        int64                    `json:"maxConvertibleDeposit,omitempty"`
	RentWithOwner                bool                     `json:"rentWithOwner,omitempty"`
	HouseInfo                    HouseInfo                `json:"houseInfo"`
	Addresses                    []PropertyAddress        `json:"addresses"`
	Media                        []PropertyMedia          `json:"media"`
	VaultIDs                     []string                 `json:"vaultIds,omitempty"`
	VaultPlacements              []PropertyVaultPlacement `json:"vaultPlacements,omitempty"`
	BusinessCommissionPercent    float64                  `json:"businessCommissionPercent"`
	OwnerCommissionPercent       float64                  `json:"ownerCommissionPercent"`
	SharingHistory               []PropertySharingHistory `json:"sharingHistory,omitempty"`
	CreatedAt                    time.Time                `json:"createdAt"`
	UpdatedAt                    time.Time                `json:"updatedAt"`
}

type PropertyShareRequestStatus string

const (
	PropertySharePending  PropertyShareRequestStatus = "pending"
	PropertyShareApproved PropertyShareRequestStatus = "approved"
	PropertyShareRejected PropertyShareRequestStatus = "rejected"
	PropertyShareReceived PropertyShareRequestStatus = "received"
)

type PropertyShareRequest struct {
	ID                string                     `json:"id"`
	BusinessID        string                     `json:"businessId"`
	PropertyFileID    string                     `json:"propertyFileId"`
	PropertyTitle     string                     `json:"propertyTitle"`
	OwnerUserID       string                     `json:"ownerUserId"`
	RequesterUserID   string                     `json:"requesterUserId"`
	RequesterName     string                     `json:"requesterName,omitempty"`
	RequesterPhone    string                     `json:"requesterPhone,omitempty"`
	CommissionPercent float64                    `json:"commissionPercent"`
	Status            PropertyShareRequestStatus `json:"status"`
	SharedCopyFileID  string                     `json:"sharedCopyFileId,omitempty"`
	CreatedAt         time.Time                  `json:"createdAt"`
	UpdatedAt         time.Time                  `json:"updatedAt"`
	DecidedAt         time.Time                  `json:"decidedAt,omitempty"`
	ReceivedAt        time.Time                  `json:"receivedAt,omitempty"`
}

type PropertyOfferStatus string

const (
	PropertyOfferCandidate         PropertyOfferStatus = "candidate"
	PropertyOfferSent              PropertyOfferStatus = "sent"
	PropertyOfferRequesterApproved PropertyOfferStatus = "requester_approved"
	PropertyOfferApproved          PropertyOfferStatus = "approved"
	PropertyOfferRejected          PropertyOfferStatus = "rejected"
)

type PropertyOfferHistoryEntry struct {
	ID         string              `json:"id"`
	ActorID    string              `json:"actorId,omitempty"`
	Action     string              `json:"action"`
	FromStatus PropertyOfferStatus `json:"fromStatus,omitempty"`
	ToStatus   PropertyOfferStatus `json:"toStatus,omitempty"`
	Note       string              `json:"note,omitempty"`
	CreatedAt  time.Time           `json:"createdAt"`
}

type PropertyOffer struct {
	ID                string                      `json:"id"`
	DedupKey          string                      `json:"dedupKey,omitempty"`
	BusinessID        string                      `json:"businessId"`
	PropertyFileID    string                      `json:"propertyFileId"`
	PropertyTitle     string                      `json:"propertyTitle"`
	OwnerUserID       string                      `json:"ownerUserId"`
	OwnerName         string                      `json:"ownerName,omitempty"`
	Owner             *User                       `json:"owner,omitempty"`
	RequesterUserID   string                      `json:"requesterUserId"`
	RequesterName     string                      `json:"requesterName,omitempty"`
	ContactID         string                      `json:"contactId,omitempty"`
	ContactName       string                      `json:"contactName,omitempty"`
	RequestID         string                      `json:"requestId,omitempty"`
	RequestTitle      string                      `json:"requestTitle,omitempty"`
	CommissionPercent float64                     `json:"commissionPercent"`
	Score             int                         `json:"score,omitempty"`
	Tier              string                      `json:"tier,omitempty"`
	Status            PropertyOfferStatus         `json:"status"`
	ChatChannelID     string                      `json:"chatChannelId,omitempty"`
	SharedCopyFileID  string                      `json:"sharedCopyFileId,omitempty"`
	PropertyFile      *PropertyFile               `json:"propertyFile,omitempty"`
	History           []PropertyOfferHistoryEntry `json:"history,omitempty"`
	CreatedAt         time.Time                   `json:"createdAt"`
	UpdatedAt         time.Time                   `json:"updatedAt"`
}

type PropertySharingHistory struct {
	RequestID         string                     `json:"requestId"`
	UserID            string                     `json:"userId"`
	UserName          string                     `json:"userName,omitempty"`
	UserPhone         string                     `json:"userPhone,omitempty"`
	CommissionPercent float64                    `json:"commissionPercent"`
	Status            PropertyShareRequestStatus `json:"status"`
	SharedCopyFileID  string                     `json:"sharedCopyFileId,omitempty"`
	CreatedAt         time.Time                  `json:"createdAt"`
	UpdatedAt         time.Time                  `json:"updatedAt"`
}

type PropertyVaultPlacement struct {
	VaultID           string  `json:"vaultId"`
	CommissionPercent float64 `json:"commissionPercent"`
}

type HouseInfo struct {
	PropertyType         string   `json:"propertyType"`
	AreaM2               int      `json:"areaM2"`
	Bedrooms             int      `json:"bedrooms"`
	Floor                int      `json:"floor"`
	TotalFloors          int      `json:"totalFloors"`
	AgeYears             int      `json:"ageYears"`
	Renovated            bool     `json:"renovated"`
	Parking              bool     `json:"parking"`
	Elevator             bool     `json:"elevator"`
	Storage              bool     `json:"storage"`
	TerraceCount         int      `json:"terraceCount"`
	BackyardAreaM2       int      `json:"backyardAreaM2"`
	SeparateEntrance     bool     `json:"separateEntrance"`
	DocumentType         string   `json:"documentType"`
	Directions           []string `json:"directions,omitempty"`
	Flooring             string   `json:"flooring,omitempty"`
	Heating              string   `json:"heating,omitempty"`
	CabinetType          string   `json:"cabinetType,omitempty"`
	Cooling              []string `json:"cooling,omitempty"`
	WallCovering         string   `json:"wallCovering,omitempty"`
	KitchenType          string   `json:"kitchenType,omitempty"`
	StoveTop             bool     `json:"stoveTop"`
	ParkingDisturbance   bool     `json:"parkingDisturbance"`
	ParkingInDeed        bool     `json:"parkingInDeed"`
	ParkingCommon        bool     `json:"parkingCommon"`
	MasterServiceCount   int      `json:"masterServiceCount"`
	DoubleGlazedWindows  bool     `json:"doubleGlazedWindows"`
	GardenBuilding       bool     `json:"gardenBuilding"`
	GardenBuildingAreaM2 int      `json:"gardenBuildingAreaM2"`
	GardenBuildingFloors int      `json:"gardenBuildingFloors"`
	Pool                 bool     `json:"pool"`
	WaterUtility         bool     `json:"waterUtility"`
	ElectricityUtility   bool     `json:"electricityUtility"`
	GasUtility           bool     `json:"gasUtility"`
	WaterRight           bool     `json:"waterRight"`
	Permit               bool     `json:"permit"`
	LandType             string   `json:"landType,omitempty"`
}

type PropertyAddress struct {
	AreaID             string `json:"areaId"`
	AreaName           string `json:"areaName"`
	StreetID           string `json:"streetId"`
	StreetName         string `json:"streetName"`
	NeighborhoodID     string `json:"neighborhoodId"`
	NeighborhoodName   string `json:"neighborhoodName"`
	ManualExactAddress string `json:"manualExactAddress"`
}

type PropertyMedia struct {
	ID          string    `json:"id"`
	FileID      string    `json:"fileId"`
	Kind        string    `json:"kind"`
	URL         string    `json:"url"`
	ContentType string    `json:"contentType"`
	Size        int64     `json:"size"`
	Width       int       `json:"width,omitempty"`
	Height      int       `json:"height,omitempty"`
	DurationSec int       `json:"durationSec,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
}
