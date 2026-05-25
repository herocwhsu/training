package interfaces

import (
	"context"

	"github.com/herocwhsu/training/utexample/internal/user/domain"
)

//go:generate mockgen -source=interfaces.go -destination mock/interfaces.go -package=mock

type UserRepository interface {
	Save(ctx context.Context, user *domain.User) (*domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	List(ctx context.Context) ([]*domain.User, error)
	Remove(ctx context.Context, id string) error
}

type UserDAO interface {
	Insert(ctx context.Context, email, name string) (id string, err error)
	FindByID(ctx context.Context, id string) (*UserRow, error)
	List(ctx context.Context) ([]*UserRow, error)
	DeleteByID(ctx context.Context, id string) error
}

type UserRow struct {
	ID    string
	Email string
	Name  string
}

type UserService interface {
	Create(ctx context.Context, email, name string) (string, error)
	Get(ctx context.Context, id string) (*domain.User, error)
	List(ctx context.Context) ([]*domain.User, error)
	Remove(ctx context.Context, id string) error
}
