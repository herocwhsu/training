package application_test

import (
	"errors"
	"testing"

	"github.com/herocwhsu/training/utexample/internal/user/application"
	"github.com/herocwhsu/training/utexample/internal/user/domain"
	"github.com/herocwhsu/training/utexample/internal/user/interfaces/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type userServiceDeps struct {
	repo *mock.MockUserRepository
	svc  *application.UserService
}

func setupUserServiceTest(t *testing.T) *userServiceDeps {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockUserRepository(ctrl)
	return &userServiceDeps{
		repo: repo,
		svc:  application.NewUserService(repo),
	}
}

func TestUserService_Create(t *testing.T) {
	t.Run("ShouldReturnID_WhenInputIsValid", func(t *testing.T) {
		d := setupUserServiceTest(t)
		d.repo.EXPECT().
			Save(t.Context(), gomock.Any()).
			DoAndReturn(func(_ any, u *domain.User) error {
				u.ID = "usr_1"
				return nil
			})

		id, err := d.svc.Create(t.Context(), "a@b.com", "Alice")
		require.NoError(t, err)
		assert.Equal(t, "usr_1", id)
	})

	t.Run("ShouldReturnError_WhenEmailIsEmpty", func(t *testing.T) {
		d := setupUserServiceTest(t)

		_, err := d.svc.Create(t.Context(), "", "Alice")
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInvalidEmail)
	})

	t.Run("ShouldReturnError_WhenRepoFails", func(t *testing.T) {
		d := setupUserServiceTest(t)
		d.repo.EXPECT().Save(t.Context(), gomock.Any()).Return(errors.New("db error"))

		_, err := d.svc.Create(t.Context(), "a@b.com", "Alice")
		require.Error(t, err)
	})
}

func TestUserService_Get(t *testing.T) {
	t.Run("ShouldReturnUser_WhenFound", func(t *testing.T) {
		d := setupUserServiceTest(t)
		expected := &domain.User{ID: "usr_1", Email: "a@b.com", Name: "Alice"}
		d.repo.EXPECT().FindByID(t.Context(), "usr_1").Return(expected, nil)

		got, err := d.svc.Get(t.Context(), "usr_1")
		require.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("ShouldReturnError_WhenNotFound", func(t *testing.T) {
		d := setupUserServiceTest(t)
		d.repo.EXPECT().FindByID(t.Context(), "usr_404").Return(nil, domain.ErrUserNotFound)

		got, err := d.svc.Get(t.Context(), "usr_404")
		require.Error(t, err)
		assert.Nil(t, got)
		assert.ErrorIs(t, err, domain.ErrUserNotFound)
	})
}

func TestUserService_Remove(t *testing.T) {
	t.Run("ShouldReturnNil_WhenRemoveSucceeds", func(t *testing.T) {
		d := setupUserServiceTest(t)
		d.repo.EXPECT().Remove(t.Context(), "usr_1").Return(nil)

		err := d.svc.Remove(t.Context(), "usr_1")
		require.NoError(t, err)
	})
}
