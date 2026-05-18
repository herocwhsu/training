# Training Repo Architecture Update Design

**Date:** 2026-05-18  
**Goal:** Rewrite `utexample/` to follow vortex-backend clean architecture conventions, using `company` and `user` as training domains.

---

## 1. Motivation

The current `utexample/` uses a layer-first layout, a flat `mocks/` folder, reflect-mode mockgen, and outdated test patterns (`defer ctrl.Finish()`, `context.Background()`, separate `tests/` folder). This design brings it in line with the vortex-backend conventions that developers will encounter in production.

---

## 2. Directory Structure

Replace the existing `utexample/` with the following domain-first layout:

```
utexample/
└── internal/
    ├── company/
    │   ├── domain/
    │   │   ├── company.go              # Company entity + Validate(), business methods
    │   │   └── error.go                # ErrCompanyNotFound, ErrInvalidEmail, ErrInvalidName
    │   ├── application/
    │   │   ├── company_service.go      # CompanyService struct + Create/Get/List/Remove
    │   │   └── company_service_test.go # Co-located unit tests
    │   ├── adapter/
    │   │   ├── controller/
    │   │   │   ├── company_controller.go      # Controller struct + DTOs
    │   │   │   └── company_controller_test.go
    │   │   └── repo/
    │   │       ├── company_repo.go            # companyDoc struct, doc↔entity mapping, RDS stub impl
    │   │       └── company_repo_test.go
    │   └── interfaces/
    │       ├── interfaces.go           # CompanyRepository, CompanyDAO interfaces + go:generate
    │       └── mock/
    │           └── interfaces.go       # Generated mocks (mockgen source mode)
    └── user/
        ├── domain/
        │   ├── user.go
        │   └── error.go
        ├── application/
        │   ├── user_service.go
        │   └── user_service_test.go
        ├── adapter/
        │   ├── controller/
        │   │   ├── user_controller.go
        │   │   └── user_controller_test.go
        │   └── repo/
        │       ├── user_repo.go
        │       └── user_repo_test.go
        └── interfaces/
            ├── interfaces.go
            └── mock/
                └── interfaces.go
```

Tests are **co-located** with source files, not in a separate `tests/` folder.

---

## 3. Interfaces & Mocks

### Convention
- All interfaces for a domain live in a single `interfaces/interfaces.go`
- Mocks are generated from that file using **mockgen source mode**
- Mocks live in `interfaces/mock/interfaces.go`

### go:generate directive
```go
//go:generate mockgen -source=interfaces.go -destination mock/interfaces.go -package=mock
```

### Company interfaces example
```go
package interfaces

import (
    "context"
    "github.com/herocwhsu/training/utexample/internal/company/domain"
)

type CompanyRepository interface {
    Save(ctx context.Context, company *domain.Company) error
    FindByID(ctx context.Context, id string) (*domain.Company, error)
    List(ctx context.Context) ([]*domain.Company, error)
    Remove(ctx context.Context, id string) error
}

type CompanyDAO interface {
    Insert(ctx context.Context, email, name string) (id string, err error)
    FindByID(ctx context.Context, id string) (email, name string, err error)
    List(ctx context.Context) ([]struct{ ID, Email, Name string }, error)
    DeleteByID(ctx context.Context, id string) error
}
```

---

## 4. Doc↔Entity Mapping

The repo adapter owns a DB doc struct with DB tags. Domain entities have no DB tags.

```go
// adapter/repo/company_repo.go

type companyDoc struct {
    ID    string `db:"id"`
    Email string `db:"email"`
    Name  string `db:"name"`
}

func docToEntity(doc *companyDoc) *domain.Company {
    return &domain.Company{ID: doc.ID, Email: doc.Email, Name: doc.Name}
}

func entityToDoc(e *domain.Company) *companyDoc {
    return &companyDoc{ID: e.ID, Email: e.Email, Name: e.Name}
}
```

The application and domain layers never import DB-specific types.

---

## 5. Domain Errors

Each domain has a dedicated `error.go`:

```go
// domain/error.go
package domain

import "errors"

var (
    ErrCompanyNotFound = errors.New("company not found")
    ErrInvalidEmail    = errors.New("invalid email")
    ErrInvalidName     = errors.New("invalid name")
)
```

Callers use `errors.Is()` for comparison. No string matching.

---

## 6. Test Conventions

### Setup pattern — deps struct
```go
type companyServiceDeps struct {
    repo *mock.MockCompanyRepository
    svc  *CompanyService
}

func setupCompanyServiceTest(t *testing.T) *companyServiceDeps {
    ctrl := gomock.NewController(t)
    repo := mock.NewMockCompanyRepository(ctrl)
    return &companyServiceDeps{
        repo: repo,
        svc:  NewCompanyService(repo),
    }
}
```

### Naming convention
`Should{Result}_When{Condition}` — e.g. `ShouldReturnID_WhenInputIsValid`, `ShouldReturnError_WhenRepoFails`

### Rules
| Old pattern | New pattern |
|---|---|
| `defer ctrl.Finish()` | Not needed (auto-cleanup via `t`) |
| `context.Background()` | `t.Context()` |
| Separate `tests/` folder | Co-located `*_test.go` |
| Inline mock setup per test | `setupTest()` returning deps struct |
| `ShouldSuccess_WhenRepoOK` | `ShouldReturnID_WhenInputIsValid` |

### Example test
```go
func TestCompanyService_Create(t *testing.T) {
    t.Run("ShouldReturnID_WhenInputIsValid", func(t *testing.T) {
        d := setupCompanyServiceTest(t)
        d.repo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil)

        id, err := d.svc.Create(t.Context(), "acme@co.com", "Acme")
        require.NoError(t, err)
        assert.NotEmpty(t, id)
    })

    t.Run("ShouldReturnError_WhenRepoFails", func(t *testing.T) {
        d := setupCompanyServiceTest(t)
        d.repo.EXPECT().Save(gomock.Any(), gomock.Any()).Return(errors.New("db error"))

        _, err := d.svc.Create(t.Context(), "acme@co.com", "Acme")
        require.Error(t, err)
    })
}
```

---

## 7. What Is Excluded

- **Google Wire DI** — too complex for training; manual `New()` constructors are used
- **GraphQL / REST routing** — out of scope; training focuses on internal layers only
- **Real DB connections** — DAO implementations remain stubs

---

## 8. Regenerating Mocks

After changing any interface, re-run:

```bash
go generate ./utexample/internal/company/interfaces/...
go generate ./utexample/internal/user/interfaces/...
# or all at once
go generate ./utexample/...
```

---

## 9. Running Tests

```bash
go test ./utexample/...
go test -v ./utexample/...
```
