package domain

import "time"

const (
	ContactTagRealEstateConsultant = "مشاور املاک"
	ContactTagRealEstateOffice     = "املاک"
	ContactTagPropertyOwner        = "مالک املاک"
	ContactTagBuilder              = "سازنده"
	ContactTagOwner                = "مالک"
	ContactTagTenant               = "مستاجر"
)

type ContactPhone struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

type ContactPropertyRef struct {
	PropertyID string `json:"propertyId,omitempty"`
	Title      string `json:"title"`
	Note       string `json:"note,omitempty"`
}

type ContactRequestLocation struct {
	ID          string `json:"id,omitempty"`
	Level       string `json:"level"`
	Name        string `json:"name"`
	AreaID      string `json:"areaId,omitempty"`
	StreetID    string `json:"streetId,omitempty"`
	IncludeAll  bool   `json:"includeAll"`
	Preference  string `json:"preference"`
	Description string `json:"description,omitempty"`
}

type ContactRequestFloorRule struct {
	FloorMin   int    `json:"floorMin,omitempty"`
	FloorMax   int    `json:"floorMax,omitempty"`
	Elevator   bool   `json:"elevator"`
	Preference string `json:"preference"`
}

type ContactRequestOptionFilter struct {
	Key        string   `json:"key"`
	Values     []string `json:"values"`
	Preference string   `json:"preference"`
}

type ContactRequestBooleanFilter struct {
	Key        string `json:"key"`
	Value      bool   `json:"value"`
	Preference string `json:"preference"`
}

type ContactRequestNumberFilter struct {
	Key        string `json:"key"`
	Min        int    `json:"min,omitempty"`
	Preference string `json:"preference"`
}

type ContactRequestHistoryChange struct {
	Field string `json:"field"`
	From  string `json:"from,omitempty"`
	To    string `json:"to,omitempty"`
}

type ContactRequestHistoryEntry struct {
	ID          string                        `json:"id"`
	ChangedByID string                        `json:"changedById,omitempty"`
	ChangedAt   time.Time                     `json:"changedAt"`
	Description string                        `json:"description"`
	Changes     []ContactRequestHistoryChange `json:"changes"`
}

type ContactRequest struct {
	ID                    string                        `json:"id"`
	Title                 string                        `json:"title"`
	Type                  string                        `json:"type,omitempty"`
	Status                string                        `json:"status,omitempty"`
	BudgetMin             int64                         `json:"budgetMin,omitempty"`
	BudgetMax             int64                         `json:"budgetMax,omitempty"`
	PurchaseMin           int64                         `json:"purchaseMin,omitempty"`
	PurchaseMax           int64                         `json:"purchaseMax,omitempty"`
	SuggestedPurchaseMin  int64                         `json:"suggestedPurchaseMin,omitempty"`
	SuggestedPurchaseMax  int64                         `json:"suggestedPurchaseMax,omitempty"`
	PartnershipMin        int64                         `json:"partnershipMin,omitempty"`
	PartnershipMax        int64                         `json:"partnershipMax,omitempty"`
	ShareMin              int                           `json:"shareMin,omitempty"`
	ShareMax              int                           `json:"shareMax,omitempty"`
	DepositMin            int64                         `json:"depositMin,omitempty"`
	DepositMax            int64                         `json:"depositMax,omitempty"`
	SuggestedDepositMin   int64                         `json:"suggestedDepositMin,omitempty"`
	SuggestedDepositMax   int64                         `json:"suggestedDepositMax,omitempty"`
	RentMin               int64                         `json:"rentMin,omitempty"`
	RentMax               int64                         `json:"rentMax,omitempty"`
	SuggestedRentMin      int64                         `json:"suggestedRentMin,omitempty"`
	SuggestedRentMax      int64                         `json:"suggestedRentMax,omitempty"`
	MinAreaM2             int                           `json:"minAreaM2,omitempty"`
	MaxAgeYears           int                           `json:"maxAgeYears,omitempty"`
	Convertible           bool                          `json:"convertible,omitempty"`
	MaxConvertibleDeposit int64                         `json:"maxConvertibleDeposit,omitempty"`
	RentWithOwner         bool                          `json:"rentWithOwner,omitempty"`
	LandMinAreaM2         int                           `json:"landMinAreaM2,omitempty"`
	BuildingMinAreaM2     int                           `json:"buildingMinAreaM2,omitempty"`
	PermitFloorsMin       int                           `json:"permitFloorsMin,omitempty"`
	Locations             []ContactRequestLocation      `json:"locations,omitempty"`
	FloorRules            []ContactRequestFloorRule     `json:"floorRules,omitempty"`
	OptionFilters         []ContactRequestOptionFilter  `json:"optionFilters,omitempty"`
	BooleanFilters        []ContactRequestBooleanFilter `json:"booleanFilters,omitempty"`
	NumberFilters         []ContactRequestNumberFilter  `json:"numberFilters,omitempty"`
	Note                  string                        `json:"note,omitempty"`
	ChangeDescription     string                        `json:"changeDescription,omitempty"`
	History               []ContactRequestHistoryEntry  `json:"history,omitempty"`
	CreatedAt             time.Time                     `json:"createdAt"`
}

type Contact struct {
	ID          string               `json:"id"`
	BusinessID  string               `json:"businessId"`
	CreatedByID string               `json:"createdById"`
	FirstName   string               `json:"firstName"`
	LastName    string               `json:"lastName"`
	DisplayName string               `json:"displayName"`
	Company     string               `json:"company,omitempty"`
	Phones      []ContactPhone       `json:"phones"`
	Tags        []string             `json:"tags"`
	Properties  []ContactPropertyRef `json:"properties"`
	Requests    []ContactRequest     `json:"requests"`
	Note        string               `json:"note,omitempty"`
	CreatedAt   time.Time            `json:"createdAt"`
	UpdatedAt   time.Time            `json:"updatedAt"`
}

type PropertyMatchResult struct {
	PropertyFile   PropertyFile          `json:"propertyFile"`
	Score          int                   `json:"score"`
	Tier           string                `json:"tier"`
	MatchedReasons []string              `json:"matchedReasons"`
	MissedReasons  []string              `json:"missedReasons"`
	Access         []PropertyMatchAccess `json:"access,omitempty"`
}

type PropertyMatchAccess struct {
	Source            string  `json:"source"`
	VaultID           string  `json:"vaultId,omitempty"`
	VaultTitle        string  `json:"vaultTitle,omitempty"`
	CommissionPercent float64 `json:"commissionPercent,omitempty"`
	Collaboration     bool    `json:"collaboration,omitempty"`
}

func SystemContactTags() []string {
	return []string{
		ContactTagRealEstateConsultant,
		ContactTagRealEstateOffice,
		ContactTagPropertyOwner,
		ContactTagBuilder,
		ContactTagOwner,
		ContactTagTenant,
	}
}
