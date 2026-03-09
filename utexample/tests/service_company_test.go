package tests

import (
	"context"
	"errors"
	"testing"

	"github.com/herocwhsu/training/utexample/internal/service/companyservice"
	"github.com/herocwhsu/training/utexample/mocks"

	"github.com/herocwhsu/training/utexample/internal/domain"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestService_CreateCompany(t *testing.T) {
	t.Run("ShouldSuccess_WhenRepoOK", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockCompanyRepository(ctrl)
		svc := companyservice.New(mockRepo)

		mockRepo.EXPECT().Create(gomock.Any(), "mail@co.com", "Co").Return("cmp_2", nil)

		id, err := svc.CreateCompany(context.Background(), "mail@co.com", "Co")
		assert.Nil(t, err)
		assert.Equal(t, "cmp_2", id)
	})

	t.Run("ShouldFail_WhenRepoErrors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockCompanyRepository(ctrl)
		svc := companyservice.New(mockRepo)

		mockRepo.EXPECT().Create(gomock.Any(), "bad@co.com", "Bad").Return("", errors.New("repo error"))

		id, err := svc.CreateCompany(context.Background(), "bad@co.com", "Bad")
		assert.Error(t, err)
		assert.Empty(t, id)
	})
}

func TestService_GetCompany(t *testing.T) {
	t.Run("ShouldSuccess_WhenRepoReturnsEntity", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockCompanyRepository(ctrl)
		svc := companyservice.New(mockRepo)

		ent := &domain.Company{ID: "cmp_3", Email: "z@z.com", Name: "Zed"}
		mockRepo.EXPECT().Get(gomock.Any(), "cmp_3").Return(ent, nil)

		got, err := svc.GetCompany(context.Background(), "cmp_3")
		assert.Nil(t, err)
		assert.Equal(t, ent, got)
	})

	t.Run("ShouldFail_WhenRepoErrors", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockRepo := mocks.NewMockCompanyRepository(ctrl)
		svc := companyservice.New(mockRepo)

		mockRepo.EXPECT().Get(gomock.Any(), "cmp_404").Return(nil, errors.New("repo error"))

		got, err := svc.GetCompany(context.Background(), "cmp_404")
		assert.Error(t, err)
		assert.Nil(t, got)
	})
}
