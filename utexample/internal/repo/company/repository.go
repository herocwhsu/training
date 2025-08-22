package companyrepo

import (
	"context"

	companydao "github.com/herocwhsu/training/utexample/internal/dao/company"
	"github.com/herocwhsu/training/utexample/internal/domain"
)

//go:generate mockgen -destination=../../../mocks/mock_repository.go -package=mocks github.com/herocwhsu/training/utexample/internal/repo/company CompanyRepository

type CompanyRepository interface {
	Create(ctx context.Context, email, name string) (string, error)
	Get(ctx context.Context, id string) (*domain.Company, error)
}

type companyRepository struct {
	dao companydao.CompanyDAO
}

func New(dao companydao.CompanyDAO) CompanyRepository {
	return &companyRepository{dao: dao}
}

func (r *companyRepository) Create(ctx context.Context, email, name string) (string, error) {
	return r.dao.Insert(ctx, email, name)
}

func (r *companyRepository) Get(ctx context.Context, id string) (*domain.Company, error) {
	email, name, err := r.dao.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	c := &domain.Company{ID: id, Email: email, Name: name}
	if err := c.Validate(); err != nil {
		return nil, err
	}
	return c, nil
}
