package domain

import "time"

type PlatformSettings struct {
	ID               string    `json:"id" bson:"id"`
	OTPAPIKey        string    `json:"-" bson:"otpApiKey"`
	OTPAPIKeyMasked  string    `json:"otpApiKeyMasked" bson:"-"`
	ServiceSMSAPIKey string    `json:"-" bson:"serviceSmsApiKey"`
	ServiceSMSMasked string    `json:"serviceSmsApiKeyMasked" bson:"-"`
	UpdatedByAdminID string    `json:"updatedByAdminId,omitempty" bson:"updatedByAdminId,omitempty"`
	UpdatedAt        time.Time `json:"updatedAt" bson:"updatedAt"`
}
