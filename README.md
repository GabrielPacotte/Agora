# Agora

Agora is an open-source discussion platform backend.

It is designed as a long-term project with a strong focus on correctness, clarity, and maintainability, rather than short-term feature velocity.

This repository currently hosts **version 0.1**, which is a foundation release.

---

## What is Agora?

Agora is a backend for a discussion platform centered around:
- posts with explicit stances
- structured comments
- user preferences
- content curation
- moderation-ready foundations

The project explores how to build a social platform backend that is:
- explicit rather than implicit
- auditable rather than opaque
- evolvable rather than fragile

Agora is not built to “move fast and break things”.  
It is built to **move deliberately and understand things**.

---

## Current State (v0.1)

Version 0.1 is **not a product release**.  
It is a **technical milestone**.

At this stage, Agora provides:
- a complete domain model
- a working HTTP API
- authentication and session handling
- persistence via PostgreSQL
- a clean separation between layers

Some features are intentionally simple or naïve.  
This is by design.

The goal of v0.1 is to establish:
- the architectural direction
- the coding standards
- the mental model of the system

---

## What This Version Is (and Is Not)

This version **is**:
- a usable backend
- testable and structured
- suitable as a foundation for further development
- representative of the project’s philosophy

This version **is not**:
- production-ready
- feature-complete
- optimized for scale or UX
- a finished product

---

## Philosophy

Agora is guided by a few core ideas:

- Business rules belong in the domain, not scattered across layers
- SQL should be explicit and readable, not hidden behind ORMs
- Security mechanisms should be extensible from day one
- Code should be understandable by humans first
- Contributors should know exactly where a change belongs

This project intentionally favors clarity over cleverness.

---

## Documentation

Detailed documentation is provided separately using **Docusaurus**, including:
- architectural overview
- domain concepts
- HTTP API reference
- contribution guidelines

The README only provides a high-level view.

---

## Contributing

Agora is open-source, but it is not a “drive-by PR” project.

Contributions are welcome from developers who:
- care about clean architecture
- are comfortable reasoning about domain logic
- value explicit design choices
- are open to review and discussion

A detailed CONTRIBUTING.md will describe expectations and workflow.

---

## Roadmap

The next major milestone (v1) will focus on:
- improving feed and search logic
- introducing moderation workflows
- hardening authentication and security
- adding observability and monitoring
- preparing production-grade migrations

---

## Why This Project Exists

Even if Agora never reaches a large audience, it serves as:
- a reference backend architecture
- a concrete example of disciplined software design
- a serious open-source codebase meant to be read and learned from

If it grows beyond that, even better.

---

## License

Agora is licensed under the GNU Affero General Public License v3 (GNU AGPLv3).
This choice is intentional: any deployment of a modified version of Agora as a service must remain open-source.