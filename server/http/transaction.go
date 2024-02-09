package http

import (
	"net/http"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/http/request"
	"github.com/bradenrayhorn/beans/server/http/response"
	"github.com/go-chi/chi/v5"
)

func responseFromTransaction(transaction beans.TransactionWithRelations) response.Transaction {
	var category *response.AssociatedCategory
	if c, ok := transaction.Category.Value(); ok {
		category = &response.AssociatedCategory{
			ID:   c.ID,
			Name: c.Name,
		}

	}

	var payee *response.AssociatedPayee
	if p, ok := transaction.Payee.Value(); ok {
		payee = &response.AssociatedPayee{
			ID:   p.ID,
			Name: p.Name,
		}
	}

	return response.Transaction{
		ID: transaction.ID,
		Account: response.AssociatedAccount{
			ID:   transaction.AccountID,
			Name: transaction.Account.Name,
		},
		Category: category,
		Payee:    payee,
		Amount:   transaction.Amount,
		Date:     transaction.Date,
		Notes:    transaction.Notes,
	}
}

func (s *Server) handleTransactionCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req request.CreateTransactionRequest
		if err := decodeRequest(r, &req); err != nil {
			Error(w, err)
			return
		}

		transactionID, err := s.contracts.Transaction.Create(r.Context(), getBudgetAuth(r), beans.TransactionCreateParams{
			TransactionParams: beans.TransactionParams{
				AccountID:  req.AccountID,
				CategoryID: req.CategoryID,
				PayeeID:    req.PayeeID,
				Amount:     req.Amount,
				Date:       req.Date,
				Notes:      req.Notes,
			},
		})

		if err != nil {
			Error(w, err)
			return
		}

		jsonResponse(w, response.CreateTransactionResponse{
			Data: response.ID{ID: transactionID},
		}, http.StatusOK)
	}
}

func (s *Server) handleTransactionUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req request.UpdateTransactionRequest
		if err := decodeRequest(r, &req); err != nil {
			Error(w, err)
			return
		}

		transactionID, err := beans.BeansIDFromString(chi.URLParam(r, "transactionID"))
		if err != nil {
			Error(w, err)
			return
		}

		err = s.contracts.Transaction.Update(r.Context(), getBudgetAuth(r), beans.TransactionUpdateParams{
			ID: transactionID,
			TransactionParams: beans.TransactionParams{
				AccountID:  req.AccountID,
				CategoryID: req.CategoryID,
				PayeeID:    req.PayeeID,
				Amount:     req.Amount,
				Date:       req.Date,
				Notes:      req.Notes,
			},
		})

		if err != nil {
			Error(w, err)
			return
		}
	}
}

func (s *Server) handleTransactionDelete() http.HandlerFunc {
	type request struct {
		TransactionIDs []beans.ID `json:"ids"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := decodeRequest(r, &req); err != nil {
			Error(w, err)
			return
		}

		err := s.contracts.Transaction.Delete(r.Context(), getBudgetAuth(r), req.TransactionIDs)
		if err != nil {
			Error(w, err)
			return
		}
	}
}

func (s *Server) handleTransactionGetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transactions, err := s.contracts.Transaction.GetAll(r.Context(), getBudgetAuth(r))
		if err != nil {
			Error(w, err)
			return
		}

		res := response.ListTransactionsResponse{Data: make([]response.Transaction, len(transactions))}
		for i, t := range transactions {
			res.Data[i] = responseFromTransaction(t)
		}

		jsonResponse(w, res, http.StatusOK)
	}
}

func (s *Server) handleTransactionGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := beans.BeansIDFromString(chi.URLParam(r, "transactionID"))
		if err != nil {
			Error(w, beans.WrapError(err, beans.ErrorNotFound))
			return
		}

		transaction, err := s.contracts.Transaction.Get(r.Context(), getBudgetAuth(r), id)
		if err != nil {
			Error(w, err)
			return
		}

		jsonResponse(w, responseFromTransaction(transaction), http.StatusOK)
	}
}
