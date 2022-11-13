package logic_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/mocks"
	"github.com/bradenrayhorn/beans/internal/testutils"
	"github.com/bradenrayhorn/beans/logic"
	mockassert "github.com/derision-test/go-mockgen/testutil/assert"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrUpdate(t *testing.T) {

	existing := &beans.MonthCategory{
		ID:         beans.NewBeansID(),
		MonthID:    beans.NewBeansID(),
		CategoryID: beans.NewBeansID(),
		Amount:     beans.NewAmount(5, -1),
	}

	t.Run("creates if month category does not exist", func(t *testing.T) {
		repository := mocks.NewMockMonthCategoryRepository()
		svc := logic.NewMonthCategoryService(repository)
		repository.GetByMonthAndCategoryFunc.SetDefaultReturn(nil, beans.ErrorNotFound)

		err := svc.CreateOrUpdate(context.Background(), beans.NewBeansID(), beans.NewBeansID(), beans.NewAmount(5, 0))
		assert.Nil(t, err)
		mockassert.NotCalled(t, repository.UpdateAmountFunc)
		mockassert.CalledOnce(t, repository.CreateFunc)
	})

	t.Run("updates existing month category", func(t *testing.T) {
		repository := mocks.NewMockMonthCategoryRepository()
		svc := logic.NewMonthCategoryService(repository)
		repository.GetByMonthAndCategoryFunc.SetDefaultReturn(existing, nil)

		err := svc.CreateOrUpdate(context.Background(), beans.NewBeansID(), beans.NewBeansID(), beans.NewAmount(5, 0))
		assert.Nil(t, err)
		mockassert.CalledOnce(t, repository.UpdateAmountFunc)
		mockassert.NotCalled(t, repository.CreateFunc)
	})

	t.Run("does not update if failed to fetch", func(t *testing.T) {
		repository := mocks.NewMockMonthCategoryRepository()
		svc := logic.NewMonthCategoryService(repository)
		repository.GetByMonthAndCategoryFunc.SetDefaultReturn(nil, errors.New("no"))

		err := svc.CreateOrUpdate(context.Background(), beans.NewBeansID(), beans.NewBeansID(), beans.NewAmount(5, 0))
		assert.NotNil(t, err)
		mockassert.NotCalled(t, repository.UpdateAmountFunc)
		mockassert.NotCalled(t, repository.CreateFunc)
	})

	t.Run("cannot set amount to negative", func(t *testing.T) {
		repository := mocks.NewMockMonthCategoryRepository()
		svc := logic.NewMonthCategoryService(repository)
		repository.GetByMonthAndCategoryFunc.SetDefaultReturn(nil, beans.ErrorNotFound)

		err := svc.CreateOrUpdate(context.Background(), beans.NewBeansID(), beans.NewBeansID(), beans.NewAmount(-5, 0))
		testutils.AssertError(t, err, "Amount must be positive.")
		mockassert.NotCalled(t, repository.UpdateAmountFunc)
		mockassert.NotCalled(t, repository.CreateFunc)
	})

	t.Run("cannot set amount to zero", func(t *testing.T) {
		repository := mocks.NewMockMonthCategoryRepository()
		svc := logic.NewMonthCategoryService(repository)
		repository.GetByMonthAndCategoryFunc.SetDefaultReturn(nil, beans.ErrorNotFound)

		err := svc.CreateOrUpdate(context.Background(), beans.NewBeansID(), beans.NewBeansID(), beans.NewAmount(0, 0))
		testutils.AssertError(t, err, "Amount must not be zero.")
		mockassert.NotCalled(t, repository.UpdateAmountFunc)
		mockassert.NotCalled(t, repository.CreateFunc)
	})
}
