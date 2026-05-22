package controller_test

import (
	"errors"
	"testing"

	"github.com/herocwhsu/training/utexample/internal/membership/adapter/controller"
	"github.com/herocwhsu/training/utexample/internal/membership/domain"
	"github.com/herocwhsu/training/utexample/internal/membership/interfaces/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type membershipControllerDeps struct {
	svc  *mock.MockMembershipService
	ctrl *controller.MembershipController
}

func setupMembershipControllerTest(t *testing.T) *membershipControllerDeps {
	c := gomock.NewController(t)
	svc := mock.NewMockMembershipService(c)
	return &membershipControllerDeps{
		svc:  svc,
		ctrl: controller.NewMembershipController(svc),
	}
}

func TestMembershipController_Add(t *testing.T) {
	t.Run("ShouldReturnID_WhenServiceSucceeds", func(t *testing.T) {
		d := setupMembershipControllerTest(t)
		d.svc.EXPECT().Add(t.Context(), "cmp_1", "usr_1", "member").Return("mem_1", nil)

		id, err := d.ctrl.Add(t.Context(), controller.AddMemberInput{CompanyID: "cmp_1", UserID: "usr_1", Role: "member"})
		require.NoError(t, err)
		assert.Equal(t, "mem_1", id)
	})

	t.Run("ShouldReturnError_WhenServiceFails", func(t *testing.T) {
		d := setupMembershipControllerTest(t)
		d.svc.EXPECT().Add(t.Context(), "cmp_1", "usr_1", "member").Return("", errors.New("service error"))

		_, err := d.ctrl.Add(t.Context(), controller.AddMemberInput{CompanyID: "cmp_1", UserID: "usr_1", Role: "member"})
		require.Error(t, err)
	})
}

func TestMembershipController_Remove(t *testing.T) {
	t.Run("ShouldReturnNil_WhenServiceSucceeds", func(t *testing.T) {
		d := setupMembershipControllerTest(t)
		d.svc.EXPECT().Remove(t.Context(), "mem_1").Return(nil)

		err := d.ctrl.Remove(t.Context(), "mem_1")
		require.NoError(t, err)
	})

	t.Run("ShouldReturnError_WhenServiceFails", func(t *testing.T) {
		d := setupMembershipControllerTest(t)
		d.svc.EXPECT().Remove(t.Context(), "mem_1").Return(errors.New("service error"))

		err := d.ctrl.Remove(t.Context(), "mem_1")
		require.Error(t, err)
	})
}

func TestMembershipController_ListByCompany(t *testing.T) {
	t.Run("ShouldReturnOutputs_WhenServiceSucceeds", func(t *testing.T) {
		d := setupMembershipControllerTest(t)
		d.svc.EXPECT().ListByCompany(t.Context(), "cmp_1").Return([]*domain.Membership{
			{ID: "mem_1", CompanyID: "cmp_1", UserID: "usr_1", Role: "member"},
		}, nil)

		out, err := d.ctrl.ListByCompany(t.Context(), "cmp_1")
		require.NoError(t, err)
		assert.Equal(t, []*controller.MembershipOutput{
			{ID: "mem_1", CompanyID: "cmp_1", UserID: "usr_1", Role: "member"},
		}, out)
	})

	t.Run("ShouldReturnError_WhenServiceFails", func(t *testing.T) {
		d := setupMembershipControllerTest(t)
		d.svc.EXPECT().ListByCompany(t.Context(), "cmp_1").Return(nil, errors.New("service error"))

		out, err := d.ctrl.ListByCompany(t.Context(), "cmp_1")
		require.Error(t, err)
		assert.Nil(t, out)
	})
}

func TestMembershipController_ListByUser(t *testing.T) {
	t.Run("ShouldReturnOutputs_WhenServiceSucceeds", func(t *testing.T) {
		d := setupMembershipControllerTest(t)
		d.svc.EXPECT().ListByUser(t.Context(), "usr_1").Return([]*domain.Membership{
			{ID: "mem_1", CompanyID: "cmp_1", UserID: "usr_1", Role: "member"},
		}, nil)

		out, err := d.ctrl.ListByUser(t.Context(), "usr_1")
		require.NoError(t, err)
		assert.Equal(t, []*controller.MembershipOutput{
			{ID: "mem_1", CompanyID: "cmp_1", UserID: "usr_1", Role: "member"},
		}, out)
	})

	t.Run("ShouldReturnError_WhenServiceFails", func(t *testing.T) {
		d := setupMembershipControllerTest(t)
		d.svc.EXPECT().ListByUser(t.Context(), "usr_1").Return(nil, errors.New("service error"))

		out, err := d.ctrl.ListByUser(t.Context(), "usr_1")
		require.Error(t, err)
		assert.Nil(t, out)
	})
}
