package beans

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAccountToRelated(t *testing.T) {
	account := Account{
		ID:        NewID(),
		Name:      "Charlie",
		OffBudget: true,
	}

	assert.Equal(t,
		RelatedAccount{
			ID:        account.ID,
			Name:      "Charlie",
			OffBudget: true,
		},
		account.ToRelated(),
	)
}
