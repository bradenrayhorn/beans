package response

import "github.com/bradenrayhorn/beans/server/beans"

type Transaction struct {
	ID       beans.ID                 `json:"id"`
	Variant  beans.TransactionVariant `json:"variant"`
	Account  AssociatedAccount        `json:"account"`
	Category *AssociatedCategory      `json:"category"`
	Payee    *AssociatedPayee         `json:"payee"`
	Amount   beans.Amount             `json:"amount"`
	Date     beans.Date               `json:"date"`
	Notes    beans.TransactionNotes   `json:"notes"`

	TransferID      beans.ID           `json:"transferID"`
	TransferAccount *AssociatedAccount `json:"transferAccount"`
}

type Split struct {
	ID       beans.ID               `json:"id"`
	Category AssociatedCategory     `json:"category"`
	Amount   beans.Amount           `json:"amount"`
	Notes    beans.TransactionNotes `json:"notes"`
}

type CreateTransactionResponse Data[ID]

type ListTransactionsResponse Data[[]Transaction]

type GetTransactionResponse Data[Transaction]

type GetSplitsResponse Data[[]Split]
