package tests

import (
	"testing"

	"github.com/herocwhsu/training/utexample/internal/domain"

	"github.com/stretchr/testify/assert"
)

func TestCompany_Validate(t *testing.T) {
	t.Run("ShouldSuccess_WhenValid", func(t *testing.T) {
		c := &domain.Company{ID: "cmp_1", Email: "ok@co.com", Name: "OK Co"}
		err := c.Validate()
		assert.Nil(t, err)
	})

	t.Run("ShouldFail_WhenEmailMissing", func(t *testing.T) {
		c := &domain.Company{ID: "cmp_2", Email: "", Name: "No Email Co"}
		err := c.Validate()
		assert.Error(t, err)
	})
}
