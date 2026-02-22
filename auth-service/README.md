# Auth Service

Authentication and authorization microservice for the Ad System platform. Built with Go using Domain-Driven Design (DDD) and Hexagonal (Ports & Adapters) architecture.

---

## Table of Contents

- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Project Structure](#project-structure)
- [Database Schema](#database-schema)
- [Features](#features)
- [API Endpoints](#api-endpoints)
- [Environment Variables](#environment-variables)
- [Running the Service](#running-the-service)
- [Migrations](#migrations)
- [Error Reference](#error-reference)
- [Architecture Decisions](#architecture-decisions)
- [Planned Features](#planned-features)

---

## Architecture

This service follows **Hexagonal Architecture** (also called Ports & Adapters), combined with **Domain-Driven Design** principles.

```
HTTP Request
    ↓
Adapter (Inbound) — HTTP handler, parses request, calls port
    ↓
Port (Inbound) — AuthService interface, the contract
    ↓
Application — Business logic, orchestrates domain + outbound ports
    ↓
Port (Outbound) — UserRepository / TokenMaker / OAuthProvider interfaces
    ↓
Adapter (Outbound) — PostgreSQL / JWT / Google OAuth implementations
```

**Key principle:** The application layer (business logic) knows nothing about HTTP, databases, or JWT libraries. It only speaks the language of domain types and port interfaces. This means you can swap PostgreSQL for MongoDB, or JWT for Paseto, without touching a single line of business logic.

---

## Tech Stack

| Concern          | Technology                            |
| ---------------- | ------------------------------------- |
| Language         | Go 1.25.1                             |
| HTTP Router      | `net/http` stdlib mux (no framework)  |
| Database         | PostgreSQL (pgx/v5 + pgxpool)         |
| Password Hashing | Argon2id                              |
| JWT              | golang-jwt/v5 (HMAC SHA-256)          |
| OAuth2           | golang.org/x/oauth2 + Google provider |
| Logging          | `log/slog` (structured)               |
| Migrations       | golang-migrate/migrate                |
| Docs             | Swagger (swaggo/swag)                 |
| Live Reload      | Air                                   |
| Linting          | golangci-lint                         |

---

## Project Structure

```
auth-service/
├── cmd/
│   └── main.go                          # Entry point
├── config/
│   ├── config.go                        # App config loader
│   ├── authConfig.go                    # JWT secrets, token durations
│   └── dbConfig.go                      # Database connection config
├── internal/
│   ├── domain/                          # Core business types — no imports from other layers
│   │   ├── user.go                      # User struct, Role, Status types
│   │   └── errors.go                    # Sentinel domain errors
│   ├── port/
│   │   ├── inbound/
│   │   │   └── auth_service.go          # AuthService interface + DTOs (RegisterInput, LoginResponse, etc.)
│   │   └── outbound/
│   │       ├── user_repository.go       # UserRepository interface
│   │       ├── token_maker.go           # TokenMaker interface
│   │       └── oauth_provider.go        # OAuthProvider interface + GoogleUserInfo
│   ├── application/
│   │   └── auth_service.go              # Business logic implementation
│   ├── adapter/
│   │   ├── inbound/
│   │   │   └── http/
│   │   │       ├── server.go            # Wires all dependencies, starts server
│   │   │       ├── router.go            # Route registration
│   │   │       ├── handler/
│   │   │       │   ├── authHandler.go   # All auth HTTP handlers
│   │   │       │   └── healthHandler.go # Health check
│   │   │       └── middleware/
│   │   │           ├── auth.go          # JWT verification middleware
│   │   │           ├── Logger.go        # Request/response logging
│   │   │           ├── recovery.go      # Panic recovery
│   │   │           ├── RequestId.go     # Request ID injection
│   │   │           └── setUpMiddleware.go
│   │   └── outbound/
│   │       ├── postgres/
│   │       │   ├── db.go                # pgxpool connection
│   │       │   ├── queries.go           # Raw SQL constants
│   │       │   └── user_repository.go   # PostgreSQL implementation of UserRepository
│   │       ├── jwt/
│   │       │   └── token_maker.go       # JWT implementation of TokenMaker
│   │       └── oauth/
│   │           └── google.go            # Google OAuth2 implementation of OAuthProvider
│   └── utils/
│       ├── context.go                   # Extract claims from request context
│       ├── cookies.go                   # Refresh token cookie helpers
│       ├── errorHandler.go              # Maps domain errors → HTTP status codes
│       ├── password.go                  # Argon2id hash + compare
│       ├── readJson.go                  # Decode request body
│       ├── response.go                  # WriteJSON, WriteMessage, WriteError helpers
│       ├── validator.go                 # Struct validation wrapper
│       └── middleware.go                # Middleware utilities
├── migrations/
│   ├── 000001_create_users.up.sql       # Initial users table
│   ├── 000001_create_users.down.sql
│   ├── 000002_add_oauth_and_profile.up.sql  # OAuth columns, nullable password
│   └── 000002_add_oauth_and_profile.down.sql
├── docs/                                # Swagger generated files
├── Makefile
└── go.mod
```

---

## Database Schema

### `users` table (after all migrations)

| Column          | Type             | Nullable | Notes                              |
| --------------- | ---------------- | -------- | ---------------------------------- |
| `id`            | UUID             | NOT NULL | Primary key, gen_random_uuid()     |
| `email`         | VARCHAR(255)     | NOT NULL | Unique                             |
| `name`          | VARCHAR(255)     | NULL     | Added in migration 002             |
| `password_hash` | TEXT             | NULL     | NULL for OAuth users               |
| `provider`      | VARCHAR(50)      | NOT NULL | `'local'` or `'google'`            |
| `provider_id`   | TEXT             | NULL     | Google Sub ID for OAuth users      |
| `avatar_url`    | TEXT             | NULL     | Google profile picture             |
| `role`          | user_role ENUM   | NOT NULL | `advertiser`, `publisher`, `admin` |
| `status`        | user_status ENUM | NOT NULL | `active`, `suspended`              |
| `created_at`    | TIMESTAMPTZ      | NOT NULL |                                    |
| `updated_at`    | TIMESTAMPTZ      | NOT NULL |                                    |

**Indexes:**

- `idx_users_email` — on `email`
- `idx_users_provider_id` — unique on `(provider, provider_id)` where `provider_id IS NOT NULL`

**Why `password_hash` is nullable:** OAuth users (Google) have no password. A NOT NULL constraint would force a fake value. NULL is the honest representation: "this user has no password."

---

## Features

### ✅ Local Authentication

- **Register** — email + password + role. Argon2id hashing. Rejects duplicate emails.
- **Login** — verifies password with `ComparePassword` (salt-aware). Returns access token in body, refresh token in `HttpOnly` cookie.
- **Refresh Token** — validates refresh token from cookie, issues new access token.
- **Get Me** — protected route, returns current user profile from JWT claims email.

### ✅ Google OAuth2

- **Initiate** — generates cryptographically random state, builds Google consent URL, sets 3 short-lived `HttpOnly` cookies (`oauth_state`, `oauth_flow`, `oauth_role`).
- **Callback** — CSRF check (state cookie vs query param), exchanges code for Google user info via userinfo endpoint, handles login vs register flow, prevents email conflicts between local and OAuth accounts.

### ✅ Change Password

- Protected route — requires valid JWT.
- Verifies current password before allowing change (using `ComparePassword`, not string comparison).
- Blocks OAuth users from using this endpoint (`ErrOAuthUser`).
- Detects identical new password via `ComparePassword` before hashing (`ErrSamePassword`).
- Updates `password_hash` + `updated_at` in a single `UPDATE` query.

---

## API Endpoints

| Method | Path                    | Auth            | Description                           |
| ------ | ----------------------- | --------------- | ------------------------------------- |
| `POST` | `/auth/register`        | Public          | Register with email + password + role |
| `POST` | `/auth/login`           | Public          | Login, get tokens                     |
| `POST` | `/auth/refresh-token`   | Public (cookie) | Refresh access token                  |
| `GET`  | `/auth/me`              | Protected       | Get current user                      |
| `POST` | `/auth/change-password` | Protected       | Change password                       |
| `GET`  | `/auth/google`          | Public          | Initiate Google OAuth2                |
| `GET`  | `/auth/google/callback` | Public          | Google OAuth2 callback                |
| `GET`  | `/health`               | Public          | Health check                          |

### Request / Response examples

**Register**

```json
POST /auth/register
{
  "email": "user@example.com",
  "password": "securepassword",
  "name": "John Doe",
  "role": "advertiser"
}
```

**Login**

```json
POST /auth/login
{
  "email": "user@example.com",
  "password": "securepassword"
}
// Response body contains access_token
// Response sets HttpOnly cookie: refresh_token
```

**Change Password**

```json
POST /auth/change-password
Authorization: Bearer <access_token>
{
  "current_password": "oldpassword",
  "new_password": "newpassword"
}
```

**Google OAuth — Login flow**

```
GET /auth/google?flow=login
→ redirects to Google consent screen
→ Google redirects to /auth/google/callback
→ returns tokens
```

**Google OAuth — Register flow**

```
GET /auth/google?flow=register&role=advertiser
→ redirects to Google consent screen
→ Google redirects to /auth/google/callback
→ creates account + returns tokens
```

---

## Environment Variables

```env
# Server
APP_PORT=8080
APP_ENV=development           # or production

# Database
DB_URL=postgres://user:password@localhost:5432/auth_db?sslmode=disable

# JWT
ACCESS_TOKEN_SECRET=your-access-secret
REFRESH_TOKEN_SECRET=your-refresh-secret
ACCESS_TOKEN_DURATION=15m
REFRESH_TOKEN_DURATION=168h   # 7 days

# Google OAuth2
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/auth/google/callback
```

---

## Running the Service

**Prerequisites:** Go 1.21+, PostgreSQL, golang-migrate

```bash
# Install dependencies
make tidy

# Run migrations
make migrate-up

# Run in development (with live reload via Air)
make dev

# Run directly
make run

# Build binary
make build

# Generate Swagger docs
make swagger

# Lint
make lint
```

---

## Migrations

```bash
# Apply all pending migrations
make migrate-up

# Rollback all migrations
make migrate-down

# Create a new migration
make migrate-create name=your_migration_name

# Check current version
make migrate-version

# Force a version (use after failed migration)
make migrate-force version=1
```

---

## Error Reference

| Domain Error            | HTTP Status | Meaning                                  |
| ----------------------- | ----------- | ---------------------------------------- |
| `ErrUserNotFound`       | 404         | No user with that email/provider         |
| `ErrEmailAlreadyExists` | 409         | Email taken (local registration)         |
| `ErrEmailConflict`      | 409         | Email registered with different provider |
| `ErrInvalidCredentials` | 401         | Wrong password or bad OAuth state        |
| `ErrUserSuspended`      | 403         | Account suspended                        |
| `ErrInvalidRole`        | 400         | Role must be advertiser or publisher     |
| `ErrOAuthUser`          | 400         | Action not allowed for OAuth accounts    |
| `ErrSamePassword`       | 400         | New password same as current             |

---

## Architecture Decisions

**Why `*string` for nullable columns?**
Go's `string` zero value is `""`, which is indistinguishable from an empty string in a database NULL column. `*string` can be `nil` (representing NULL) or point to an actual string. The repository layer bridges this gap: on write, `""` → `nil` → NULL; on read, NULL → `nil` → don't dereference, leave the field as `""`.

**Why two JWT secrets?**
Access tokens and refresh tokens have different security profiles. If one secret leaks, the other remains safe. Refresh tokens are long-lived and stored in `HttpOnly` cookies — they need stronger isolation from the short-lived access tokens sent as Bearer headers.

**Why state cookie for OAuth?**
This is CSRF protection. Without verifying that the `state` param in the callback matches what was sent in the initiation request, an attacker could forge a callback request and trick your server into authenticating a malicious code. The cookie is `HttpOnly` and scoped to 5 minutes.

**Why `flow` and `role` cookies instead of URL params in callback?**
Google's callback URL only includes `code` and `state` — you cannot add extra query params and have Google pass them through. Cookies are the correct mechanism to carry your own state across the redirect.

**Why `provider_id` (Sub) not email for OAuth user lookup?**
Google's Sub ID is permanent and unique per user per app. Email can change. If a user changes their Google email and logs in again, Sub still matches — email lookup would create a duplicate account.

---

## Planned Features

- [ ] **Forgot Password** — password reset via email token (requires migration 000003 for `password_reset_tokens` table + email sender adapter)
- [ ] **Email Verification** — verify email on registration before allowing login
- [ ] **Token Invalidation** — `token_version` column approach (increment on password change/logout to invalidate old JWTs without Redis)
- [ ] **Soft Delete** — `deleted_at` column, filter in all queries
- [ ] **Audit Trail** — separate table logging sensitive actions (login, password change, etc.)
- [ ] **Additional OAuth Providers** — GitHub, etc. (OAuthProvider interface already supports multiple providers)
