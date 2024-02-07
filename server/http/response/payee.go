package response

import "github.com/bradenrayhorn/beans/server/beans"

type AssociatedPayee struct {
	ID   beans.ID   `json:"id"`
	Name beans.Name `json:"name"`
}

type Payee struct {
	ID   beans.ID   `json:"id"`
	Name beans.Name `json:"name"`
}

type CreatePayeeResponse Data[ID]
type ListPayeesResponse Data[[]Payee]
type GetPayeeResponse Data[Payee]
