# Smart Logistics Hub

Smart Logistics Hub is a logistics operations system modeling the warehouse and delivery workflows of Shopee Express (SPX). It is built as a **Go modular monolith** — domain-first, dependency-injected — and serves as a learning vehicle for Go backend development, production-oriented architecture, and microservices readiness.

---

## Tech Stack

| Layer | Technology |
|---|---|
| **Language** | Go 1.26 |
| **HTTP Framework** | Gin v1.12 |
| **Database** | MariaDB 11 (`database/sql` + `go-sql-driver/mysql`) |
| **Auth** | Keycloak (JWT RS256, JWKS verification) |
| **Cache (optional)** | Redis v9 |
| **Metrics** | Prometheus (`prometheus/client_golang`) on a separate internal port |
| **Migration** | golang-migrate CLI (`cmd/migrate`) |
| **Orchestration** | Docker Compose (MariaDB + Keycloak + Backend) |

---

## Project Structure

```
backend/
├── cmd/api/main.go              # Entry point (API + metrics server)
├── cmd/migrate/main.go          # Migration CLI (golang-migrate)
├── migrations/                  # SQL migration files (13 tables)
├── docker-entrypoint.sh         # Container startup (migrate + server)
└── internal/
    ├── driver/                  # Driver domain — FULLY IMPLEMENTED
    ├── order/                   # Order domain — FULLY IMPLEMENTED
    ├── inventory/               # Inventory domain — FULLY IMPLEMENTED
    ├── tracking/                # Tracking domain — FULLY IMPLEMENTED
    ├── product/                 # Product domain — FULLY IMPLEMENTED
    ├── warehouse/               # Warehouse domain — FULLY IMPLEMENTED
    ├── trip/                    # Trip domain — FULLY IMPLEMENTED
    ├── inbound/                 # Inbound domain — FULLY IMPLEMENTED
    ├── billing/                 # Billing domain — FULLY IMPLEMENTED
    ├── ai/                      # AI events (ANPR) domain — FULLY IMPLEMENTED
    ├── common/errors/           # Shared APIError sentinels
    └── infrastructure/
        ├── config/              # Env-based config
        ├── database/            # MariaDB connection/pool
        ├── redis/               # Redis client (optional)
        ├── keycloak/            # JWT verifier (JWKS)
        ├── middleware/          # auth, cors, rbac, error, request_id, metrics
        └── logger/              # slog wrapper
```

**Domain pattern** (interface injection: handler → service → repository → MariaDB):

```
internal/{domain}/
├── entity/{entity}.go
├── dto/request.go, dto/response.go
├── handler/handler.go
├── service/service.go, service_test.go
├── repository/mariadb.go
└── routes.go
```

Services depend on consumer-side repository **interfaces** (never `*sql.DB` directly), which makes unit testing trivial via mocks. Cross-domain orchestration is also interface-based: `order` → `product`/`inventory` (stock reservation), `trip` → `driver` (availability check), `inbound` → `inventory` (stock top-up on completion), `billing` → `order` (validation).

---

## Dev Environment

### 1. Start infrastructure with Docker

```bash
docker compose up -d
```

This starts MariaDB (host port `3307`), Keycloak (host port `8180`), and optionally the backend. For backend-only local development you can start just the infrastructure:

```bash
docker compose up -d mariadb postgres-keycloak keycloak
```

### 2. Run migrations

```bash
# From the backend directory
cd backend
go run ./cmd/migrate up
```

Migrations are also applied automatically inside the container when `AUTO_MIGRATE=true` (the development compose override).

### 3. Run the server

```bash
go run ./cmd/api/main.go
```

The API listens on `http://localhost:8000` and exposes:

- **API** — `/api/v1/...` (all routes JWT-protected; DELETE requires the `admin` role)
- **Liveness** — `/healthz` (no DB check)
- **Readiness** — `/readyz` (DB ping)
- **Prometheus metrics** — `http://localhost:9090/metrics` (separate internal port)

### 4. Run unit tests

```bash
go test ./internal/...
```

Unit tests use mocked repositories and do not require a database.

---

## API Endpoints Summary

All `/api/v1` endpoints require a valid JWT. `DELETE` endpoints marked `admin` require the `admin` role via `RequireRole`.

| Domain | Base Path | Operations |
|---|---|---|
| **Orders** | `/api/v1/orders` | CRUD (admin DELETE) |
| **Drivers** | `/api/v1/drivers` | CRUD |
| **Inventory** | `/api/v1/inventory` | CRUD (admin DELETE) |
| **Tracking** | `/api/v1/tracking-logs` | CRUD + `GET /order/:order_code` (admin DELETE) |
| **Products** | `/api/v1/products` | CRUD (admin DELETE) |
| **Warehouses** | `/api/v1/warehouses` | CRUD (admin DELETE) |
| **Trips** | `/api/v1/trips` | CRUD + `POST /:id/assign-driver` (admin DELETE) |
| **Inbounds** | `/api/v1/inbounds` | CRUD (admin DELETE) |
| **Billing** | `/api/v1/billing` | CRUD + `GET /code/:billing_code` + `GET /order/:order_code` (admin DELETE) |
| **AI Events** | `/api/v1/ai-events` | CRUD (admin DELETE) |

**Notable behaviors**:

- `POST /api/v1/orders` requires `warehouse_id`; the order service validates each product and **reserves stock** (decreases `available_qty`, increases `reserved_qty`), rejecting insufficient inventory with `409 Conflict`.
- `POST /api/v1/inbounds` requires `warehouse_id`; completing an inbound (`status = COMPLETED`) adds `received_qty` to inventory `available_qty`.
- `POST /api/v1/trips` requires `driver_code`; non-`AVAILABLE` drivers are rejected with `409 Conflict`.
- `POST /api/v1/billing` validates the `order_code` exists and rejects duplicate `PAID` invoices with `409 Conflict`.
- `POST /api/v1/ai-events` records ANPR gate events and flags low-confidence readings (`confidence_score < 0.7`).
- Public endpoints: `/healthz`, `/readyz`, `/api/v1/orders/health`. `/metrics` is served on a separate internal port.

---

## Running Tests

```bash
# Unit tests (mock repos, no DB needed)
go test ./internal/...

# Integration tests (requires running MariaDB)
# Start the compose stack or set MARIADB_* env vars pointing to a test DB:
MARIADB_HOST=localhost MARIADB_PORT=3306 MARIADB_USER=root \
MARIADB_PASSWORD=root MARIADB_DB_NAME=smart_logistics \
go test ./test/integration/...
```

Integration tests apply migrations automatically (down → up) to ensure a clean schema per run.

---

## Database Migrations

Managed by [golang-migrate](https://github.com/golang-migrate/migrate) via `cmd/migrate`:

```bash
# From host (using compose)
docker compose exec backend /migrate up
docker compose exec backend /migrate version

# From local Go (in backend/)
go run ./cmd/migrate up
go run ./cmd/migrate down
go run ./cmd/migrate version
```

`AUTO_MIGRATE` (default `false` in `docker-compose.yml`) is set to `true` in development overrides to auto-run migrations on every container start. In CI, migrations run as an explicit step before `go test`.

---

## Configuration

| Env | Purpose |
|---|---|
| `APP_ENV` | `development` / `production` |
| `SERVER_HOST` / `SERVER_PORT` | API listen address |
| `MARIADB_HOST` / `PORT` / `USER` / `PASSWORD` / `DB_NAME` | Database connection |
| `METRICS_ENABLED` / `HOST` / `PORT` | Prometheus metrics server |
| `AUTO_MIGRATE` | Run migrations on container start |
| `KEYCLOAK_SERVER_URL` / `REALM` / `CLIENT_ID` | Keycloak auth |
| `DEV_SKIP_AUTH` | Skip JWT verification (development only) |
| `REDIS_ENABLED` | Enable Redis caching (optional) |

See `backend/.env.example` for the full template. `docker-compose.yml` provides suitable development defaults.

---

## CI

GitHub Actions (`.github/workflows/ci.yml`) runs on push/PR to `main`/`master`:

1. gofmt check → 2. go vet → 3. migrate up (MariaDB service container) → 4. go test → 5. go build

---

## Stub Directories

- **`ai_service/`** — placeholder for a future AI plate-recognition service (YOLOv8 + OCR). No implementation yet.
- **`data_pipeline/`** — placeholder for future Spark/ETL + Airflow DAGs. No implementation yet.

---

## Why MariaDB?

MariaDB was chosen over MongoDB for the business domain because the majority of logistics data (orders, inventory, billing, trips) is inherently relational — requiring joins, transactional integrity, and referential constraints between entities like orders and order_items, inventory levels per warehouse, driver-to-trip assignments, and tracking event timelines. All 10 business domains use MariaDB-compatible `db` struct tags and `database/sql` repositories.
