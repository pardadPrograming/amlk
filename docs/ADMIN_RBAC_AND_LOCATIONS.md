# Admin RBAC and System Locations

This document captures the intended design for global administration roles and
city-based system locations.

## Problem

The current application has business-scoped roles:

- `owner`
- `manager`
- `consultant`

These roles control access inside one real estate business. They are not enough
for the platform-level management panel because platform admins need access
across users, businesses, visitors, sales, finance, training, and shared
location data.

The current location model is also business-scoped:

- `Area`
- `Street`
- `Neighborhood`

The target model needs shared, city-based system locations. If an area, street,
or neighborhood does not exist in the system catalog, a user can suggest it.
The suggestion appears in the global admin panel, and an authorized admin can
approve or reject it. Approved items become part of the system catalog.

## Location Model

Use two layers:

1. System catalog locations
2. User/business submitted location suggestions

### System Catalog

System locations should be city-based and reusable by every business.

Recommended entities:

- `City`
- `SystemArea`
- `SystemStreet`
- `SystemNeighborhood`

Recommended fields:

- `id`
- `cityId`
- `parentId` where relevant
- `name`
- `normalizedName`
- `status`: `active`, `disabled`, `merged`
- `createdAt`
- `updatedAt`
- `createdByAdminId`
- `mergedIntoId`

Hierarchy:

```text
City
  Area
    Street
      Neighborhood
```

### Suggestions

When a user cannot find a location, the frontend should let them submit a
suggestion instead of directly creating a system location.

Recommended entity: `LocationSuggestion`

Fields:

- `id`
- `cityId`
- `type`: `area`, `street`, `neighborhood`
- `name`
- `normalizedName`
- `parentAreaId`
- `parentStreetId`
- `manualParentName`
- `submittedByUserId`
- `businessId`
- `sourcePropertyId`
- `status`: `pending`, `approved`, `rejected`, `duplicate`, `merged`
- `reviewedByAdminId`
- `reviewNote`
- `approvedSystemLocationId`
- `createdAt`
- `reviewedAt`

Approval behavior:

- `pending`: visible in the admin panel.
- `approved`: creates a system area/street/neighborhood or links to an existing
  matching item.
- `duplicate`: links the suggestion to an already existing system location.
- `rejected`: remains in history but is not visible to regular users.
- `merged`: used when an admin consolidates near-duplicate names.

## Frontend Flow

In property/location forms:

1. User selects city.
2. User searches system area/street/neighborhood.
3. If the item does not exist, user clicks "add new suggestion".
4. The suggested item can be used temporarily in that file with a `pending`
   badge.
5. Admin reviews the suggestion in the global admin panel.
6. If approved, future users see it as a normal system location.

Business-specific custom locations should remain possible only when they are
needed for private internal labeling, but they must not silently pollute the
shared system catalog.

## Global Admin Roles

Add platform-scoped admin roles separate from business membership roles.

Recommended entity: `AdminAccount`

Fields:

- `id`
- `userId`
- `roles`
- `status`: `active`, `disabled`
- `createdAt`
- `updatedAt`
- `createdByAdminId`

Recommended global roles:

- `super_admin`
- `sales_manager`
- `finance_manager`
- `training_manager`

## Permission Matrix

| Capability | Super Admin | Sales Manager | Finance Manager | Training Manager |
| --- | --- | --- | --- | --- |
| Manage admin accounts | yes | no | no | no |
| Manage system locations | yes | no | no | limited review if assigned |
| View all businesses/accounts | yes | yes | finance summary only | yes |
| View visitors | yes | yes | no | yes |
| View visitor sales amounts | yes | yes | yes where financial | no |
| View customer details | yes | no | no | yes |
| Edit customer details | yes | no | no | yes if post-sale support |
| View user/property files | yes | no | no | yes if support/training related |
| Edit user/property files | yes | no | no | limited/support-only |
| Financial reports | yes | no | yes | no |
| Profit/loss reports | yes | no | yes | no |
| Pre-sales pipeline | yes | yes | no | read-only or no |
| Post-sales support/training | yes | no | no | yes |

Important separation:

- Sales manager is focused on pre-sales: visitors, leads, formed accounts, and
  sales pipeline visibility.
- Training manager is focused on post-sales: visitor/customer onboarding,
  account details, support/training state, and customer/user details, but not
  sales amounts or financial performance.
- Finance manager sees financial reports and profit/loss, but should not get
  operational edit access to customer files unless explicitly granted.
- Super admin bypasses all platform restrictions.

## Permission Names

Suggested platform permissions:

- `platform:admins:manage`
- `platform:users:read`
- `platform:users:update`
- `platform:businesses:read`
- `platform:businesses:update`
- `platform:visitors:read`
- `platform:visitors:sales_read`
- `platform:customers:read`
- `platform:customers:update`
- `platform:files:read`
- `platform:files:update`
- `platform:finance:read`
- `platform:profit_loss:read`
- `platform:locations:read`
- `platform:locations:review`
- `platform:locations:manage`
- `platform:training:read`
- `platform:training:update`

Default role permissions:

```text
super_admin:
  all platform permissions

sales_manager:
  platform:users:read
  platform:businesses:read
  platform:visitors:read
  platform:visitors:sales_read

finance_manager:
  platform:businesses:read
  platform:finance:read
  platform:profit_loss:read

training_manager:
  platform:users:read
  platform:businesses:read
  platform:visitors:read
  platform:customers:read
  platform:customers:update
  platform:files:read
  platform:training:read
  platform:training:update
```

## Backend API Draft

System locations:

- `GET /api/v1/catalog/cities`
- `GET /api/v1/catalog/cities/{cityId}/locations`
- `GET /api/v1/catalog/cities/{cityId}/locations/search?q=...`
- `POST /api/v1/location-suggestions`

Admin panel:

- `GET /api/v1/admin/location-suggestions?status=pending`
- `POST /api/v1/admin/location-suggestions/{suggestionId}/approve`
- `POST /api/v1/admin/location-suggestions/{suggestionId}/reject`
- `POST /api/v1/admin/location-suggestions/{suggestionId}/mark-duplicate`
- `GET /api/v1/admin/users`
- `GET /api/v1/admin/businesses`
- `GET /api/v1/admin/visitors`
- `GET /api/v1/admin/customers`
- `GET /api/v1/admin/finance/reports`
- `GET /api/v1/admin/profit-loss`
- `GET /api/v1/admin/training`

Admin routes must use a platform admin middleware, not the existing
business-member middleware.

## Implementation Notes

- Keep business roles and platform admin roles separate.
- Do not attach global admin privileges to `BusinessMember`.
- Add a platform permission checker, for example:
  `RequirePlatformPermission("platform:locations:review")`.
- Normalize Persian names before duplicate checks:
  trim spaces, normalize Arabic/Persian Yeh and Kaf, remove duplicate spaces.
- Store rejected and duplicate suggestions for audit and future deduplication.
- Do not expose sales amount fields to `training_manager` responses. Prefer
  response shaping at service/query layer, not only hiding fields in the UI.
- Do not expose customer/file details to `sales_manager` responses.
- Super admin can access every route, but actions should still be audit logged.

## Suggested Delivery Order

1. Add platform admin roles, permissions, middleware, and audit logging.
2. Add city/system location catalog models and repository methods.
3. Add location suggestion model and review workflow.
4. Update property/location forms to search catalog first and suggest missing
   locations.
5. Build the global admin panel shell.
6. Add location suggestion review screens.
7. Add role-specific admin dashboards for sales, finance, and training.
