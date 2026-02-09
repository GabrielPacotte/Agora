---
sidebar_position: 1
sidebar_title: "Overview"
---

# Code Overview

Agora is not a side project, a prototype, or a playground.

It is a deliberately engineered backend whose primary goal is to demonstrate
how a production-ready system can be designed with strict separation of concerns,
strong domain boundaries, and long-term maintainability in mind.

This documentation is intended for contributors who are already comfortable with
backend development concepts and want to understand how Agora is structured,
why certain architectural decisions were made, and how to contribute without
breaking those guarantees.

---

## Core Philosophy

Agora follows a strict interpretation of Clean Architecture:

- The **domain layer is sovereign**
- Dependencies always point **inward**
- Infrastructure is replaceable
- SQL is treated as an **internal API**, not an implementation detail
- HTTP is a delivery mechanism, not the core of the system

Each layer has a single responsibility and a clearly defined role.

---

## Request Lifecycle

A typical request flows through the system as follows:

1. **HTTP handler**
   - Parses input (path, query, body)
   - Maps DTOs to domain inputs
   - Handles authentication and authorization context
2. **Service**
   - Applies business rules
   - Orchestrates use cases
   - Enforces permissions and invariants
3. **Repository**
   - Acts as a gateway to persistence
   - Maps domain objects to database calls
4. **SQL function**
   - Encapsulates data access logic
   - Guarantees shape and constraints of returned data

The response then flows back in reverse order.

At no point should:
- HTTP handlers contain business logic
- Services know about SQL
- Domain objects depend on infrastructure concerns

---

## What Agora Is Not

- Not a CRUD application
- Not a framework
- Not optimized for rapid iteration at the expense of clarity
- Not permissive toward architectural shortcuts

This project favors **clarity over speed**, **explicitness over convenience**, and
**long-term maintainability over short-term productivity**.

---

## Stability Guarantees

- The **domain layer** is considered highly stable
- Service interfaces evolve cautiously
- SQL functions are versioned intentionally
- HTTP contracts are documented explicitly via OpenAPI

Breaking changes are deliberate, documented, and justified.