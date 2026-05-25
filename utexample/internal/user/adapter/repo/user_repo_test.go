package repo_test

import (
	"errors"
	"testing"

	"github.com/herocwhsu/training/utexample/internal/user/adapter/repo"
	"github.com/herocwhsu/training/utexample/internal/user/domain"
	"github.com/herocwhsu/training/utexample/internal/user/interfaces"
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
	t.Run("ShouldReturnSavedUser_WhenDAOInsertSucceeds", func(t *testing.T) {
		d := setupUserRepoTest(t)
		user := &domain.User{Email: "a@b.com", Name: "Alice"}
		d.dao.EXPECT().Insert(t.Context(), "a@b.com", "Alice").Return("usr_1", nil)

		got, err := d.repo.Save(t.Context(), user)
		require.NoError(t, err)
		assert.Equal(t, &domain.User{ID: "usr_1", Email: "a@b.com", Name: "Alice"}, got)
		assert.Empty(t, user.ID, "input entity must not be mutated")
	})

	t.Run("ShouldReturnError_WhenDAOInsertFails", func(t *testing.T) {
		d := setupUserRepoTest(t)
		user := &domain.User{Email: "a@b.com", Name: "Alice"}
		d.dao.EXPECT().Insert(t.Context(), "a@b.com", "Alice").Return("", errors.New("db error"))

		got, err := d.repo.Save(t.Context(), user)
		require.Error(t, err)
		assert.Nil(t, got)
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

func TestUserRepository_List(t *testing.T) {
	t.Run("ShouldReturnUsers_WhenDAOListSucceeds", func(t *testing.T) {
		d := setupUserRepoTest(t)
		rows := []*interfaces.UserRow{
			{ID: "usr_1", Email: "a@b.com", Name: "Alice"},
			{ID: "usr_2", Email: "b@b.com", Name: "Bob"},
		}
		d.dao.EXPECT().List(t.Context()).Return(rows, nil)

		got, err := d.repo.List(t.Context())
		require.NoError(t, err)
		assert.Equal(t, []*domain.User{
			{ID: "usr_1", Email: "a@b.com", Name: "Alice"},
			{ID: "usr_2", Email: "b@b.com", Name: "Bob"},
		}, got)
	})

	t.Run("ShouldReturnError_WhenDAOListFails", func(t *testing.T) {
		d := setupUserRepoTest(t)
		d.dao.EXPECT().List(t.Context()).Return(nil, errors.New("db error"))

		got, err := d.repo.List(t.Context())
		require.Error(t, err)
		assert.Nil(t, got)
	})
}
