package http

import (
	"net/http"

	"github.com/bradenrayhorn/beans/beans"
)

type transactionResponse struct {
	ID      string                 `json:"id"`
	Account responseAccount        `json:"account"`
	Amount  beans.Amount           `json:"amount"`
	Date    string                 `json:"date"`
	Notes   beans.TransactionNotes `json:"notes"`
}

func responseFromTransaction(transaction *beans.Transaction) transactionResponse {
	return transactionResponse{
		ID:      transaction.ID.String(),
		Account: responseFromAccount(transaction.Account),
		Amount:  transaction.Amount,
		Date:    transaction.Date.String(),
		Notes:   transaction.Notes,
	}
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

		res := response{Data: responseFromTransaction(transaction)}
		jsonResponse(w, res, http.StatusOK)
	}
}

func (s *Server) handleTransactionGetAll() http.HandlerFunc {
	type response struct {
		Data []transactionResponse `json:"data"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		transactions, err := s.transactionRepository.GetForBudget(r.Context(), getBudget(r).ID)
		if err != nil {
			Error(w, err)
			return
		}

		res := response{Data: make([]transactionResponse, len(transactions))}
		for i, t := range transactions {
			res.Data[i] = responseFromTransaction(t)
		}

		jsonResponse(w, res, http.StatusOK)
	}
}
