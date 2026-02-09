# Project Structure

Agora’s codebase is organized around responsibility boundaries rather than technical layers.

Below is an overview of the most important directories and their roles.

---

## `internal/domain`

The heart of the system.

Contains:
- Core business entities
- Domain invariants and validations
- Domain-specific errors

Characteristics:
- No dependency on infrastructure
- No HTTP concepts
- No database concerns
- No JSON or serialization logic

This package should remain usable in complete isolation.

---

## `internal/services`

Application use cases.

Contains:
- Business workflows
- Authorization rules
- Coordination between repositories

Characteristics:
- Depends on domain and repository interfaces
- Does not know about SQL or HTTP
- Owns all business decision-making

Services are where “what is allowed” and “what is forbidden” is defined.

---

## `internal/repositories`

Persistence contracts and implementations.

Contains:
- Repository interfaces
- PostgreSQL implementations
- Mapping between database rows and domain objects

Characteristics:
- SQL is accessed exclusively via functions
- No business rules
- No authorization logic

---

## `internal/http`

HTTP delivery layer.

Contains:
- Handlers
- DTOs (input and output)
- Middlewares
- Routing

Characteristics:
- Responsible for HTTP semantics only
- Translates HTTP ↔ domain
- No business logic

---

## `db/`

Database layer.

Contains:
- Schemas
- SQL functions
- Seeds
- (Later) migrations

SQL functions are considered a stable API for repositories.

---

## `cmd/api`

Application bootstrap.

Contains:
- Dependency wiring
- Configuration loading
- HTTP server startup

No business logic should exist here.

---

## Design Rule

If you are unsure where code belongs, ask:

> “Can this be tested without HTTP or SQL?”

If the answer is yes, it probably does not belong in `internal/http` or `db/`.