package application

import (
	"context"
	"fmt"

	"github.com/herocwhsu/training/utexample/internal/user/domain"
	"github.com/herocwhsu/training/utexample/internal/user/interfaces"
)

type UserService struct {
	repo interfaces.UserRepository
}

func NewUserService(repo interfaces.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, email, name string) (string, error) {
	u := &domain.User{Email: email, Name: name}
	if err := u.Validate(); err != nil {
		return "", err
	}
	if err := s.repo.Save(ctx, u); err != nil {
		return "", fmt.Errorf("save user: %w", err)
	}
	return u.ID, nil
}

func (s *UserService) Get(ctx context.Context, id string) (*domain.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *UserService) List(ctx context.Context) ([]*domain.User, error) {
	return s.repo.List(ctx)
}

func (s *UserService) Remove(ctx context.Context, id string) error {
	return s.repo.Remove(ctx, id)
}
