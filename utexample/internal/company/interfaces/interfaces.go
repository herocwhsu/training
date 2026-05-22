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

type CompanyService interface {
	Create(ctx context.Context, email, name string) (string, error)
	Get(ctx context.Context, id string) (*domain.Company, error)
	List(ctx context.Context) ([]*domain.Company, error)
	Remove(ctx context.Context, id string) error
}
