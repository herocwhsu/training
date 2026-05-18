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
	c, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find company by id: %w", err)
	}
	return c, nil
}

func (s *CompanyService) List(ctx context.Context) ([]*domain.Company, error) {
	return s.repo.List(ctx)
}

func (s *CompanyService) Remove(ctx context.Context, id string) error {
	if err := s.repo.Remove(ctx, id); err != nil {
		return fmt.Errorf("remove company: %w", err)
	}
	return nil
}
