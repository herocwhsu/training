package repo_test

import (
	"errors"
	"testing"

	"github.com/herocwhsu/training/utexample/internal/user/adapter/repo"
	"github.com/herocwhsu/training/utexample/internal/user/domain"
	"github.com/herocwhsu/training/utexample/internal/user/interfaces/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type userRepoDeps struct {
	dao  *mock.MockUserDAO
	repo *repo.UserRepository
}

func setupUserRepoTest(t *testing.T) *userRepoDeps {
	ctrl := gomock.NewController(t)
	dao := mock.NewMockUserDAO(ctrl)
	return &userRepoDeps{
		dao:  dao,
		repo: repo.NewUserRepository(dao),
	}
}

func TestUserRepository_Save(t *testing.T) {
	t.Run("ShouldSetID_WhenDAOInsertSucceeds", func(t *testing.T) {
		d := setupUserRepoTest(t)
		user := &domain.User{Email: "a@b.com", Name: "Alice"}
		d.dao.EXPECT().Insert(t.Context(), "a@b.com", "Alice").Return("usr_1", nil)

		err := d.repo.Save(t.Context(), user)
		require.NoError(t, err)
		assert.Equal(t, "usr_1", user.ID)
	})

	t.Run("ShouldReturnError_WhenDAOInsertFails", func(t *testing.T) {
		d := setupUserRepoTest(t)
		user := &domain.User{Email: "a@b.com", Name: "Alice"}
		d.dao.EXPECT().Insert(t.Context(), "a@b.com", "Alice").Return("", errors.New("db error"))

		err := d.repo.Save(t.Context(), user)
		require.Error(t, err)
	})
}

func TestUserRepository_FindByID(t *testing.T) {
	t.Run("ShouldReturnUser_WhenDAOFindsRecord", func(t *testing.T) {
		d := setupUserRepoTest(t)
		d.dao.EXPECT().FindByID(t.Context(), "usr_1").Return("a@b.com", "Alice", nil)

		got, err := d.repo.FindByID(t.Context(), "usr_1")
		require.NoError(t, err)
		assert.Equal(t, &domain.User{ID: "usr_1", Email: "a@b.com", Name: "Alice"}, got)
	})

	t.Run("ShouldReturnError_WhenDAOReturnsNotFound", func(t *testing.T) {
		d := setupUserRepoTest(t)
		d.dao.EXPECT().FindByID(t.Context(), "usr_404").Return("", "", domain.ErrUserNotFound)

		got, err := d.repo.FindByID(t.Context(), "usr_404")
		require.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestUserRepository_Remove(t *testing.T) {
	t.Run("ShouldReturnNil_WhenDAODeleteSucceeds", func(t *testing.T) {
		d := setupUserRepoTest(t)
		d.dao.EXPECT().DeleteByID(t.Context(), "usr_1").Return(nil)

		err := d.repo.Remove(t.Context(), "usr_1")
		require.NoError(t, err)
	})
}
