# Amlak CRM

Initial Flutter + Go foundation for a Persian real estate CRM.

## Product design notes

- Admin RBAC and city-based system locations: `docs/ADMIN_RBAC_AND_LOCATIONS.md`

## Run backend

```powershell
cd C:\projects\amlak\backend
go run .\cmd\api
```

Health check:

```powershell
Invoke-RestMethod http://localhost:8080/healthz
```

OTP uses a mock provider in development and returns `developmentCode` in the response.

Latest OTP for test/development without sending a phone number:

```powershell
Invoke-RestMethod http://localhost:8080/api/v1/test/latest-otp
```

## Run frontend

```powershell
cd C:\projects\amlak\frontend
flutter run -d chrome --dart-define=API_BASE_URL=http://localhost:8080/api/v1
```

## Local infrastructure

```powershell
cd C:\projects\amlak
docker compose up -d
```

The current primary backend adapter is in-memory for fast development. Platform admin accounts, city catalog locations, system locations, and location suggestions are persisted in MongoDB when `MONGO_URI` is reachable; otherwise that platform catalog falls back to memory. User/session reads are wrapped by a 60-second write-through identity cache backed by Redis when `REDIS_ADDR` is reachable, otherwise by an in-process TTL cache. User/session changes publish RabbitMQ domain events when `RABBITMQ_URL` is available. MongoDB, Redis, MinIO, and RabbitMQ are wired in Docker and ready for the remaining persistence adapters.
