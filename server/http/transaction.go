package http

import (
	"net/http"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/http/response"
	"github.com/go-chi/chi/v5"
)

func responseFromTransaction(transaction *beans.Transaction) response.Transaction {
	var category *response.AssociatedCategory
	if !transaction.CategoryID.Empty() && !transaction.CategoryName.Empty() {
		category = &response.AssociatedCategory{
			ID:   transaction.CategoryID,
			Name: beans.Name(transaction.CategoryName.String()),
		}
	}

	var payee *response.AssociatedPayee
	if !transaction.PayeeID.Empty() && !transaction.PayeeName.Empty() {
		payee = &response.AssociatedPayee{
			ID:   transaction.PayeeID,
			Name: beans.Name(transaction.PayeeName.String()),
		}
	}

	return response.Transaction{
		ID: transaction.ID.String(),
		Account: response.AssociatedAccount{
			ID:   transaction.AccountID,
			Name: transaction.Account.Name,
		},
		Category: category,
		Payee:    payee,
		Amount:   transaction.Amount,
		Date:     transaction.Date.String(),
		Notes:    transaction.Notes,
	}
}

func (s *Server) handleTransactionCreate() http.HandlerFunc {
	type request struct {
		AccountID  beans.ID               `json:"account_id"`
		CategoryID beans.ID               `json:"category_id"`
		PayeeID    beans.ID               `json:"payee_id"`
		Amount     beans.Amount           `json:"amount"`
		Date       beans.Date             `json:"date"`
		Notes      beans.TransactionNotes `json:"notes"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := decodeRequest(r, &req); err != nil {
			Error(w, err)
			return
		}

		transaction, err := s.transactionContract.Create(r.Context(), getBudgetAuth(r), beans.TransactionCreateParams{
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
			Data: response.ID{ID: transaction.ID},
		}, http.StatusOK)
	}
}

func (s *Server) handleTransactionUpdate() http.HandlerFunc {
	type request struct {
		AccountID  beans.ID               `json:"account_id"`
		CategoryID beans.ID               `json:"category_id"`
		PayeeID    beans.ID               `json:"payee_id"`
		Amount     beans.Amount           `json:"amount"`
		Date       beans.Date             `json:"date"`
		Notes      beans.TransactionNotes `json:"notes"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := decodeRequest(r, &req); err != nil {
			Error(w, err)
			return
		}

		transactionID, err := beans.BeansIDFromString(chi.URLParam(r, "transactionID"))
		if err != nil {
			Error(w, err)
			return
		}

		err = s.transactionContract.Update(r.Context(), getBudgetAuth(r), beans.TransactionUpdateParams{
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

		err := s.transactionContract.Delete(r.Context(), getBudgetAuth(r), req.TransactionIDs)
		if err != nil {
			Error(w, err)
			return
		}
	}
}

func (s *Server) handleTransactionGetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		transactions, err := s.transactionContract.GetAll(r.Context(), getBudgetAuth(r))
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
