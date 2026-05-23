package searchengine

import (
	"testing"

	"amlakcrm/backend/internal/domain"
)

func TestRentLeaseFinancialMatching(t *testing.T) {
	tests := []struct {
		name        string
		request     domain.ContactRequest
		file        domain.PropertyFile
		wantTier    string
		wantMatched bool
	}{
		{
			name: "direct point matches preferred range",
			request: domain.ContactRequest{
				Type:       "rent_lease",
				DepositMin: 500_000_000,
				DepositMax: 700_000_000,
				RentMax:    20_000_000,
			},
			file: domain.PropertyFile{
				Type:         domain.PropertyFileRentLease,
				Types:        []domain.PropertyFileType{domain.PropertyFileRentLease},
				DepositPrice: 600_000_000,
				RentPrice:    15_000_000,
			},
			wantTier:    "green",
			wantMatched: true,
		},
		{
			name: "convertible point matches using 100m to 3m rule",
			request: domain.ContactRequest{
				Type:       "rent_lease",
				DepositMin: 590_000_000,
				DepositMax: 610_000_000,
				RentMax:    17_000_000,
			},
			file: domain.PropertyFile{
				Type:                  domain.PropertyFileRentLease,
				Types:                 []domain.PropertyFileType{domain.PropertyFileRentLease},
				DepositPrice:          500_000_000,
				RentPrice:             20_000_000,
				Convertible:           true,
				MaxConvertibleDeposit: 600_000_000,
			},
			wantTier:    "green",
			wantMatched: true,
		},
		{
			name: "convertible range does not over-discount rent",
			request: domain.ContactRequest{
				Type:       "rent_lease",
				DepositMin: 590_000_000,
				DepositMax: 610_000_000,
				RentMax:    15_000_000,
			},
			file: domain.PropertyFile{
				Type:                  domain.PropertyFileRentLease,
				Types:                 []domain.PropertyFileType{domain.PropertyFileRentLease},
				DepositPrice:          500_000_000,
				RentPrice:             20_000_000,
				Convertible:           true,
				MaxConvertibleDeposit: 600_000_000,
			},
			wantMatched: false,
		},
		{
			name: "suggested range matches as yellow",
			request: domain.ContactRequest{
				Type:                "rent_lease",
				DepositMin:          500_000_000,
				DepositMax:          700_000_000,
				RentMax:             20_000_000,
				SuggestedDepositMin: 500_000_000,
				SuggestedDepositMax: 700_000_000,
				SuggestedRentMax:    25_000_000,
			},
			file: domain.PropertyFile{
				Type:         domain.PropertyFileRentLease,
				Types:        []domain.PropertyFileType{domain.PropertyFileRentLease},
				DepositPrice: 600_000_000,
				RentPrice:    25_000_000,
			},
			wantTier:    "yellow",
			wantMatched: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := matchProperty(tt.request, tt.file)
			if !tt.wantMatched {
				if result.Tier != "" {
					t.Fatalf("expected no match, got tier %q", result.Tier)
				}
				return
			}
			if result.Tier != tt.wantTier {
				t.Fatalf("expected tier %q, got %q", tt.wantTier, result.Tier)
			}
		})
	}
}
