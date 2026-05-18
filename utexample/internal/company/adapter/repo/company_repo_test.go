package repo_test

import (
	"errors"
	"testing"

	"github.com/herocwhsu/training/utexample/internal/company/adapter/repo"
	"github.com/herocwhsu/training/utexample/internal/company/domain"
	"github.com/herocwhsu/training/utexample/internal/company/interfaces/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type companyRepoDeps struct {
	dao  *mock.MockCompanyDAO
	repo *repo.CompanyRepository
}

func setupCompanyRepoTest(t *testing.T) *companyRepoDeps {
	ctrl := gomock.NewController(t)
	dao := mock.NewMockCompanyDAO(ctrl)
	return &companyRepoDeps{
		dao:  dao,
		repo: repo.NewCompanyRepository(dao),
	}
}

func TestCompanyRepository_Save(t *testing.T) {
	t.Run("ShouldReturnNil_WhenDAOInsertSucceeds", func(t *testing.T) {
		d := setupCompanyRepoTest(t)
		company := &domain.Company{Email: "a@b.com", Name: "Acme"}
		d.dao.EXPECT().Insert(t.Context(), "a@b.com", "Acme").Return("cmp_1", nil)

		err := d.repo.Save(t.Context(), company)
		require.NoError(t, err)
		assert.Equal(t, "cmp_1", company.ID)
	})

	t.Run("ShouldReturnError_WhenDAOInsertFails", func(t *testing.T) {
		d := setupCompanyRepoTest(t)
		company := &domain.Company{Email: "a@b.com", Name: "Acme"}
		d.dao.EXPECT().Insert(t.Context(), "a@b.com", "Acme").Return("", errors.New("db error"))

		err := d.repo.Save(t.Context(), company)
		require.Error(t, err)
	})
}

func TestCompanyRepository_FindByID(t *testing.T) {
	t.Run("ShouldReturnCompany_WhenDAOFindsRecord", func(t *testing.T) {
		d := setupCompanyRepoTest(t)
		d.dao.EXPECT().FindByID(t.Context(), "cmp_1").Return("a@b.com", "Acme", nil)

		got, err := d.repo.FindByID(t.Context(), "cmp_1")
		require.NoError(t, err)
		assert.Equal(t, &domain.Company{ID: "cmp_1", Email: "a@b.com", Name: "Acme"}, got)
	})

	t.Run("ShouldReturnError_WhenDAOReturnsNotFound", func(t *testing.T) {
		d := setupCompanyRepoTest(t)
		d.dao.EXPECT().FindByID(t.Context(), "cmp_404").Return("", "", domain.ErrCompanyNotFound)

		got, err := d.repo.FindByID(t.Context(), "cmp_404")
		require.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestCompanyRepository_Remove(t *testing.T) {
	t.Run("ShouldReturnNil_WhenDAODeleteSucceeds", func(t *testing.T) {
		d := setupCompanyRepoTest(t)
		d.dao.EXPECT().DeleteByID(t.Context(), "cmp_1").Return(nil)

		err := d.repo.Remove(t.Context(), "cmp_1")
		require.NoError(t, err)
	})
}
