# Membership Domain Design

## Goal

Add a `membership` domain to `utexample` that demonstrates cross-domain communication in clean architecture. A user can be added to or removed from a company. The membership domain coordinates between company and user without either knowing about the other.

---

## Architecture

The membership domain follows the same four-layer structure as `company` and `user`. Cross-domain access is achieved through narrow interfaces owned by the membership domain — `CompanyReader` and `UserReader` — which the other domains' repositories satisfy at wire-up time. Neither `company` nor `user` imports from `membership`.

```
membership/interfaces  defines CompanyReader, UserReader
                           ↑ satisfied by
company/adapter/repo   CompanyRepository.FindByID
user/adapter/repo      UserRepository.FindByID
```

---

## Entity

```go
type Membership struct {
    ID        string
    CompanyID string
    UserID    string
    Role      string // "member" | "admin"
}
```

`Validate()` returns an error if `CompanyID`, `UserID`, or `Role` is empty, or if `Role` is not `"member"` or `"admin"`.

---

## Errors

```go
var (
    ErrMembershipNotFound = errors.New("membership not found")
    ErrAlreadyMember      = errors.New("already a member")
    ErrInvalidRole        = errors.New("invalid role")
)
```

---

## Interfaces (`membership/interfaces/interfaces.go`)

```go
// Cross-domain: narrow read interfaces owned by membership
type CompanyReader interface {
    FindByID(ctx context.Context, id string) (*companydomain.Company, error)
}
type UserReader interface {
    FindByID(ctx context.Context, id string) (*userdomain.User, error)
}

// Membership-owned persistence
type MembershipRepository interface {
    Save(ctx context.Context, m *domain.Membership) error
    FindByID(ctx context.Context, id string) (*domain.Membership, error)
    FindByCompanyID(ctx context.Context, companyID string) ([]*domain.Membership, error)
    FindByUserID(ctx context.Context, userID string) ([]*domain.Membership, error)
    Remove(ctx context.Context, id string) error
}

type MembershipDAO interface {
    Insert(ctx context.Context, companyID, userID, role string) (string, error)
    FindByID(ctx context.Context, id string) (companyID, userID, role string, err error)
    FindByCompanyID(ctx context.Context, companyID string) ([]MembershipRow, error)
    FindByUserID(ctx context.Context, userID string) ([]MembershipRow, error)
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

---

## Service (`membership/application/membership_service.go`)

`MembershipService` takes `MembershipRepository`, `CompanyReader`, `UserReader`.

**`Add(ctx, companyID, userID, role)`**
1. Call `companyReader.FindByID` — return `ErrCompanyNotFound` (wrapped) if not found
2. Call `userReader.FindByID` — return `ErrUserNotFound` (wrapped) if not found
3. Build `Membership{CompanyID, UserID, Role}`, call `Validate()`
4. Call `repo.Save` — return error if fails
5. Return `membership.ID`

**`Remove(ctx, membershipID)`** — delegates to `repo.Remove`, wraps error.

**`ListByCompany(ctx, companyID)`** — delegates to `repo.FindByCompanyID`, wraps error.

**`ListByUser(ctx, userID)`** — delegates to `repo.FindByUserID`, wraps error.

---

## Adapter: Repo (`membership/adapter/repo/membership_repo.go`)

`membershipDoc` struct mirrors `MembershipRow`. `docToEntity` maps doc → `domain.Membership`. `MembershipRepository` wraps `MembershipDAO`.

---

## Adapter: Controller (`membership/adapter/controller/membership_controller.go`)

DTOs: `AddMemberInput{CompanyID, UserID, Role}`, `MembershipOutput{ID, CompanyID, UserID, Role}`.

`MembershipController` wraps `MembershipService`. Methods: `Add`, `Remove`, `ListByCompany`, `ListByUser`.

---

## Directory Layout

```
utexample/internal/membership/
├── domain/
│   ├── membership.go
│   └── error.go
├── application/
│   └── membership_service.go
├── adapter/
│   ├── repo/
│   │   └── membership_repo.go
│   └── controller/
│       └── membership_controller.go
└── interfaces/
    ├── interfaces.go
    └── mock/
        └── interfaces.go
```

---

## Testing

All tests follow existing conventions: `setupMembershipServiceTest(t)` deps struct, `t.Context()`, `Should{Result}_When{Condition}` naming, no `defer ctrl.Finish()`.

**Service tests** mock `MembershipRepository`, `CompanyReader`, `UserReader`.

Key cases:
- `Add` succeeds when both company and user exist
- `Add` returns `ErrCompanyNotFound` when company missing
- `Add` returns `ErrUserNotFound` when user missing
- `Add` returns `ErrInvalidRole` when role is invalid
- `Remove`, `ListByCompany`, `ListByUser` happy path + repo failure

**Repo tests** mock `MembershipDAO`.

---

## Key Conventions (unchanged from existing domains)

- `go:generate mockgen -source=interfaces.go -destination mock/interfaces.go -package=mock`
- Domain errors via `errors.Is()`, never string matching
- Error wrapping: `fmt.Errorf("context: %w", err)`
- No DB tags on domain entities
- `Validate()` called on write path only (not in `FindByID`)
