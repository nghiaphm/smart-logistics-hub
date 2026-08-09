# Smart Logistic Hub — Production-Oriented Architecture Refactor

You are refactoring the Smart Logistic Hub repository toward a clean, maintainable, production-oriented Go backend architecture.

IMPORTANT:

This is a REAL refactoring task, not a theoretical architecture exercise.

The goal is NOT to maximize the number of folders, abstractions, design patterns, services, or technologies.

The goal is to produce a practical architecture that:

- is understandable by a Junior/Mid-level Go developer
- follows idiomatic Go practices
- has clear domain boundaries
- has clear dependency direction
- is testable
- is maintainable
- can evolve toward microservices later
- is suitable for a real product
- avoids over-engineering
- preserves existing business behavior where possible

==================================================
0. SOURCE OF TRUTH
==================================================

Before modifying anything:

1. Read `ARCHITECTURE.md`.
2. Inspect the actual repository.
3. Verify important findings against the source code.
4. Do not assume that ARCHITECTURE.md is always correct if the code has changed.
5. If ARCHITECTURE.md and the source code disagree, use the actual source code as the final source of truth and report the discrepancy.

The repository currently contains:

- backend/
- ai_service/
- data_pipeline/
- agents/
- infrastructure/deployment configuration
- CI/CD configuration

This refactor primarily targets:

backend/

Do NOT unnecessarily redesign `ai_service/` or `data_pipeline/` if they are currently empty or experimental.

==================================================
1. TARGET ARCHITECTURE
==================================================

The backend should become a modular monolith with clear domain boundaries.

Target concept:

backend/
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── auth/
│   ├── driver/
│   ├── warehouse/
│   ├── product/
│   ├── inventory/
│   ├── inbound/
│   ├── order/
│   ├── trip/
│   ├── tracking/
│   ├── billing/
│   ├── shipment/
│   ├── ai/
│   │
│   ├── common/
│   │
│   └── infrastructure/
│       ├── config/
│       ├── database/
│       ├── redis/
│       ├── keycloak/
│       ├── logger/
│       ├── middleware/
│       └── server/
│
├── migrations/
├── deployments/
├── test/
├── Dockerfile
├── go.mod
└── go.sum

This is a MODULAR MONOLITH.

Do NOT split every domain into a separate microservice.

Do NOT create:
- order-service
- inventory-service
- warehouse-service
- driver-service
etc.

as independent services during this refactor.

The architecture must be designed so that domains can be extracted into microservices in the future without prematurely introducing distributed-system complexity.

==================================================
2. DATABASE DECISION — MARIA DB
==================================================

IMPORTANT:

MariaDB is now the PRIMARY BUSINESS DATABASE.

The current repository uses MongoDB in parts of the existing implementation.

The refactor must migrate the core backend persistence from MongoDB to MariaDB where appropriate.

Do NOT simply rename MongoDB repository files to MariaDB files.

Perform a real persistence-layer migration.

First identify all current MongoDB usage:

- mongo-driver dependencies
- MongoDB clients
- database initialization
- collections
- models/documents
- repository implementations
- queries
- indexes
- ObjectID usage
- aggregation pipelines
- transactions
- references between documents

Then map each business entity to a relational model.

Do not blindly create one table per Go struct.

Design tables based on actual business relationships.

Core relational domains include, where applicable:

- warehouses
- warehouse locations
- products
- inventory
- inventory transactions
- orders
- order items
- drivers
- trips
- shipments
- tracking events
- inbound records
- billing
- payments

Preserve business semantics during migration.

==================================================
3. MARIA DB ACCESS
==================================================

Use a practical Go MariaDB approach.

Prefer the standard `database/sql` ecosystem unless the existing project has a strong reason to use another library.

Choose ONE database access strategy.

Do not introduce multiple SQL libraries.

The database layer must provide:

- connection initialization
- connection pooling
- configuration
- health checking
- context-aware queries
- transaction support
- graceful shutdown
- migration support

Do not place SQL/database implementation inside handlers.

Do not place SQL/database implementation inside domain services.

==================================================
4. DATABASE MIGRATIONS
==================================================

Use the project's migration strategy consistently.

If `golang-migrate` is already selected, use it consistently.

Create proper:

migrations/
├── 000001_....up.sql
├── 000001_....down.sql
├── 000002_....up.sql
└── 000002_....down.sql

Migration files must contain real schema definitions.

Include where justified:

- primary keys
- foreign keys
- unique constraints
- indexes
- NOT NULL constraints
- sensible defaults

Do not add constraints merely for theoretical purity if they conflict with existing business behavior.

Avoid excessive normalization where it makes the business model unnecessarily difficult.

==================================================
5. DOMAIN STRUCTURE
==================================================

Each domain should own its business logic.

Recommended structure:

internal/order/
├── entity/
│   └── order.go
├── dto/
│   ├── request.go
│   └── response.go
├── handler/
│   └── handler.go
├── service/
│   └── service.go
├── repository/
│   ├── repository.go
│   └── mariadb.go
└── routes.go

Apply the same structure ONLY where it is actually useful.

Do not create empty files just to make every domain look identical.

For simple domains, use fewer files.

The important principle is:

domain code stays inside its domain.

Avoid:

internal/
├── handlers/
├── services/
├── repositories/
├── models/
└── schemas/

for business-domain code.

==================================================
6. REMOVE OLD LAYER-FIRST ARCHITECTURE
==================================================

The old global folders currently include concepts such as:

- handlers/
- services/
- repositories/
- models/
- schemas/

These are leftovers from the previous architecture.

Migrate their responsibilities into the appropriate domain.

Examples:

old:

internal/models/order.go

new:

internal/order/entity/order.go

old:

internal/repositories/order_repo.go

new:

internal/order/repository/mariadb.go

old:

internal/schemas/orders/request.go

new:

internal/order/dto/request.go

old:

internal/handlers/order.go

new:

internal/order/handler/handler.go

Do NOT simply move files mechanically.

Update:
- package names
- imports
- dependencies
- interfaces
- constructors
- tests
- route registration

After migration, remove obsolete global domain-layer folders if they are no longer used.

Do not leave duplicate implementations.

==================================================
7. ENTITY VS DTO
==================================================

Separate persistence/domain entities from HTTP DTOs.

Entity:

internal/order/entity/order.go

should represent the domain/persistence model.

DTO:

internal/order/dto/request.go
internal/order/dto/response.go

should represent API input/output.

Do not expose database-specific types directly through HTTP responses.

Do not reuse request DTOs as database entities just to reduce files.

However, do not create unnecessary mappings for trivial values.

Use clear boundaries where they provide real value.

==================================================
8. HANDLER RESPONSIBILITY
==================================================

Handlers should be thin.

Handler responsibilities:

- parse HTTP request
- validate request shape
- extract path/query parameters
- call service
- map service result to HTTP response
- map known application errors to HTTP status codes

Handlers must NOT contain:

- SQL queries
- MongoDB queries
- business rules
- transaction orchestration
- complex domain logic

Example:

HTTP
 ↓
Handler
 ↓
Service
 ↓
Repository
 ↓
MariaDB

==================================================
9. SERVICE RESPONSIBILITY
==================================================

Services own business use cases.

Examples:

OrderService:
- create order
- update order
- cancel order
- assign driver
- complete order

InventoryService:
- reserve inventory
- release inventory
- adjust inventory
- check stock

WarehouseService:
- create warehouse
- update warehouse
- manage locations

Services may coordinate multiple repositories when a business use case requires it.

Do NOT put business logic into repositories.

Do NOT make services simple pass-through wrappers around repository calls if there is no actual use-case logic.

==================================================
10. REPOSITORY RESPONSIBILITY
==================================================

Repositories are persistence adapters.

They should handle:

- SQL queries
- scanning rows
- persistence mapping
- database-specific behavior

They should NOT contain business decisions.

Avoid generic "BaseRepository" abstractions unless there is a concrete, repeated behavior that genuinely benefits from them.

Do NOT create a generic repository such as:

Repository[T]

just for architectural style.

Prefer domain-specific repositories.

Example:

type OrderRepository interface {
    Create(...)
    GetByID(...)
    UpdateStatus(...)
}

Implementation:

type MariaDBRepository struct {
    db *sql.DB
}

The interface should exist only when it provides useful abstraction, such as:

- mocking
- dependency inversion
- multiple implementations
- testability

Do not create interfaces purely because "production code should always use interfaces".

==================================================
11. DEPENDENCY DIRECTION
==================================================

Maintain this general dependency direction:

Handler
   ↓
Service
   ↓
Repository interface
   ↓
MariaDB repository implementation
   ↓
Infrastructure database

Infrastructure must NOT depend on business handlers.

Repositories should not import HTTP handlers.

Services should not import router packages.

Domain packages should not depend on `cmd/`.

Avoid circular dependencies.

Verify imports after refactoring.

Run:

go list ./...

and:

go build ./...

to detect dependency problems.

==================================================
12. DEPENDENCY INJECTION
==================================================

Use explicit dependency injection.

Prefer constructor functions:

NewOrderRepository(...)
NewOrderService(...)
NewOrderHandler(...)

Wire dependencies in:

cmd/api/main.go

The entry point should perform bootstrap only:

- load configuration
- initialize logger
- initialize database
- initialize Redis
- initialize Keycloak
- initialize repositories
- initialize services
- initialize handlers
- register routes
- start server
- handle graceful shutdown

Do not put business logic in main.go.

Do not introduce a dependency injection framework.

Use manual dependency injection.

==================================================
13. CONFIGURATION
==================================================

Centralize application configuration.

Use environment variables.

Maintain:

.env.example

Do NOT commit actual secrets.

Configuration should include, where applicable:

- application environment
- server port
- MariaDB connection
- Redis connection
- Keycloak configuration
- logging configuration

Do not let business domains directly read environment variables.

Instead:

config
 ↓
application bootstrap
 ↓
dependencies

Avoid global mutable configuration such as:

AppConfig.SomeValue

when dependency injection can reasonably be used.

==================================================
14. DATABASE CONFIGURATION
==================================================

Create a clear database infrastructure package.

Example:

internal/infrastructure/database/

Responsibilities:

- MariaDB connection
- pool configuration
- ping/health
- close/shutdown
- transaction support

Do not put repositories in this package.

Correct:

infrastructure/database
    ↓
domain repositories

Incorrect:

infrastructure/database
    └── order repository

==================================================
15. REDIS
==================================================

Redis is infrastructure.

Keep Redis initialization and client configuration under:

internal/infrastructure/redis/

Do not duplicate Redis initialization elsewhere.

Do not create a second Redis client package.

Domain services may depend on a Redis abstraction when there is a real business use case.

Do not make every domain depend on Redis.

==================================================
16. KEYCLOAK / AUTHENTICATION
==================================================

Authentication is a high-priority security area.

Inspect the current Keycloak/JWKS implementation.

JWT tokens MUST be cryptographically verified.

Do NOT use ParseUnverified for authentication.

The authentication flow should conceptually be:

HTTP Request
 ↓
Bearer Token
 ↓
JWT parsing
 ↓
Signature verification
 ↓
JWKS public key
 ↓
Issuer validation
 ↓
Audience validation where applicable
 ↓
Expiry validation
 ↓
Claims extraction
 ↓
Authorization/RBAC

JWKS retrieval should:

- use Keycloak configuration
- have timeout
- cache keys
- support key rotation
- handle failures safely

Do not store Keycloak JWKS logic inside generic config code.

Keep Keycloak-specific integration under:

internal/infrastructure/keycloak/

Authentication middleware should use the Keycloak/JWT infrastructure.

==================================================
17. AUTHORIZATION
==================================================

Separate authentication from authorization.

Authentication answers:

"Who is this user?"

Authorization answers:

"Is this user allowed to perform this action?"

Keep authorization logic explicit.

Do not scatter role checks randomly throughout handlers.

If the current project does not have a complete authorization model, implement only what is required by existing business behavior.

Do not invent a complicated RBAC/ABAC system.

==================================================
18. ERROR HANDLING
==================================================

Create a consistent application error strategy.

Errors should distinguish:

- validation errors
- authentication errors
- authorization errors
- not found
- conflict
- database errors
- internal errors

Do not expose raw SQL/database errors to clients.

Use error wrapping:

fmt.Errorf("...: %w", err)

Preserve the underlying error.

Map application errors to HTTP responses centrally where practical.

Do not create an enormous global error hierarchy.

Keep domain-specific errors close to the domain when appropriate.

==================================================
19. VALIDATION
==================================================

Validate HTTP input at the boundary.

Examples:

- required fields
- string length
- numeric ranges
- enum values
- IDs
- pagination parameters

Business validation belongs in services/domain logic.

Example:

HTTP validation:
"name is required"

Business validation:
"warehouse cannot be deleted while inventory exists"

Do not duplicate the same validation in handler and service unless there is a clear reason.

==================================================
20. LOGGING
==================================================

Use structured logging.

Prefer Go's `log/slog` unless there is a concrete reason to retain another logger.

Logs should contain useful context:

- request ID
- operation
- domain
- error
- relevant entity ID

Do NOT log:

- passwords
- access tokens
- refresh tokens
- API keys
- database credentials
- sensitive personal data unnecessarily

Do not scatter random `fmt.Println` debugging throughout production code.

==================================================
21. HTTP SERVER
==================================================

Implement a proper HTTP server lifecycle.

Requirements:

- configurable address/port
- reasonable read/write/idle timeouts
- graceful shutdown
- SIGINT/SIGTERM handling
- context cancellation
- dependency cleanup

The server package should own HTTP server configuration.

main.go should orchestrate startup/shutdown.

==================================================
22. MIDDLEWARE
==================================================

Shared HTTP middleware should live in infrastructure/middleware.

Potential middleware:

- request ID
- logging
- recovery
- CORS
- authentication
- rate limiting

Only implement middleware that is actually required.

Do not add middleware simply to fill a folder.

==================================================
23. API ROUTING
==================================================

Each domain should be responsible for registering its own routes.

For example:

order/routes.go
inventory/routes.go
warehouse/routes.go

Central application router:

- creates router
- applies global middleware
- calls domain route registration

Do not duplicate route registration.

Do not leave obsolete route systems in place.

==================================================
24. RESPONSE FORMAT
==================================================

If the existing API uses a consistent response envelope, preserve it unless there is a strong reason to change it.

Do not introduce a new response format only for architectural cleanliness.

If no consistent format exists, establish a simple production-friendly format.

Do not over-engineer generic response wrappers.

==================================================
25. TESTING
==================================================

Refactor tests along with source code.

Unit tests should remain close to the code:

*_test.go

Use:

test/integration/
test/e2e/

only for tests that genuinely require those scopes.

Prioritize:

1. service/business logic tests
2. repository integration tests
3. HTTP handler tests
4. authentication tests
5. database integration tests

Do not generate hundreds of meaningless tests.

Mocks should only be introduced where useful.

If using mockery, generate mocks only for interfaces that actually need mocking.

==================================================
26. DATABASE TESTING
==================================================

Because MariaDB becomes the primary database, test repository behavior against MariaDB where practical.

Do not rely exclusively on mocks for SQL behavior.

At minimum, test:

- insert
- query
- update
- delete
- constraints
- transactions
- important indexes/query behavior

Prefer an isolated test database/container if the project infrastructure supports it.

Do not use production credentials.

==================================================
27. MARIADB MIGRATION FROM MONGODB
==================================================

This is a critical phase.

Before deleting MongoDB code:

1. inventory all MongoDB collections
2. identify all repository operations
3. identify document relationships
4. identify embedded documents
5. identify ObjectID usage
6. identify indexes
7. identify unique constraints
8. identify transaction requirements
9. identify fields that need normalization
10. map them to MariaDB tables

Create a migration plan internally before making destructive changes.

Do NOT delete MongoDB code before equivalent MariaDB functionality exists.

Do NOT invent new business behavior during migration.

Once all required functionality is migrated and verified:

- remove MongoDB dependencies
- remove MongoDB initialization
- remove obsolete Mongo repositories
- remove obsolete Mongo models
- update configuration
- update documentation

The final backend must not contain unused MongoDB infrastructure.

==================================================
28. DATABASE DESIGN PRINCIPLES
==================================================

Use relational modeling where relationships matter.

For example:

orders
order_items
products

inventory
inventory_transactions

warehouses
warehouse_locations

drivers
trips
shipments
tracking_events

Use foreign keys where they represent real ownership/relationships.

Use unique constraints for actual business uniqueness.

Examples:

order_code
product_code
warehouse_code

Do not add foreign keys blindly if the domain architecture intentionally allows independent lifecycle.

Document any deliberate exceptions.

==================================================
29. TRANSACTIONS
==================================================

Use database transactions for operations that require atomicity.

Example:

Create order
+
reserve inventory
+
create related records

should be transactional when these operations must succeed/fail together.

Do not wrap every database operation in a transaction.

Transaction boundaries belong to business use cases.

The service/use-case layer should control transaction orchestration where multiple repository operations must be atomic.

==================================================
30. PAGINATION
==================================================

Create reusable pagination behavior only where it is genuinely shared.

Pagination should support:

- page/limit or equivalent
- validation
- total count where required
- consistent response metadata

Do not create a giant generic pagination abstraction that leaks into every layer.

==================================================
31. COMMON / SHARED CODE
==================================================

Review:

internal/common/
pkg/

Remove or consolidate unused/duplicate packages.

Keep business logic out of common/shared packages.

Do not move everything into common just because multiple domains use it.

A good shared package should represent truly cross-cutting functionality.

If `pkg/` has no legitimate external consumers, do not force code into it.

==================================================
32. PACKAGE NAMING
==================================================

Use idiomatic Go package names.

Avoid mismatches such as:

directory:
infrastructure/database/

but source:
package repositories

Package name should reflect the package responsibility.

Avoid unnecessary plural package names.

Avoid generic names such as:

utils
helpers
misc

unless there is a clear reason.

==================================================
33. REMOVE DEAD CODE
==================================================

After migration:

Remove:

- unused files
- obsolete MongoDB code
- obsolete global repositories
- obsolete global handlers
- obsolete schemas
- dead routes
- duplicate database initialization
- duplicate Redis clients
- unused dependencies

But only delete code after verifying it is unused.

Use compiler/static analysis to confirm.

==================================================
34. GO TOOLING
==================================================

After each major phase, run:

go fmt ./...
go vet ./...
go test ./...
go build ./...

If available, also run:

go test -race ./...

Do not ignore compiler errors.

Do not hide failures.

Do not leave the repository in a knowingly broken state between phases unless the current migration step explicitly requires it and can be completed immediately.

==================================================
35. CI/CD
==================================================

The current CI configuration may still reference the previous Python backend.

Update it to match the Go backend.

Minimum pipeline:

- gofmt check
- go vet
- go test
- go build

If Docker is implemented:

- Docker build

Do not introduce Kubernetes deployment at this stage.

Do not introduce unnecessary CI complexity.

==================================================
36. DOCKER
==================================================

Create a production-appropriate Dockerfile for the Go backend.

Prefer a multi-stage build:

builder
 ↓
small runtime image

Do not run the application as root if practical.

Do not embed secrets into the image.

Docker Compose should support local development infrastructure.

Do not confuse development Compose with production orchestration.

==================================================
37. HEALTH CHECKS
==================================================

Implement:

/health
/readiness

or equivalent endpoints if they do not already exist.

Health should indicate application liveness.

Readiness should verify required dependencies where appropriate.

Do not make health endpoints perform expensive operations.

==================================================
38. OBSERVABILITY
==================================================

Implement only a practical baseline:

- structured logs
- request ID
- health endpoint
- readiness endpoint

Do not add Prometheus, Grafana, OpenTelemetry, ELK, etc. unless there is already a concrete requirement.

The goal is a clean foundation, not infrastructure theater.

==================================================
39. SECURITY
==================================================

Fix high-priority security problems identified in ARCHITECTURE.md.

Especially:

- JWT signature verification
- JWKS handling
- secret handling
- environment configuration
- authentication bypass behavior
- error exposure
- logging sensitive information

Development bypasses must never silently apply to production.

Validate environment configuration on startup.

Fail fast when required production configuration is missing.

==================================================
40. SECRETS
==================================================

Never commit secrets.

Inspect:

- .env files
- Docker Compose
- GitHub Actions
- config files

Use placeholders in `.env.example`.

If an actual credential was previously committed:

- do NOT print it
- report that credential rotation is required
- remove it from active configuration
- do not assume deleting the current file removes it from Git history

==================================================
41. AI SERVICE
==================================================

Do not force AI logic into the Go backend.

Keep:

ai_service/

as a separate service boundary.

The backend should communicate with it through an explicit API/event contract when implementation begins.

Do not create microservice infrastructure before the AI service actually exists.

==================================================
42. DATA PIPELINE
==================================================

Keep:

data_pipeline/

separate from backend.

Do not import backend business code directly into the data pipeline.

When implemented, prefer explicit data contracts or event/API boundaries.

Do not implement Kafka/Spark/Airflow merely for architectural appearance.

Only introduce them when the corresponding use case exists.

==================================================
43. MICROservice READINESS
==================================================

The backend should be designed as a modular monolith that can evolve toward microservices.

Each domain should have:

- clear business responsibility
- limited dependencies
- repository boundary
- service/use-case boundary
- API boundary where appropriate

Avoid direct cross-domain database access.

Prefer domain service/use-case communication.

Do not extract services yet.

==================================================
44. CROSS-DOMAIN DEPENDENCIES
==================================================

Inspect dependencies such as:

order → inventory
order → driver
shipment → tracking
billing → order

Avoid direct access to another domain's repository.

Prefer explicit service/use-case interfaces where appropriate.

Example:

OrderService
    ↓
InventoryService

rather than:

OrderService
    ↓
InventoryRepository

unless there is a clear architectural reason.

Do not create circular domain dependencies.

If a circular dependency is found, report it and redesign the dependency direction.

==================================================
45. DO NOT OVER-ENGINEER
==================================================

Do NOT introduce:

- dependency injection frameworks
- generic repository frameworks
- generic service frameworks
- CQRS
- event sourcing
- Kafka
- RabbitMQ
- gRPC
- Kubernetes
- service mesh
- distributed tracing
- API gateway
- excessive abstractions

unless the existing project has a concrete requirement for them.

Production-ready does NOT mean technology-heavy.

==================================================
46. REFACTORING ORDER
==================================================

Perform the refactor in this order:

PHASE 0
Security and secrets review.

PHASE 1
Make the Go project structurally consistent.

PHASE 2
Fix package names and imports.

PHASE 3
Complete domain-first migration.

PHASE 4
Remove obsolete global handlers/services/repositories/models/schemas.

PHASE 5
Establish infrastructure packages.

PHASE 6
Implement MariaDB infrastructure.

PHASE 7
Implement MariaDB schema and migrations.

PHASE 8
Migrate repositories from MongoDB to MariaDB.

PHASE 9
Update services and transaction boundaries.

PHASE 10
Fix authentication/JWKS verification.

PHASE 11
Implement consistent errors and validation.

PHASE 12
Implement structured logging and request ID.

PHASE 13
Implement graceful shutdown.

PHASE 14
Update tests.

PHASE 15
Update Docker.

PHASE 16
Update CI/CD for Go.

PHASE 17
Add health/readiness.

PHASE 18
Final architecture cleanup.

After every significant phase:

go fmt ./...
go vet ./...
go test ./...
go build ./...

==================================================
47. PRESERVE BUSINESS BEHAVIOR
==================================================

Do not change business rules merely to make architecture cleaner.

If a business rule is unclear:

- inspect current implementation
- inspect tests
- inspect API contracts
- preserve existing behavior where possible
- report uncertainty

Do not invent requirements.

==================================================
48. API COMPATIBILITY
==================================================

Preserve existing API routes and request/response contracts where practical.

If a route must change because of the refactor:

- document the change
- update route registration
- update tests
- update API documentation if present

Do not silently break APIs.

==================================================
49. FILE CREATION RULE
==================================================

Do not create files merely because the architecture diagram contains them.

Create a file only when there is actual code responsibility for it.

For example, do not create:

repository/interface.go

if there is no interface needed.

Do not create:

service/errors.go

if no domain-specific errors exist.

Architecture should reflect real code, not empty templates.

==================================================
50. FINAL TARGET
==================================================

The final backend should approximately follow:

HTTP
 ↓
Middleware
 ↓
Handler
 ↓
Service / Use Case
 ↓
Repository Interface
 ↓
MariaDB Repository
 ↓
MariaDB

Cross-cutting infrastructure:

Config
Logger
Database
Redis
Keycloak
Middleware
Server

Domain structure:

internal/
├── auth/
├── driver/
├── warehouse/
├── product/
├── inventory/
├── inbound/
├── order/
├── trip/
├── tracking/
├── shipment/
├── billing/
└── ai/

The exact final structure may differ if the actual code demonstrates a better, simpler design.

==================================================
51. FINAL VALIDATION
==================================================

Before finishing the refactor:

1. `go build ./...` passes.
2. `go test ./...` passes, or all remaining failures are explicitly documented.
3. `go vet ./...` passes or remaining issues are documented.
4. No obsolete MongoDB dependency remains unless explicitly justified.
5. MariaDB is the primary backend database.
6. No duplicate database initialization exists.
7. No duplicate Redis initialization exists.
8. No obsolete global handlers/services/repositories/models/schemas remain.
9. Domain dependencies are clear.
10. No circular imports exist.
11. JWT signature verification is enabled.
12. Secrets are not committed.
13. `.env.example` exists and contains no real secrets.
14. Graceful shutdown works.
15. Docker build works if Dockerfile is implemented.
16. CI uses Go instead of the previous Python workflow.
17. Health/readiness endpoints work.
18. Existing business behavior is preserved where possible.
19. No unnecessary frameworks or infrastructure were introduced.

==================================================
52. FINAL REPORT
==================================================

After completing the refactor, provide a concise report containing:

### Changed
- major architecture changes
- domain migrations
- database migration
- infrastructure changes

### Removed
- obsolete folders
- obsolete files
- obsolete dependencies

### Added
- new infrastructure
- migrations
- tests
- configuration

### Database
- MongoDB components removed
- MariaDB schema created
- migrations created
- repositories migrated

### Security
- authentication changes
- JWKS changes
- secret handling

### Validation
- gofmt
- go vet
- go test
- go build
- Docker build if applicable

### Remaining Issues
Only list issues that genuinely remain.

Do NOT claim production-ready if important issues remain.

Use the term:

"Production-oriented foundation"

rather than falsely claiming:

"Production-ready"

unless all relevant requirements are actually satisfied.