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
