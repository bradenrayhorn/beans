package response

import "github.com/bradenrayhorn/beans/server/beans"

type AssociatedAccount struct {
	ID        beans.ID   `json:"id"`
	Name      beans.Name `json:"name"`
	OffBudget bool       `json:"offBudget"`
}

type Account struct {
	ID        beans.ID `json:"id"`
	Name      string   `json:"name"`
	OffBudget bool     `json:"offBudget"`
}

type ListAccount struct {
	ID        beans.ID     `json:"id"`
	Name      string       `json:"name"`
	Balance   beans.Amount `json:"balance"`
	OffBudget bool         `json:"offBudget"`
}

type CreateAccountResponse Data[ID]
type ListAccountResponse Data[[]ListAccount]
type GetAccountResponse Data[Account]
type GetTransactableAccounts Data[[]Account]
