package companyctl

import (
	"context"

	"github.com/herocwhsu/training/utexample/internal/service/companyservice"
)

type CompanyInfo struct {
	ID    string
	Email string
	Name  string
}

type Controller struct {
	svc companyservice.CompanyService
}

func New(svc companyservice.CompanyService) *Controller {
	return &Controller{svc: svc}
}

func (c *Controller) CreateCompany(ctx context.Context, input CompanyInfo) (string, error) {
	// map input DTO to service call
	return c.svc.CreateCompany(ctx, input.Email, input.Name)
}

func (c *Controller) GetCompany(ctx context.Context, companyID string) (*CompanyInfo, error) {
	ent, err := c.svc.GetCompany(ctx, companyID)
	if err != nil {
		return nil, err
	}
	return &CompanyInfo{ID: ent.ID, Email: ent.Email, Name: ent.Name}, nil
}
