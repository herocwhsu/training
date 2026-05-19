package controller

import (
	"context"

	"github.com/herocwhsu/training/utexample/internal/membership/interfaces"
)

type AddMemberInput struct {
	CompanyID string
	UserID    string
	Role      string
}

type MembershipOutput struct {
	ID        string
	CompanyID string
	UserID    string
	Role      string
}

type MembershipController struct {
	svc interfaces.MembershipService
}

func NewMembershipController(svc interfaces.MembershipService) *MembershipController {
	return &MembershipController{svc: svc}
}

func (c *MembershipController) Add(ctx context.Context, input AddMemberInput) (string, error) {
	return c.svc.Add(ctx, input.CompanyID, input.UserID, input.Role)
}

func (c *MembershipController) Remove(ctx context.Context, membershipID string) error {
	return c.svc.Remove(ctx, membershipID)
}

func (c *MembershipController) ListByCompany(ctx context.Context, companyID string) ([]*MembershipOutput, error) {
	memberships, err := c.svc.ListByCompany(ctx, companyID)
	if err != nil {
		return nil, err
	}
	out := make([]*MembershipOutput, 0, len(memberships))
	for _, m := range memberships {
		out = append(out, &MembershipOutput{ID: m.ID, CompanyID: m.CompanyID, UserID: m.UserID, Role: m.Role})
	}
	return out, nil
}

func (c *MembershipController) ListByUser(ctx context.Context, userID string) ([]*MembershipOutput, error) {
	memberships, err := c.svc.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	out := make([]*MembershipOutput, 0, len(memberships))
	for _, m := range memberships {
		out = append(out, &MembershipOutput{ID: m.ID, CompanyID: m.CompanyID, UserID: m.UserID, Role: m.Role})
	}
	return out, nil
}
