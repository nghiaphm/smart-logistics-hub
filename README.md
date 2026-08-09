# Smart Logistics Hub

Smart Logistics Hub is a Go-based logistics operations system modeling the warehouse and delivery workflows of Shopee Express (SPX). It serves as a learning vehicle for Go backend development, production-oriented architecture, and microservices readiness.

**Tech stack (current)**: Go 1.26 + Gin v1.12, MariaDB 11, Keycloak (JWT RS256), Prometheus metrics, Docker Compose.

---

## Quick Start

```bash
# 1. Start infrastructure (MariaDB + Keycloak + backend)
docker compose up -d

# 2. Run database migrations (if not using AUTO_MIGRATE)
docker compose exec backend /migrate up

# 3. Verify
curl http://localhost:8000/healthz
curl http://localhost:8000/readyz
curl http://localhost:9090/metrics
```

The backend exposes:
- **API** on `localhost:8000` (`/api/v1/orders`, `/api/v1/drivers`, `/api/v1/inventory`, `/api/v1/tracking-logs`)
- **Liveness** at `/healthz` (no DB check)
- **Readiness** at `/readyz` (DB ping)
- **Prometheus metrics** on `localhost:9090` (separate internal port)

---

## Architecture

**Modular monolith** with domain-first structure under `backend/internal/`:

```
backend/
├── cmd/api/main.go              # Entry point (API + metrics server)
├── cmd/migrate/main.go          # Migration CLI (golang-migrate)
├── migrations/                  # SQL migration files
├── docker-entrypoint.sh         # Container startup (migrate + server)
├── internal/
│   ├── driver/                  # Driver domain (CRUD)
│   ├── order/                   # Order domain (CRUD + RBAC)
│   ├── inventory/               # Inventory domain (CRUD)
│   ├── tracking/                # Tracking domain (CRUD)
│   ├── ai/                      # AI event entity (planned)
│   ├── billing/                 # Billing entity (planned)
│   ├── inbound/                 # Inbound entity (planned)
│   ├── product/                 # Product entity (planned)
│   ├── trip/                    # Trip entity (planned)
│   ├── warehouse/               # Warehouse entity (planned)
│   ├── common/errors/           # Shared APIError sentinels
│   └── infrastructure/
│       ├── config/              # Env-based config
│       ├── database/            # MariaDB connection
│       ├── redis/               # Redis client (optional)
│       ├── keycloak/            # JWT verifier (JWKS)
│       ├── middleware/          # auth, cors, rbac, error, request_id, metrics
│       └── logger/              # slog wrapper
└── test/integration/            # Repository integration tests
```

**Domain pattern** (handler → service → repository → MariaDB):
```
driver/
├── entity/driver.go
├── dto/request.go, response.go
├── handler/handler.go
├── service/service.go, service_test.go
├── repository/mariadb.go
└── routes.go
```

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

Tests apply migrations automatically (down → up) to ensure a clean schema per run.

---

## Database Migrations

Managed by [golang-migrate](https://github.com/golang-migrate/migrate) (`cmd/migrate`):

```bash
# From host (using compose)
docker compose exec backend /migrate up
docker compose exec backend /migrate version

# From local Go
go run ./cmd/migrate up
go run ./cmd/migrate down
go run ./cmd/migrate version
```

**AUTO_MIGRATE** (default `false` in `docker-compose.yml`): set `true` in development overrides (`docker-compose.override.yml`) to auto-run migrations on every container start. In CI, migrations run as an explicit step before `go test`.

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

See `backend/.env.example` for full template. `docker-compose.yml` provides suitable development defaults.

---

## CI

GitHub Actions (`.github/workflows/ci.yml`) runs on push/PR:

1. gofmt check → 2. go vet → 3. migrate up (MariaDB service) → 4. go test → 5. go build

---

## Stub Directories

- **`ai_service/`** — placeholder for future AI-based plate recognition (YOLOv8 + OCR). No implementation yet.
- **`data_pipeline/`** — placeholder for future Spark/ETL + Airflow DAGs. No implementation yet.

---

## Why MariaDB?

MariaDB was chosen over MongoDB for the business domain because the majority of logistics data (orders, inventory, billing, trips) is inherently relational — requiring joins, transactional integrity, and referential constraints between entities like orders and order_items, inventory levels per warehouse, driver-to-trip assignments, and tracking event timelines.

