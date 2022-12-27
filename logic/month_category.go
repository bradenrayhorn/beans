package logic

import (
	"context"
	"errors"

	"github.com/bradenrayhorn/beans/beans"
)

type monthCategoryService struct {
	monthCategoryRepository beans.MonthCategoryRepository
}

func NewMonthCategoryService(mcr beans.MonthCategoryRepository) *monthCategoryService {
	return &monthCategoryService{monthCategoryRepository: mcr}
}

func (s *monthCategoryService) CreateOrUpdate(ctx context.Context, monthID beans.ID, categoryID beans.ID, amount beans.Amount) error {
	if err := beans.ValidateFields(
		beans.Field("Amount", beans.NonZero(amount), beans.Positive(amount)),
	); err != nil {
		return err
	}

	monthCategory, err := s.monthCategoryRepository.GetByMonthAndCategory(ctx, monthID, categoryID)
	if err == nil {
		return s.monthCategoryRepository.UpdateAmount(ctx, monthCategory.ID, amount)
	}

	if errors.Is(err, beans.ErrorNotFound) {
		monthCategory = &beans.MonthCategory{
			ID:         beans.NewBeansID(),
			MonthID:    monthID,
			CategoryID: categoryID,
			Amount:     amount,
		}

		return s.monthCategoryRepository.Create(ctx, monthCategory)
	}

	return err
}

func (s *monthCategoryService) CreateIfNotExists(ctx context.Context, monthID beans.ID, categoryID beans.ID) error {
	monthCategory, err := s.monthCategoryRepository.GetByMonthAndCategory(ctx, monthID, categoryID)
	if err == nil {
		return nil
	}

	if errors.Is(err, beans.ErrorNotFound) {
		monthCategory = &beans.MonthCategory{
			ID:         beans.NewBeansID(),
			MonthID:    monthID,
			CategoryID: categoryID,
			Amount:     beans.NewAmount(0, 0),
		}

		return s.monthCategoryRepository.Create(ctx, monthCategory)
	}

	return err
}
