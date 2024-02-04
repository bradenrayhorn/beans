package response

import "github.com/bradenrayhorn/beans/server/beans"

type Budget struct {
	ID   beans.ID   `json:"id"`
	Name beans.Name `json:"name"`
}

type CreateBudgetResponse Data[ID]

type ListBudgetsResposne Data[[]Budget]

type GetBudgetResponse Data[Budget]
