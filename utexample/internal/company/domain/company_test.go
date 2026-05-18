package domain_test

import (
	"testing"

	"github.com/herocwhsu/training/utexample/internal/company/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompany_Validate(t *testing.T) {
	t.Run("ShouldReturnNil_WhenAllFieldsAreValid", func(t *testing.T) {
		c := &domain.Company{ID: "cmp_1", Email: "a@b.com", Name: "Acme"}
		require.NoError(t, c.Validate())
	})

	t.Run("ShouldReturnError_WhenEmailIsEmpty", func(t *testing.T) {
		c := &domain.Company{ID: "cmp_1", Email: "", Name: "Acme"}
		err := c.Validate()
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInvalidEmail)
	})

	t.Run("ShouldReturnError_WhenNameIsEmpty", func(t *testing.T) {
		c := &domain.Company{ID: "cmp_1", Email: "a@b.com", Name: ""}
		err := c.Validate()
		require.Error(t, err)
		assert.ErrorIs(t, err, domain.ErrInvalidName)
	})
}
