package tests

import (
	"context"
	"testing"

	companyrepo "github.com/herocwhsu/training/utexample/internal/repo/company"
	"github.com/herocwhsu/training/utexample/mocks"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRepository_Create(t *testing.T) {
	t.Run("ShouldSuccess_WhenDAOOK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDAO := mocks.NewMockCompanyDAO(ctrl)
		repo := companyrepo.New(mockDAO)

		mockDAO.EXPECT().Insert(gomock.Any(), "a@b.com", "Acme").Return("cmp_X", nil)

		id, err := repo.Create(context.Background(), "a@b.com", "Acme")
		assert.Nil(t, err)
		assert.Equal(t, "cmp_X", id)
	})

	t.Run("ShouldFail_WhenDAOErrors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDAO := mocks.NewMockCompanyDAO(ctrl)
		repo := companyrepo.New(mockDAO)

		mockDAO.EXPECT().Insert(gomock.Any(), "bad@b.com", "Bad").Return("", context.DeadlineExceeded)

		id, err := repo.Create(context.Background(), "bad@b.com", "Bad")
		assert.Error(t, err)
		assert.Empty(t, id)
	})
}

func TestRepository_Get(t *testing.T) {
	t.Run("ShouldSuccess_WhenFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDAO := mocks.NewMockCompanyDAO(ctrl)
		repo := companyrepo.New(mockDAO)

		mockDAO.EXPECT().FindByID(gomock.Any(), "cmp_Y").Return("x@y.com", "X Corp", nil)

		ent, err := repo.Get(context.Background(), "cmp_Y")
		assert.Nil(t, err)
		assert.Equal(t, "cmp_Y", ent.ID)
		assert.Equal(t, "x@y.com", ent.Email)
		assert.Equal(t, "X Corp", ent.Name)
	})

	t.Run("ShouldFail_WhenDAOReturnsError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockDAO := mocks.NewMockCompanyDAO(ctrl)
		repo := companyrepo.New(mockDAO)

		mockDAO.EXPECT().FindByID(gomock.Any(), "cmp_404").Return("", "", context.Canceled)

		ent, err := repo.Get(context.Background(), "cmp_404")
		assert.Error(t, err)
		assert.Nil(t, ent)
	})
}
