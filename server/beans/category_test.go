package beans

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCategoryToRelated(t *testing.T) {
	category := Category{
		ID:       NewID(),
		BudgetID: NewID(),
		GroupID:  NewID(),
		Name:     "Charlie",
	}

	assert.Equal(t,
		RelatedCategory{
			ID:   category.ID,
			Name: "Charlie",
		},
		category.ToRelated(),
	)
}
