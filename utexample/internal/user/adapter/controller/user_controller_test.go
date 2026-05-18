package controller_test

import (
	"errors"
	"testing"

	"github.com/herocwhsu/training/utexample/internal/user/adapter/controller"
	"github.com/herocwhsu/training/utexample/internal/user/domain"
	"github.com/herocwhsu/training/utexample/internal/user/interfaces/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type userControllerDeps struct {
	svc  *mock.MockUserService
	ctrl *controller.UserController
}

func setupUserControllerTest(t *testing.T) *userControllerDeps {
	c := gomock.NewController(t)
	svc := mock.NewMockUserService(c)
	return &userControllerDeps{
		svc:  svc,
		ctrl: controller.NewUserController(svc),
	}
}

func TestUserController_Create(t *testing.T) {
	t.Run("ShouldReturnID_WhenServiceSucceeds", func(t *testing.T) {
		d := setupUserControllerTest(t)
		d.svc.EXPECT().Create(t.Context(), "a@b.com", "Alice").Return("usr_1", nil)

		id, err := d.ctrl.Create(t.Context(), controller.CreateUserInput{Email: "a@b.com", Name: "Alice"})
		require.NoError(t, err)
		assert.Equal(t, "usr_1", id)
	})

	t.Run("ShouldReturnError_WhenServiceFails", func(t *testing.T) {
		d := setupUserControllerTest(t)
		d.svc.EXPECT().Create(t.Context(), "bad@b.com", "Bad").Return("", errors.New("service error"))

		_, err := d.ctrl.Create(t.Context(), controller.CreateUserInput{Email: "bad@b.com", Name: "Bad"})
		require.Error(t, err)
	})
}

func TestUserController_Get(t *testing.T) {
	t.Run("ShouldReturnUserOutput_WhenFound", func(t *testing.T) {
		d := setupUserControllerTest(t)
		d.svc.EXPECT().Get(t.Context(), "usr_1").Return(
			&domain.User{ID: "usr_1", Email: "a@b.com", Name: "Alice"}, nil,
		)

		out, err := d.ctrl.Get(t.Context(), "usr_1")
		require.NoError(t, err)
		assert.Equal(t, &controller.UserOutput{ID: "usr_1", Email: "a@b.com", Name: "Alice"}, out)
	})

	t.Run("ShouldReturnError_WhenNotFound", func(t *testing.T) {
		d := setupUserControllerTest(t)
		d.svc.EXPECT().Get(t.Context(), "usr_404").Return(nil, domain.ErrUserNotFound)

		out, err := d.ctrl.Get(t.Context(), "usr_404")
		require.Error(t, err)
		assert.Nil(t, out)
	})
}

func TestUserController_List(t *testing.T) {
	t.Run("ShouldReturnUserOutputs_WhenServiceSucceeds", func(t *testing.T) {
		d := setupUserControllerTest(t)
		d.svc.EXPECT().List(t.Context()).Return([]*domain.User{
			{ID: "usr_1", Email: "a@b.com", Name: "Alice"},
			{ID: "usr_2", Email: "b@b.com", Name: "Bob"},
		}, nil)

		out, err := d.ctrl.List(t.Context())
		require.NoError(t, err)
		assert.Equal(t, []*controller.UserOutput{
			{ID: "usr_1", Email: "a@b.com", Name: "Alice"},
			{ID: "usr_2", Email: "b@b.com", Name: "Bob"},
		}, out)
	})

	t.Run("ShouldReturnError_WhenServiceFails", func(t *testing.T) {
		d := setupUserControllerTest(t)
		d.svc.EXPECT().List(t.Context()).Return(nil, errors.New("service error"))

		out, err := d.ctrl.List(t.Context())
		require.Error(t, err)
		assert.Nil(t, out)
	})
}

func TestUserController_Remove(t *testing.T) {
	t.Run("ShouldReturnNil_WhenServiceSucceeds", func(t *testing.T) {
		d := setupUserControllerTest(t)
		d.svc.EXPECT().Remove(t.Context(), "usr_1").Return(nil)

		err := d.ctrl.Remove(t.Context(), "usr_1")
		require.NoError(t, err)
	})

	t.Run("ShouldReturnError_WhenServiceFails", func(t *testing.T) {
		d := setupUserControllerTest(t)
		d.svc.EXPECT().Remove(t.Context(), "usr_404").Return(errors.New("not found"))

		err := d.ctrl.Remove(t.Context(), "usr_404")
		require.Error(t, err)
	})
}
