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
