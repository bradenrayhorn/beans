package request

import "github.com/bradenrayhorn/beans/server/beans"

type CreateAccount struct {
	Name      beans.Name `json:"name"`
	OffBudget bool       `json:"off_budget"`
}
