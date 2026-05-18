package repo

import (
	"context"

	"github.com/herocwhsu/training/utexample/internal/user/domain"
	"github.com/herocwhsu/training/utexample/internal/user/interfaces"
)

type userDoc struct {
	ID    string
	Email string
	Name  string
}

func docToEntity(doc *userDoc) *domain.User {
	return &domain.User{ID: doc.ID, Email: doc.Email, Name: doc.Name}
}

func entityToDoc(u *domain.User) *userDoc {
	return &userDoc{ID: u.ID, Email: u.Email, Name: u.Name}
}

type UserRepository struct {
	dao interfaces.UserDAO
}

func NewUserRepository(dao interfaces.UserDAO) *UserRepository {
	return &UserRepository{dao: dao}
}

func (r *UserRepository) Save(ctx context.Context, user *domain.User) error {
	id, err := r.dao.Insert(ctx, user.Email, user.Name)
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*domain.User, error) {
	email, name, err := r.dao.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	doc := &userDoc{ID: id, Email: email, Name: name}
	entity := docToEntity(doc)
	if err := entity.Validate(); err != nil {
		return nil, err
	}
	return entity, nil
}

func (r *UserRepository) List(ctx context.Context) ([]*domain.User, error) {
	rows, err := r.dao.List(ctx)
	if err != nil {
		return nil, err
	}
	result := make([]*domain.User, 0, len(rows))
	for _, row := range rows {
		result = append(result, docToEntity(&userDoc{ID: row.ID, Email: row.Email, Name: row.Name}))
	}
	return result, nil
}

func (r *UserRepository) Remove(ctx context.Context, id string) error {
	return r.dao.DeleteByID(ctx, id)
}
