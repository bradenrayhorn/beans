package request

import "github.com/bradenrayhorn/beans/server/beans"

type UpdateMonth struct {
	Carryover beans.Amount `json:"carryover"`
}

type UpdateMonthCategory struct {
	CategoryID beans.ID     `json:"category_id"`
	Amount     beans.Amount `json:"amount"`
}
