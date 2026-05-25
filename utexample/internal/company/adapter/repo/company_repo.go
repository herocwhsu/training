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


type CompanyRepository struct {
	dao interfaces.CompanyDAO
}

func NewCompanyRepository(dao interfaces.CompanyDAO) *CompanyRepository {
	return &CompanyRepository{dao: dao}
}

func (r *CompanyRepository) Save(ctx context.Context, company *domain.Company) (*domain.Company, error) {
	id, err := r.dao.Insert(ctx, company.Email, company.Name)
	if err != nil {
		return nil, err
	}
	return &domain.Company{ID: id, Email: company.Email, Name: company.Name}, nil
}

func (r *CompanyRepository) FindByID(ctx context.Context, id string) (*domain.Company, error) {
	email, name, err := r.dao.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	doc := &companyDoc{ID: id, Email: email, Name: name}
	return docToEntity(doc), nil
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
