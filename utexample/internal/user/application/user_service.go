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
	saved, err := s.repo.Save(ctx, u)
	if err != nil {
		return "", fmt.Errorf("save user: %w", err)
	}
	return saved.ID, nil
}

func (s *UserService) Get(ctx context.Context, id string) (*domain.User, error) {
	u, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return u, nil
}

func (s *UserService) List(ctx context.Context) ([]*domain.User, error) {
	return s.repo.List(ctx)
}

func (s *UserService) Remove(ctx context.Context, id string) error {
	if err := s.repo.Remove(ctx, id); err != nil {
		return fmt.Errorf("remove user: %w", err)
	}
	return nil
}
