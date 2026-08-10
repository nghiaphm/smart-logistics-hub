# Smart Logistic Project — Architecture

*Audit date: 2026-08-10. Based on actual source code, imports, and configuration. This audit includes a refactor that standardized dependency injection across all implemented domains, verified auth/RBAC coverage on every private endpoint, completed all 10 domains, and added inter-domain orchestration to the `order`, `trip`, `inbound`, and `billing` services.*

---

## 1. Project Overview

**Goal**: A Smart Logistics Hub simulating warehouse operations and delivery management (modeled after Shopee Express / SPX). Built as a learning vehicle for Go backend development, microservices architecture, and production-oriented system design.

**Current state**: The project is a **modular monolith** in Go. The primary database is MariaDB. The README is **outdated** — it describes the original Python/FastAPI + MongoDB architecture which no longer exists in `backend/`.

**Build status**: `go build ./cmd/api/`, `go vet ./...`, and `go fmt ./...` all pass.

**Key technologies in actual use**:

| Technology | Usage |
|---|---|
| **Go 1.26.4** | Backend language |
| **Gin v1.12** | HTTP framework |
| **MariaDB 11** | Primary relational database (`database/sql` + `go-sql-driver/mysql`) |
| **Redis v9** | Optional cache/session store (`go-redis/v9`) |
| **Keycloak** | Authentication (JWT via RS256) |
| **AWS S3** | Optional object storage (not used by business logic) |
| **Docker Compose** | Local dev infrastructure (MariaDB, Keycloak, backend) |

---

## 2. Repository Structure

```
smart-logistic-project/
├── .env.development              # Dev env vars (MariaDB, Keycloak, Redis disabled)
├── .env.production               # Production env template (CHANGE_ME placeholders, no secrets)
├── .github/workflows/ci.yml      # Go CI: gofmt, vet, migrate, test (with MariaDB), build
├── .gitattributes                # *.go text eol=lf normalization
├── .kilo/                        # Kilo CLI local config (untracked)
├── .vscode/settings.json         # Editor settings
├── docker-compose.yml            # MariaDB + Backend + Keycloak + Postgres
├── diff_output.txt               # Leftover diff artifact (tracked)
├── README.md                     # Go + Gin + MariaDB + Keycloak architecture guide
├── ARCHITECTURE.md               # This document
├── docs/
│   └── refactor-prompts/         # Historical refactor task prompts
│       ├── architecture_audit.md
│       ├── architecture-refactor.md
│       └── refactor-prompt-phase1/2/3.md
├── ai_service/                   # Stub — PLANNED (no implementation)
│   ├── .env, .env.example, main.py, requirements.txt   # all 0 bytes
│   └── model/                    # Empty directory
├── data_pipeline/                # Stub — PLANNED (no implementation)
│   ├── dags/                     # Empty
│   └── spark_jobs/               # Empty
└── backend/
    ├── .env                      # Empty
    ├── .env.example              # Template: MariaDB, Redis, S3, Keycloak, server, metrics, AUTO_MIGRATE
    ├── Dockerfile                # Multi-stage Go build (golang:1.26-alpine → alpine:3.20)
    ├── docker-entrypoint.sh      # Runs /migrate up (unless AUTO_MIGRATE=false) then execs /server
    ├── go.mod / go.sum
    ├── cmd/api/main.go           # Application entry point (API + metrics server)
    ├── cmd/migrate/main.go       # golang-migrate CLI (up / down / version)
    ├── migrations/
    │   ├── 000001_initial_schema.up.sql    # 13 tables (golang-migrate format)
    │   └── 000001_initial_schema.down.sql  # Drop all tables
    ├── docs/                     # Empty
    ├── test/integration/         # Repository integration tests (driver, order, inventory, tracking)
    └── internal/
        ├── common/
        │   └── errors/errors.go  # APIError struct + 6 sentinel errors
        ├── infrastructure/
        │   ├── config/config.go          # Env-based config loading (incl. metrics)
        │   ├── database/mariadb.go       # MariaDB connection + pool
        │   ├── redis/client.go           # Redis client factory
        │   ├── keycloak/verifier.go      # JWTVerifier with JWKS + RSA
        │   ├── middleware/auth.go        # JWT auth middleware (accepts verifier interface)
        │   ├── middleware/cors.go        # CORS middleware
        │   ├── middleware/rbac.go        # RequireRole middleware
        │   ├── middleware/error_handler.go  # Centralized error rendering (c.Error)
        │   ├── middleware/request_id.go  # Request ID middleware + request-scoped logger
        │   ├── middleware/metrics.go     # Prometheus HTTP metrics middleware
        │   └── logger/logger.go         # slog wrapper
        ├── driver/                       # FULLY IMPLEMENTED
        ├── order/                        # FULLY IMPLEMENTED
        ├── inventory/                    # FULLY IMPLEMENTED
        ├── tracking/                     # FULLY IMPLEMENTED
        ├── product/                      # FULLY IMPLEMENTED
        ├── warehouse/                    # FULLY IMPLEMENTED
        ├── trip/                         # FULLY IMPLEMENTED
        ├── inbound/                      # FULLY IMPLEMENTED
        ├── billing/                      # FULLY IMPLEMENTED
        └── ai/                           # FULLY IMPLEMENTED
```

---

## 3. System Components

### 3.1 Backend (`backend/`)

| Field | Detail |
|---|---|
| **Responsibility** | REST API for logistics operations |
| **Language** | Go 1.26.4 |
| **Module** | `my-web-app.com/smart-logistic-hub` |
| **HTTP Framework** | Gin v1.12.0 |
| **Entry Point** | `backend/cmd/api/main.go` — `func main()` |
| **Primary Database** | MariaDB 11 (`database/sql` + `go-sql-driver/mysql v1.10.0`) |
| **Optional Infrastructure** | Redis v9.20.1, AWS S3 (v1.42), Keycloak |
| **Architectural Style** | **Modular monolith** — domain-first, dependency injection |
| **Dockerfile** | Multi-stage (golang:1.24-alpine → alpine:3.21, non-root user, healthcheck) |
| **Compiles?** | **Yes** — `go build`, `go vet`, `go fmt` all pass |

### 3.2 AI Service (`ai_service/`)

| Field | Detail |
|---|---|
| **Current State** | **EMPTY** — `main.py`, `requirements.txt`, `.env.example`, `.env` all 0 bytes |
| **Status** | PLANNED — no implementation |

### 3.3 Data Pipeline (`data_pipeline/`)

| Field | Detail |
|---|---|
| **Current State** | **EMPTY** — `dags/` and `spark_jobs/` contain no files |
| **Status** | PLANNED — no implementation |

---

## 4. Backend Architecture

### 4.1 Entry Point — `backend/cmd/api/main.go`

```
package main
```

The entry point performs bootstrap in order:

1. `config.LoadConfig()` — loads env vars from `.env.{APP_ENV}`, returns `*Config`
2. `logger.New(cfg.Environment)` — creates structured `*slog.Logger`, sets as default
3. `database.Connect(cfg)` — opens MariaDB connection (25 max open, 5 idle, 5min lifetime), pings
4. `infraredis.NewRedisClient(ctx, cfg)` — optional Redis (only if `REDIS_ENABLED`)
5. `keycloak.NewJWTVerifier(cfg)` — creates JWT verifier with JWKS support
6. `middleware.AuthMiddleware(cfg, devSkipAuth, verifier)` — JWT auth middleware
7. `middleware.CORSMiddleware(frontendURL)` — CORS middleware
8. Gin router setup:
   - Global middleware chain: `RequestIDMiddleware()` → `gin.Recovery()` → CORS → `MetricsMiddleware()` → `ErrorHandler()`
   - `/healthz` — public liveness, returns 200 while the process is alive (no DB check)
   - `/readyz` — public readiness, pings MariaDB (503 when DB unreachable)
   - `/api/v1/orders/*` — auth-protected CRUD (admin-only DELETE)
   - `/api/v1/drivers/*` — auth-protected CRUD
   - `/api/v1/inventory/*` — auth-protected CRUD (admin-only DELETE)
   - `/api/v1/tracking-logs/*` — auth-protected CRUD (admin-only DELETE)
   - `/api/v1/products/*` — auth-protected CRUD (admin-only DELETE)
   - `/api/v1/warehouses/*` — auth-protected CRUD (admin-only DELETE)
   - `/api/v1/trips/*` — auth-protected CRUD + assign-driver (admin-only DELETE)
   - `/api/v1/inbounds/*` — auth-protected CRUD (admin-only DELETE)
   - `/api/v1/billing/*` — auth-protected CRUD + lookup by code/order (admin-only DELETE)
   - `/api/v1/ai-events/*` — auth-protected CRUD (admin-only DELETE)
9. Optional internal Prometheus metrics server on a separate port (`METRICS_PORT`, default 9090, controlled by `METRICS_ENABLED`) serving `/metrics`
10. HTTP server with graceful shutdown (SIGINT/SIGTERM, 10s drain); both API and metrics servers are shut down cleanly

**Dependency injection**: All dependencies wired manually in `main()` — no framework used.

### 4.2 Domain Implementation Status

| Domain | Entity | DTO | Handler | Service | Repository | Routes | Auth | DB | Status |
|---|---|---|---|---|---|---|---|---|---|
| **driver** | ✓ | ✓ | ✓ | ✓ | ✓ (MariaDB) | ✓ | AuthMw | MariaDB | **FULL** |
| **order** | ✓ | ✓ | ✓ | ✓ | ✓ (MariaDB) | ✓ | AuthMw + RBAC | MariaDB | **FULL** |
| **inventory** | ✓ | ✓ | ✓ | ✓ | ✓ (MariaDB) | ✓ | AuthMw + RBAC | MariaDB | **FULL** |
| **tracking** | ✓ | ✓ | ✓ | ✓ | ✓ (MariaDB) | ✓ | AuthMw + RBAC | MariaDB | **FULL** |
| **product** | ✓ | ✓ | ✓ | ✓ | ✓ (MariaDB) | ✓ | AuthMw + RBAC | MariaDB | **FULL** |
| **warehouse** | ✓ | ✓ | ✓ | ✓ | ✓ (MariaDB) | ✓ | AuthMw + RBAC | MariaDB | **FULL** |
| **trip** | ✓ | ✓ | ✓ | ✓ | ✓ (MariaDB) | ✓ | AuthMw + RBAC | MariaDB | **FULL** |
| **inbound** | ✓ | ✓ | ✓ | ✓ | ✓ (MariaDB) | ✓ | AuthMw + RBAC | MariaDB | **FULL** |
| **billing** | ✓ | ✓ | ✓ | ✓ | ✓ (MariaDB) | ✓ | AuthMw + RBAC | MariaDB | **FULL** |
| **ai** | ✓ | ✓ | ✓ | ✓ | ✓ (MariaDB) | ✓ | AuthMw + RBAC | MariaDB | **FULL** |

### 4.3 Domain Architecture Patterns

**All implemented domains** (driver, order, inventory, tracking, product, warehouse, trip, inbound, billing, ai) follow the same subdirectory structure:

```
internal/{domain}/
├── entity/{entity}.go
├── dto/request.go, dto/response.go
├── handler/handler.go
├── service/service.go, service_test.go
├── repository/mariadb.go
└── routes.go
```

Pattern (Interface Injection):
```
Handler (Gin context → typed DTO binding → calls service)
  ↓
Service (business logic, validation, use-case orchestration)
  ↓
Repository (interface) — SQL queries via database/sql
  ↓
MariaDB
```

Wiring flow in each domain's `RegisterRoutes`:
```
db → repository.New(db) → service.New(repo) → handler.New(service)
```

The service layer never receives `*sql.DB` directly — it depends only on its consumer-side repository interface (`DriverRepository`, `OrderRepository`, `InventoryRepository`, `TrackingRepository`, `ProductRepository`, `WarehouseRepository`, `TripRepository`, `InboundRepository`, `BillingRepository`, `AIRepository`), enabling mocking in unit tests.

**Inter-domain orchestration**: `order.Service` orchestrates across domains during order creation. It depends on two extra injected consumer-side interfaces — `ProductRepository` (validates each item's product exists via `GetByID`) and `InventoryRepository` (checks and reserves stock via `GetByProductWarehouse`/`Update`). `order/routes.go` wires `db → productRepository.New(db)` and `db → inventoryRepository.New(db)` and passes both into `service.NewService(orderRepo, productRepo, inventoryRepo)`. Cross-domain calls are strictly interface-based (no concrete repository types, no direct `*sql.DB`, no circular imports).

`trip.Service` depends on a `DriverRepository` consumer interface to validate driver existence and `AVAILABLE` status before assignment (`Create` and `AssignDriver`). `inbound.Service` depends on an `InventoryRepository` consumer interface; completing an inbound (`status = COMPLETED`) adds `received_qty` to `available_qty` for each item at the inbound's warehouse. `billing.Service` depends on an `OrderRepository` consumer interface — it verifies the order exists before creating an invoice and rejects a second billing record for an already-`PAID` order. `ai.Service` is self-contained (single `AIRepository`); it records gate/ANPR events and flags low-confidence readings (`confidence_score < 0.7`).

### 4.4 Stub Domains

**None remain.** All 10 business domains are fully implemented.

---

## 5. Infrastructure

### 5.1 Config (`internal/infrastructure/config/`)

- **File**: `config.go` — `package config`
- **Exports**: `type Config struct` (31 fields), `func LoadConfig() *Config`
- **No global variable** — `LoadConfig()` returns `*Config`, passed via DI
- **Loading**: `.env.{APP_ENV}` via `godotenv`, OS env fallback
- **Fields**: ProjectName, Version, Environment, ServerHost/Port, S3 (enabled/region/bucket/endpoint/pathstyle), Redis (enabled/URI/host/port/password/DB/ttl), MariaDB (enabled/host/port/user/password/name/URI), FrontendURL, AIServiceURL, KeycloakServerURL/Realm/ClientID, DevSkipAuth

### 5.2 Database — MariaDB (`internal/infrastructure/database/`)

- **File**: `mariadb.go` — `package database`
- **Exports**: `func Connect(cfg *config.Config) (*sql.DB, error)`, `func Close(db *sql.DB)`
- **Driver**: `github.com/go-sql-driver/mysql` (imported via blank import)
- **Pool**: MaxOpenConns=25, MaxIdleConns=5, ConnMaxLifetime=5min
- **Connection**: Supports `MARIADB_URI` (full DSN) or individual host/port/user/password/name fields
- **Startup validation**: Pings DB on `Connect()`

### 5.3 Redis (`internal/infrastructure/redis/`)

- **File**: `client.go` — `package redis`
- **Exports**: `func NewRedisClient(ctx, cfg) (*redis.Client, error)`
- **Conditional**: Only connects if `REDIS_ENABLED=true`
- **URI or fields**: Supports `REDIS_URI` (parsed) or individual host/port/password/DB
- **Usage**: Client created in `main.go` but **not currently used** by any business logic

### 5.4 Keycloak / Authentication (`internal/infrastructure/keycloak/`)

**Two files exist, one is orphaned:**

| File | Package | Used? | Purpose |
|---|---|---|---|
| `verifier.go` | `keycloak` | **Yes** (imported by main.go) | `JWTVerifier` struct with JWKS fetching, RSA key parsing, RS256 token verification with issuer validation |
| `client.go` | `keycloak` | **No** (orphaned) | Package-level `FetchJWKS()` function — never called |

**JWTVerifier** (`verifier.go`):
- `func NewJWTVerifier(cfg *config.Config) *JWTVerifier`
- `func (v *JWTVerifier) VerifyToken(ctx, tokenString) (jwt.MapClaims, error)`:
  - Parses RS256 JWT
  - Validates issuer: `{KeycloakServerURL}/realms/{KeycloakRealm}`
  - Fetches JWKS from Keycloak `/certs` endpoint
  - Caches JWKS in-memory with `sync.RWMutex`
  - Extracts RSA public key from `x5c` certificate chain
  - Returns parsed claims on success

### 5.5 Middleware (`internal/infrastructure/middleware/`)

| File | Exports | Description |
|---|---|---|
| `auth.go` | `type JWTVerifier interface`, `func AuthMiddleware(cfg, devSkipAuth, verifier) gin.HandlerFunc` | Bearer token extraction, dev mode mock user, JWT verification via interface |
| `cors.go` | `func CORSMiddleware(frontendURL string) gin.HandlerFunc` | CORS for configured frontend origin, common methods/headers, credentials, 12h max age |
| `rbac.go` | `func RequireRole(allowedRoles ...string) gin.HandlerFunc` | Extracts roles from JWT `realm_access.roles` and `resource_access.*.roles` |
| `error_handler.go` | `func ErrorHandler() gin.HandlerFunc` | Global error middleware registered via `r.Use(...)`. Renders errors recorded with `c.Error(err)` as `{ "error": { "code", "message" } }`, maps `*APIError` to its status code, masks unknown errors as generic 500 (logs detail server-side via the request-scoped logger) |
| `request_id.go` | `func RequestIDMiddleware() gin.HandlerFunc`, `func LoggerFromContext(c) *slog.Logger` | First middleware in the chain. Generates/reuses a request ID, sets gin `request_id`, attaches it to a request-scoped `slog` logger, and sets the `X-Request-ID` response header |
| `metrics.go` | `func MetricsMiddleware() gin.HandlerFunc` | Records `http_requests_total{method,path,status}` and `http_request_duration_seconds{method,path}` for every request |

### 5.6 Logger (`internal/infrastructure/logger/`)

- **File**: `logger.go` — `package logger`
- **Exports**: `func New(level string) *slog.Logger`, `func WithRequestID(logger, requestID) *slog.Logger`
- **Implementation**: Go standard library `log/slog` with JSON output to stdout
- **Levels**: debug, info, warn, error

---

## 6. API & Routing

### 6.1 Route Registration

Route registration follows a consistent pattern. Each domain exports a `RegisterRoutes` function called from `main.go`:

```go
order.RegisterRoutes(api, db, authMw)
driver.RegisterRoutes(api, db, authMw)
inventory.RegisterRoutes(api, db, authMw)
tracking.RegisterRoutes(api, db, authMw)
product.RegisterRoutes(api, db, authMw)
warehouse.RegisterRoutes(api, db, authMw)
trip.RegisterRoutes(api, db, authMw)
inbound.RegisterRoutes(api, db, authMw)
billing.RegisterRoutes(api, db, authMw)
ai.RegisterRoutes(api, db, authMw)
```

All ten implemented domains receive the same `authMw` (JWT auth) parameter.

### 6.2 API Endpoints

| Method | Path | Auth | Role | Domain |
|---|---|---|---|---|
| GET | `/healthz` | Public | — | Global (liveness) |
| GET | `/readyz` | Public | — | Global (readiness) |
| GET | `/metrics` | Internal port | — | Global (Prometheus) |
| GET | `/api/v1/orders/health` | Public | — | Order |
| POST | `/api/v1/orders` | JWT | — | Order |
| GET | `/api/v1/orders` | JWT | — | Order |
| GET | `/api/v1/orders/:id` | JWT | — | Order |
| PUT | `/api/v1/orders/:id` | JWT | — | Order |
| DELETE | `/api/v1/orders/:id` | JWT | admin | Order |
| POST | `/api/v1/drivers` | JWT | — | Driver |
| GET | `/api/v1/drivers` | JWT | — | Driver |
| GET | `/api/v1/drivers/:id` | JWT | — | Driver |
| PUT | `/api/v1/drivers/:id` | JWT | — | Driver |
| DELETE | `/api/v1/drivers/:id` | JWT | — | Driver |
| POST | `/api/v1/inventory` | JWT | — | Inventory |
| GET | `/api/v1/inventory` | JWT | — | Inventory |
| GET | `/api/v1/inventory/:id` | JWT | — | Inventory |
| PATCH | `/api/v1/inventory/:id` | JWT | — | Inventory |
| DELETE | `/api/v1/inventory/:id` | JWT | admin | Inventory |
| POST | `/api/v1/tracking-logs` | JWT | — | Tracking |
| GET | `/api/v1/tracking-logs` | JWT | — | Tracking |
| GET | `/api/v1/tracking-logs/order/:order_code` | JWT | — | Tracking |
| GET | `/api/v1/tracking-logs/:id` | JWT | — | Tracking |
| PUT | `/api/v1/tracking-logs/:id` | JWT | — | Tracking |
| DELETE | `/api/v1/tracking-logs/:id` | JWT | admin | Tracking |
| POST | `/api/v1/products` | JWT | — | Product |
| GET | `/api/v1/products` | JWT | — | Product |
| GET | `/api/v1/products/:id` | JWT | — | Product |
| PATCH | `/api/v1/products/:id` | JWT | — | Product |
| DELETE | `/api/v1/products/:id` | JWT | admin | Product |
| POST | `/api/v1/warehouses` | JWT | — | Warehouse |
| GET | `/api/v1/warehouses` | JWT | — | Warehouse |
| GET | `/api/v1/warehouses/:id` | JWT | — | Warehouse |
| PATCH | `/api/v1/warehouses/:id` | JWT | — | Warehouse |
| DELETE | `/api/v1/warehouses/:id` | JWT | admin | Warehouse |
| POST | `/api/v1/trips` | JWT | — | Trip |
| GET | `/api/v1/trips` | JWT | — | Trip |
| GET | `/api/v1/trips/:id` | JWT | — | Trip |
| PATCH | `/api/v1/trips/:id` | JWT | — | Trip |
| POST | `/api/v1/trips/:id/assign-driver` | JWT | — | Trip |
| DELETE | `/api/v1/trips/:id` | JWT | admin | Trip |
| POST | `/api/v1/inbounds` | JWT | — | Inbound |
| GET | `/api/v1/inbounds` | JWT | — | Inbound |
| GET | `/api/v1/inbounds/:id` | JWT | — | Inbound |
| PATCH | `/api/v1/inbounds/:id` | JWT | — | Inbound |
| DELETE | `/api/v1/inbounds/:id` | JWT | admin | Inbound |
| POST | `/api/v1/billing` | JWT | — | Billing |
| GET | `/api/v1/billing` | JWT | — | Billing |
| GET | `/api/v1/billing/code/:billing_code` | JWT | — | Billing |
| GET | `/api/v1/billing/order/:order_code` | JWT | — | Billing |
| GET | `/api/v1/billing/:id` | JWT | — | Billing |
| PATCH | `/api/v1/billing/:id` | JWT | — | Billing |
| DELETE | `/api/v1/billing/:id` | JWT | admin | Billing |
| POST | `/api/v1/ai-events` | JWT | — | AI |
| GET | `/api/v1/ai-events` | JWT | — | AI |
| GET | `/api/v1/ai-events/:id` | JWT | — | AI |
| PATCH | `/api/v1/ai-events/:id` | JWT | — | AI |
| DELETE | `/api/v1/ai-events/:id` | JWT | admin | AI |

**Notable**: All CRUD routes across all ten domains require a valid JWT. DELETE routes for `orders`, `inventory`, `tracking-logs`, `products`, `warehouses`, `trips`, `inbounds`, `billing`, and `ai-events` additionally require the `admin` role via `RequireRole("admin")`. `POST /api/v1/orders` requires a `warehouse_id`; on creation the order service validates each item's product and reserves stock (decreases `available_qty`, increases `reserved_qty`) in the inventory domain, rejecting insufficient stock with 409 Conflict. `POST /api/v1/inbounds` requires a `warehouse_id`; completing an inbound adds `received_qty` to inventory `available_qty`. `POST /api/v1/trips` requires a `driver_code` and rejects non-`AVAILABLE` drivers with 409 Conflict. `POST /api/v1/billing` validates the `order_code` exists and rejects duplicate `PAID` invoices with 409 Conflict. `POST /api/v1/ai-events` records ANPR gate events and flags low-confidence readings (`confidence_score < 0.7`). Only `/healthz`, `/readyz`, `/api/v1/orders/health` are public; `/metrics` is served on a separate internal port and is not part of the public API surface.

### 6.3 Inter-Service Communication

| Path | Status |
|---|---|
| Backend → AI Service | **N/A** — config has `AI_SERVICE_URL` but no HTTP calls exist |
| Backend → Data Pipeline | **N/A** — pipeline has no code |
| Backend → Keycloak | **Implemented** — JWKS fetching for JWT verification |
| Backend → Redis | **Connectable** — client created but not used by business logic |
| Backend → S3 | **Disabled** — no code uses S3 |
| Frontend → Backend | **Configured** — CORS allows `FRONTEND_URL` origin |

---

## 7. Database Architecture

### 7.1 MariaDB

**13 tables** defined in `migrations/000001_initial_schema.up.sql`:

| Table | Key Columns | Constraints | Used By |
|---|---|---|---|
| `warehouses` | id, warehouse_code, name, address, lat, lng, contact_phone, manager_name, is_active | UNIQUE(warehouse_code) | warehouse domain |
| `products` | id, sku, name, category, price, weight_gram, length_cm, width_cm, height_cm | UNIQUE(sku) | product domain |
| `drivers` | id, driver_code, full_name, phone, vehicle_type, license_plate, status, current_lat/lng, warehouse_id | UNIQUE(driver_code), INDEX(status) | driver domain |
| `orders` | id, order_code, sender_*, receiver_*, status, assigned_driver_id | UNIQUE(order_code), INDEX(status, driver_id) | order domain |
| `order_items` | id, order_id, product_id, product_name, quantity, weight_gram | FK→orders ON DELETE CASCADE | order domain |
| `inventory` | id, product_id, warehouse_id, available_qty, reserved_qty, damaged_qty, hold_qty | UNIQUE(product_id, warehouse_id), INDEXES | inventory domain |
| `trips` | id, trip_code, driver_id, vehicle_license_plate, status, distance, duration | UNIQUE(trip_code), INDEX(driver_id) | trip domain |
| `trip_stops` | id, trip_id, order_code, stop_type, address, lat, lng, status, timestamps | FK→trips ON DELETE CASCADE | trip domain |
| `tracking_events` | id, order_code, driver_code, status_update, lat, lng, note, timestamp | INDEX(order_code, driver_code, timestamp) | tracking domain |
| `billing` | id, billing_code, order_code, amount, currency, payment_method/status | UNIQUE(billing_code), INDEX(order_code) | billing domain |
| `inbounds` | id, receipt_code, supplier_name, warehouse_id, status, completed_at | UNIQUE(receipt_code), INDEX(warehouse_id) | inbound domain |
| `inbound_items` | id, inbound_id, product_id, expected/received/rejected/qc_passed qty | FK→inbounds ON DELETE CASCADE | inbound domain |
| `ai_events` | id, event_code, license_plate, confidence_score, event_type, gate_id | UNIQUE(event_code) | ai domain |

**All tables**: InnoDB engine, utf8mb4 charset, timestamps with CURRENT_TIMESTAMP defaults.

**Migration tool**: **golang-migrate** (`github.com/golang-migrate/migrate/v4` with the mysql driver and file source). Migration files follow golang-migrate convention (`{version}_{description}.up.sql` / `.down.sql`, e.g. `000001_initial_schema.up.sql`). A dedicated CLI at `backend/cmd/migrate/main.go` supports `up`, `down`, and `version` subcommands, reading the same `MARIADB_*` env vars as the app.

**Migration strategy — explicit CLI, not auto-run:**

| Context | `AUTO_MIGRATE` | Who runs migration | Rationale |
|---|---|---|---|
| **docker-compose.yml** (default) | `false` | Manual: `docker compose exec backend /migrate up` | Production-like safety — no surprise schema changes on restart. The separate `cmd/migrate` CLI was chosen precisely to decouple migration from server startup. |
| **docker-compose.override.yml** (dev convenience) | `true` (override) | `docker-entrypoint.sh` auto-runs on container start | Developer convenience — migration runs automatically so `docker compose up -d` just works. |
| **CI workflow** (`.github/workflows/ci.yml`) | N/A | Explicit `go run ./cmd/migrate up` step before `go test` | CI always runs migrations explicitly via the CLI, independent of `AUTO_MIGRATE`. |
| **`test/integration/setup_test.go`** | N/A | Applies raw SQL directly (down + up) as part of `TestMain` | Self-contained for hermetic tests — does not depend on golang-migrate versioning. |

The `docker-entrypoint.sh` script gates behind `AUTO_MIGRATE={true|false}`. Because golang-migrate uses MySQL advisory locks (`GET_LOCK`), concurrent startup of multiple replicas with `AUTO_MIGRATE=true` would not race — but the explicit `false` default in compose ensures migration intent is always deliberate outside of dev/CI.

### 7.2 MongoDB (Removed)

MongoDB is **no longer the primary database**. The `go.mongodb.org/mongo-driver/v2` remains in `go.mod` as an **indirect** dependency (transitive through another package). No Go code references MongoDB. All repositories use MariaDB. All entity structs in stub domains have been migrated to MariaDB-compatible `db` struct tags.

---

## 8. Dependency Flow

### 8.1 Runtime Dependency Graph

```
main.go
├── config.LoadConfig()
├── logger.New()
├── database.Connect()        → MariaDB
├── redis.NewRedisClient()    → Redis (optional)
├── keycloak.NewJWTVerifier() → Keycloak (JWKS)
├── middleware
│   ├── CORSMiddleware()
│   └── AuthMiddleware()     → JWTVerifier interface
├── domain routes
│   ├── order.RegisterRoutes()
│   │   ├── handler → service → order repository (interface) → MariaDB
│   │   └── order service → product repository (interface) & inventory repository (interface) — validation + stock reservation
│   ├── driver.RegisterRoutes()
│   │   └── handler → service → repository (interface) → MariaDB
│   ├── inventory.RegisterRoutes()
│   │   └── handler → service → repository (interface) → MariaDB
│   ├── tracking.RegisterRoutes()
│   │   └── handler → service → repository (interface) → MariaDB
│   ├── product.RegisterRoutes()
│   │   └── handler → service → repository (interface) → MariaDB
│   ├── warehouse.RegisterRoutes()
│   │   └── handler → service → repository (interface) → MariaDB
│   ├── trip.RegisterRoutes()
│   │   ├── handler → service → trip repository (interface) → MariaDB
│   │   └── trip service → driver repository (interface) — driver availability check on assignment
│   ├── inbound.RegisterRoutes()
│   │   ├── handler → service → inbound repository (interface) → MariaDB
│   │   └── inbound service → inventory repository (interface) — adds received_qty on completion
│   ├── billing.RegisterRoutes()
│   │   ├── handler → service → billing repository (interface) → MariaDB
│   │   └── billing service → order repository (interface) — order existence + duplicate-PAID check
│   └── ai.RegisterRoutes()
│       └── handler → service → ai repository (interface) → MariaDB
└── health/readiness         → MariaDB ping
```

### 8.2 Mermaid Diagram

```mermaid
graph TB
    MAIN[main.go]

    subgraph Infrastructure
        CFG[config]
        LOG[logger/slog]
        DB[database - MariaDB]
        RDS[redis - optional]
        KC[keycloak - JWTVerifier]
        MID[middleware - auth/cors/rbac]
    end

    subgraph "Implemented Domains"
        direction LR
        ORDH[order handler] --> ORDS[order service] --> ORDR[order repository<br/>(interface)] --> MDB[(MariaDB)]
        DRVH[driver handler] --> DRVS[driver service] --> DRVR[driver repository<br/>(interface)] --> MDB
        INVH[inventory handler] --> INVS[inventory service] --> INVR[inventory repository<br/>(interface)] --> MDB
        TRKH[tracking handler] --> TRKS[tracking service] --> TRKR[tracking repository<br/>(interface)] --> MDB
        PRDH[product handler] --> PRDS[product service] --> PRDR[product repository<br/>(interface)] --> MDB
        WHH[warehouse handler] --> WHS[warehouse service] --> WHR[warehouse repository<br/>(interface)] --> MDB
        TRPH[trip handler] --> TRPS[trip service] --> TRPR[trip repository<br/>(interface)] --> MDB
        INBH[inbound handler] --> INBS[inbound service] --> INBR[inbound repository<br/>(interface)] --> MDB
        BILH[billing handler] --> BILS[billing service] --> BILR[billing repository<br/>(interface)] --> MDB
        AIH[ai handler] --> AIS[ai service] --> AIR[ai repository<br/>(interface)] --> MDB

        ORDS -.-> PRDR
        ORDS -.-> INVR
        TRPS -.-> DRVR
        INBS -.-> INVR
        BILS -.-> ORDR
    end

    MAIN --> CFG
    MAIN --> LOG
    MAIN --> DB
    MAIN --> RDS
    MAIN --> KC
    MAIN --> MID
    MAIN --> ORDH
    MAIN --> DRVH
    MAIN --> INVH
    MAIN --> TRKH
    MAIN --> PRDH
    MAIN --> WHH
    MAIN --> TRPH
    MAIN --> INBH
    MAIN --> BILH
    MAIN --> AIH
    KC --> KEYCLOAK[Keycloak Server]

    subgraph External
        KEYCLOAK
    end
```

### 8.3 Concrete Import Analysis

Verified by inspecting actual `import` statements:

| Package | Imported By |
|---|---|
| `internal/infrastructure/config` | main.go, database, redis, keycloak, middleware/auth |
| `internal/infrastructure/database` | main.go |
| `internal/infrastructure/redis` | main.go |
| `internal/infrastructure/keycloak` | main.go (verifier.go only; client.go is orphaned) |
| `internal/infrastructure/middleware` | main.go, order/handler |
| `internal/infrastructure/logger` | main.go |
| `internal/common/errors` | driver/repository, order/repository, inventory/repository, tracking/repository, product/repository, warehouse/repository, trip/repository, inbound/repository, billing/repository, ai/repository |
| `internal/driver/{handler,service,repository,dto,entity}` | driver/routes.go, handler, service, repository (internal to domain) |
| `internal/order/{handler,service,repository,dto,entity}` | order/routes.go, handler, service, repository (internal to domain) |
| `internal/product/{handler,service,repository,dto,entity}` | product/routes.go, handler, service, repository (internal to domain) |
| `internal/warehouse/{handler,service,repository,dto,entity}` | warehouse/routes.go, handler, service, repository (internal to domain) |
| `internal/trip/{handler,service,repository,dto,entity}` | trip/routes.go, handler, service, repository (internal to domain) |
| `internal/inbound/{handler,service,repository,dto,entity}` | inbound/routes.go, handler, service, repository (internal to domain) |
| `internal/billing/{handler,service,repository,dto,entity}` | billing/routes.go, handler, service, repository (internal to domain) |
| `internal/ai/{handler,service,repository,dto,entity}` | ai/routes.go, handler, service, repository (internal to domain) |
| `internal/inventory/entity`, `internal/product/entity` | order/service (consumer-side interface types for stock reservation) |
| `internal/inventory/repository` | order/routes.go (injected into order service) |
| `internal/product/repository` | order/routes.go (injected into order service) |
| `internal/driver/entity` | trip/service (consumer-side interface type for driver assignment) |
| `internal/driver/repository` | trip/routes.go (injected into trip service) |
| `internal/inventory/entity` | inbound/service (consumer-side interface type for stock top-up) |
| `internal/inventory/repository` | inbound/routes.go (injected into inbound service) |
| `internal/order/entity` | billing/service (consumer-side interface type for order validation) |
| `internal/order/repository` | billing/routes.go (injected into billing service) |

**No circular dependencies detected.** `order` depends on `product`/`inventory`, `trip` depends on `driver`, `inbound` depends on `inventory`, and `billing` depends on `order` — always via consumer-side interfaces and leaf `entity` packages only; the reverse directions are absent.

---

## 9. Configuration

### 9.1 Environment Files

| File | Status | Contents |
|---|---|---|
| `backend/.env.example` | Template | MariaDB, Redis, S3, Keycloak, server settings. No real secrets. |
| `backend/.env` | Empty | — |
| `.env.development` (root) | In use | MariaDB localhost, Redis disabled, Keycloak, DEV_SKIP_AUTH=true |
| `.env.production` (root) | Empty | — |

### 9.2 Configuration Loading

`config.LoadConfig()` in `internal/infrastructure/config/config.go`:
1. Reads `APP_ENV` (default: `development`)
2. Loads `.env.{APP_ENV}` via `godotenv`
3. Falls back to OS environment variables
4. Returns `*Config` — **no global state**

### 9.3 Secrets Handling

- `.env.development` uses localhost credentials only (no secrets exposed)
- `.env.example` uses placeholder values (`root`, `localhost`)
- `docker-compose.yml` credentials are development-only defaults
- No real cloud credentials found in the repository

---

## 10. Error Handling

- **Shared errors**: `internal/common/errors/errors.go` (`package errors`)
  - `APIError{StatusCode, Message}` with `Error()` method
  - Sentinel errors: `ErrBadRequest` (400), `ErrUnauthorized` (401), `ErrForbidden` (403), `ErrNotFound` (404), `ErrConflict` (409), `ErrInternal` (500)
- **Repository errors**: Return `apierrors.ErrNotFound` for `sql.ErrNoRows`, wrap DB errors with `fmt.Errorf`
- **Service errors**: Propagate repository errors, add business validation
- **Centralized error middleware**: `internal/infrastructure/middleware/error_handler.go` registers a global `ErrorHandler()` middleware (`r.Use(...)`). Handlers record errors via `c.Error(err)` and return; the middleware:
  - Extracts `*APIError` from the error chain via `errors.As` and renders `{ "error": { "code": <status>, "message": "<msg>" } }`
  - Falls back to HTTP 500 with a generic `"Internal server error"` message for unknown errors — the underlying error is only logged server-side via `slog`, never leaked to the client
  - Binding/validation errors are wrapped with `ErrBadRequest` by handlers before being passed to `c.Error`
- **No per-handler `resolveError()` helper remains**: all ten implemented domains (driver, order, inventory, tracking, product, warehouse, trip, inbound, billing, ai) use `c.Error(err)` + the global middleware

---

## 11. Observability

| Capability | Status |
|---|---|
| Structured logging | **Yes** — `log/slog` with JSON output, configurable levels |
| Request ID | **Yes** — `RequestIDMiddleware()` (first in the chain) generates/reuses a request ID, stores it in gin context, attaches it to a request-scoped slog logger (`request_id` field), and echoes `X-Request-ID` in the response header |
| Health check (liveness) | **Yes** — `GET /healthz` returns 200 while the process is alive; no DB dependency |
| Readiness check | **Yes** — `GET /readyz` pings MariaDB, returns 503 when DB is unreachable |
| Tracing | **No** |
| Metrics | **Yes** — Prometheus (`prometheus/client_golang`), exposed on a **separate internal port** (`METRICS_PORT`, default 9090) via `promhttp.Handler()`, not on the public API port |
| Prometheus / Grafana | Prometheus client library only; no Grafana setup |
| Error monitoring | **No** |

**HTTP metrics** (`internal/infrastructure/middleware/metrics.go`, applied globally):
- `http_requests_total{method,path,status}` — request counter (5xx errors are derivable from the `status` label)
- `http_request_duration_seconds{method,path}` — request latency histogram

The metrics server is started alongside the API server in `cmd/api/main.go` and shuts down gracefully with it. `METRICS_ENABLED` (default `true`) controls whether it starts.

---

## 12. Testing

| Artifact | Status |
|---|---|
| Service unit tests | **Implemented** — `driver/service/service_test.go`, `order/service/service_test.go`, `inventory/service_test.go`, `tracking/service_test.go`, `product/service/service_test.go`, `warehouse/service/service_test.go`, `trip/service/service_test.go`, `inbound/service/service_test.go`, `billing/service/service_test.go`, `ai/service/service_test.go` |
| Middleware tests | **Implemented** — `infrastructure/middleware/error_handler_test.go`, `infrastructure/middleware/request_id_test.go` |
| Repository integration tests | **Implemented** — `test/integration/` (driver, order, inventory, tracking) |
| Migration CLI | `cmd/migrate` — verified locally (`up`/`down`/`version` round-trip against real MariaDB) |
| E2E tests | **None** |
| Test database | **MariaDB** — real database via `MARIADB_*` env vars (CI provides a MariaDB service container) |

**Service unit tests** use mocked repositories (consumer-side interfaces `DriverRepository`, `OrderRepository`, `InventoryRepository`, `TrackingRepository`, `ProductRepository`, `WarehouseRepository`, `TripRepository`, `InboundRepository`, `BillingRepository`, `AIRepository` defined in each service package). Mock repositories implement these interfaces and are injected via `NewServiceWithRepo(...)` — no real DB is touched. `order`, `trip`, `inbound`, and `billing` tests also mock the injected cross-domain interfaces (`ProductRepository`, `InventoryRepository`, `DriverRepository`, `OrderRepository`).

**Repository integration tests** in `test/integration/` run against a real MariaDB. `TestMain` connects using `MARIADB_HOST/PORT/USER/PASSWORD/DB_NAME` env vars (defaults: localhost/3306/root/root/smart_logistics), resets the schema via the down/up migration files, and each test truncates tables for isolation. They verify CRUD for the driver, order, inventory, and tracking repositories plus `sql.ErrNoRows` → `ErrNotFound` behavior.

**CI**: `.github/workflows/ci.yml` runs `go run ./cmd/migrate up` (applying the schema) then `go test ./...` with a MariaDB 11 service container, so integration tests execute in CI against the migrated schema.

---

## 13. Deployment

### 13.1 Docker

| Component | File | Status |
|---|---|---|
| Backend | `backend/Dockerfile` | **Implemented** — multi-stage (golang:1.26-alpine builder → alpine:3.20 runtime), builds `/server` (cmd/api) + `/migrate` (cmd/migrate), non-root user, `HEALTHCHECK` against `/healthz` |
| Backend entrypoint | `backend/docker-entrypoint.sh` | Runs `/migrate up` before starting `/server` (skippable via `AUTO_MIGRATE=false`) |
| AI Service | None | — |
| Data Pipeline | None | — |

### 13.2 Docker Compose

`docker-compose.yml` defines 4 services + 1 network:

1. **mariadb** — MariaDB 11, host port 3307 → container 3306, healthcheck enabled
2. **backend** — Go app, builds from `./backend/Dockerfile`, ports 8000 (API) + 9090 (metrics), `AUTO_MIGRATE=true`, depends on healthy mariadb + started keycloak
3. **postgres-keycloak** — PostgreSQL 16 for Keycloak
4. **keycloak** — Keycloak latest (dev mode), host port 8180 → 8080

### 13.3 CI/CD (`.github/workflows/ci.yml`)

- **Triggers**: Push/PR to `main`/`master`
- **Working directory**: `backend/`
- **Service**: MariaDB 11 container
- **Steps**:
  1. Checkout
  2. Set up Go 1.24
  3. `gofmt` format check
  4. `go vet ./...`
  5. `go run ./cmd/migrate up` (applies migrations before tests)
  6. `go test ./...` (with MariaDB env vars)
  7. `go build ./cmd/api/`

### 13.4 Independent Deployability

| Component | Independently Deployable? |
|---|---|
| Backend | **Yes** — Dockerfile + docker-compose |
| AI Service | **No** — no code |
| Data Pipeline | **No** — no code |
| Keycloak | **Yes** — via Docker Compose |
| MariaDB | **Yes** — via Docker Compose |

---

## 14. Orphaned / Dead Code

**Cleared in Phase 3.** The previously orphaned files (`pkg/utils/logger.go`, `keycloak/client.go`, `auth/handler/auth_handler.go`) have been removed after confirming zero imports via grep across the entire codebase. The empty `pkg/` and `auth/` directories were also removed.

---

## 15. Architectural Observations

### 15.1 Structural Inconsistencies

1. ~~**Dual file structure styles**~~ — **Resolved.** All ten implemented domains now use the same subdirectory layout (entity/, dto/, handler/, service/, repository/ + routes.go at root).

2. ~~**Unused bson tags**~~ — **Resolved.** All stub domains have been migrated from `bson` struct tags to MariaDB-compatible `db` struct tags matching the migration schema. No stub domains remain — all 10 domains are fully implemented.

3. ~~**Orphaned packages**~~ — **Resolved.** `pkg/utils/`, `infrastructure/keycloak/client.go`, and `auth/` have been removed.

4. ~~**Service constructor asymmetry**~~ — **Resolved.** All ten implemented domains (`driver`, `order`, `inventory`, `tracking`, `product`, `warehouse`, `trip`, `inbound`, `billing`, `ai`) now strictly follow the `NewService(repo)` pattern via **Interface Injection**. Each domain's `RegisterRoutes` wires `db → repository.New(db) → service.New(repo) → handler.New(service)`, and no service package depends on `*sql.DB` directly. `NewServiceWithRepo` remains only as an alias used by unit tests for repository mocking.

5. ~~**Inconsistent auth/RBAC coverage**~~ — **Resolved.** All CRUD endpoints in `inventory` and `tracking` are protected by the JWT `authMw` (previously documented as `None`), and their DELETE routes additionally require `RequireRole("admin")` — identical to `order` and `driver`. `cmd/api/main.go` passes the same `authMw` to all ten `RegisterRoutes` calls.

6. ~~**No inter-domain communication**~~ — **Resolved.** Inter-domain orchestration is now implemented via injected consumer-side interfaces: `order.Service` validates each item's product and reserves stock (decrease `available_qty`, increase `reserved_qty`) on creation, rejecting insufficient stock with `apierrors.ErrConflict` and unknown products with `apierrors.ErrBadRequest`; `trip.Service` validates driver existence and `AVAILABLE` status before assignment (unknown → `ErrBadRequest`, unavailable → `ErrConflict`); `inbound.Service` adds `received_qty` to inventory `available_qty` when an inbound transitions to `COMPLETED`; `billing.Service` validates the order exists before invoicing and rejects a second billing record for an already-`PAID` order. All cross-domain dependencies are interface-based and wired in each domain's `routes.go` — no circular imports. **Known limitation**: repositories own separate `*sql.DB` connections, so cross-domain mutations (e.g., order creation + inventory reservation, inbound completion + stock top-up) are not wrapped in a single DB transaction; a shared transaction boundary would be needed for atomic rollback.

7. ~~**Stub domains**~~ — **Resolved.** The final stub domains (`billing`, `ai`) are fully implemented. `billing.Service` provides invoicing with `PENDING`/`PAID`/`FAILED` status flow, order validation, and duplicate-`PAID` prevention. `ai.Service` records ANPR gate events and flags low-confidence readings (`confidence_score < 0.7`). All 10/10 domains are complete.

### 15.2 Dependency Direction

Correct direction maintained:
```
Handler → Service → Repository (Interface) → MariaDB
order.Service → product Repository (Interface) & inventory Repository (Interface)
trip.Service → driver Repository (Interface)
inbound.Service → inventory Repository (Interface)
billing.Service → order Repository (Interface)
Infrastructure (config, database, redis, keycloak, middleware)
```
No circular dependencies. Infrastructure packages do not import domain packages.

### 15.3 Microservice Readiness

The modular monolith structure supports future extraction:
- Each domain has clear boundaries (entity, service, repository)
- No cross-domain database access in current code
- Cross-domain calls are interface-based (`order.Service` → `product`/`inventory`, `trip.Service` → `driver`, `inbound.Service` → `inventory`, `billing.Service` → `order` consumer interfaces), so domains can be extracted into separate services with database per service
- Domains could be extracted into separate services with database per service

---

## 16. Current Architecture Diagram

```mermaid
graph TB
    subgraph "External"
        CLIENT[HTTP Client]
        KC[Keycloak :8180]
    end

    subgraph "Backend :8000"
        GIN[Gin Router]
        CORS[CORS Middleware]
        AUTH[Auth Middleware]
        HEALTH["/health<br/>/readiness"]

        subgraph "Domains"
            ORD[Order<br/>auth+RBACon DELETE]
            DRV[Driver<br/>auth]
            INV[Inventory<br/>auth+RBACon DELETE]
            TRK[Tracking<br/>auth+RBACon DELETE]
            PRD[Product<br/>auth+RBACon DELETE]
            WH[Warehouse<br/>auth+RBACon DELETE]
            TRP[Trip<br/>auth+RBACon DELETE]
            INB[Inbound<br/>auth+RBACon DELETE]
            BIL[Billing<br/>auth+RBACon DELETE]
            AIE[AI Events<br/>auth+RBACon DELETE]
        end
    end

    subgraph "Infrastructure"
        MDB[(MariaDB :3306)]
        REDIS[(Redis<br/>optional)]
        PG[(PostgreSQL<br/>Keycloak DB)]
    end

    CLIENT --> GIN
    GIN --> CORS
    GIN --> AUTH
    AUTH --> KC
    GIN --> HEALTH
    GIN --> ORD
    GIN --> DRV
    GIN --> INV
    GIN --> TRK
    GIN --> PRD
    GIN --> WH
    GIN --> TRP
    GIN --> INB
    GIN --> BIL
    GIN --> AIE
    ORD --> MDB
    DRV --> MDB
    INV --> MDB
    TRK --> MDB
    PRD --> MDB
    WH --> MDB
    TRP --> MDB
    INB --> MDB
    BIL --> MDB
    AIE --> MDB
    KC --> PG
```

---

## 17. Summary

| Metric | Value |
|---|---|
| Total Go files | 103 (86 source + 17 test) |
| Fully implemented domains | 10 (driver, order, inventory, tracking, product, warehouse, trip, inbound, billing, ai) — all consistent subdirectory structure |
| Stub domains (entity/DTO only) | 0 |
| Broken stub (1 line) | 0 (auth/ removed) |
| MariaDB tables | 13 (in migration via golang-migrate) |
| API endpoints | 56 (plus /healthz, /readyz; /metrics on internal port) |
| External services used | Keycloak (active), Redis (optional), S3 (optional, disabled) |
| Orphaned files | 0 (all confirmed dead code removed in Phase 3) |
| Empty stubs (`ai_service/`) | 4 files (PLANNED) |
| Empty stubs (`data_pipeline/`) | 2 directories (PLANNED) |
| Build status | `go build`, `go vet`, `go fmt` — all **pass** |
| Test status | **Implemented** — 10 service unit test suites (mocked repos), 2 middleware test suites (error handler + request ID), 12 repository integration tests (real MariaDB) |

*Audit completed: 2026-08-10. Based on actual source code verification. Source code was modified as part of this audit to standardize dependency injection across all implemented domains, verify auth/RBAC on all private routes, complete all 10 domains (including `product`, `warehouse`, `trip`, `inbound`, `billing`, and `ai`), and add inter-domain orchestration (`order` → `product`/`inventory`, `trip` → `driver`, `inbound` → `inventory`, `billing` → `order`).*
