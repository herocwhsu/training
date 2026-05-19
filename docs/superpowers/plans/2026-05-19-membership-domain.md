# Membership Domain Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Add a `membership` domain that demonstrates cross-domain communication — a user can be added to or removed from a company, coordinated by membership without company or user knowing about each other.

**Architecture:** Membership owns narrow `CompanyReader` and `UserReader` interfaces; the existing `CompanyRepository` and `UserRepository` satisfy them at wire-up time. Membership has its own four-layer structure (domain → interfaces → application → adapter) identical to the existing `company` and `user` domains.

**Tech Stack:** Go 1.25, `github.com/golang/mock`, `github.com/stretchr/testify`

---

## File Map

| File | Action |
|------|--------|
| `utexample/internal/membership/domain/membership.go` | Create — entity + Validate() |
| `utexample/internal/membership/domain/error.go` | Create — domain error vars |
| `utexample/internal/membership/interfaces/interfaces.go` | Create — all interfaces incl. CompanyReader/UserReader |
| `utexample/internal/membership/interfaces/mock/interfaces.go` | Generate via go generate |
| `utexample/internal/membership/adapter/repo/membership_repo.go` | Create — DAO-backed repo |
| `utexample/internal/membership/adapter/repo/membership_repo_test.go` | Create — repo tests |
| `utexample/internal/membership/application/membership_service.go` | Create — service orchestration |
| `utexample/internal/membership/application/membership_service_test.go` | Create — service tests |
| `utexample/internal/membership/adapter/controller/membership_controller.go` | Create — DTOs + controller |
| `utexample/internal/membership/adapter/controller/membership_controller_test.go` | Create — controller tests |
| `README.md` | Modify — add membership section |

---

## Task 1: Domain Entity and Errors

**Files:**
- Create: `utexample/internal/membership/domain/membership.go`
- Create: `utexample/internal/membership/domain/error.go`

- [ ] **Step 1: Create the domain error file**

```go
// utexample/internal/membership/domain/error.go
package domain

import "errors"

var (
	ErrMembershipNotFound = errors.New("membership not found")
	ErrAlreadyMember      = errors.New("already a member")
	ErrInvalidRole        = errors.New("invalid role")
)
```

- [ ] **Step 2: Create the entity file**

```go
// utexample/internal/membership/domain/membership.go
package domain

type Membership struct {
	ID        string
	CompanyID string
	UserID    string
	Role      string // "member" | "admin"
}

func (m *Membership) Validate() error {
	if m.Role != "member" && m.Role != "admin" {
		return ErrInvalidRole
	}
	return nil
}
```

- [ ] **Step 3: Verify it compiles**

Run: `go build ./utexample/internal/membership/domain/...`
Expected: no output (success)

- [ ] **Step 4: Commit**

```bash
git add utexample/internal/membership/domain/
git commit -m "feat: add membership domain entity and errors"
```

---

## Task 2: Interfaces

**Files:**
- Create: `utexample/internal/membership/interfaces/interfaces.go`

- [ ] **Step 1: Create the interfaces file**

```go
// utexample/internal/membership/interfaces/interfaces.go
package interfaces

import (
	"context"

	companydomain "github.com/herocwhsu/training/utexample/internal/company/domain"
	"github.com/herocwhsu/training/utexample/internal/membership/domain"
	userdomain "github.com/herocwhsu/training/utexample/internal/user/domain"
)

//go:generate mockgen -source=interfaces.go -destination mock/interfaces.go -package=mock

// Cross-domain: narrow read interfaces owned by membership.
// Satisfied at wire-up by *companyrepo.CompanyRepository and *userrepo.UserRepository.
type CompanyReader interface {
	FindByID(ctx context.Context, id string) (*companydomain.Company, error)
}

type UserReader interface {
	FindByID(ctx context.Context, id string) (*userdomain.User, error)
}

type MembershipRepository interface {
	Save(ctx context.Context, m *domain.Membership) error
	FindByID(ctx context.Context, id string) (*domain.Membership, error)
	FindByCompanyID(ctx context.Context, companyID string) ([]*domain.Membership, error)
	FindByUserID(ctx context.Context, userID string) ([]*domain.Membership, error)
	ExistsByCompanyAndUser(ctx context.Context, companyID, userID string) (bool, error)
	Remove(ctx context.Context, id string) error
}

type MembershipDAO interface {
	Insert(ctx context.Context, companyID, userID, role string) (string, error)
	FindByID(ctx context.Context, id string) (companyID, userID, role string, err error)
	FindByCompanyID(ctx context.Context, companyID string) ([]MembershipRow, error)
	FindByUserID(ctx context.Context, userID string) ([]MembershipRow, error)
	ExistsByCompanyAndUser(ctx context.Context, companyID, userID string) (bool, error)
	DeleteByID(ctx context.Context, id string) error
}

type MembershipRow struct {
	ID        string
	CompanyID string
	UserID    string
	Role      string
}

type MembershipService interface {
	Add(ctx context.Context, companyID, userID, role string) (string, error)
	Remove(ctx context.Context, membershipID string) error
	ListByCompany(ctx context.Context, companyID string) ([]*domain.Membership, error)
	ListByUser(ctx context.Context, userID string) ([]*domain.Membership, error)
}
```

- [ ] **Step 2: Verify it compiles**

Run: `go build ./utexample/internal/membership/interfaces/...`
Expected: no output (success)

- [ ] **Step 3: Commit**

```bash
git add utexample/internal/membership/interfaces/interfaces.go
git commit -m "feat: add membership interfaces"
```

---

## Task 3: Generate Mocks

**Files:**
- Create: `utexample/internal/membership/interfaces/mock/interfaces.go` (generated)

- [ ] **Step 1: Run go generate**

Run: `go generate ./utexample/internal/membership/interfaces/...`
Expected: creates `utexample/internal/membership/interfaces/mock/interfaces.go`

- [ ] **Step 2: Verify the mock file exists**

Run: `ls utexample/internal/membership/interfaces/mock/`
Expected: `interfaces.go`

- [ ] **Step 3: Verify it compiles**

Run: `go build ./utexample/internal/membership/interfaces/mock/...`
Expected: no output (success)

- [ ] **Step 4: Commit**

```bash
git add utexample/internal/membership/interfaces/mock/
git commit -m "feat: generate membership mocks"
```

---

## Task 4: Repository Adapter

**Files:**
- Create: `utexample/internal/membership/adapter/repo/membership_repo.go`
- Create: `utexample/internal/membership/adapter/repo/membership_repo_test.go`

- [ ] **Step 1: Write the failing repo tests**

```go
// utexample/internal/membership/adapter/repo/membership_repo_test.go
package repo_test

import (
	"errors"
	"testing"

	"github.com/herocwhsu/training/utexample/internal/membership/adapter/repo"
	"github.com/herocwhsu/training/utexample/internal/membership/domain"
	"github.com/herocwhsu/training/utexample/internal/membership/interfaces"
	"github.com/herocwhsu/training/utexample/internal/membership/interfaces/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type membershipRepoDeps struct {
	dao  *mock.MockMembershipDAO
	repo *repo.MembershipRepository
}

func setupMembershipRepoTest(t *testing.T) *membershipRepoDeps {
	ctrl := gomock.NewController(t)
	dao := mock.NewMockMembershipDAO(ctrl)
	return &membershipRepoDeps{
		dao:  dao,
		repo: repo.NewMembershipRepository(dao),
	}
}

func TestMembershipRepository_Save(t *testing.T) {
	t.Run("ShouldSetID_WhenDAOInsertSucceeds", func(t *testing.T) {
		d := setupMembershipRepoTest(t)
		m := &domain.Membership{CompanyID: "cmp_1", UserID: "usr_1", Role: "member"}
		d.dao.EXPECT().Insert(t.Context(), "cmp_1", "usr_1", "member").Return("mem_1", nil)

		err := d.repo.Save(t.Context(), m)
		require.NoError(t, err)
		assert.Equal(t, "mem_1", m.ID)
	})

	t.Run("ShouldReturnError_WhenDAOInsertFails", func(t *testing.T) {
		d := setupMembershipRepoTest(t)
		m := &domain.Membership{CompanyID: "cmp_1", UserID: "usr_1", Role: "member"}
		d.dao.EXPECT().Insert(t.Context(), "cmp_1", "usr_1", "member").Return("", errors.New("db error"))

		err := d.repo.Save(t.Context(), m)
		require.Error(t, err)
	})
}

func TestMembershipRepository_FindByID(t *testing.T) {
	t.Run("ShouldReturnMembership_WhenDAOFindsRecord", func(t *testing.T) {
		d := setupMembershipRepoTest(t)
		d.dao.EXPECT().FindByID(t.Context(), "mem_1").Return("cmp_1", "usr_1", "member", nil)

		got, err := d.repo.FindByID(t.Context(), "mem_1")
		require.NoError(t, err)
		assert.Equal(t, &domain.Membership{ID: "mem_1", CompanyID: "cmp_1", UserID: "usr_1", Role: "member"}, got)
	})

	t.Run("ShouldReturnError_WhenDAOReturnsNotFound", func(t *testing.T) {
		d := setupMembershipRepoTest(t)
		d.dao.EXPECT().FindByID(t.Context(), "mem_404").Return("", "", "", domain.ErrMembershipNotFound)

		got, err := d.repo.FindByID(t.Context(), "mem_404")
		require.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestMembershipRepository_FindByCompanyID(t *testing.T) {
	t.Run("ShouldReturnMemberships_WhenDAOSucceeds", func(t *testing.T) {
		d := setupMembershipRepoTest(t)
		rows := []interfaces.MembershipRow{
			{ID: "mem_1", CompanyID: "cmp_1", UserID: "usr_1", Role: "member"},
			{ID: "mem_2", CompanyID: "cmp_1", UserID: "usr_2", Role: "admin"},
		}
		d.dao.EXPECT().FindByCompanyID(t.Context(), "cmp_1").Return(rows, nil)

		got, err := d.repo.FindByCompanyID(t.Context(), "cmp_1")
		require.NoError(t, err)
		assert.Equal(t, []*domain.Membership{
			{ID: "mem_1", CompanyID: "cmp_1", UserID: "usr_1", Role: "member"},
			{ID: "mem_2", CompanyID: "cmp_1", UserID: "usr_2", Role: "admin"},
		}, got)
	})

	t.Run("ShouldReturnError_WhenDAOFails", func(t *testing.T) {
		d := setupMembershipRepoTest(t)
		d.dao.EXPECT().FindByCompanyID(t.Context(), "cmp_1").Return(nil, errors.New("db error"))

		got, err := d.repo.FindByCompanyID(t.Context(), "cmp_1")
		require.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestMembershipRepository_FindByUserID(t *testing.T) {
	t.Run("ShouldReturnMemberships_WhenDAOSucceeds", func(t *testing.T) {
		d := setupMembershipRepoTest(t)
		rows := []interfaces.MembershipRow{
			{ID: "mem_1", CompanyID: "cmp_1", UserID: "usr_1", Role: "member"},
		}
		d.dao.EXPECT().FindByUserID(t.Context(), "usr_1").Return(rows, nil)

		got, err := d.repo.FindByUserID(t.Context(), "usr_1")
		require.NoError(t, err)
		assert.Equal(t, []*domain.Membership{
			{ID: "mem_1", CompanyID: "cmp_1", UserID: "usr_1", Role: "member"},
		}, got)
	})

	t.Run("ShouldReturnError_WhenDAOFails", func(t *testing.T) {
		d := setupMembershipRepoTest(t)
		d.dao.EXPECT().FindByUserID(t.Context(), "usr_1").Return(nil, errors.New("db error"))

		got, err := d.repo.FindByUserID(t.Context(), "usr_1")
		require.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestMembershipRepository_ExistsByCompanyAndUser(t *testing.T) {
	t.Run("ShouldReturnTrue_WhenDAOReturnsTrue", func(t *testing.T) {
		d := setupMembershipRepoTest(t)
		d.dao.EXPECT().ExistsByCompanyAndUser(t.Context(), "cmp_1", "usr_1").Return(true, nil)

		exists, err := d.repo.ExistsByCompanyAndUser(t.Context(), "cmp_1", "usr_1")
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("ShouldReturnError_WhenDAOFails", func(t *testing.T) {
		d := setupMembershipRepoTest(t)
		d.dao.EXPECT().ExistsByCompanyAndUser(t.Context(), "cmp_1", "usr_1").Return(false, errors.New("db error"))

		_, err := d.repo.ExistsByCompanyAndUser(t.Context(), "cmp_1", "usr_1")
		require.Error(t, err)
	})
}

func TestMembershipRepository_Remove(t *testing.T) {
	t.Run("ShouldReturnNil_WhenDAODeleteSucceeds", func(t *testing.T) {
		d := setupMembershipRepoTest(t)
		d.dao.EXPECT().DeleteByID(t.Context(), "mem_1").Return(nil)

		err := d.repo.Remove(t.Context(), "mem_1")
		require.NoError(t, err)
	})
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./utexample/internal/membership/adapter/repo/...`
Expected: FAIL — package `repo` does not exist yet

- [ ] **Step 3: Implement the repository**

```go
// utexample/internal/membership/adapter/repo/membership_repo.go
package repo

import (
	"context"

	"github.com/herocwhsu/training/utexample/internal/membership/domain"
	"github.com/herocwhsu/training/utexample/internal/membership/interfaces"
)

type membershipDoc struct {
	ID        string
	CompanyID string
	UserID    string
	Role      string
}

func docToEntity(doc *membershipDoc) *domain.Membership {
	return &domain.Membership{ID: doc.ID, CompanyID: doc.CompanyID, UserID: doc.UserID, Role: doc.Role}
}

type MembershipRepository struct {
	dao interfaces.MembershipDAO
}

func NewMembershipRepository(dao interfaces.MembershipDAO) *MembershipRepository {
	return &MembershipRepository{dao: dao}
}

func (r *MembershipRepository) Save(ctx context.Context, m *domain.Membership) error {
	id, err := r.dao.Insert(ctx, m.CompanyID, m.UserID, m.Role)
	if err != nil {
		return err
	}
	m.ID = id
	return nil
}

func (r *MembershipRepository) FindByID(ctx context.Context, id string) (*domain.Membership, error) {
	companyID, userID, role, err := r.dao.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return docToEntity(&membershipDoc{ID: id, CompanyID: companyID, UserID: userID, Role: role}), nil
}

func (r *MembershipRepository) FindByCompanyID(ctx context.Context, companyID string) ([]*domain.Membership, error) {
	rows, err := r.dao.FindByCompanyID(ctx, companyID)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.Membership, 0, len(rows))
	for _, row := range rows {
		result = append(result, docToEntity(&membershipDoc{ID: row.ID, CompanyID: row.CompanyID, UserID: row.UserID, Role: row.Role}))
	}
	return result, nil
}

func (r *MembershipRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.Membership, error) {
	rows, err := r.dao.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.Membership, 0, len(rows))
	for _, row := range rows {
		result = append(result, docToEntity(&membershipDoc{ID: row.ID, CompanyID: row.CompanyID, UserID: row.UserID, Role: row.Role}))
	}
	return result, nil
}

func (r *MembershipRepository) ExistsByCompanyAndUser(ctx context.Context, companyID, userID string) (bool, error) {
	return r.dao.ExistsByCompanyAndUser(ctx, companyID, userID)
}

func (r *MembershipRepository) Remove(ctx context.Context, id string) error {
	return r.dao.DeleteByID(ctx, id)
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./utexample/internal/membership/adapter/repo/...`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add utexample/internal/membership/adapter/repo/
git commit -m "feat: add membership repository adapter"
```

---

## Task 5: Application Service

**Files:**
- Create: `utexample/internal/membership/application/membership_service.go`
- Create: `utexample/internal/membership/application/membership_service_test.go`

- [ ] **Step 1: Write the failing service tests**

```go
// utexample/internal/membership/application/membership_service_test.go
package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/herocwhsu/training/utexample/internal/membership/application"
	companydomain "github.com/herocwhsu/training/utexample/internal/company/domain"
	"github.com/herocwhsu/training/utexample/internal/membership/domain"
	userdomain "github.com/herocwhsu/training/utexample/internal/user/domain"
	"github.com/herocwhsu/training/utexample/internal/membership/interfaces/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type membershipServiceDeps struct {
	repo          *mock.MockMembershipRepository
	companyReader *mock.MockCompanyReader
	userReader    *mock.MockUserReader
	svc           *application.MembershipService
}

func setupMembershipServiceTest(t *testing.T) *membershipServiceDeps {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockMembershipRepository(ctrl)
	companyReader := mock.NewMockCompanyReader(ctrl)
	userReader := mock.NewMockUserReader(ctrl)
	return &membershipServiceDeps{
		repo:          repo,
		companyReader: companyReader,
		userReader:    userReader,
		svc:           application.NewMembershipService(repo, companyReader, userReader),
	}
}

func TestMembershipService_Add(t *testing.T) {
	t.Run("ShouldReturnID_WhenBothExistAndNotAlreadyMember", func(t *testing.T) {
		d := setupMembershipServiceTest(t)
		d.companyReader.EXPECT().FindByID(t.Context(), "cmp_1").Return(&companydomain.Company{ID: "cmp_1"}, nil)
		d.userReader.EXPECT().FindByID(t.Context(), "usr_1").Return(&userdomain.User{ID: "usr_1"}, nil)
		d.repo.EXPECT().ExistsByCompanyAndUser(t.Context(), "cmp_1", "usr_1").Return(false, nil)
		d.repo.EXPECT().
			Save(t.Context(), gomock.Any()).
			DoAndReturn(func(_ context.Context, m *domain.Membership) error {
				m.ID = "mem_1"
				return nil
			})

		id, err := d.svc.Add(t.Context(), "cmp_1", "usr_1", "member")
		require.NoError(t, err)
		assert.Equal(t, "mem_1", id)
	})

	t.Run("ShouldReturnError_WhenCompanyNotFound", func(t *testing.T) {
		d := setupMembershipServiceTest(t)
		d.companyReader.EXPECT().FindByID(t.Context(), "cmp_404").Return(nil, companydomain.ErrCompanyNotFound)

		_, err := d.svc.Add(t.Context(), "cmp_404", "usr_1", "member")
		require.Error(t, err)
		assert.ErrorIs(t, err, companydomain.ErrCompanyNotFound)
	})

	t.Run("ShouldReturnError_WhenUserNotFound", func(t *testing.T) {
		d := setupMembershipServiceTest(t)
		d.companyReader.EXPECT().FindByID(t.Context(), "cmp_1").Return(&companydomain.Company{ID: "cmp_1"}, nil)
		d.userReader.EXPECT().FindByID(t.Context(), "usr_404").Return(nil, userdomain.ErrUserNotFound)

		_, err := d.svc.Add(t.Context(), "cmp_1", "usr_404", "member")
		require.Error(t, err)
		assert.ErrorIs(t, err, userdomain.ErrUserNotFound)
	})

	t.Run("ShouldReturnError_WhenAlreadyMember", func(t *testing.T) {
		d := setupMembershipServiceTest(t)
		d.companyReader.EXPECT().FindByID(t.Context(), "cmp_1").Return(&companydomain.Company{ID: "cmp_1"}, nil)
		d.userReader.EXPECT().FindByID(t.Context(), "usr_1").Return(&userdomain.User{ID: "usr_1"}, nil)
		d.repo.EXPECT().ExistsByCompanyAndUser(t.Context(), "cmp_1", "usr_1").Return(true, nil)

		_, err := d.svc.Add(t.Context(), "cmp_1", "usr_1", "member")
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrAlreadyMember)
	})

	t.Run("ShouldReturnError_WhenRoleIsInvalid", func(t *testing.T) {
		d := setupMembershipServiceTest(t)
		d.companyReader.EXPECT().FindByID(t.Context(), "cmp_1").Return(&companydomain.Company{ID: "cmp_1"}, nil)
		d.userReader.EXPECT().FindByID(t.Context(), "usr_1").Return(&userdomain.User{ID: "usr_1"}, nil)
		d.repo.EXPECT().ExistsByCompanyAndUser(t.Context(), "cmp_1", "usr_1").Return(false, nil)

		_, err := d.svc.Add(t.Context(), "cmp_1", "usr_1", "superuser")
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInvalidRole)
	})
}

func TestMembershipService_Remove(t *testing.T) {
	t.Run("ShouldReturnNil_WhenRemoveSucceeds", func(t *testing.T) {
		d := setupMembershipServiceTest(t)
		d.repo.EXPECT().Remove(t.Context(), "mem_1").Return(nil)

		err := d.svc.Remove(t.Context(), "mem_1")
		require.NoError(t, err)
	})

	t.Run("ShouldReturnError_WhenRepoFails", func(t *testing.T) {
		d := setupMembershipServiceTest(t)
		d.repo.EXPECT().Remove(t.Context(), "mem_1").Return(errors.New("db error"))

		err := d.svc.Remove(t.Context(), "mem_1")
		require.Error(t, err)
	})
}

func TestMembershipService_ListByCompany(t *testing.T) {
	t.Run("ShouldReturnMemberships_WhenRepoSucceeds", func(t *testing.T) {
		d := setupMembershipServiceTest(t)
		expected := []*domain.Membership{
			{ID: "mem_1", CompanyID: "cmp_1", UserID: "usr_1", Role: "member"},
		}
		d.repo.EXPECT().FindByCompanyID(t.Context(), "cmp_1").Return(expected, nil)

		got, err := d.svc.ListByCompany(t.Context(), "cmp_1")
		require.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("ShouldReturnError_WhenRepoFails", func(t *testing.T) {
		d := setupMembershipServiceTest(t)
		d.repo.EXPECT().FindByCompanyID(t.Context(), "cmp_1").Return(nil, errors.New("db error"))

		got, err := d.svc.ListByCompany(t.Context(), "cmp_1")
		require.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestMembershipService_ListByUser(t *testing.T) {
	t.Run("ShouldReturnMemberships_WhenRepoSucceeds", func(t *testing.T) {
		d := setupMembershipServiceTest(t)
		expected := []*domain.Membership{
			{ID: "mem_1", CompanyID: "cmp_1", UserID: "usr_1", Role: "member"},
		}
		d.repo.EXPECT().FindByUserID(t.Context(), "usr_1").Return(expected, nil)

		got, err := d.svc.ListByUser(t.Context(), "usr_1")
		require.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("ShouldReturnError_WhenRepoFails", func(t *testing.T) {
		d := setupMembershipServiceTest(t)
		d.repo.EXPECT().FindByUserID(t.Context(), "usr_1").Return(nil, errors.New("db error"))

		got, err := d.svc.ListByUser(t.Context(), "usr_1")
		require.Error(t, err)
		assert.Nil(t, got)
	})
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./utexample/internal/membership/application/...`
Expected: FAIL — package `application` does not exist yet

- [ ] **Step 3: Implement the service**

```go
// utexample/internal/membership/application/membership_service.go
package application

import (
	"context"
	"fmt"

	"github.com/herocwhsu/training/utexample/internal/membership/domain"
	"github.com/herocwhsu/training/utexample/internal/membership/interfaces"
)

type MembershipService struct {
	repo          interfaces.MembershipRepository
	companyReader interfaces.CompanyReader
	userReader    interfaces.UserReader
}

func NewMembershipService(
	repo interfaces.MembershipRepository,
	companyReader interfaces.CompanyReader,
	userReader interfaces.UserReader,
) *MembershipService {
	return &MembershipService{repo: repo, companyReader: companyReader, userReader: userReader}
}

func (s *MembershipService) Add(ctx context.Context, companyID, userID, role string) (string, error) {
	if _, err := s.companyReader.FindByID(ctx, companyID); err != nil {
		return "", fmt.Errorf("find company: %w", err)
	}
	if _, err := s.userReader.FindByID(ctx, userID); err != nil {
		return "", fmt.Errorf("find user: %w", err)
	}
	exists, err := s.repo.ExistsByCompanyAndUser(ctx, companyID, userID)
	if err != nil {
		return "", fmt.Errorf("check membership: %w", err)
	}
	if exists {
		return "", domain.ErrAlreadyMember
	}
	m := &domain.Membership{CompanyID: companyID, UserID: userID, Role: role}
	if err := m.Validate(); err != nil {
		return "", err
	}
	if err := s.repo.Save(ctx, m); err != nil {
		return "", fmt.Errorf("save membership: %w", err)
	}
	return m.ID, nil
}

func (s *MembershipService) Remove(ctx context.Context, membershipID string) error {
	if err := s.repo.Remove(ctx, membershipID); err != nil {
		return fmt.Errorf("remove membership: %w", err)
	}
	return nil
}

func (s *MembershipService) ListByCompany(ctx context.Context, companyID string) ([]*domain.Membership, error) {
	memberships, err := s.repo.FindByCompanyID(ctx, companyID)
	if err != nil {
		return nil, fmt.Errorf("list by company: %w", err)
	}
	return memberships, nil
}

func (s *MembershipService) ListByUser(ctx context.Context, userID string) ([]*domain.Membership, error) {
	memberships, err := s.repo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("list by user: %w", err)
	}
	return memberships, nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./utexample/internal/membership/application/...`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add utexample/internal/membership/application/
git commit -m "feat: add membership service"
```

---

## Task 6: Controller Adapter

**Files:**
- Create: `utexample/internal/membership/adapter/controller/membership_controller.go`
- Create: `utexample/internal/membership/adapter/controller/membership_controller_test.go`

- [ ] **Step 1: Write the failing controller tests**

```go
// utexample/internal/membership/adapter/controller/membership_controller_test.go
package controller_test

import (
	"errors"
	"testing"

	"github.com/herocwhsu/training/utexample/internal/membership/adapter/controller"
	"github.com/herocwhsu/training/utexample/internal/membership/domain"
	"github.com/herocwhsu/training/utexample/internal/membership/interfaces/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type membershipControllerDeps struct {
	svc  *mock.MockMembershipService
	ctrl *controller.MembershipController
}

func setupMembershipControllerTest(t *testing.T) *membershipControllerDeps {
	c := gomock.NewController(t)
	svc := mock.NewMockMembershipService(c)
	return &membershipControllerDeps{
		svc:  svc,
		ctrl: controller.NewMembershipController(svc),
	}
}

func TestMembershipController_Add(t *testing.T) {
	t.Run("ShouldReturnID_WhenServiceSucceeds", func(t *testing.T) {
		d := setupMembershipControllerTest(t)
		d.svc.EXPECT().Add(t.Context(), "cmp_1", "usr_1", "member").Return("mem_1", nil)

		id, err := d.ctrl.Add(t.Context(), controller.AddMemberInput{CompanyID: "cmp_1", UserID: "usr_1", Role: "member"})
		require.NoError(t, err)
		assert.Equal(t, "mem_1", id)
	})

	t.Run("ShouldReturnError_WhenServiceFails", func(t *testing.T) {
		d := setupMembershipControllerTest(t)
		d.svc.EXPECT().Add(t.Context(), "cmp_1", "usr_1", "member").Return("", errors.New("service error"))

		_, err := d.ctrl.Add(t.Context(), controller.AddMemberInput{CompanyID: "cmp_1", UserID: "usr_1", Role: "member"})
		require.Error(t, err)
	})
}

func TestMembershipController_Remove(t *testing.T) {
	t.Run("ShouldReturnNil_WhenServiceSucceeds", func(t *testing.T) {
		d := setupMembershipControllerTest(t)
		d.svc.EXPECT().Remove(t.Context(), "mem_1").Return(nil)

		err := d.ctrl.Remove(t.Context(), "mem_1")
		require.NoError(t, err)
	})

	t.Run("ShouldReturnError_WhenServiceFails", func(t *testing.T) {
		d := setupMembershipControllerTest(t)
		d.svc.EXPECT().Remove(t.Context(), "mem_1").Return(errors.New("service error"))

		err := d.ctrl.Remove(t.Context(), "mem_1")
		require.Error(t, err)
	})
}

func TestMembershipController_ListByCompany(t *testing.T) {
	t.Run("ShouldReturnOutputs_WhenServiceSucceeds", func(t *testing.T) {
		d := setupMembershipControllerTest(t)
		d.svc.EXPECT().ListByCompany(t.Context(), "cmp_1").Return([]*domain.Membership{
			{ID: "mem_1", CompanyID: "cmp_1", UserID: "usr_1", Role: "member"},
		}, nil)

		out, err := d.ctrl.ListByCompany(t.Context(), "cmp_1")
		require.NoError(t, err)
		assert.Equal(t, []*controller.MembershipOutput{
			{ID: "mem_1", CompanyID: "cmp_1", UserID: "usr_1", Role: "member"},
		}, out)
	})

	t.Run("ShouldReturnError_WhenServiceFails", func(t *testing.T) {
		d := setupMembershipControllerTest(t)
		d.svc.EXPECT().ListByCompany(t.Context(), "cmp_1").Return(nil, errors.New("service error"))

		out, err := d.ctrl.ListByCompany(t.Context(), "cmp_1")
		require.Error(t, err)
		assert.Nil(t, out)
	})
}

func TestMembershipController_ListByUser(t *testing.T) {
	t.Run("ShouldReturnOutputs_WhenServiceSucceeds", func(t *testing.T) {
		d := setupMembershipControllerTest(t)
		d.svc.EXPECT().ListByUser(t.Context(), "usr_1").Return([]*domain.Membership{
			{ID: "mem_1", CompanyID: "cmp_1", UserID: "usr_1", Role: "member"},
		}, nil)

		out, err := d.ctrl.ListByUser(t.Context(), "usr_1")
		require.NoError(t, err)
		assert.Equal(t, []*controller.MembershipOutput{
			{ID: "mem_1", CompanyID: "cmp_1", UserID: "usr_1", Role: "member"},
		}, out)
	})

	t.Run("ShouldReturnError_WhenServiceFails", func(t *testing.T) {
		d := setupMembershipControllerTest(t)
		d.svc.EXPECT().ListByUser(t.Context(), "usr_1").Return(nil, errors.New("service error"))

		out, err := d.ctrl.ListByUser(t.Context(), "usr_1")
		require.Error(t, err)
		assert.Nil(t, out)
	})
}
```

- [ ] **Step 2: Run tests to verify they fail**

Run: `go test ./utexample/internal/membership/adapter/controller/...`
Expected: FAIL — package `controller` does not exist yet

- [ ] **Step 3: Implement the controller**

```go
// utexample/internal/membership/adapter/controller/membership_controller.go
package controller

import (
	"context"

	"github.com/herocwhsu/training/utexample/internal/membership/interfaces"
)

type AddMemberInput struct {
	CompanyID string
	UserID    string
	Role      string
}

type MembershipOutput struct {
	ID        string
	CompanyID string
	UserID    string
	Role      string
}

type MembershipController struct {
	svc interfaces.MembershipService
}

func NewMembershipController(svc interfaces.MembershipService) *MembershipController {
	return &MembershipController{svc: svc}
}

func (c *MembershipController) Add(ctx context.Context, input AddMemberInput) (string, error) {
	return c.svc.Add(ctx, input.CompanyID, input.UserID, input.Role)
}

func (c *MembershipController) Remove(ctx context.Context, membershipID string) error {
	return c.svc.Remove(ctx, membershipID)
}

func (c *MembershipController) ListByCompany(ctx context.Context, companyID string) ([]*MembershipOutput, error) {
	memberships, err := c.svc.ListByCompany(ctx, companyID)
	if err != nil {
		return nil, err
	}
	out := make([]*MembershipOutput, 0, len(memberships))
	for _, m := range memberships {
		out = append(out, &MembershipOutput{ID: m.ID, CompanyID: m.CompanyID, UserID: m.UserID, Role: m.Role})
	}
	return out, nil
}

func (c *MembershipController) ListByUser(ctx context.Context, userID string) ([]*MembershipOutput, error) {
	memberships, err := c.svc.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]*MembershipOutput, 0, len(memberships))
	for _, m := range memberships {
		out = append(out, &MembershipOutput{ID: m.ID, CompanyID: m.CompanyID, UserID: m.UserID, Role: m.Role})
	}
	return out, nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

Run: `go test ./utexample/internal/membership/adapter/controller/...`
Expected: PASS

- [ ] **Step 5: Run all membership tests**

Run: `go test ./utexample/internal/membership/...`
Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add utexample/internal/membership/adapter/controller/
git commit -m "feat: add membership controller adapter"
```

---

## Task 7: Update README

**Files:**
- Modify: `README.md`

- [ ] **Step 1: Add membership to the project structure section**

In `README.md`, find the project structure block and add `membership/` alongside `company/` and `user/`:

```
utexample/
└── internal/
    ├── company/
    │   ├── domain/
    │   ├── application/
    │   ├── adapter/
    │   │   ├── controller/
    │   │   └── repo/
    │   └── interfaces/
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

- [ ] **Step 2: Add cross-domain communication section to Key Conventions**

Add after the existing conventions list:

```markdown
## Cross-Domain Communication

The `membership` domain demonstrates how domains communicate without direct coupling:

- `membership/interfaces` defines narrow `CompanyReader` and `UserReader` interfaces
- At wire-up, `*companyrepo.CompanyRepository` and `*userrepo.UserRepository` satisfy these interfaces
- Neither `company` nor `user` imports from `membership`
- `MembershipService` depends only on interfaces it owns — not on other domains' packages
```

- [ ] **Step 3: Add mock regeneration command for membership**

In the "Regenerating Mocks" section, add:

```bash
go generate ./utexample/internal/membership/interfaces/...
```

- [ ] **Step 4: Run all tests to confirm nothing is broken**

Run: `go test ./utexample/...`
Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add README.md
git commit -m "docs: update README with membership domain"
```
