package domain

const (
	PermBusinessRead        = "business:read"
	PermBusinessUpdate      = "business:update"
	PermBusinessDelete      = "business:delete"
	PermLicenseManage       = "license:manage"
	PermInvitationsManage   = "invitations:manage"
	PermMembersRead         = "members:read"
	PermMembersManage       = "members:manage"
	PermMembersPromote      = "members:promote"
	PermLocationsManage     = "locations:manage"
	PermPropertiesManageOwn = "properties:manage_own"
	PermCRMManageOwn        = "crm:manage_own"
)

func DefaultPermissions(role Role) []string {
	switch role {
	case RoleOwner:
		return []string{
			PermBusinessRead, PermBusinessUpdate, PermBusinessDelete, PermLicenseManage,
			PermInvitationsManage, PermMembersRead, PermMembersManage, PermMembersPromote,
			PermLocationsManage, PermPropertiesManageOwn, PermCRMManageOwn,
		}
	case RoleManager:
		return []string{
			PermBusinessRead, PermBusinessUpdate, PermInvitationsManage,
			PermMembersRead, PermMembersManage, PermLocationsManage, PermPropertiesManageOwn, PermCRMManageOwn,
		}
	default:
		return []string{PermBusinessRead, PermPropertiesManageOwn, PermCRMManageOwn}
	}
}

func HasPermission(member BusinessMember, permission string) bool {
	for _, item := range member.Permissions {
		if item == permission {
			return true
		}
	}
	return false
}
