package request

import "github.com/bradenrayhorn/beans/server/beans"

type CreateTransaction struct {
	AccountID  beans.ID               `json:"account_id"`
	CategoryID beans.ID               `json:"category_id"`
	PayeeID    beans.ID               `json:"payee_id"`
	Amount     beans.Amount           `json:"amount"`
	Date       beans.Date             `json:"date"`
	Notes      beans.TransactionNotes `json:"notes"`

	Splits []Split `json:"splits"`

	TransferAccountID beans.ID `json:"transferAccountID"`
}

type UpdateTransaction struct {
	AccountID  beans.ID               `json:"account_id"`
	CategoryID beans.ID               `json:"category_id"`
	PayeeID    beans.ID               `json:"payee_id"`
	Amount     beans.Amount           `json:"amount"`
	Date       beans.Date             `json:"date"`
	Notes      beans.TransactionNotes `json:"notes"`

	Splits []Split `json:"splits"`
}

type Split struct {
	CategoryID beans.ID               `json:"category_id"`
	Amount     beans.Amount           `json:"amount"`
	Notes      beans.TransactionNotes `json:"notes"`
}

type DeleteTransaction struct {
	IDs []beans.ID `json:"ids"`
}
