package companyctl

import (
	"context"

	"github.com/herocwhsu/training/utexample/internal/service/companysvc"
)

type CompanyInfo struct {
	Email string
	Name  string
}

type Controller struct {
	svc companysvc.CompanyService
}

func New(svc companysvc.CompanyService) *Controller {
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
	return &CompanyInfo{Email: ent.Email, Name: ent.Name}, nil
}
