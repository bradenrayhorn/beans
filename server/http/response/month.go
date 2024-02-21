package response

import "github.com/bradenrayhorn/beans/server/beans"

type MonthCategory struct {
	ID         beans.ID     `json:"id"`
	Assigned   beans.Amount `json:"assigned"`
	Activity   beans.Amount `json:"activity"`
	Available  beans.Amount `json:"available"`
	CategoryID beans.ID     `json:"category_id"`
}
type Month struct {
	ID          beans.ID        `json:"id"`
	Date        beans.MonthDate `json:"date"`
	Budgetable  beans.Amount    `json:"budgetable"`
	Carryover   beans.Amount    `json:"carryover"`
	Income      beans.Amount    `json:"income"`
	Assigned    beans.Amount    `json:"assigned"`
	CarriedOver beans.Amount    `json:"carried_over"`
	Categories  []MonthCategory `json:"categories"`
}

type GetMonthResponse Data[Month]
