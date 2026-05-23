# Amlak CRM Architecture

## Product scope

This monorepo contains an initial production-oriented foundation for a Persian, RTL real estate CRM:

- `backend`: Go modular monolith, designed around Clean Architecture boundaries.
- `frontend`: Flutter app with GetX, responsive Material UI, and centralized theme/localization.
- `docker-compose.yml`: local MongoDB, Redis, MinIO, RabbitMQ, and backend wiring.
- `backend/cmd/realtime`: lightweight WebSocket microservice for realtime delivery.

The first release focuses on OTP authentication, user profile bootstrap, business creation, invitations, consultant/member management, file metadata, and dashboard placeholders.

Platform-level admin roles and city-based system location approval are tracked in
`docs/ADMIN_RBAC_AND_LOCATIONS.md`. These are intentionally separate from
business-scoped member roles and business-specific locations.

## Backend structure

```text
backend/
  cmd/api                 application entrypoint
  cmd/realtime            realtime WebSocket service entrypoint
  internal/config         environment configuration
  internal/domain         business models, roles, permissions
  internal/repository     persistence ports and current in-memory adapter
  internal/service        usecases
  internal/transport/http HTTP handlers, middleware, routing
```

The backend is a modular monolith. Each feature is isolated behind service and repository contracts so it can later move behind gRPC without rewriting handlers or domain rules.

## Service communication decision

Use both, but at different levels:

- RabbitMQ for async domain events such as `consultant.invited`, `otp.requested`, `license.expiring`, and future CRM reminders.
- gRPC for future synchronous service-to-service APIs once modules become independent services.

The MVP keeps these modules in-process and exposes HTTP REST under `/api/v1`.

## Identity cache and events

User and session reads now go through a reusable write-through cache wrapper:

- `internal/cache`: generic TTL cache abstraction with Redis and in-process memory adapters.
- `internal/repository/CachedStore`: wraps the primary store and caches hot `User` and active `Session` records for `IDENTITY_CACHE_TTL_SECONDS`, default `60`.
- User/session writes update the primary store first and the cache immediately after.
- Revoked or expired sessions are evicted from cache but remain in the primary database/store for history and audit.
- Protected requests validate the `sessionId`, touch `lastSeenAt`, and refresh the cached session.

When `REDIS_ADDR` is configured and reachable, the API uses Redis for this identity cache. If Redis is unavailable in development, it falls back to the in-process TTL cache and logs the condition. Cache data is disposable and short-lived; the primary database/store remains the source of truth.

Inter-service communication uses RabbitMQ via `internal/events`:

- Exchange: `EVENT_EXCHANGE`, default `amlak.events`.
- Published events include `user.upserted`, `user.updated`, `session.saved`, `session.touched`, and `session.revoked`.
- In development, if `RABBITMQ_URL` is empty or unavailable, the API falls back to a no-op publisher and logs the condition.

## Core MongoDB collections

- `users`: phone, display name, profile status, devices, audit timestamps.
- `otp_attempts`: phone, code hash, expiry, send count, failed attempts, rate limit metadata.
- `sessions`: user, refresh token hash, device metadata, expiry, revoked state.
- `businesses`: owner, phones, address, working hours, license status, logo file id.
- `business_members`: user, business, role, permissions, commission percent, active status.
- `invitations`: business, invitee phone/user, role, commission, status, expiry, inviter.
- `files`: object storage provider, bucket/key, content type, size, public/private URL.
- `audit_logs`: actor, action, resource, IP, user agent, metadata.

Future platform/admin collections:

- `admin_accounts`: platform admin users, global roles, status, audit metadata.
- `cities`: shared city catalog.
- `system_areas`: approved shared area catalog per city.
- `system_streets`: approved shared street catalog per city and area.
- `system_neighborhoods`: approved shared neighborhood catalog per street.
- `location_suggestions`: user/business submitted location suggestions pending admin review.

These platform/admin collections are already wired to MongoDB through a hybrid
store adapter. The rest of the MVP store remains in-memory until its persistence
adapter is added.

## Main API contracts

All responses use:

```json
{ "data": {}, "error": null }
```

Errors use:

```json
{ "data": null, "error": { "code": "validation_error", "message": "..." } }
```

### Auth

- `POST /api/v1/auth/request-otp`
- `POST /api/v1/auth/verify-otp`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`
- `GET /api/v1/auth/me`
- `PATCH /api/v1/auth/profile`
- `GET /api/v1/auth/security`
- `PATCH /api/v1/auth/privacy`
- `GET /api/v1/auth/sessions`
- `DELETE /api/v1/auth/sessions/{sessionId}`

Access tokens include `userId` and `sessionId`. The API validates the session on protected requests, updates `lastSeenAt`, and lets the user revoke other active sessions from the profile/security UI.

### Realtime

- `GET /healthz` on the realtime service returns service health and active connection count.
- `GET /ws?token={accessToken}` upgrades to WebSocket after JWT validation.
- `POST /internal/publish` publishes a JSON message to all active connections for one user.

The realtime service is intentionally small: no business logic, no database dependency, sharded in-memory connection maps, bounded per-client queues, disabled compression, ping keepalive, and a clear path to Redis Pub/Sub or NATS for multi-node fanout. It is designed to be horizontally replicated behind a load balancer; sticky sessions are helpful but not required once broker fanout is enabled.

### Business

- `POST /api/v1/businesses`
- `GET /api/v1/businesses`
- `GET /api/v1/businesses/{businessId}/dashboard`
- `PATCH /api/v1/businesses/{businessId}`

### Invitations and members

- `POST /api/v1/businesses/{businessId}/invitations`
- `GET /api/v1/businesses/{businessId}/invitations`
- `GET /api/v1/invitations/inbox`
- `POST /api/v1/invitations/{invitationId}/accept`
- `POST /api/v1/invitations/{invitationId}/reject`
- `GET /api/v1/businesses/{businessId}/members`
- `PATCH /api/v1/businesses/{businessId}/members/{memberId}`

### Files

- `POST /api/v1/files/business-logo`
- `DELETE /api/v1/files/{fileId}`

## Authorization model

Roles:

- `owner`: full control.
- `manager`: operational admin, cannot manage owner/license/destructive business actions.
- `consultant`: CRM operations for assigned scope.

Permissions are stored separately from role names. Role defaults are applied on membership creation, then can be extended later.

Platform admin authorization is a separate model from business authorization.
Global roles such as `super_admin`, `sales_manager`, `finance_manager`, and
`training_manager` must be checked by platform permissions and must not be
implemented as `BusinessMember` roles.
