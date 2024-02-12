package request

import "github.com/bradenrayhorn/beans/server/beans"

type transactionRequest struct {
	AccountID  beans.ID               `json:"account_id"`
	CategoryID beans.ID               `json:"category_id"`
	PayeeID    beans.ID               `json:"payee_id"`
	Amount     beans.Amount           `json:"amount"`
	Date       beans.Date             `json:"date"`
	Notes      beans.TransactionNotes `json:"notes"`
}

type CreateTransactionRequest transactionRequest
type UpdateTransactionRequest transactionRequest