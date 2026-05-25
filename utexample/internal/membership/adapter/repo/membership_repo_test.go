package repo_test

import (
	"errors"
	"testing"

	"github.com/herocwhsu/training/utexample/internal/membership/adapter/repo"
	"github.com/herocwhsu/training/utexample/internal/membership/domain"
	"github.com/herocwhsu/training/utexample/internal/membership/interfaces"
	"github.com/herocwhsu/training/utexample/internal/membership/interfaces/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type membershipRepoDeps struct {
	dao  *mock.MockMembershipDAO
	repo *repo.MembershipRepository
}

func setupMembershipRepoTest(t *testing.T) *membershipRepoDeps {
	ctrl := gomock.NewController(t)
	dao := mock.NewMockMembershipDAO(ctrl)
	return &membershipRepoDeps{
		dao:  dao,
		repo: repo.NewMembershipRepository(dao),
	}
}

func TestMembershipRepository_Save(t *testing.T) {
	t.Run("ShouldReturnSavedMembership_WhenDAOInsertSucceeds", func(t *testing.T) {
		d := setupMembershipRepoTest(t)
		m := &domain.Membership{CompanyID: "cmp_1", UserID: "usr_1", Role: "member"}
		d.dao.EXPECT().Insert(t.Context(), "cmp_1", "usr_1", "member").Return("mem_1", nil)

		got, err := d.repo.Save(t.Context(), m)
		require.NoError(t, err)
		assert.Equal(t, &domain.Membership{ID: "mem_1", CompanyID: "cmp_1", UserID: "usr_1", Role: "member"}, got)
		assert.Empty(t, m.ID, "input entity must not be mutated")
	})

	t.Run("ShouldReturnError_WhenDAOInsertFails", func(t *testing.T) {
		d := setupMembershipRepoTest(t)
		m := &domain.Membership{CompanyID: "cmp_1", UserID: "usr_1", Role: "member"}
		d.dao.EXPECT().Insert(t.Context(), "cmp_1", "usr_1", "member").Return("", errors.New("db error"))

		got, err := d.repo.Save(t.Context(), m)
		require.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestMembershipRepository_FindByID(t *testing.T) {
	t.Run("ShouldReturnMembership_WhenDAOFindsRecord", func(t *testing.T) {
		d := setupMembershipRepoTest(t)
		d.dao.EXPECT().FindByID(t.Context(), "mem_1").Return("cmp_1", "usr_1", "member", nil)

		got, err := d.repo.FindByID(t.Context(), "mem_1")
		require.NoError(t, err)
		assert.Equal(t, &domain.Membership{ID: "mem_1", CompanyID: "cmp_1", UserID: "usr_1", Role: "member"}, got)
	})

	t.Run("ShouldReturnError_WhenDAOReturnsNotFound", func(t *testing.T) {
		d := setupMembershipRepoTest(t)
		d.dao.EXPECT().FindByID(t.Context(), "mem_404").Return("", "", "", domain.ErrMembershipNotFound)

		got, err := d.repo.FindByID(t.Context(), "mem_404")
		require.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestMembershipRepository_FindByCompanyID(t *testing.T) {
	t.Run("ShouldReturnMemberships_WhenDAOSucceeds", func(t *testing.T) {
		d := setupMembershipRepoTest(t)
		rows := []interfaces.MembershipRow{
			{ID: "mem_1", CompanyID: "cmp_1", UserID: "usr_1", Role: "member"},
			{ID: "mem_2", CompanyID: "cmp_1", UserID: "usr_2", Role: "admin"},
		}
		d.dao.EXPECT().FindByCompanyID(t.Context(), "cmp_1").Return(rows, nil)

		got, err := d.repo.FindByCompanyID(t.Context(), "cmp_1")
		require.NoError(t, err)
		assert.Equal(t, []*domain.Membership{
			{ID: "mem_1", CompanyID: "cmp_1", UserID: "usr_1", Role: "member"},
			{ID: "mem_2", CompanyID: "cmp_1", UserID: "usr_2", Role: "admin"},
		}, got)
	})

	t.Run("ShouldReturnError_WhenDAOFails", func(t *testing.T) {
		d := setupMembershipRepoTest(t)
		d.dao.EXPECT().FindByCompanyID(t.Context(), "cmp_1").Return(nil, errors.New("db error"))

		got, err := d.repo.FindByCompanyID(t.Context(), "cmp_1")
		require.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestMembershipRepository_FindByUserID(t *testing.T) {
	t.Run("ShouldReturnMemberships_WhenDAOSucceeds", func(t *testing.T) {
		d := setupMembershipRepoTest(t)
		rows := []interfaces.MembershipRow{
			{ID: "mem_1", CompanyID: "cmp_1", UserID: "usr_1", Role: "member"},
		}
		d.dao.EXPECT().FindByUserID(t.Context(), "usr_1").Return(rows, nil)

		got, err := d.repo.FindByUserID(t.Context(), "usr_1")
		require.NoError(t, err)
		assert.Equal(t, []*domain.Membership{
			{ID: "mem_1", CompanyID: "cmp_1", UserID: "usr_1", Role: "member"},
		}, got)
	})

	t.Run("ShouldReturnError_WhenDAOFails", func(t *testing.T) {
		d := setupMembershipRepoTest(t)
		d.dao.EXPECT().FindByUserID(t.Context(), "usr_1").Return(nil, errors.New("db error"))

		got, err := d.repo.FindByUserID(t.Context(), "usr_1")
		require.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestMembershipRepository_ExistsByCompanyAndUser(t *testing.T) {
	t.Run("ShouldReturnTrue_WhenDAOReturnsTrue", func(t *testing.T) {
		d := setupMembershipRepoTest(t)
		d.dao.EXPECT().ExistsByCompanyAndUser(t.Context(), "cmp_1", "usr_1").Return(true, nil)

		exists, err := d.repo.ExistsByCompanyAndUser(t.Context(), "cmp_1", "usr_1")
		require.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("ShouldReturnError_WhenDAOFails", func(t *testing.T) {
		d := setupMembershipRepoTest(t)
		d.dao.EXPECT().ExistsByCompanyAndUser(t.Context(), "cmp_1", "usr_1").Return(false, errors.New("db error"))

		_, err := d.repo.ExistsByCompanyAndUser(t.Context(), "cmp_1", "usr_1")
		require.Error(t, err)
	})
}

func TestMembershipRepository_Remove(t *testing.T) {
	t.Run("ShouldReturnNil_WhenDAODeleteSucceeds", func(t *testing.T) {
		d := setupMembershipRepoTest(t)
		d.dao.EXPECT().DeleteByID(t.Context(), "mem_1").Return(nil)

		err := d.repo.Remove(t.Context(), "mem_1")
		require.NoError(t, err)
	})

	t.Run("ShouldReturnError_WhenDAODeleteFails", func(t *testing.T) {
		d := setupMembershipRepoTest(t)
		d.dao.EXPECT().DeleteByID(t.Context(), "mem_1").Return(errors.New("db error"))

		err := d.repo.Remove(t.Context(), "mem_1")
		require.Error(t, err)
	})
}
