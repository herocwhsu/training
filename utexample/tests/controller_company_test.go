package tests

import (
	"context"
	"testing"

	"github.com/herocwhsu/training/utexample/mocks"

	"github.com/golang/mock/gomock"
	"github.com/herocwhsu/training/utexample/internal/controller/companyctl"
	"github.com/herocwhsu/training/utexample/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestController_CreateCompany(t *testing.T) {
	t.Run("ShouldSuccess_WhenInputIsValid", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSvc := mocks.NewMockCompanyService(ctrl)
		c := companyctl.New(mockSvc)

		in := companyctl.CompanyInfo{Email: "a@b.com", Name: "Acme"}
		mockSvc.EXPECT().CreateCompany(gomock.Any(), "a@b.com", "Acme").Return("cmp_1", nil)

		id, err := c.CreateCompany(context.Background(), in)
		assert.Nil(t, err)
		assert.Equal(t, "cmp_1", id)
	})

	t.Run("ShouldFail_WhenServiceReturnsError", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSvc := mocks.NewMockCompanyService(ctrl)
		c := companyctl.New(mockSvc)

		in := companyctl.CompanyInfo{Email: "bad@b.com", Name: "BadCo"}
		mockSvc.EXPECT().CreateCompany(gomock.Any(), "bad@b.com", "BadCo").Return("", context.DeadlineExceeded)

		id, err := c.CreateCompany(context.Background(), in)
		assert.Error(t, err)
		assert.Empty(t, id)
	})
}

func TestController_GetCompany(t *testing.T) {
	t.Run("ShouldSuccess_WhenFound", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSvc := mocks.NewMockCompanyService(ctrl)
		c := companyctl.New(mockSvc)

		ent := &domain.Company{ID: "cmp_2", Email: "x@y.com", Name: "X Corp"}
		mockSvc.EXPECT().GetCompany(gomock.Any(), "cmp_2").Return(ent, nil)

		out, err := c.GetCompany(context.Background(), "cmp_2")
		assert.Nil(t, err)
		assert.Equal(t, &companyctl.CompanyInfo{Email: "x@y.com", Name: "X Corp"}, out)
	})

	t.Run("ShouldFail_WhenServiceErrors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockSvc := mocks.NewMockCompanyService(ctrl)
		c := companyctl.New(mockSvc)

		mockSvc.EXPECT().GetCompany(gomock.Any(), "cmp_404").Return(nil, context.Canceled)

		out, err := c.GetCompany(context.Background(), "cmp_404")
		assert.Error(t, err)
		assert.Nil(t, out)
	})
}
