package postgres

import (
	"context"
	"errors"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/db"
	"github.com/bradenrayhorn/beans/server/postgres/mapper"
)

type monthRepository struct {
	repository
}

func NewMonthRepository(pool *DbPool) *monthRepository {
	return &monthRepository{repository{pool}}
}

func (r *monthRepository) Create(ctx context.Context, tx beans.Tx, month *beans.Month) error {
	return r.DB(tx).CreateMonth(ctx, db.CreateMonthParams{
		ID:        month.ID.String(),
		BudgetID:  month.BudgetID.String(),
		Date:      mapper.MonthDateToPg(month.Date),
		Carryover: mapper.AmountToNumeric(month.Carryover.OrZero()),
	})
}

func (r *monthRepository) Get(ctx context.Context, id beans.ID) (*beans.Month, error) {
	res, err := r.DB(nil).GetMonthByID(ctx, id.String())
	if err != nil {
		return nil, mapPostgresError(err)
	}

	return mapper.Month(res)
}

func (r *monthRepository) Update(ctx context.Context, month *beans.Month) error {
	return r.DB(nil).UpdateMonth(ctx, db.UpdateMonthParams{
		ID:        month.ID.String(),
		Carryover: mapper.AmountToNumeric(month.Carryover.OrZero()),
	})
}

func (r *monthRepository) GetOrCreate(ctx context.Context, tx beans.Tx, budgetID beans.ID, date beans.MonthDate) (*beans.Month, error) {
	res, err := r.DB(tx).GetMonthByDate(ctx, db.GetMonthByDateParams{BudgetID: budgetID.String(), Date: mapper.MonthDateToPg(date)})
	if err != nil {
		err = mapPostgresError(err)

		if errors.Is(err, beans.ErrorNotFound) {
			month := &beans.Month{
				ID:        beans.NewBeansID(),
				BudgetID:  budgetID,
				Date:      date,
				Carryover: beans.NewAmount(0, 0),
			}
			return month, r.Create(ctx, tx, month)
		}
	}

	return mapper.Month(res)
}

func (r *monthRepository) GetForBudget(ctx context.Context, budgetID beans.ID) ([]*beans.Month, error) {
	res, err := r.DB(nil).GetMonthsByBudget(ctx, budgetID.String())
	if err != nil {
		return nil, mapPostgresError(err)
	}

	return mapper.MapSlice(res, mapper.Month)
}
