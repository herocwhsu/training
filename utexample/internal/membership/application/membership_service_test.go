package application_test

import (
	"context"
	"errors"
	"testing"

	"github.com/herocwhsu/training/utexample/internal/membership/application"
	companydomain "github.com/herocwhsu/training/utexample/internal/company/domain"
	"github.com/herocwhsu/training/utexample/internal/membership/domain"
	userdomain "github.com/herocwhsu/training/utexample/internal/user/domain"
	"github.com/herocwhsu/training/utexample/internal/membership/interfaces/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type membershipServiceDeps struct {
	repo          *mock.MockMembershipRepository
	companyReader *mock.MockCompanyReader
	userReader    *mock.MockUserReader
	svc           *application.MembershipService
}

func setupMembershipServiceTest(t *testing.T) *membershipServiceDeps {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockMembershipRepository(ctrl)
	companyReader := mock.NewMockCompanyReader(ctrl)
	userReader := mock.NewMockUserReader(ctrl)
	return &membershipServiceDeps{
		repo:          repo,
		companyReader: companyReader,
		userReader:    userReader,
		svc:           application.NewMembershipService(repo, companyReader, userReader),
	}
}

func TestMembershipService_Add(t *testing.T) {
	t.Run("ShouldReturnID_WhenBothExistAndNotAlreadyMember", func(t *testing.T) {
		d := setupMembershipServiceTest(t)
		d.companyReader.EXPECT().FindByID(t.Context(), "cmp_1").Return(&companydomain.Company{ID: "cmp_1"}, nil)
		d.userReader.EXPECT().FindByID(t.Context(), "usr_1").Return(&userdomain.User{ID: "usr_1"}, nil)
		d.repo.EXPECT().ExistsByCompanyAndUser(t.Context(), "cmp_1", "usr_1").Return(false, nil)
		d.repo.EXPECT().
			Save(t.Context(), gomock.Any()).
			DoAndReturn(func(_ context.Context, m *domain.Membership) (*domain.Membership, error) {
				return &domain.Membership{ID: "mem_1", CompanyID: m.CompanyID, UserID: m.UserID, Role: m.Role}, nil
			})

		id, err := d.svc.Add(t.Context(), "cmp_1", "usr_1", "member")
		require.NoError(t, err)
		assert.Equal(t, "mem_1", id)
	})

	t.Run("ShouldReturnError_WhenCompanyNotFound", func(t *testing.T) {
		d := setupMembershipServiceTest(t)
		d.companyReader.EXPECT().FindByID(t.Context(), "cmp_404").Return(nil, companydomain.ErrCompanyNotFound)

		_, err := d.svc.Add(t.Context(), "cmp_404", "usr_1", "member")
		require.Error(t, err)
		assert.ErrorIs(t, err, companydomain.ErrCompanyNotFound)
	})

	t.Run("ShouldReturnError_WhenUserNotFound", func(t *testing.T) {
		d := setupMembershipServiceTest(t)
		d.companyReader.EXPECT().FindByID(t.Context(), "cmp_1").Return(&companydomain.Company{ID: "cmp_1"}, nil)
		d.userReader.EXPECT().FindByID(t.Context(), "usr_404").Return(nil, userdomain.ErrUserNotFound)

		_, err := d.svc.Add(t.Context(), "cmp_1", "usr_404", "member")
		require.Error(t, err)
		assert.ErrorIs(t, err, userdomain.ErrUserNotFound)
	})

	t.Run("ShouldReturnError_WhenAlreadyMember", func(t *testing.T) {
		d := setupMembershipServiceTest(t)
		d.companyReader.EXPECT().FindByID(t.Context(), "cmp_1").Return(&companydomain.Company{ID: "cmp_1"}, nil)
		d.userReader.EXPECT().FindByID(t.Context(), "usr_1").Return(&userdomain.User{ID: "usr_1"}, nil)
		d.repo.EXPECT().ExistsByCompanyAndUser(t.Context(), "cmp_1", "usr_1").Return(true, nil)

		_, err := d.svc.Add(t.Context(), "cmp_1", "usr_1", "member")
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrAlreadyMember)
	})

	t.Run("ShouldReturnError_WhenRoleIsInvalid", func(t *testing.T) {
		d := setupMembershipServiceTest(t)
		d.companyReader.EXPECT().FindByID(t.Context(), "cmp_1").Return(&companydomain.Company{ID: "cmp_1"}, nil)
		d.userReader.EXPECT().FindByID(t.Context(), "usr_1").Return(&userdomain.User{ID: "usr_1"}, nil)
		d.repo.EXPECT().ExistsByCompanyAndUser(t.Context(), "cmp_1", "usr_1").Return(false, nil)

		_, err := d.svc.Add(t.Context(), "cmp_1", "usr_1", "superuser")
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInvalidRole)
	})

	t.Run("ShouldReturnError_WhenExistsCheckFails", func(t *testing.T) {
		d := setupMembershipServiceTest(t)
		d.companyReader.EXPECT().FindByID(t.Context(), "cmp_1").Return(&companydomain.Company{ID: "cmp_1"}, nil)
		d.userReader.EXPECT().FindByID(t.Context(), "usr_1").Return(&userdomain.User{ID: "usr_1"}, nil)
		d.repo.EXPECT().ExistsByCompanyAndUser(t.Context(), "cmp_1", "usr_1").Return(false, errors.New("db error"))

		_, err := d.svc.Add(t.Context(), "cmp_1", "usr_1", "member")
		require.Error(t, err)
	})

	t.Run("ShouldReturnError_WhenSaveFails", func(t *testing.T) {
		d := setupMembershipServiceTest(t)
		d.companyReader.EXPECT().FindByID(t.Context(), "cmp_1").Return(&companydomain.Company{ID: "cmp_1"}, nil)
		d.userReader.EXPECT().FindByID(t.Context(), "usr_1").Return(&userdomain.User{ID: "usr_1"}, nil)
		d.repo.EXPECT().ExistsByCompanyAndUser(t.Context(), "cmp_1", "usr_1").Return(false, nil)
		d.repo.EXPECT().Save(t.Context(), gomock.Any()).Return(nil, errors.New("db error"))

		_, err := d.svc.Add(t.Context(), "cmp_1", "usr_1", "member")
		require.Error(t, err)
	})
}

func TestMembershipService_Remove(t *testing.T) {
	t.Run("ShouldReturnNil_WhenRemoveSucceeds", func(t *testing.T) {
		d := setupMembershipServiceTest(t)
		d.repo.EXPECT().Remove(t.Context(), "mem_1").Return(nil)

		err := d.svc.Remove(t.Context(), "mem_1")
		require.NoError(t, err)
	})

	t.Run("ShouldReturnError_WhenRepoFails", func(t *testing.T) {
		d := setupMembershipServiceTest(t)
		d.repo.EXPECT().Remove(t.Context(), "mem_1").Return(errors.New("db error"))

		err := d.svc.Remove(t.Context(), "mem_1")
		require.Error(t, err)
	})
}

func TestMembershipService_ListByCompany(t *testing.T) {
	t.Run("ShouldReturnMemberships_WhenRepoSucceeds", func(t *testing.T) {
		d := setupMembershipServiceTest(t)
		expected := []*domain.Membership{
			{ID: "mem_1", CompanyID: "cmp_1", UserID: "usr_1", Role: "member"},
		}
		d.repo.EXPECT().FindByCompanyID(t.Context(), "cmp_1").Return(expected, nil)

		got, err := d.svc.ListByCompany(t.Context(), "cmp_1")
		require.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("ShouldReturnError_WhenRepoFails", func(t *testing.T) {
		d := setupMembershipServiceTest(t)
		d.repo.EXPECT().FindByCompanyID(t.Context(), "cmp_1").Return(nil, errors.New("db error"))

		got, err := d.svc.ListByCompany(t.Context(), "cmp_1")
		require.Error(t, err)
		assert.Nil(t, got)
	})
}

func TestMembershipService_ListByUser(t *testing.T) {
	t.Run("ShouldReturnMemberships_WhenRepoSucceeds", func(t *testing.T) {
		d := setupMembershipServiceTest(t)
		expected := []*domain.Membership{
			{ID: "mem_1", CompanyID: "cmp_1", UserID: "usr_1", Role: "member"},
		}
		d.repo.EXPECT().FindByUserID(t.Context(), "usr_1").Return(expected, nil)

		got, err := d.svc.ListByUser(t.Context(), "usr_1")
		require.NoError(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("ShouldReturnError_WhenRepoFails", func(t *testing.T) {
		d := setupMembershipServiceTest(t)
		d.repo.EXPECT().FindByUserID(t.Context(), "usr_1").Return(nil, errors.New("db error"))

		got, err := d.svc.ListByUser(t.Context(), "usr_1")
		require.Error(t, err)
		assert.Nil(t, got)
	})
}
