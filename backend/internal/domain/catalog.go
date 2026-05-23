package domain

import "time"

type SystemLocationStatus string

const (
	SystemLocationActive   SystemLocationStatus = "active"
	SystemLocationDisabled SystemLocationStatus = "disabled"
	SystemLocationMerged   SystemLocationStatus = "merged"
)

type LocationSuggestionType string

const (
	LocationSuggestionArea         LocationSuggestionType = "area"
	LocationSuggestionStreet       LocationSuggestionType = "street"
	LocationSuggestionNeighborhood LocationSuggestionType = "neighborhood"
)

type LocationSuggestionStatus string

const (
	LocationSuggestionPending   LocationSuggestionStatus = "pending"
	LocationSuggestionApproved  LocationSuggestionStatus = "approved"
	LocationSuggestionRejected  LocationSuggestionStatus = "rejected"
	LocationSuggestionDuplicate LocationSuggestionStatus = "duplicate"
	LocationSuggestionMerged    LocationSuggestionStatus = "merged"
)

type City struct {
	ID             string               `json:"id" bson:"id"`
	Name           string               `json:"name" bson:"name"`
	NormalizedName string               `json:"normalizedName" bson:"normalizedName"`
	Status         SystemLocationStatus `json:"status" bson:"status"`
	CreatedAt      time.Time            `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time            `json:"updatedAt" bson:"updatedAt"`
}

type SystemArea struct {
	ID             string               `json:"id" bson:"id"`
	CityID         string               `json:"cityId" bson:"cityId"`
	Name           string               `json:"name" bson:"name"`
	NormalizedName string               `json:"normalizedName" bson:"normalizedName"`
	Status         SystemLocationStatus `json:"status" bson:"status"`
	MergedIntoID   string               `json:"mergedIntoId,omitempty" bson:"mergedIntoId,omitempty"`
	CreatedAt      time.Time            `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time            `json:"updatedAt" bson:"updatedAt"`
}

type SystemStreet struct {
	ID             string               `json:"id" bson:"id"`
	CityID         string               `json:"cityId" bson:"cityId"`
	AreaID         string               `json:"areaId" bson:"areaId"`
	Name           string               `json:"name" bson:"name"`
	NormalizedName string               `json:"normalizedName" bson:"normalizedName"`
	Status         SystemLocationStatus `json:"status" bson:"status"`
	MergedIntoID   string               `json:"mergedIntoId,omitempty" bson:"mergedIntoId,omitempty"`
	CreatedAt      time.Time            `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time            `json:"updatedAt" bson:"updatedAt"`
}

type SystemNeighborhood struct {
	ID             string               `json:"id" bson:"id"`
	CityID         string               `json:"cityId" bson:"cityId"`
	AreaID         string               `json:"areaId" bson:"areaId"`
	StreetID       string               `json:"streetId" bson:"streetId"`
	Name           string               `json:"name" bson:"name"`
	NormalizedName string               `json:"normalizedName" bson:"normalizedName"`
	Status         SystemLocationStatus `json:"status" bson:"status"`
	MergedIntoID   string               `json:"mergedIntoId,omitempty" bson:"mergedIntoId,omitempty"`
	CreatedAt      time.Time            `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time            `json:"updatedAt" bson:"updatedAt"`
}

type LocationSuggestion struct {
	ID                       string                   `json:"id" bson:"id"`
	CityID                   string                   `json:"cityId" bson:"cityId"`
	Type                     LocationSuggestionType   `json:"type" bson:"type"`
	Name                     string                   `json:"name" bson:"name"`
	NormalizedName           string                   `json:"normalizedName" bson:"normalizedName"`
	ParentAreaID             string                   `json:"parentAreaId,omitempty" bson:"parentAreaId,omitempty"`
	ParentStreetID           string                   `json:"parentStreetId,omitempty" bson:"parentStreetId,omitempty"`
	ManualParentName         string                   `json:"manualParentName,omitempty" bson:"manualParentName,omitempty"`
	SubmittedByUserID        string                   `json:"submittedByUserId" bson:"submittedByUserId"`
	BusinessID               string                   `json:"businessId,omitempty" bson:"businessId,omitempty"`
	SourcePropertyID         string                   `json:"sourcePropertyId,omitempty" bson:"sourcePropertyId,omitempty"`
	Status                   LocationSuggestionStatus `json:"status" bson:"status"`
	ReviewedByAdminID        string                   `json:"reviewedByAdminId,omitempty" bson:"reviewedByAdminId,omitempty"`
	ReviewNote               string                   `json:"reviewNote,omitempty" bson:"reviewNote,omitempty"`
	ApprovedSystemLocationID string                   `json:"approvedSystemLocationId,omitempty" bson:"approvedSystemLocationId,omitempty"`
	CreatedAt                time.Time                `json:"createdAt" bson:"createdAt"`
	ReviewedAt               time.Time                `json:"reviewedAt,omitempty" bson:"reviewedAt,omitempty"`
}
