package application_test

import (
	"errors"
	"testing"

	"github.com/herocwhsu/training/utexample/internal/company/application"
	"github.com/herocwhsu/training/utexample/internal/company/domain"
	"github.com/herocwhsu/training/utexample/internal/company/interfaces/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type companyServiceDeps struct {
	repo *mock.MockCompanyRepository
	svc  *application.CompanyService
}

func setupCompanyServiceTest(t *testing.T) *companyServiceDeps {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockCompanyRepository(ctrl)
	return &companyServiceDeps{
		repo: repo,
		svc:  application.NewCompanyService(repo),
	}
}

func TestCompanyService_Create(t *testing.T) {
	t.Run("ShouldReturnID_WhenInputIsValid", func(t *testing.T) {
		d := setupCompanyServiceTest(t)
		d.repo.EXPECT().
			Save(t.Context(), gomock.Any()).
			DoAndReturn(func(_ any, c *domain.Company) error {
				c.ID = "cmp_1"
				return nil
			})

		id, err := d.svc.Create(t.Context(), "a@b.com", "Acme")
		require.NoError(t, err)
		assert.Equal(t, "cmp_1", id)
	})

	t.Run("ShouldReturnError_WhenEmailIsEmpty", func(t *testing.T) {
		d := setupCompanyServiceTest(t)

		_, err := d.svc.Create(t.Context(), "", "Acme")
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInvalidEmail)
	})

	t.Run("ShouldReturnError_WhenRepoFails", func(t *testing.T) {
		d := setupCompanyServiceTest(t)
		d.repo.EXPECT().Save(t.Context(), gomock.Any()).Return(errors.New("db error"))

		_, err := d.svc.Create(t.Context(), "a@b.com", "Acme")
		require.Error(t, err)
	})
}

func TestCompanyService_Get(t *testing.T) {
	t.Run("ShouldReturnCompany_WhenFound", func(t *testing.T) {
		d := setupCompanyServiceTest(t)
		expected := &domain.Company{ID: "cmp_1", Email: "a@b.com", Name: "Acme"}
		d.repo.EXPECT().FindByID(t.Context(), "cmp_1").Return(expected, nil)

		got, err := d.svc.Get(t.Context(), "cmp_1")
		require.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("ShouldReturnError_WhenNotFound", func(t *testing.T) {
		d := setupCompanyServiceTest(t)
		d.repo.EXPECT().FindByID(t.Context(), "cmp_404").Return(nil, domain.ErrCompanyNotFound)

		got, err := d.svc.Get(t.Context(), "cmp_404")
		require.Error(t, err)
		assert.Nil(t, got)
		assert.ErrorIs(t, err, domain.ErrCompanyNotFound)
	})
}

func TestCompanyService_Remove(t *testing.T) {
	t.Run("ShouldReturnNil_WhenRemoveSucceeds", func(t *testing.T) {
		d := setupCompanyServiceTest(t)
		d.repo.EXPECT().Remove(t.Context(), "cmp_1").Return(nil)

		err := d.svc.Remove(t.Context(), "cmp_1")
		require.NoError(t, err)
	})
}
