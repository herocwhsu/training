# Training Repo Architecture Update Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Rewrite `utexample/` to follow vortex-backend clean architecture conventions using `company` and `user` as training domains.

**Architecture:** Domain-first layout under `utexample/internal/{domain}/` with four sub-layers: `domain/`, `application/`, `adapter/`, `interfaces/`. Interfaces are centralized per domain; mocks are generated via mockgen source mode and co-located with interfaces. Tests live next to source files.

**Tech Stack:** Go 1.25, github.com/golang/mock v1.6.0, github.com/stretchr/testify v1.10.0

---

## File Map

### Company domain
| File | Responsibility |
|---|---|
| `utexample/internal/company/domain/company.go` | Company entity, Validate(), business methods |
| `utexample/internal/company/domain/error.go` | Domain error vars |
| `utexample/internal/company/interfaces/interfaces.go` | CompanyRepository + CompanyDAO interfaces |
| `utexample/internal/company/interfaces/mock/interfaces.go` | Generated mocks |
| `utexample/internal/company/adapter/repo/company_repo.go` | companyDoc, doc↔entity mapping, RDS stub |
| `utexample/internal/company/adapter/repo/company_repo_test.go` | Repo unit tests |
| `utexample/internal/company/application/company_service.go` | CompanyService struct + methods |
| `utexample/internal/company/application/company_service_test.go` | Service unit tests |
| `utexample/internal/company/adapter/controller/company_controller.go` | Controller struct + DTOs |
| `utexample/internal/company/adapter/controller/company_controller_test.go` | Controller unit tests |

### User domain
| File | Responsibility |
|---|---|
| `utexample/internal/user/domain/user.go` | User entity, Validate() |
| `utexample/internal/user/domain/error.go` | Domain error vars |
| `utexample/internal/user/interfaces/interfaces.go` | UserRepository + UserDAO interfaces |
| `utexample/internal/user/interfaces/mock/interfaces.go` | Generated mocks |
| `utexample/internal/user/adapter/repo/user_repo.go` | userDoc, doc↔entity mapping, RDS stub |
| `utexample/internal/user/adapter/repo/user_repo_test.go` | Repo unit tests |
| `utexample/internal/user/application/user_service.go` | UserService struct + methods |
| `utexample/internal/user/application/user_service_test.go` | Service unit tests |
| `utexample/internal/user/adapter/controller/user_controller.go` | Controller struct + DTOs |
| `utexample/internal/user/adapter/controller/user_controller_test.go` | Controller unit tests |

---
## Task 1: Remove old utexample structure

**Files:**
- Delete: `utexample/` (entire directory)

- [ ] **Step 1: Remove old code**

```bash
rm -rf utexample/
```

- [ ] **Step 2: Verify removal**

```bash
ls utexample 2>&1
```
Expected: `ls: utexample: No such file or directory`

- [ ] **Step 3: Commit**

```bash
git add -A
git commit -m "chore: remove old utexample structure"
```

---

## Task 2: Company domain layer

**Files:**
- Create: `utexample/internal/company/domain/company.go`
- Create: `utexample/internal/company/domain/error.go`
- Create: `utexample/internal/company/domain/company_test.go`

- [ ] **Step 1: Write the failing domain test**

Create `utexample/internal/company/domain/company_test.go`:

```go
package domain_test

import (
	"testing"

	"github.com/herocwhsu/training/utexample/internal/company/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompany_Validate(t *testing.T) {
	t.Run("ShouldReturnNil_WhenAllFieldsAreValid", func(t *testing.T) {
		c := &domain.Company{ID: "cmp_1", Email: "a@b.com", Name: "Acme"}
		require.NoError(t, c.Validate())
	})

	t.Run("ShouldReturnError_WhenEmailIsEmpty", func(t *testing.T) {
		c := &domain.Company{ID: "cmp_1", Email: "", Name: "Acme"}
		err := c.Validate()
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInvalidEmail)
	})

	t.Run("ShouldReturnError_WhenNameIsEmpty", func(t *testing.T) {
		c := &domain.Company{ID: "cmp_1", Email: "a@b.com", Name: ""}
		err := c.Validate()
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInvalidName)
	})
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./utexample/internal/company/domain/...
```
Expected: compile error — package not found

- [ ] **Step 3: Create domain error file**

Create `utexample/internal/company/domain/error.go`:

```go
package domain

import "errors"

var (
	ErrCompanyNotFound = errors.New("company not found")
	ErrInvalidEmail    = errors.New("invalid email")
	ErrInvalidName     = errors.New("invalid name")
)
```

- [ ] **Step 4: Create domain entity**

Create `utexample/internal/company/domain/company.go`:

```go
package domain

type Company struct {
	ID    string
	Email string
	Name  string
}

func (c *Company) Validate() error {
	if c.Email == "" {
		return ErrInvalidEmail
	}
	if c.Name == "" {
		return ErrInvalidName
	}
	return nil
}
```

- [ ] **Step 5: Run test to verify it passes**

```bash
go test ./utexample/internal/company/domain/...
```
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add utexample/internal/company/domain/
git commit -m "feat: add company domain entity and errors"
```

---
## Task 3: Company interfaces and mocks

**Files:**
- Create: `utexample/internal/company/interfaces/interfaces.go`
- Create: `utexample/internal/company/interfaces/mock/interfaces.go` (generated)

- [ ] **Step 1: Create interfaces file**

Create `utexample/internal/company/interfaces/interfaces.go`:

```go
package interfaces

import (
	"context"

	"github.com/herocwhsu/training/utexample/internal/company/domain"
)

//go:generate mockgen -source=interfaces.go -destination mock/interfaces.go -package=mock

type CompanyRepository interface {
	Save(ctx context.Context, company *domain.Company) error
	FindByID(ctx context.Context, id string) (*domain.Company, error)
	List(ctx context.Context) ([]*domain.Company, error)
	Remove(ctx context.Context, id string) error
}

type CompanyDAO interface {
	Insert(ctx context.Context, email, name string) (id string, err error)
	FindByID(ctx context.Context, id string) (email, name string, err error)
	List(ctx context.Context) ([]*CompanyRow, error)
	DeleteByID(ctx context.Context, id string) error
}

type CompanyRow struct {
	ID    string
	Email string
	Name  string
}
```

- [ ] **Step 2: Generate mocks**

```bash
cd utexample/internal/company/interfaces && go generate ./...
```
Expected: creates `mock/interfaces.go`

- [ ] **Step 3: Verify mock file exists**

```bash
ls utexample/internal/company/interfaces/mock/
```
Expected: `interfaces.go`

- [ ] **Step 4: Commit**

```bash
git add utexample/internal/company/interfaces/
git commit -m "feat: add company interfaces and generated mocks"
```

---

## Task 4: Company repo adapter

**Files:**
- Create: `utexample/internal/company/adapter/repo/company_repo.go`
- Create: `utexample/internal/company/adapter/repo/company_repo_test.go`

- [ ] **Step 1: Write the failing repo test**

Create `utexample/internal/company/adapter/repo/company_repo_test.go`:

```go
package repo_test

import (
	"errors"
	"testing"

	"github.com/herocwhsu/training/utexample/internal/company/adapter/repo"
	"github.com/herocwhsu/training/utexample/internal/company/domain"
	"github.com/herocwhsu/training/utexample/internal/company/interfaces/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type companyRepoDeps struct {
	dao  *mock.MockCompanyDAO
	repo *repo.CompanyRepository
}

func setupCompanyRepoTest(t *testing.T) *companyRepoDeps {
	ctrl := gomock.NewController(t)
	dao := mock.NewMockCompanyDAO(ctrl)
	return &companyRepoDeps{
		dao:  dao,
		repo: repo.NewCompanyRepository(dao),
	}
}

func TestCompanyRepository_Save(t *testing.T) {
	t.Run("ShouldReturnNil_WhenDAOInsertSucceeds", func(t *testing.T) {
		d := setupCompanyRepoTest(t)
		company := &domain.Company{Email: "a@b.com", Name: "Acme"}
		d.dao.EXPECT().Insert(t.Context(), "a@b.com", "Acme").Return("cmp_1", nil)

		err := d.repo.Save(t.Context(), company)
		require.NoError(t, err)
		assert.Equal(t, "cmp_1", company.ID)
	})

	t.Run("ShouldReturnError_WhenDAOInsertFails", func(t *testing.T) {
		d := setupCompanyRepoTest(t)
		company := &domain.Company{Email: "a@b.com", Name: "Acme"}
		d.dao.EXPECT().Insert(t.Context(), "a@b.com", "Acme").Return("", errors.New("db error"))

		err := d.repo.Save(t.Context(), company)
		require.Error(t, err)
	})
}

func TestCompanyRepository_FindByID(t *testing.T) {
	t.Run("ShouldReturnCompany_WhenDAOFindsRecord", func(t *testing.T) {
		d := setupCompanyRepoTest(t)
		d.dao.EXPECT().FindByID(t.Context(), "cmp_1").Return("a@b.com", "Acme", nil)

		got, err := d.repo.FindByID(t.Context(), "cmp_1")
		require.NoError(t, err)
		assert.Equal(t, &domain.Company{ID: "cmp_1", Email: "a@b.com", Name: "Acme"}, got)
	})

	t.Run("ShouldReturnError_WhenDAOReturnsNotFound", func(t *testing.T) {
		d := setupCompanyRepoTest(t)
		d.dao.EXPECT().FindByID(t.Context(), "cmp_404").Return("", "", domain.ErrCompanyNotFound)

		got, err := d.repo.FindByID(t.Context(), "cmp_404")
		require.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestCompanyRepository_Remove(t *testing.T) {
	t.Run("ShouldReturnNil_WhenDAODeleteSucceeds", func(t *testing.T) {
		d := setupCompanyRepoTest(t)
		d.dao.EXPECT().DeleteByID(t.Context(), "cmp_1").Return(nil)

		err := d.repo.Remove(t.Context(), "cmp_1")
		require.NoError(t, err)
	})
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./utexample/internal/company/adapter/repo/...
```
Expected: compile error — package not found

- [ ] **Step 3: Create repo implementation**

Create `utexample/internal/company/adapter/repo/company_repo.go`:

```go
package repo

import (
	"context"

	"github.com/herocwhsu/training/utexample/internal/company/domain"
	"github.com/herocwhsu/training/utexample/internal/company/interfaces"
)

type companyDoc struct {
	ID    string
	Email string
	Name  string
}

func docToEntity(doc *companyDoc) *domain.Company {
	return &domain.Company{ID: doc.ID, Email: doc.Email, Name: doc.Name}
}

func entityToDoc(e *domain.Company) *companyDoc {
	return &companyDoc{ID: e.ID, Email: e.Email, Name: e.Name}
}

type CompanyRepository struct {
	dao interfaces.CompanyDAO
}

func NewCompanyRepository(dao interfaces.CompanyDAO) *CompanyRepository {
	return &CompanyRepository{dao: dao}
}

func (r *CompanyRepository) Save(ctx context.Context, company *domain.Company) error {
	id, err := r.dao.Insert(ctx, company.Email, company.Name)
	if err != nil {
		return err
	}
	company.ID = id
	return nil
}

func (r *CompanyRepository) FindByID(ctx context.Context, id string) (*domain.Company, error) {
	email, name, err := r.dao.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	doc := &companyDoc{ID: id, Email: email, Name: name}
	entity := docToEntity(doc)
	if err := entity.Validate(); err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *CompanyRepository) List(ctx context.Context) ([]*domain.Company, error) {
	rows, err := r.dao.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.Company, 0, len(rows))
	for _, row := range rows {
		result = append(result, docToEntity(&companyDoc{ID: row.ID, Email: row.Email, Name: row.Name}))
	}
	return result, nil
}

func (r *CompanyRepository) Remove(ctx context.Context, id string) error {
	return r.dao.DeleteByID(ctx, id)
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test ./utexample/internal/company/adapter/repo/...
```
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add utexample/internal/company/adapter/repo/
git commit -m "feat: add company repo adapter with doc-entity mapping"
```

---
## Task 5: Company application service

**Files:**
- Create: `utexample/internal/company/application/company_service.go`
- Create: `utexample/internal/company/application/company_service_test.go`

- [ ] **Step 1: Write the failing service test**

Create `utexample/internal/company/application/company_service_test.go`:

```go
package application_test

import (
	"errors"
	"testing"

	"github.com/herocwhsu/training/utexample/internal/company/application"
	"github.com/herocwhsu/training/utexample/internal/company/domain"
	"github.com/herocwhsu/training/utexample/internal/company/interfaces/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type companyServiceDeps struct {
	repo *mock.MockCompanyRepository
	svc  *application.CompanyService
}

func setupCompanyServiceTest(t *testing.T) *companyServiceDeps {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockCompanyRepository(ctrl)
	return &companyServiceDeps{
		repo: repo,
		svc:  application.NewCompanyService(repo),
	}
}

func TestCompanyService_Create(t *testing.T) {
	t.Run("ShouldReturnID_WhenInputIsValid", func(t *testing.T) {
		d := setupCompanyServiceTest(t)
		d.repo.EXPECT().
			Save(t.Context(), gomock.Any()).
			DoAndReturn(func(_ interface{}, c *domain.Company) error {
				c.ID = "cmp_1"
				return nil
			})

		id, err := d.svc.Create(t.Context(), "a@b.com", "Acme")
		require.NoError(t, err)
		assert.Equal(t, "cmp_1", id)
	})

	t.Run("ShouldReturnError_WhenEmailIsEmpty", func(t *testing.T) {
		d := setupCompanyServiceTest(t)

		_, err := d.svc.Create(t.Context(), "", "Acme")
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInvalidEmail)
	})

	t.Run("ShouldReturnError_WhenRepoFails", func(t *testing.T) {
		d := setupCompanyServiceTest(t)
		d.repo.EXPECT().Save(t.Context(), gomock.Any()).Return(errors.New("db error"))

		_, err := d.svc.Create(t.Context(), "a@b.com", "Acme")
		require.Error(t, err)
	})
}

func TestCompanyService_Get(t *testing.T) {
	t.Run("ShouldReturnCompany_WhenFound", func(t *testing.T) {
		d := setupCompanyServiceTest(t)
		expected := &domain.Company{ID: "cmp_1", Email: "a@b.com", Name: "Acme"}
		d.repo.EXPECT().FindByID(t.Context(), "cmp_1").Return(expected, nil)

		got, err := d.svc.Get(t.Context(), "cmp_1")
		require.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("ShouldReturnError_WhenNotFound", func(t *testing.T) {
		d := setupCompanyServiceTest(t)
		d.repo.EXPECT().FindByID(t.Context(), "cmp_404").Return(nil, domain.ErrCompanyNotFound)

		got, err := d.svc.Get(t.Context(), "cmp_404")
		require.Error(t, err)
		assert.Nil(t, got)
		assert.ErrorIs(t, err, domain.ErrCompanyNotFound)
	})
}

func TestCompanyService_Remove(t *testing.T) {
	t.Run("ShouldReturnNil_WhenRemoveSucceeds", func(t *testing.T) {
		d := setupCompanyServiceTest(t)
		d.repo.EXPECT().Remove(t.Context(), "cmp_1").Return(nil)

		err := d.svc.Remove(t.Context(), "cmp_1")
		require.NoError(t, err)
	})
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./utexample/internal/company/application/...
```
Expected: compile error — package not found

- [ ] **Step 3: Create service implementation**

Create `utexample/internal/company/application/company_service.go`:

```go
package application

import (
	"context"
	"fmt"

	"github.com/herocwhsu/training/utexample/internal/company/domain"
	"github.com/herocwhsu/training/utexample/internal/company/interfaces"
)

type CompanyService struct {
	repo interfaces.CompanyRepository
}

func NewCompanyService(repo interfaces.CompanyRepository) *CompanyService {
	return &CompanyService{repo: repo}
}

func (s *CompanyService) Create(ctx context.Context, email, name string) (string, error) {
	c := &domain.Company{Email: email, Name: name}
	if err := c.Validate(); err != nil {
		return "", err
	}
	if err := s.repo.Save(ctx, c); err != nil {
		return "", fmt.Errorf("save company: %w", err)
	}
	return c.ID, nil
}

func (s *CompanyService) Get(ctx context.Context, id string) (*domain.Company, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *CompanyService) List(ctx context.Context) ([]*domain.Company, error) {
	return s.repo.List(ctx)
}

func (s *CompanyService) Remove(ctx context.Context, id string) error {
	return s.repo.Remove(ctx, id)
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test ./utexample/internal/company/application/...
```
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add utexample/internal/company/application/
git commit -m "feat: add company application service"
```

---

## Task 6: Company controller

**Files:**
- Create: `utexample/internal/company/adapter/controller/company_controller.go`
- Create: `utexample/internal/company/adapter/controller/company_controller_test.go`

- [ ] **Step 1: Write the failing controller test**

Create `utexample/internal/company/adapter/controller/company_controller_test.go`:

```go
package controller_test

import (
	"errors"
	"testing"

	"github.com/herocwhsu/training/utexample/internal/company/adapter/controller"
	"github.com/herocwhsu/training/utexample/internal/company/domain"
	"github.com/herocwhsu/training/utexample/internal/company/interfaces/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type companyControllerDeps struct {
	svc  *mock.MockCompanyService
	ctrl *controller.CompanyController
}

func setupCompanyControllerTest(t *testing.T) *companyControllerDeps {
	c := gomock.NewController(t)
	svc := mock.NewMockCompanyService(c)
	return &companyControllerDeps{
		svc:  svc,
		ctrl: controller.NewCompanyController(svc),
	}
}

func TestCompanyController_Create(t *testing.T) {
	t.Run("ShouldReturnID_WhenServiceSucceeds", func(t *testing.T) {
		d := setupCompanyControllerTest(t)
		d.svc.EXPECT().Create(t.Context(), "a@b.com", "Acme").Return("cmp_1", nil)

		id, err := d.ctrl.Create(t.Context(), controller.CreateCompanyInput{Email: "a@b.com", Name: "Acme"})
		require.NoError(t, err)
		assert.Equal(t, "cmp_1", id)
	})

	t.Run("ShouldReturnError_WhenServiceFails", func(t *testing.T) {
		d := setupCompanyControllerTest(t)
		d.svc.EXPECT().Create(t.Context(), "bad@b.com", "Bad").Return("", errors.New("service error"))

		_, err := d.ctrl.Create(t.Context(), controller.CreateCompanyInput{Email: "bad@b.com", Name: "Bad"})
		require.Error(t, err)
	})
}

func TestCompanyController_Get(t *testing.T) {
	t.Run("ShouldReturnCompanyOutput_WhenFound", func(t *testing.T) {
		d := setupCompanyControllerTest(t)
		d.svc.EXPECT().Get(t.Context(), "cmp_1").Return(
			&domain.Company{ID: "cmp_1", Email: "a@b.com", Name: "Acme"}, nil,
		)

		out, err := d.ctrl.Get(t.Context(), "cmp_1")
		require.NoError(t, err)
		assert.Equal(t, &controller.CompanyOutput{ID: "cmp_1", Email: "a@b.com", Name: "Acme"}, out)
	})

	t.Run("ShouldReturnError_WhenNotFound", func(t *testing.T) {
		d := setupCompanyControllerTest(t)
		d.svc.EXPECT().Get(t.Context(), "cmp_404").Return(nil, domain.ErrCompanyNotFound)

		out, err := d.ctrl.Get(t.Context(), "cmp_404")
		require.Error(t, err)
		assert.Nil(t, out)
	})
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./utexample/internal/company/adapter/controller/...
```
Expected: compile error — package not found

- [ ] **Step 3: Add CompanyService interface to interfaces.go**

Append to `utexample/internal/company/interfaces/interfaces.go` (inside the file, after CompanyDAO):

```go
type CompanyService interface {
	Create(ctx context.Context, email, name string) (string, error)
	Get(ctx context.Context, id string) (*domain.Company, error)
	List(ctx context.Context) ([]*domain.Company, error)
	Remove(ctx context.Context, id string) error
}
```

- [ ] **Step 4: Regenerate mocks**

```bash
cd utexample/internal/company/interfaces && go generate ./...
```

- [ ] **Step 5: Create controller implementation**

Create `utexample/internal/company/adapter/controller/company_controller.go`:

```go
package controller

import (
	"context"

	"github.com/herocwhsu/training/utexample/internal/company/interfaces"
)

type CreateCompanyInput struct {
	Email string
	Name  string
}

type CompanyOutput struct {
	ID    string
	Email string
	Name  string
}

type CompanyController struct {
	svc interfaces.CompanyService
}

func NewCompanyController(svc interfaces.CompanyService) *CompanyController {
	return &CompanyController{svc: svc}
}

func (c *CompanyController) Create(ctx context.Context, input CreateCompanyInput) (string, error) {
	return c.svc.Create(ctx, input.Email, input.Name)
}

func (c *CompanyController) Get(ctx context.Context, id string) (*CompanyOutput, error) {
	company, err := c.svc.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &CompanyOutput{ID: company.ID, Email: company.Email, Name: company.Name}, nil
}
```

- [ ] **Step 6: Run test to verify it passes**

```bash
go test ./utexample/internal/company/...
```
Expected: PASS

- [ ] **Step 7: Commit**

```bash
git add utexample/internal/company/
git commit -m "feat: add company controller with DTOs"
```

---
## Task 7: User domain layer

**Files:**
- Create: `utexample/internal/user/domain/user.go`
- Create: `utexample/internal/user/domain/error.go`
- Create: `utexample/internal/user/domain/user_test.go`

- [ ] **Step 1: Write the failing domain test**

Create `utexample/internal/user/domain/user_test.go`:

```go
package domain_test

import (
	"testing"

	"github.com/herocwhsu/training/utexample/internal/user/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser_Validate(t *testing.T) {
	t.Run("ShouldReturnNil_WhenAllFieldsAreValid", func(t *testing.T) {
		u := &domain.User{ID: "usr_1", Email: "a@b.com", Name: "Alice"}
		require.NoError(t, u.Validate())
	})

	t.Run("ShouldReturnError_WhenEmailIsEmpty", func(t *testing.T) {
		u := &domain.User{ID: "usr_1", Email: "", Name: "Alice"}
		err := u.Validate()
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInvalidEmail)
	})

	t.Run("ShouldReturnError_WhenNameIsEmpty", func(t *testing.T) {
		u := &domain.User{ID: "usr_1", Email: "a@b.com", Name: ""}
		err := u.Validate()
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInvalidName)
	})
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./utexample/internal/user/domain/...
```
Expected: compile error — package not found

- [ ] **Step 3: Create domain error file**

Create `utexample/internal/user/domain/error.go`:

```go
package domain

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidEmail = errors.New("invalid email")
	ErrInvalidName  = errors.New("invalid name")
)
```

- [ ] **Step 4: Create domain entity**

Create `utexample/internal/user/domain/user.go`:

```go
package domain

type User struct {
	ID    string
	Email string
	Name  string
}

func (u *User) Validate() error {
	if u.Email == "" {
		return ErrInvalidEmail
	}
	if u.Name == "" {
		return ErrInvalidName
	}
	return nil
}
```

- [ ] **Step 5: Run test to verify it passes**

```bash
go test ./utexample/internal/user/domain/...
```
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add utexample/internal/user/domain/
git commit -m "feat: add user domain entity and errors"
```

---

## Task 8: User interfaces and mocks

**Files:**
- Create: `utexample/internal/user/interfaces/interfaces.go`
- Create: `utexample/internal/user/interfaces/mock/interfaces.go` (generated)

- [ ] **Step 1: Create interfaces file**

Create `utexample/internal/user/interfaces/interfaces.go`:

```go
package interfaces

import (
	"context"

	"github.com/herocwhsu/training/utexample/internal/user/domain"
)

//go:generate mockgen -source=interfaces.go -destination mock/interfaces.go -package=mock

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) error
	FindByID(ctx context.Context, id string) (*domain.User, error)
	List(ctx context.Context) ([]*domain.User, error)
	Remove(ctx context.Context, id string) error
}

type UserDAO interface {
	Insert(ctx context.Context, email, name string) (id string, err error)
	FindByID(ctx context.Context, id string) (email, name string, err error)
	List(ctx context.Context) ([]*UserRow, error)
	DeleteByID(ctx context.Context, id string) error
}

type UserRow struct {
	ID    string
	Email string
	Name  string
}

type UserService interface {
	Create(ctx context.Context, email, name string) (string, error)
	Get(ctx context.Context, id string) (*domain.User, error)
	List(ctx context.Context) ([]*domain.User, error)
	Remove(ctx context.Context, id string) error
}
```

- [ ] **Step 2: Generate mocks**

```bash
cd utexample/internal/user/interfaces && go generate ./...
```
Expected: creates `mock/interfaces.go`

- [ ] **Step 3: Verify mock file exists**

```bash
ls utexample/internal/user/interfaces/mock/
```
Expected: `interfaces.go`

- [ ] **Step 4: Commit**

```bash
git add utexample/internal/user/interfaces/
git commit -m "feat: add user interfaces and generated mocks"
```

---
## Task 9: User repo adapter

**Files:**
- Create: `utexample/internal/user/adapter/repo/user_repo.go`
- Create: `utexample/internal/user/adapter/repo/user_repo_test.go`

- [ ] **Step 1: Write the failing repo test**

Create `utexample/internal/user/adapter/repo/user_repo_test.go`:

```go
package repo_test

import (
	"errors"
	"testing"

	"github.com/herocwhsu/training/utexample/internal/user/adapter/repo"
	"github.com/herocwhsu/training/utexample/internal/user/domain"
	"github.com/herocwhsu/training/utexample/internal/user/interfaces/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type userRepoDeps struct {
	dao  *mock.MockUserDAO
	repo *repo.UserRepository
}

func setupUserRepoTest(t *testing.T) *userRepoDeps {
	ctrl := gomock.NewController(t)
	dao := mock.NewMockUserDAO(ctrl)
	return &userRepoDeps{
		dao:  dao,
		repo: repo.NewUserRepository(dao),
	}
}

func TestUserRepository_Save(t *testing.T) {
	t.Run("ShouldSetID_WhenDAOInsertSucceeds", func(t *testing.T) {
		d := setupUserRepoTest(t)
		user := &domain.User{Email: "a@b.com", Name: "Alice"}
		d.dao.EXPECT().Insert(t.Context(), "a@b.com", "Alice").Return("usr_1", nil)

		err := d.repo.Save(t.Context(), user)
		require.NoError(t, err)
		assert.Equal(t, "usr_1", user.ID)
	})

	t.Run("ShouldReturnError_WhenDAOInsertFails", func(t *testing.T) {
		d := setupUserRepoTest(t)
		user := &domain.User{Email: "a@b.com", Name: "Alice"}
		d.dao.EXPECT().Insert(t.Context(), "a@b.com", "Alice").Return("", errors.New("db error"))

		err := d.repo.Save(t.Context(), user)
		require.Error(t, err)
	})
}

func TestUserRepository_FindByID(t *testing.T) {
	t.Run("ShouldReturnUser_WhenDAOFindsRecord", func(t *testing.T) {
		d := setupUserRepoTest(t)
		d.dao.EXPECT().FindByID(t.Context(), "usr_1").Return("a@b.com", "Alice", nil)

		got, err := d.repo.FindByID(t.Context(), "usr_1")
		require.NoError(t, err)
		assert.Equal(t, &domain.User{ID: "usr_1", Email: "a@b.com", Name: "Alice"}, got)
	})

	t.Run("ShouldReturnError_WhenDAOReturnsNotFound", func(t *testing.T) {
		d := setupUserRepoTest(t)
		d.dao.EXPECT().FindByID(t.Context(), "usr_404").Return("", "", domain.ErrUserNotFound)

		got, err := d.repo.FindByID(t.Context(), "usr_404")
		require.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestUserRepository_Remove(t *testing.T) {
	t.Run("ShouldReturnNil_WhenDAODeleteSucceeds", func(t *testing.T) {
		d := setupUserRepoTest(t)
		d.dao.EXPECT().DeleteByID(t.Context(), "usr_1").Return(nil)

		err := d.repo.Remove(t.Context(), "usr_1")
		require.NoError(t, err)
	})
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./utexample/internal/user/adapter/repo/...
```
Expected: compile error — package not found

- [ ] **Step 3: Create repo implementation**

Create `utexample/internal/user/adapter/repo/user_repo.go`:

```go
package repo

import (
	"context"

	"github.com/herocwhsu/training/utexample/internal/user/domain"
	"github.com/herocwhsu/training/utexample/internal/user/interfaces"
)

type userDoc struct {
	ID    string
	Email string
	Name  string
}

func docToEntity(doc *userDoc) *domain.User {
	return &domain.User{ID: doc.ID, Email: doc.Email, Name: doc.Name}
}

func entityToDoc(u *domain.User) *userDoc {
	return &userDoc{ID: u.ID, Email: u.Email, Name: u.Name}
}

type UserRepository struct {
	dao interfaces.UserDAO
}

func NewUserRepository(dao interfaces.UserDAO) *UserRepository {
	return &UserRepository{dao: dao}
}

func (r *UserRepository) Save(ctx context.Context, user *domain.User) error {
	id, err := r.dao.Insert(ctx, user.Email, user.Name)
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	email, name, err := r.dao.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	doc := &userDoc{ID: id, Email: email, Name: name}
	entity := docToEntity(doc)
	if err := entity.Validate(); err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *UserRepository) List(ctx context.Context) ([]*domain.User, error) {
	rows, err := r.dao.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.User, 0, len(rows))
	for _, row := range rows {
		result = append(result, docToEntity(&userDoc{ID: row.ID, Email: row.Email, Name: row.Name}))
	}
	return result, nil
}

func (r *UserRepository) Remove(ctx context.Context, id string) error {
	return r.dao.DeleteByID(ctx, id)
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test ./utexample/internal/user/adapter/repo/...
```
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add utexample/internal/user/adapter/repo/
git commit -m "feat: add user repo adapter with doc-entity mapping"
```

---

## Task 10: User application service

**Files:**
- Create: `utexample/internal/user/application/user_service.go`
- Create: `utexample/internal/user/application/user_service_test.go`

- [ ] **Step 1: Write the failing service test**

Create `utexample/internal/user/application/user_service_test.go`:

```go
package application_test

import (
	"errors"
	"testing"

	"github.com/herocwhsu/training/utexample/internal/user/application"
	"github.com/herocwhsu/training/utexample/internal/user/domain"
	"github.com/herocwhsu/training/utexample/internal/user/interfaces/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type userServiceDeps struct {
	repo *mock.MockUserRepository
	svc  *application.UserService
}

func setupUserServiceTest(t *testing.T) *userServiceDeps {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockUserRepository(ctrl)
	return &userServiceDeps{
		repo: repo,
		svc:  application.NewUserService(repo),
	}
}

func TestUserService_Create(t *testing.T) {
	t.Run("ShouldReturnID_WhenInputIsValid", func(t *testing.T) {
		d := setupUserServiceTest(t)
		d.repo.EXPECT().
			Save(t.Context(), gomock.Any()).
			DoAndReturn(func(_ interface{}, u *domain.User) error {
				u.ID = "usr_1"
				return nil
			})

		id, err := d.svc.Create(t.Context(), "a@b.com", "Alice")
		require.NoError(t, err)
		assert.Equal(t, "usr_1", id)
	})

	t.Run("ShouldReturnError_WhenEmailIsEmpty", func(t *testing.T) {
		d := setupUserServiceTest(t)

		_, err := d.svc.Create(t.Context(), "", "Alice")
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInvalidEmail)
	})

	t.Run("ShouldReturnError_WhenRepoFails", func(t *testing.T) {
		d := setupUserServiceTest(t)
		d.repo.EXPECT().Save(t.Context(), gomock.Any()).Return(errors.New("db error"))

		_, err := d.svc.Create(t.Context(), "a@b.com", "Alice")
		require.Error(t, err)
	})
}

func TestUserService_Get(t *testing.T) {
	t.Run("ShouldReturnUser_WhenFound", func(t *testing.T) {
		d := setupUserServiceTest(t)
		expected := &domain.User{ID: "usr_1", Email: "a@b.com", Name: "Alice"}
		d.repo.EXPECT().FindByID(t.Context(), "usr_1").Return(expected, nil)

		got, err := d.svc.Get(t.Context(), "usr_1")
		require.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("ShouldReturnError_WhenNotFound", func(t *testing.T) {
		d := setupUserServiceTest(t)
		d.repo.EXPECT().FindByID(t.Context(), "usr_404").Return(nil, domain.ErrUserNotFound)

		got, err := d.svc.Get(t.Context(), "usr_404")
		require.Error(t, err)
		assert.Nil(t, got)
		assert.ErrorIs(t, err, domain.ErrUserNotFound)
	})
}

func TestUserService_Remove(t *testing.T) {
	t.Run("ShouldReturnNil_WhenRemoveSucceeds", func(t *testing.T) {
		d := setupUserServiceTest(t)
		d.repo.EXPECT().Remove(t.Context(), "usr_1").Return(nil)

		err := d.svc.Remove(t.Context(), "usr_1")
		require.NoError(t, err)
	})
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./utexample/internal/user/application/...
```
Expected: compile error — package not found

- [ ] **Step 3: Create service implementation**

Create `utexample/internal/user/application/user_service.go`:

```go
package application

import (
	"context"
	"fmt"

	"github.com/herocwhsu/training/utexample/internal/user/domain"
	"github.com/herocwhsu/training/utexample/internal/user/interfaces"
)

type UserService struct {
	repo interfaces.UserRepository
}

func NewUserService(repo interfaces.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, email, name string) (string, error) {
	u := &domain.User{Email: email, Name: name}
	if err := u.Validate(); err != nil {
		return "", err
	}
	if err := s.repo.Save(ctx, u); err != nil {
		return "", fmt.Errorf("save user: %w", err)
	}
	return u.ID, nil
}

func (s *UserService) Get(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *UserService) List(ctx context.Context) ([]*domain.User, error) {
	return s.repo.List(ctx)
}

func (s *UserService) Remove(ctx context.Context, id string) error {
	return s.repo.Remove(ctx, id)
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test ./utexample/internal/user/application/...
```
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add utexample/internal/user/application/
git commit -m "feat: add user application service"
```

---
## Task 11: User controller

**Files:**
- Create: `utexample/internal/user/adapter/controller/user_controller.go`
- Create: `utexample/internal/user/adapter/controller/user_controller_test.go`

- [ ] **Step 1: Write the failing controller test**

Create `utexample/internal/user/adapter/controller/user_controller_test.go`:

```go
package controller_test

import (
	"errors"
	"testing"

	"github.com/herocwhsu/training/utexample/internal/user/adapter/controller"
	"github.com/herocwhsu/training/utexample/internal/user/domain"
	"github.com/herocwhsu/training/utexample/internal/user/interfaces/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type userControllerDeps struct {
	svc  *mock.MockUserService
	ctrl *controller.UserController
}

func setupUserControllerTest(t *testing.T) *userControllerDeps {
	c := gomock.NewController(t)
	svc := mock.NewMockUserService(c)
	return &userControllerDeps{
		svc:  svc,
		ctrl: controller.NewUserController(svc),
	}
}

func TestUserController_Create(t *testing.T) {
	t.Run("ShouldReturnID_WhenServiceSucceeds", func(t *testing.T) {
		d := setupUserControllerTest(t)
		d.svc.EXPECT().Create(t.Context(), "a@b.com", "Alice").Return("usr_1", nil)

		id, err := d.ctrl.Create(t.Context(), controller.CreateUserInput{Email: "a@b.com", Name: "Alice"})
		require.NoError(t, err)
		assert.Equal(t, "usr_1", id)
	})

	t.Run("ShouldReturnError_WhenServiceFails", func(t *testing.T) {
		d := setupUserControllerTest(t)
		d.svc.EXPECT().Create(t.Context(), "bad@b.com", "Bad").Return("", errors.New("service error"))

		_, err := d.ctrl.Create(t.Context(), controller.CreateUserInput{Email: "bad@b.com", Name: "Bad"})
		require.Error(t, err)
	})
}

func TestUserController_Get(t *testing.T) {
	t.Run("ShouldReturnUserOutput_WhenFound", func(t *testing.T) {
		d := setupUserControllerTest(t)
		d.svc.EXPECT().Get(t.Context(), "usr_1").Return(
			&domain.User{ID: "usr_1", Email: "a@b.com", Name: "Alice"}, nil,
		)

		out, err := d.ctrl.Get(t.Context(), "usr_1")
		require.NoError(t, err)
		assert.Equal(t, &controller.UserOutput{ID: "usr_1", Email: "a@b.com", Name: "Alice"}, out)
	})

	t.Run("ShouldReturnError_WhenNotFound", func(t *testing.T) {
		d := setupUserControllerTest(t)
		d.svc.EXPECT().Get(t.Context(), "usr_404").Return(nil, domain.ErrUserNotFound)

		out, err := d.ctrl.Get(t.Context(), "usr_404")
		require.Error(t, err)
		assert.Nil(t, out)
	})
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./utexample/internal/user/adapter/controller/...
```
Expected: compile error — package not found

- [ ] **Step 3: Create controller implementation**

Create `utexample/internal/user/adapter/controller/user_controller.go`:

```go
package controller

import (
	"context"

	"github.com/herocwhsu/training/utexample/internal/user/interfaces"
)

type CreateUserInput struct {
	Email string
	Name  string
}

type UserOutput struct {
	ID    string
	Email string
	Name  string
}

type UserController struct {
	svc interfaces.UserService
}

func NewUserController(svc interfaces.UserService) *UserController {
	return &UserController{svc: svc}
}

func (c *UserController) Create(ctx context.Context, input CreateUserInput) (string, error) {
	return c.svc.Create(ctx, input.Email, input.Name)
}

func (c *UserController) Get(ctx context.Context, id string) (*UserOutput, error) {
	user, err := c.svc.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &UserOutput{ID: user.ID, Email: user.Email, Name: user.Name}, nil
}
```

- [ ] **Step 4: Run all tests to verify everything passes**

```bash
go test ./utexample/...
```
Expected: PASS (all packages)

- [ ] **Step 5: Commit**

```bash
git add utexample/internal/user/adapter/controller/
git commit -m "feat: add user controller with DTOs"
```

---

## Task 12: Update README

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Update README to reflect new architecture**

Replace the contents of `README.md` with:

```markdown
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
    └── user/
        ├── domain/
        ├── application/
        ├── adapter/
        │   ├── controller/
        │   └── repo/
        └── interfaces/
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
# or all at once
go generate ./utexample/...
```

---

## Dependencies

| Package | Purpose |
|---|---|
| `github.com/golang/mock` | Mock generation and runtime |
| `github.com/stretchr/testify` | Test assertions |
```

- [ ] **Step 2: Run all tests one final time**

```bash
go test ./utexample/...
```
Expected: PASS

- [ ] **Step 3: Commit**

```bash
git add README.md
git commit -m "docs: update README for new clean architecture structure"
```

---
