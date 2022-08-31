package http

import (
	"net/http"

	"github.com/bradenrayhorn/beans/beans"
)

type amountResponse struct {
	Coefficient int64 `json:"coefficient"`
	Exponent    int64 `json:"exponent"`
}

type transactionResponse struct {
	ID        string         `json:"id"`
	AccountID string         `json:"account_id"`
	Amount    amountResponse `json:"amount"`
	Date      string         `json:"date"`
	Notes     string         `json:"notes"`
}

func (s *Server) handleTransactionCreate() http.HandlerFunc {
	type request struct {
		AccountID beans.ID               `json:"account_id"`
		Amount    beans.Amount           `json:"amount"`
		Date      beans.Date             `json:"date"`
		Notes     beans.TransactionNotes `json:"notes"`
	}

	type response struct {
		Data transactionResponse `json:"data"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := decodeRequest(r, &req); err != nil {
			Error(w, err)
			return
		}

		transaction, err := s.transactionService.Create(r.Context(), getBudget(r), beans.TransactionCreate{
			AccountID: req.AccountID,
			Amount:    req.Amount,
			Date:      req.Date,
			Notes:     req.Notes,
		})

		if err != nil {
			Error(w, err)
			return
		}

		res := response{Data: transactionResponse{
			ID:        transaction.ID.String(),
			AccountID: transaction.AccountID.String(),
			Amount:    amountResponse{Coefficient: transaction.Amount.Coefficient().Int64(), Exponent: int64(transaction.Amount.Exponent())},
			Date:      transaction.Date.String(),
			Notes:     string(transaction.Notes),
		}}

		jsonResponse(w, res, http.StatusOK)
	}
}
