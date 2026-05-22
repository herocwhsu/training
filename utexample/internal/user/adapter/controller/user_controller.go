package controller

import (
	"context"

	"github.com/herocwhsu/training/utexample/internal/user/interfaces"
)

type CreateUserInput struct {
	Email string
	Name  string
}

type UserOutput struct {
	ID    string
	Email string
	Name  string
}

type UserController struct {
	svc interfaces.UserService
}

func NewUserController(svc interfaces.UserService) *UserController {
	return &UserController{svc: svc}
}

func (c *UserController) Create(ctx context.Context, input CreateUserInput) (string, error) {
	return c.svc.Create(ctx, input.Email, input.Name)
}

func (c *UserController) Get(ctx context.Context, id string) (*UserOutput, error) {
	user, err := c.svc.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &UserOutput{ID: user.ID, Email: user.Email, Name: user.Name}, nil
}

func (c *UserController) List(ctx context.Context) ([]*UserOutput, error) {
	users, err := c.svc.List(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]*UserOutput, 0, len(users))
	for _, user := range users {
		out = append(out, &UserOutput{ID: user.ID, Email: user.Email, Name: user.Name})
	}
	return out, nil
}

func (c *UserController) Remove(ctx context.Context, id string) error {
	return c.svc.Remove(ctx, id)
}
