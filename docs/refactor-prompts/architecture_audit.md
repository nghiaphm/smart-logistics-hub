You are performing an architecture audit of the ENTIRE smart-logistic-project repository.

IMPORTANT:
This is an AUDIT AND DOCUMENTATION task only.

DO NOT refactor the project.
DO NOT move files.
DO NOT rename files or folders.
DO NOT delete files.
DO NOT create new source-code files.
DO NOT modify existing source code.
DO NOT change configuration.
DO NOT change dependencies.
DO NOT change Docker configuration.

The ONLY file you are allowed to create or modify is:

ARCHITECTURE.md

Create it at the ROOT of the repository:

smart-logistic-project/ARCHITECTURE.md

==================================================
1. PROJECT CONTEXT
==================================================

The project is a Smart Logistic Hub intended as a serious portfolio/product project.

The project is currently being developed with the goal of:

- learning Go backend development
- learning production-oriented architecture
- learning microservices architecture
- learning how real product systems are structured
- building a project suitable for a Junior Backend Engineer portfolio

The repository currently contains multiple components/services, including at least:

- backend/
- ai_service/
- data_pipeline/
- agents/
- .github/
- Docker / Docker Compose configuration
- environment configuration

The backend was previously named backend-go and has now been renamed to backend.

The backend is being migrated/refactored from a previous FastAPI/Python implementation to Go.

The current backend uses MariaDB as the primary relational database.

The project may also use technologies such as Redis, Keycloak, Docker, message brokers, AI services, data pipelines, etc.

DO NOT assume any technology is actually being used just because it appears in a folder name or dependency.

Verify everything from the actual code/configuration.

==================================================
2. PRIMARY OBJECTIVE
==================================================

Create a highly accurate document describing the CURRENT architecture of the ENTIRE repository.

The document must answer:

1. What components/services currently exist?
2. What is the responsibility of each component?
3. How is the backend currently structured?
4. How is the AI service structured?
5. How is the data pipeline structured?
6. How do these components communicate?
7. What databases and infrastructure are used?
8. What are the current dependency directions?
9. Which parts follow feature/domain-first architecture?
10. Which parts still follow layer-first architecture?
11. Where are duplicate responsibilities?
12. Where are architectural inconsistencies?
13. Which areas are unclear or potentially problematic?
14. Are there missing architectural pieces that should be reviewed?
15. Is the current structure suitable for evolving toward microservices?

IMPORTANT:

Do NOT redesign the project yet.

Do NOT decide what the architecture SHOULD be.

First document what the architecture IS TODAY.

==================================================
3. COMPLETE REPOSITORY TREE
==================================================

Scan the entire repository:

smart-logistic-project/

Include important folders and files recursively.

The tree should cover at least:

- backend/
- ai_service/
- data_pipeline/
- agents/
- .github/
- Dockerfile(s)
- docker-compose.yml
- migration files
- scripts
- tests
- configuration files
- go.mod / go.sum
- Python dependency files
- README
- API documentation
- infrastructure/deployment files

Do NOT include:

- .git/
- __pycache__/
- .pytest_cache/
- .ruff_cache/
- node_modules/
- build artifacts
- binaries
- generated files
- temporary files

Do not expose secrets.

For .env files:
- inspect variable NAMES if necessary
- DO NOT copy secret values into ARCHITECTURE.md
- never expose passwords, API keys, tokens, secrets, private keys, credentials, etc.

==================================================
4. ROOT PROJECT ARCHITECTURE
==================================================

Analyze the root structure.

For each major component, document:

### backend/
Responsibility:
Technology:
Entry point:
Main dependencies:
Database:
External services:
Communication mechanisms:
Current architectural style:

### ai_service/
Responsibility:
Technology:
Entry point:
Dependencies:
External integrations:
Communication with backend:
Current architectural style:

### data_pipeline/
Responsibility:
Technology:
Input:
Output:
Storage:
Communication:
Current architectural style:

### agents/
Determine exactly what this folder contains and how it relates to the project.

Do not assume it is part of the runtime architecture.

### .github/
Inspect CI/CD workflows and document:

- build
- test
- lint
- Docker build
- deployment
- branch/workflow assumptions

==================================================
5. BACKEND ARCHITECTURE
==================================================

Perform a DEEP audit of:

backend/

Do not only inspect the directory tree.

Read the actual Go source code.

Analyze:

backend/
├── cmd/
├── internal/
├── pkg/
├── migrations/
├── deployments/
├── scripts/
├── test/
└── configuration

Identify:

- application entrypoint
- HTTP server
- router
- middleware
- configuration
- database
- Redis
- authentication
- Keycloak
- repositories
- services
- handlers
- entities/models
- DTOs/schemas
- validation
- error handling
- health checks
- logging
- external integrations

==================================================
6. BACKEND DOMAIN ANALYSIS
==================================================

Identify every business domain currently implemented.

For example, if present:

- auth
- ai
- billing
- driver
- inbound
- inventory
- order
- product
- shipment
- tracking
- trip
- warehouse
- etc.

For EACH domain, document:

Domain:
Current path:
Responsibility:
Entity/model:
DTO/schema:
Handler:
Service:
Repository:
Routes:
Validation:
Errors:
External dependencies:

Most importantly, determine whether the domain is self-contained or whether its code is scattered across global folders such as:

internal/models/
internal/services/
internal/repositories/
internal/schemas/
internal/handlers/

If a domain is split across multiple locations, explicitly document this.

==================================================
7. FEATURE-FIRST VS LAYER-FIRST
==================================================

Determine whether the current backend is:

- Layer-first
- Feature-first
- Domain-first
- Clean Architecture
- Hexagonal Architecture
- Hybrid
- Other

Do NOT classify based only on folder names.

Use actual imports and dependencies.

For example, determine whether the project currently has:

internal/
├── handlers/
├── services/
├── repositories/
├── models/
└── schemas/

versus:

internal/
├── warehouse/
│   ├── handler/
│   ├── service/
│   ├── repository/
│   └── entity/
├── inventory/
├── order/
└── ...

If both exist, document that the project is in a transitional/hybrid state.

==================================================
8. DEPENDENCY DIRECTION
==================================================

Inspect actual imports.

Document the real dependency flow.

For example:

HTTP
 ↓
Handler
 ↓
Service
 ↓
Repository
 ↓
MariaDB

But only use this diagram if the code actually follows it.

Identify:

- handler → service
- service → repository
- repository → database
- domain → infrastructure
- infrastructure → domain
- shared package dependencies
- circular dependencies
- concrete dependency coupling
- interface usage
- dependency injection

Create a Mermaid dependency diagram based on the CURRENT code.

==================================================
9. DATABASE ARCHITECTURE
==================================================

The project is expected to use MariaDB as the primary relational database.

Verify:

- database driver
- ORM/query builder
- connection initialization
- connection pooling
- transaction handling
- migrations
- repositories
- indexes
- constraints
- database configuration
- health checks

Identify where database code currently lives.

If there are multiple database-related folders, such as:

internal/db/
internal/infrastructure/database/

document both and determine whether their responsibilities overlap.

Do NOT merge or refactor them.

==================================================
10. REDIS / CACHE
==================================================

If Redis exists:

Determine:

- client library
- initialization
- configuration
- usage
- cache responsibility
- session/token usage
- rate limiting usage
- which domains depend on Redis

Identify whether Redis is infrastructure or directly accessed by business logic.

==================================================
11. KEYCLOAK / AUTHENTICATION
==================================================

If Keycloak exists, deeply inspect:

- configuration
- JWT validation
- JWKS fetching
- token verification
- middleware
- authentication service
- authorization
- roles/permissions

IMPORTANT:

Do not classify a file as "configuration" just because it is currently located inside config/.

Determine responsibility from actual code.

For example, if a file currently located at:

config/security.go

actually performs:

- JWKS fetching
- HTTP calls to Keycloak
- token verification
- caching

document that its CURRENT LOCATION and CURRENT RESPONSIBILITY differ.

==================================================
12. AI SERVICE
==================================================

Deeply inspect:

ai_service/

Determine:

- framework
- language
- entrypoint
- API endpoints
- service layer
- AI model integration
- database usage
- queue/message broker usage
- communication with backend
- communication with data pipeline
- authentication
- configuration
- error handling
- logging
- tests

Determine whether ai_service is:

- a standalone service
- a module inside the monolith
- an independently deployable service
- an experimental component

Use actual code/configuration to determine this.

==================================================
13. DATA PIPELINE
==================================================

Deeply inspect:

data_pipeline/

Determine:

- responsibility
- data sources
- data processing
- output
- storage
- scheduling
- communication with backend
- communication with AI service
- technologies
- deployment model

Determine whether it is:

- batch pipeline
- streaming pipeline
- ETL
- ML/data preprocessing
- scheduled job
- standalone service

==================================================
14. INTER-SERVICE COMMUNICATION
==================================================

Find all communication mechanisms between components.

Look for:

- REST
- HTTP
- gRPC
- RabbitMQ
- Kafka
- Redis Pub/Sub
- Webhooks
- shared database
- file/object storage
- direct imports
- subprocess calls
- message queues

For each communication path, document:

Source → Destination
Protocol:
Purpose:
Data exchanged:
Synchronous/Asynchronous:
Current implementation:

Example:

backend → ai_service
HTTP POST
Purpose: AI prediction
Synchronous

Do NOT invent communication that is not present.

==================================================
15. SHARED CODE / COMMON / PKG
==================================================

Inspect:

backend/internal/common/
backend/pkg/

Determine:

- what is actually shared
- whether shared packages contain business logic
- whether they create coupling between domains
- whether they are appropriate for sharing
- whether packages are internal-only or intended for external reuse

Do not recommend moving anything yet.

==================================================
16. API / ROUTING
==================================================

Inspect all routes.

Document:

- route registration
- API versioning
- domain route grouping
- middleware
- authentication
- authorization
- health endpoints
- external/internal APIs

Determine whether route registration is centralized or domain-specific.

==================================================
17. CONFIGURATION MANAGEMENT
==================================================

Inspect:

- .env.development
- .env.production
- .env.example
- config files
- Docker Compose
- deployment configuration

Document:

- configuration loading
- environment separation
- configuration ownership
- secrets handling
- database config
- Redis config
- Keycloak config
- service URLs
- ports

NEVER expose actual secret values.

==================================================
18. ERROR HANDLING
==================================================

Inspect:

- exceptions
- errors
- HTTP error mapping
- domain errors
- common errors
- middleware

Determine:

- whether errors are centralized
- whether errors are domain-specific
- whether handlers expose internal errors
- whether error codes exist

==================================================
19. OBSERVABILITY
==================================================

Inspect whether the project has:

- structured logging
- request ID
- tracing
- metrics
- health checks
- readiness checks
- liveness checks
- OpenTelemetry
- Prometheus
- Grafana
- error monitoring

Document CURRENT STATE only.

==================================================
20. TESTING ARCHITECTURE
==================================================

Inspect:

- *_test.go
- test/
- integration tests
- E2E tests
- mocks
- fixtures
- test database
- test containers

Document where tests live and what layers are tested.

==================================================
21. DOCKER / DEPLOYMENT
==================================================

Inspect:

- Dockerfile(s)
- docker-compose.yml
- deployments/
- CI/CD

Determine:

- which components are independently containerized
- which databases are containers
- service dependencies
- startup order
- health checks
- environment variables
- build process

Determine whether backend, ai_service and data_pipeline can currently be deployed independently.

==================================================
22. ARCHITECTURAL DUPLICATION
==================================================

Explicitly search for duplicate or overlapping concepts.

Examples:

- tracking/
- tracking_log/
- tracking_logs/
- ai/
- ai_event/
- services/
- handlers/
- models/
- repositories/
- schemas/
- db/
- infrastructure/database/
- config/
- infrastructure/config/

For every suspicious duplication, document:

Location A:
Responsibility:

Location B:
Responsibility:

Potential overlap:
Evidence:

Do NOT delete or merge anything.

==================================================
23. ARCHITECTURAL RISKS
==================================================

Create:

## Architectural Risks

Classify observations as:

### High
Potentially affects correctness, maintainability, scalability or service boundaries.

### Medium
Structural issue that should be reviewed.

### Low
Naming/organization issue with limited impact.

Only report evidence-based observations.

Examples:

- circular dependency
- global mutable state
- business logic in handlers
- services directly accessing DB
- repositories tightly coupled to concrete infrastructure
- duplicate domain representations
- shared package containing business logic
- shared database across supposed microservices
- tightly coupled services
- unclear service boundaries

==================================================
24. MISSING / UNCLEAR COMPONENTS
==================================================

Create:

## Missing or Unclear

Identify things that cannot be determined from the current repository.

Examples:

- unclear ownership of a database
- unclear communication contract
- unclear domain boundary
- unclear authentication responsibility
- missing migration strategy
- missing health/readiness check
- missing integration tests

Do not assume they are missing if they simply could not be located.

Distinguish:

"Not found"

from:

"Exists but responsibility is unclear."

==================================================
25. CURRENT ARCHITECTURE DIAGRAM
==================================================

Create a high-level Mermaid diagram of the ENTIRE CURRENT SYSTEM.

For example:

Client
  |
  v
Backend
  |
  +----> MariaDB
  |
  +----> Redis
  |
  +----> Keycloak
  |
  +----> AI Service
  |
  +----> Data Pipeline

BUT:

This is only an example.

Build the actual diagram from the repository.

Include:

- backend
- ai_service
- data_pipeline
- database
- Redis
- Keycloak
- message broker if present
- object storage if present
- external services if relevant

==================================================
26. DOMAIN / SERVICE BOUNDARY ANALYSIS
==================================================

For each major component, determine:

- What business capability does it own?
- What data does it own?
- What data does it depend on?
- Which other components depend on it?
- Is it independently deployable?
- Does it have its own database?
- Does it directly import another service's code?
- Does it communicate through APIs/events?

This section is important because the project is intended to learn microservices architecture.

However:

Do NOT redesign service boundaries.

Only document CURRENT boundaries.

==================================================
27. PRODUCTION READINESS OBSERVATIONS
==================================================

Create:

## Production Readiness Observations

Evaluate CURRENT implementation only.

Cover:

- architecture
- dependency management
- database
- security
- authentication
- configuration
- error handling
- logging
- observability
- testing
- deployment
- graceful shutdown
- migrations
- scalability
- service boundaries

Do not turn this into a refactor plan.

==================================================
28. IMPORTANT DISTINCTION
==================================================

The entire ARCHITECTURE.md must clearly distinguish:

CURRENT STATE
from
POTENTIAL CONCERN

Do NOT write:

"Move X to Y."

Instead write:

"X is currently located at A and performs responsibility B. This may overlap with C and should be reviewed."

The next step will be a human/AI architecture review based on this document.

==================================================
29. FINAL STRUCTURE OF ARCHITECTURE.md
==================================================

Use this structure:

# Smart Logistic Project — Architecture

## 1. Project Overview

## 2. Repository Structure

## 3. System Components

## 4. Backend Architecture

## 5. Backend Domain Structure

## 6. AI Service Architecture

## 7. Data Pipeline Architecture

## 8. Infrastructure

## 9. Database Architecture

## 10. Authentication & Security

## 11. API & Communication

## 12. Dependency Flow

## 13. Configuration

## 14. Error Handling

## 15. Testing Architecture

## 16. Deployment Architecture

## 17. Observability

## 18. Shared Code

## 19. Duplications & Inconsistencies

## 20. Architectural Risks

## 21. Missing / Unclear Areas

## 22. Current Architecture Diagram

## 23. Domain / Service Boundary Analysis

## 24. Production Readiness Observations

## 25. Summary of Current State

==================================================
30. FINAL CHECK
==================================================

Before finishing:

1. Verify that every major root component was inspected.
2. Verify backend/ was inspected deeply.
3. Verify ai_service/ was inspected.
4. Verify data_pipeline/ was inspected.
5. Verify Docker Compose was inspected.
6. Verify deployment configuration was inspected.
7. Verify go.mod/go.sum and Python dependencies were inspected where applicable.
8. Verify actual imports/dependencies were inspected.
9. Verify duplicate folders were explicitly recorded.
10. Verify no secrets were written.
11. Verify no source code was modified.
12. Verify no files were moved/renamed/deleted.
13. Verify only ARCHITECTURE.md was created/modified.

At the end of your response, report ONLY:

- ARCHITECTURE.md created/updated
- Root components detected
- Number of backend domains detected
- Major architectural duplications detected
- Major dependency concerns detected
- Areas that remain unclear

Do not perform any refactoring.