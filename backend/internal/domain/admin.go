package domain

import "time"

type AdminRole string

const (
	AdminRoleSuperAdmin      AdminRole = "super_admin"
	AdminRoleSalesManager    AdminRole = "sales_manager"
	AdminRoleFinanceManager  AdminRole = "finance_manager"
	AdminRoleTrainingManager AdminRole = "training_manager"
)

const (
	PermPlatformAdminsManage     = "platform:admins:manage"
	PermPlatformUsersRead        = "platform:users:read"
	PermPlatformUsersUpdate      = "platform:users:update"
	PermPlatformBusinessesRead   = "platform:businesses:read"
	PermPlatformBusinessesUpdate = "platform:businesses:update"
	PermPlatformVisitorsRead     = "platform:visitors:read"
	PermPlatformVisitorSalesRead = "platform:visitors:sales_read"
	PermPlatformCustomersRead    = "platform:customers:read"
	PermPlatformCustomersUpdate  = "platform:customers:update"
	PermPlatformFilesRead        = "platform:files:read"
	PermPlatformFilesUpdate      = "platform:files:update"
	PermPlatformFinanceRead      = "platform:finance:read"
	PermPlatformProfitLossRead   = "platform:profit_loss:read"
	PermPlatformLocationsRead    = "platform:locations:read"
	PermPlatformLocationsReview  = "platform:locations:review"
	PermPlatformLocationsManage  = "platform:locations:manage"
	PermPlatformTrainingRead     = "platform:training:read"
	PermPlatformTrainingUpdate   = "platform:training:update"
)

type AdminAccount struct {
	ID               string      `json:"id" bson:"id"`
	UserID           string      `json:"userId" bson:"userId"`
	Roles            []AdminRole `json:"roles" bson:"roles"`
	Permissions      []string    `json:"permissions" bson:"permissions"`
	Status           string      `json:"status" bson:"status"`
	CreatedByAdminID string      `json:"createdByAdminId,omitempty" bson:"createdByAdminId,omitempty"`
	CreatedAt        time.Time   `json:"createdAt" bson:"createdAt"`
	UpdatedAt        time.Time   `json:"updatedAt" bson:"updatedAt"`
}

func DefaultPlatformPermissions(role AdminRole) []string {
	switch role {
	case AdminRoleSuperAdmin:
		return []string{"*"}
	case AdminRoleSalesManager:
		return []string{
			PermPlatformUsersRead,
			PermPlatformBusinessesRead,
			PermPlatformVisitorsRead,
			PermPlatformVisitorSalesRead,
		}
	case AdminRoleFinanceManager:
		return []string{
			PermPlatformBusinessesRead,
			PermPlatformFinanceRead,
			PermPlatformProfitLossRead,
		}
	case AdminRoleTrainingManager:
		return []string{
			PermPlatformUsersRead,
			PermPlatformBusinessesRead,
			PermPlatformVisitorsRead,
			PermPlatformCustomersRead,
			PermPlatformCustomersUpdate,
			PermPlatformFilesRead,
			PermPlatformTrainingRead,
			PermPlatformTrainingUpdate,
		}
	default:
		return nil
	}
}

func HasPlatformPermission(account AdminAccount, permission string) bool {
	if account.Status != "active" {
		return false
	}
	for _, item := range account.Permissions {
		if item == "*" || item == permission {
			return true
		}
	}
	for _, role := range account.Roles {
		for _, item := range DefaultPlatformPermissions(role) {
			if item == "*" || item == permission {
				return true
			}
		}
	}
	return false
}
