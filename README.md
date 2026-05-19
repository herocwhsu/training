# Go Clean Architecture Training Sample

A hands-on example demonstrating vortex-backend clean architecture conventions using `company` and `user` as training domains.

---

## Architecture

Each domain follows a four-layer structure. Dependencies flow inward — outer layers depend on inner layers, never the reverse.

```
┌─────────────────────────────────────────────────┐
│  adapter/controller   Maps input DTOs ↔ service │
│    depends on ▼ interfaces.XxxService            │
├─────────────────────────────────────────────────┤
│  application          Orchestrates domain +repo  │
│    depends on ▼ interfaces.XxxRepository         │
├─────────────────────────────────────────────────┤
│  adapter/repo         Doc↔entity mapping + stub  │
│    depends on ▼ interfaces.XxxDAO                │
├─────────────────────────────────────────────────┤
│  domain               Entity + errors            │
│    no external deps                              │
└─────────────────────────────────────────────────┘
```

Interfaces for each domain live in `interfaces/interfaces.go`. Mocks are generated from that file and placed in `interfaces/mock/`.

---

## Project Structure

```
utexample/
└── internal/
    ├── company/
    │   ├── domain/           # Company entity, Validate(), error vars
    │   ├── application/      # CompanyService (orchestrates repo)
    │   ├── adapter/
    │   │   ├── controller/   # Controller + DTOs
    │   │   └── repo/         # companyDoc, doc↔entity mapping, DAO stub
    │   └── interfaces/       # All interfaces + generated mocks
    ├── user/
    │   ├── domain/
    │   ├── application/
    │   ├── adapter/
    │   │   ├── controller/
    │   │   └── repo/
    │   └── interfaces/
    └── membership/
        ├── domain/           # Membership entity, Validate(), error vars
        ├── application/      # MembershipService (coordinates company + user + repo)
        ├── adapter/
        │   ├── controller/   # Controller + DTOs
        │   └── repo/         # membershipDoc, doc↔entity mapping, DAO stub
        └── interfaces/       # All interfaces + generated mocks (incl. CompanyReader, UserReader)
```

Tests are co-located with source files (`*_test.go` next to `*.go`).

---

## Key Conventions

- **Domain-first layout** — code is organized by domain, not by layer
- **Centralized interfaces** — all interfaces for a domain in one `interfaces/interfaces.go`
- **mockgen source mode** — `mockgen -source=interfaces.go`; mocks in `interfaces/mock/`
- **Doc↔entity mapping** — DB structs stay in the adapter; domain entities have no DB tags
- **Domain errors** — `errors.Is()` for comparison, never string matching
- **Test setup struct** — `setupXxxTest(t)` returns a deps struct; no `defer ctrl.Finish()`
- **Test naming** — `Should{Result}_When{Condition}`
- **Context in tests** — `t.Context()` not `context.Background()`

---

## Cross-Domain Communication

The `membership` domain demonstrates how domains communicate without direct coupling:

- `membership/interfaces` defines narrow `CompanyReader` and `UserReader` interfaces
- At wire-up, `*companyrepo.CompanyRepository` and `*userrepo.UserRepository` satisfy these interfaces
- Neither `company` nor `user` imports from `membership`
- `MembershipService` depends only on interfaces it owns — not on other domains' packages

---

## Running Tests

```bash
go test ./utexample/...
go test -v ./utexample/...
```

---

## Regenerating Mocks

```bash
go generate ./utexample/internal/company/interfaces/...
go generate ./utexample/internal/user/interfaces/...
go generate ./utexample/internal/membership/interfaces/...
# or all at once
go generate ./utexample/...
```

---

## Dependencies

| Package | Purpose |
|---|---|
| `github.com/golang/mock` | Mock generation and runtime |
| `github.com/stretchr/testify` | Test assertions |
