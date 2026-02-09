# Contributing to Agora

First of all, thank you for your interest in contributing to **Agora**.  
This project aims to be a **clean, rigorous, and production-grade open-source backend**, with strong architectural and ethical principles.

Please read this document carefully before opening an issue or submitting a pull request.

---

## Project Philosophy

Agora is built with the following principles in mind:

- **Clean Architecture** over quick hacks
- **Explicit domain modeling**
- **Strong separation of concerns**
- **Readable, testable, maintainable code**
- **Ethical considerations** around privacy, data ownership, and platform governance

This repository is **not** a playground or a tutorial project.  
Contributions are expected to meet a professional standard.

---

## License & Contributions

Agora is licensed under the **GNU Affero General Public License v3 (AGPL-3.0)**.

By contributing to this project, you agree that:

- Your contribution will be licensed under AGPL-3.0
- Any deployment of a modified version of Agora as a service must remain open-source
- You have the right to submit the code you are contributing

If you are not comfortable with AGPL, **do not contribute**.

---

## What You Can Contribute

Contributions are welcome in the following areas:

- Backend features (services, repositories, domain logic)
- SQL schema and functions
- Performance and security improvements
- API documentation
- Tests (unit, integration)
- Developer tooling (CI, scripts)
- Bug fixes and refactors with a clear rationale

Before working on a large change, **open an issue first** to discuss it.

---

## Issues Guidelines

When opening an issue:

- Be clear and precise
- Explain the **problem**, not just the solution
- Provide context (why this matters)
- Reference existing code or documentation when relevant

Low-effort issues (e.g. “it doesn’t work”, “add X please”) will be closed.

---

## Pull Request Guidelines

### General Rules

- One pull request = **one coherent change**
- Keep PRs reasonably sized
- Every PR must have a **clear description**
- Code must compile and tests must pass

### Required Quality Standards

Your contribution must:

- Follow existing architecture and naming conventions
- Avoid introducing unnecessary abstractions
- Avoid leaking infrastructure concerns into the domain layer
- Include tests when relevant
- Be formatted with standard Go tooling (`gofmt`)

If a contribution lowers the overall code quality, it will be rejected.

---

## Code Style & Architecture

### Go Code

- Idiomatic Go is expected
- Prefer explicitness over cleverness
- No global state unless strictly justified
- Errors must be handled explicitly
- Domain errors belong to the `domain` package

### Architecture

- Domain layer must remain independent of infrastructure
- Repositories must not contain business logic
- Services orchestrate use cases, not persistence details
- HTTP handlers must remain thin

If you are unsure where something belongs, **ask before coding**.

---

## Tests

- New behavior should come with tests
- Regressions must be covered
- Tests must be deterministic and isolated

The goal is **confidence**, not coverage for the sake of coverage.

---

## Commits

- Use clear, descriptive commit messages
- Avoid mixing unrelated changes
- Prefer multiple small commits over one large opaque commit

Example:
```
feat(posts): add feed ranking placeholder
fix(users): prevent duplicate preference entries
refactor(domain): simplify stance validation
```

---

## Review Process

All contributions are reviewed by the maintainer.

Expect:
- Questions
- Requests for clarification
- Requests for changes

This is normal and part of maintaining a high-quality codebase.

---

## Code of Conduct

Be respectful, constructive, and professional.

This project values:
- Technical rigor
- Honest discussion
- Clear communication

Harassment, hostility, or bad-faith behavior will not be tolerated.

---

## Final Note

Agora is a long-term project with strong foundations.  
If you enjoy **thinking deeply about architecture, trade-offs, and clean systems**, you are very welcome here.

If you are looking for quick wins or low-effort contributions, this may not be the right project.

Thank you for taking the time to read this.