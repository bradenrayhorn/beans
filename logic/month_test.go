package logic_test

import (
	"context"
	"testing"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/mocks"
	"github.com/bradenrayhorn/beans/logic"
	"github.com/stretchr/testify/assert"
)

func TestGetMonth(t *testing.T) {
	monthRepository := mocks.NewMockMonthRepository()
	svc := logic.NewMonthService(monthRepository)

	t.Run("gives not found error for wrong budget", func(t *testing.T) {
		monthRepository.GetFunc.SetDefaultReturn(&beans.Month{ID: beans.NewBeansID(), BudgetID: beans.NewBeansID()}, nil)
		_, err := svc.Get(context.Background(), beans.NewBeansID(), beans.NewBeansID())
		assert.ErrorIs(t, err, beans.ErrorNotFound)
	})
}
