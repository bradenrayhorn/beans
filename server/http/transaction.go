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

	var transferAccount *response.AssociatedAccount
	if a, ok := transaction.TransferAccount.Value(); ok {
		transferAccount = &response.AssociatedAccount{
			ID:        a.ID,
			Name:      a.Name,
			OffBudget: a.OffBudget,
		}
	}

	return response.Transaction{
		ID:      transaction.ID,
		Variant: transaction.Variant,
		Account: response.AssociatedAccount{
			ID:        transaction.Account.ID,
			Name:      transaction.Account.Name,
			OffBudget: transaction.Account.OffBudget,
		},
		Category:        category,
		Payee:           payee,
		Amount:          transaction.Amount,
		Date:            transaction.Date,
		Notes:           transaction.Notes,
		TransferAccount: transferAccount,
	}
}

func (s *Server) handleTransactionCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req request.CreateTransaction
		if err := decodeRequest(r, &req); err != nil {
			Error(w, err)
			return
		}

		splits := make([]beans.SplitParams, len(req.Splits))
		for i, s := range req.Splits {
			splits[i] = beans.SplitParams{
				Amount:     s.Amount,
				CategoryID: s.CategoryID,
				Notes:      s.Notes,
			}
		}

		transactionID, err := s.contracts.Transaction.Create(r.Context(), getBudgetAuth(r), beans.TransactionCreateParams{
			TransferAccountID: req.TransferAccountID,
			TransactionParams: beans.TransactionParams{
				AccountID:  req.AccountID,
				CategoryID: req.CategoryID,
				PayeeID:    req.PayeeID,
				Amount:     req.Amount,
				Date:       req.Date,
				Notes:      req.Notes,
				Splits:     splits,
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
		var req request.UpdateTransaction
		if err := decodeRequest(r, &req); err != nil {
			Error(w, err)
			return
		}

		transactionID, err := beans.IDFromString(chi.URLParam(r, "transactionID"))
		if err != nil {
			Error(w, err)
			return
		}

		splits := make([]beans.SplitParams, len(req.Splits))
		for i, s := range req.Splits {
			splits[i] = beans.SplitParams{
				Amount:     s.Amount,
				CategoryID: s.CategoryID,
				Notes:      s.Notes,
			}
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
				Splits:     splits,
			},
		})

		if err != nil {
			Error(w, err)
			return
		}
	}
}

func (s *Server) handleTransactionDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req request.DeleteTransaction
		if err := decodeRequest(r, &req); err != nil {
			Error(w, err)
			return
		}

		err := s.contracts.Transaction.Delete(r.Context(), getBudgetAuth(r), req.IDs)
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
		id, err := beans.IDFromString(chi.URLParam(r, "transactionID"))
		if err != nil {
			Error(w, beans.WrapError(err, beans.ErrorNotFound))
			return
		}

		transaction, err := s.contracts.Transaction.Get(r.Context(), getBudgetAuth(r), id)
		if err != nil {
			Error(w, err)
			return
		}

		jsonResponse(w,
			response.GetTransactionResponse{Data: responseFromTransaction(transaction)},
			http.StatusOK)
	}
}

func (s *Server) handleTransactionGetSplits() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := beans.IDFromString(chi.URLParam(r, "transactionID"))
		if err != nil {
			Error(w, beans.WrapError(err, beans.ErrorNotFound))
			return
		}

		splits, err := s.contracts.Transaction.GetSplits(r.Context(), getBudgetAuth(r), id)
		if err != nil {
			Error(w, err)
			return
		}

		res := response.GetSplitsResponse{Data: make([]response.Split, len(splits))}
		for i, t := range splits {
			res.Data[i] = response.Split{
				ID:       t.ID,
				Amount:   t.Amount,
				Category: response.AssociatedCategory(t.Category),
				Notes:    t.Notes,
			}
		}

		jsonResponse(w, res, http.StatusOK)
	}
}
