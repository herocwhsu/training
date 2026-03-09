# Go Unit Testing Training Sample

A hands-on example demonstrating how to write unit tests in Go using a clean layered architecture, dependency injection, `gomock`, and `testify`.

---

## Architecture

Each layer depends only on the interface defined by the layer below it, never on a concrete implementation.

```
┌─────────────────────────────────────────┐
│  Controller  (companyctl)               │  Maps HTTP/input DTOs ↔ service calls
│    depends on ▼ CompanyService          │
├─────────────────────────────────────────┤
│  Service     (companysvc)               │  Holds business rules / validation
│    depends on ▼ CompanyRepository       │
├─────────────────────────────────────────┤
│  Repository  (companyrepo)              │  Builds domain objects from raw data
│    depends on ▼ CompanyDAO              │
├─────────────────────────────────────────┤
│  DAO         (companydao)               │  Executes raw DB queries
└─────────────────────────────────────────┘
```

---

## Project Structure

```
utexample/
├── internal/
│   ├── domain/
│   │   └── company.go          # Domain model + Validate()
│   ├── dao/companydao/
│   │   └── dao.go              # CompanyDAO interface + RDS stub
│   ├── repo/companyrepo/
│   │   └── repository.go       # CompanyRepository interface + implementation
│   ├── service/companysvc/
│   │   └── service.go          # CompanyService interface + implementation
│   └── controller/companyctl/
│       └── controller.go       # Controller struct + CompanyInfo DTO
├── mocks/
│   ├── mock_dao.go             # Generated mock for CompanyDAO
│   ├── mock_repository.go      # Generated mock for CompanyRepository
│   └── mock_service.go         # Generated mock for CompanyService
└── tests/
    ├── domain_company_test.go  # Tests for domain validation logic
    ├── repo_company_test.go    # Tests for repository (mocks DAO)
    ├── service_company_test.go # Tests for service (mocks repository)
    └── controller_company_test.go # Tests for controller (mocks service)
```

---

## Key Concepts

- **Dependency Injection** — Each layer receives its dependency via constructor (`New(...)`), making it easy to swap real implementations with mocks in tests.

- **Interface-based design** — Every layer exposes an interface (`CompanyDAO`, `CompanyRepository`, `CompanyService`). The layer above depends on the interface, not the concrete struct.

- **gomock** — Used to generate and control mock implementations of each interface, letting you test each layer in complete isolation.

- **testify** — Provides readable assertions (`assert.Nil`, `assert.Equal`, `assert.Error`) that produce clear failure messages.

- **Test per layer** — Each layer has its own test file. Only the immediate dependency is mocked, so failures pinpoint exactly which layer is broken.

---

## Running Tests

```bash
go test ./utexample/...
```

To see verbose output per test case:

```bash
go test -v ./utexample/...
```

---

## Regenerating Mocks

Mocks are generated with [mockgen](https://github.com/golang/mock). After changing any interface, re-run the corresponding `go generate`:

```bash
# From the repo root
go generate ./utexample/internal/dao/companydao/...
go generate ./utexample/internal/repo/companyrepo/...
go generate ./utexample/internal/service/companysvc/...
```

Or regenerate everything at once:

```bash
go generate ./utexample/...
```

---

## Dependencies

| Package | Purpose |
|---|---|
| `github.com/golang/mock` | Mock generation and runtime |
| `github.com/stretchr/testify` | Test assertions |
