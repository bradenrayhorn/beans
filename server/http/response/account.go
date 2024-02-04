package response

import "github.com/bradenrayhorn/beans/server/beans"

type AssociatedAccount struct {
	ID   beans.ID   `json:"id"`
	Name beans.Name `json:"name"`
}

type CreateAccountResponse Data[ID]

type ListAccount struct {
	ID      beans.ID     `json:"id"`
	Name    string       `json:"name"`
	Balance beans.Amount `json:"balance"`
}

type ListAccountResponse Data[[]ListAccount]
