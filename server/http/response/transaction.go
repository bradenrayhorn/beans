package response

import "github.com/bradenrayhorn/beans/server/beans"

type Transaction struct {
	ID       string                 `json:"id"`
	Account  AssociatedAccount      `json:"account"`
	Category *AssociatedCategory    `json:"category"`
	Payee    *AssociatedPayee       `json:"payee"`
	Amount   beans.Amount           `json:"amount"`
	Date     string                 `json:"date"`
	Notes    beans.TransactionNotes `json:"notes"`
}

type CreateTransactionResponse Data[ID]

type ListTransactionsResponse Data[[]Transaction]
