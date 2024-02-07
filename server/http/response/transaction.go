package response

import "github.com/bradenrayhorn/beans/server/beans"

type Transaction struct {
	ID       beans.ID               `json:"id"`
	Account  AssociatedAccount      `json:"account"`
	Category *AssociatedCategory    `json:"category"`
	Payee    *AssociatedPayee       `json:"payee"`
	Amount   beans.Amount           `json:"amount"`
	Date     beans.Date             `json:"date"`
	Notes    beans.TransactionNotes `json:"notes"`
}

type CreateTransactionResponse Data[ID]

type ListTransactionsResponse Data[[]Transaction]

type GetTransactionResponse Data[Transaction]
