package http

import (
	"net/http"

	"github.com/bradenrayhorn/beans/beans"
)

func (s *Server) handleTransactionCreate() http.HandlerFunc {
	type request struct {
		AccountID beans.ID               `json:"account_id"`
		Amount    beans.Amount           `json:"amount"`
		Date      beans.Date             `json:"date"`
		Notes     beans.TransactionNotes `json:"notes"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := decodeRequest(r, &req); err != nil {
			Error(w, err)
			return
		}

		_, err := s.transactionService.Create(r.Context(), beans.TransactionCreate{
			AccountID: req.AccountID,
			Amount:    req.Amount,
			Date:      req.Date,
			Notes:     req.Notes,
		})

		if err != nil {
			Error(w, err)
			return
		}
	}
}
