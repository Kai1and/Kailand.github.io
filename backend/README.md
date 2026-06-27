# Backend

REST API for equipment booking.

## Run

```bash
cp .env.example .env
go mod tidy
go run ./cmd
```

Before running the app, create the PostgreSQL database and apply `migrations/001_init.sql`.

## API

Public routes:

- `GET /api/health`
- `POST /api/auth/register`
- `POST /api/auth/login`
- `GET /api/categories`
- `GET /api/equipment`
- `GET /api/equipment/{id}`

Authenticated routes:

- `GET /api/auth/me`
- `GET /api/bookings`
- `POST /api/bookings`
- `PATCH /api/bookings/{id}/cancel`

Admin routes:

- `POST /api/categories`
- `PUT /api/categories/{id}`
- `DELETE /api/categories/{id}`
- `POST /api/equipment`
- `PUT /api/equipment/{id}`
- `DELETE /api/equipment/{id}`
- `PATCH /api/bookings/{id}/status`
- `GET /api/users`
- `PATCH /api/users/{id}/role`
