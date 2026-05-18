package domain_test

import (
	"testing"

	"github.com/herocwhsu/training/utexample/internal/user/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser_Validate(t *testing.T) {
	t.Run("ShouldReturnNil_WhenAllFieldsAreValid", func(t *testing.T) {
		u := &domain.User{ID: "usr_1", Email: "a@b.com", Name: "Alice"}
		require.NoError(t, u.Validate())
	})

	t.Run("ShouldReturnError_WhenEmailIsEmpty", func(t *testing.T) {
		u := &domain.User{ID: "usr_1", Email: "", Name: "Alice"}
		err := u.Validate()
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInvalidEmail)
	})

	t.Run("ShouldReturnError_WhenNameIsEmpty", func(t *testing.T) {
		u := &domain.User{ID: "usr_1", Email: "a@b.com", Name: ""}
		err := u.Validate()
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInvalidName)
	})
}
