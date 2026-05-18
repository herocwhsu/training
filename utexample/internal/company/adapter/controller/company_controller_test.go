package controller_test

import (
	"errors"
	"testing"

	"github.com/herocwhsu/training/utexample/internal/company/adapter/controller"
	"github.com/herocwhsu/training/utexample/internal/company/domain"
	"github.com/herocwhsu/training/utexample/internal/company/interfaces/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type companyControllerDeps struct {
	svc  *mock.MockCompanyService
	ctrl *controller.CompanyController
}

func setupCompanyControllerTest(t *testing.T) *companyControllerDeps {
	c := gomock.NewController(t)
	svc := mock.NewMockCompanyService(c)
	return &companyControllerDeps{
		svc:  svc,
		ctrl: controller.NewCompanyController(svc),
	}
}

func TestCompanyController_Create(t *testing.T) {
	t.Run("ShouldReturnID_WhenServiceSucceeds", func(t *testing.T) {
		d := setupCompanyControllerTest(t)
		d.svc.EXPECT().Create(t.Context(), "a@b.com", "Acme").Return("cmp_1", nil)

		id, err := d.ctrl.Create(t.Context(), controller.CreateCompanyInput{Email: "a@b.com", Name: "Acme"})
		require.NoError(t, err)
		assert.Equal(t, "cmp_1", id)
	})

	t.Run("ShouldReturnError_WhenServiceFails", func(t *testing.T) {
		d := setupCompanyControllerTest(t)
		d.svc.EXPECT().Create(t.Context(), "bad@b.com", "Bad").Return("", errors.New("service error"))

		_, err := d.ctrl.Create(t.Context(), controller.CreateCompanyInput{Email: "bad@b.com", Name: "Bad"})
		require.Error(t, err)
	})
}

func TestCompanyController_Get(t *testing.T) {
	t.Run("ShouldReturnCompanyOutput_WhenFound", func(t *testing.T) {
		d := setupCompanyControllerTest(t)
		d.svc.EXPECT().Get(t.Context(), "cmp_1").Return(
			&domain.Company{ID: "cmp_1", Email: "a@b.com", Name: "Acme"}, nil,
		)

		out, err := d.ctrl.Get(t.Context(), "cmp_1")
		require.NoError(t, err)
		assert.Equal(t, &controller.CompanyOutput{ID: "cmp_1", Email: "a@b.com", Name: "Acme"}, out)
	})

	t.Run("ShouldReturnError_WhenNotFound", func(t *testing.T) {
		d := setupCompanyControllerTest(t)
		d.svc.EXPECT().Get(t.Context(), "cmp_404").Return(nil, domain.ErrCompanyNotFound)

		out, err := d.ctrl.Get(t.Context(), "cmp_404")
		require.Error(t, err)
		assert.Nil(t, out)
	})
}

func TestCompanyController_List(t *testing.T) {
	t.Run("ShouldReturnCompanyOutputs_WhenServiceSucceeds", func(t *testing.T) {
		d := setupCompanyControllerTest(t)
		d.svc.EXPECT().List(t.Context()).Return([]*domain.Company{
			{ID: "cmp_1", Email: "a@b.com", Name: "Acme"},
			{ID: "cmp_2", Email: "b@b.com", Name: "Beta"},
		}, nil)

		out, err := d.ctrl.List(t.Context())
		require.NoError(t, err)
		assert.Equal(t, []*controller.CompanyOutput{
			{ID: "cmp_1", Email: "a@b.com", Name: "Acme"},
			{ID: "cmp_2", Email: "b@b.com", Name: "Beta"},
		}, out)
	})

	t.Run("ShouldReturnError_WhenServiceFails", func(t *testing.T) {
		d := setupCompanyControllerTest(t)
		d.svc.EXPECT().List(t.Context()).Return(nil, errors.New("service error"))

		out, err := d.ctrl.List(t.Context())
		require.Error(t, err)
		assert.Nil(t, out)
	})
}

func TestCompanyController_Remove(t *testing.T) {
	t.Run("ShouldReturnNil_WhenServiceSucceeds", func(t *testing.T) {
		d := setupCompanyControllerTest(t)
		d.svc.EXPECT().Remove(t.Context(), "cmp_1").Return(nil)

		err := d.ctrl.Remove(t.Context(), "cmp_1")
		require.NoError(t, err)
	})

	t.Run("ShouldReturnError_WhenServiceFails", func(t *testing.T) {
		d := setupCompanyControllerTest(t)
		d.svc.EXPECT().Remove(t.Context(), "cmp_404").Return(errors.New("not found"))

		err := d.ctrl.Remove(t.Context(), "cmp_404")
		require.Error(t, err)
	})
}
