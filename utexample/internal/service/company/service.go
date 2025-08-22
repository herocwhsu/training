package companysvc

import (
	"context"

	"github.com/herocwhsu/training/utexample/internal/domain"
	companyrepo "github.com/herocwhsu/training/utexample/internal/repo/company"
)

//go:generate mockgen -destination=../../../mocks/mock_service.go -package=mocks github.com/herocwhsu/training/utexample/internal/service/company CompanyService

type CompanyService interface {
	CreateCompany(ctx context.Context, email, name string) (string, error)
	GetCompany(ctx context.Context, id string) (*domain.Company, error)
}

type companyService struct {
	repo companyrepo.CompanyRepository
}

func New(repo companyrepo.CompanyRepository) CompanyService {
	return &companyService{repo: repo}
}

func (s *companyService) CreateCompany(ctx context.Context, email, name string) (string, error) {
	// potential business rules before creating
	return s.repo.Create(ctx, email, name)
}

func (s *companyService) GetCompany(ctx context.Context, id string) (*domain.Company, error) {
	return s.repo.Get(ctx, id)
}
