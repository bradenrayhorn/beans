package http

import (
	"net/http"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/go-chi/chi/v5"
)

type transactionPayee struct {
	ID   beans.ID   `json:"id"`
	Name beans.Name `json:"name"`
}

type transactionResponse struct {
	ID       string                 `json:"id"`
	Account  responseAccount        `json:"account"`
	Category *listCategoryResponse  `json:"category"`
	Payee    *transactionPayee      `json:"payee"`
	Amount   beans.Amount           `json:"amount"`
	Date     string                 `json:"date"`
	Notes    beans.TransactionNotes `json:"notes"`
}

func responseFromTransaction(transaction *beans.Transaction) transactionResponse {
	var category *listCategoryResponse
	if !transaction.CategoryID.Empty() && !transaction.CategoryName.Empty() {
		category = &listCategoryResponse{
			ID:   transaction.CategoryID,
			Name: beans.Name(transaction.CategoryName.String()),
		}
	}

	var payee *transactionPayee
	if !transaction.PayeeID.Empty() && !transaction.PayeeName.Empty() {
		payee = &transactionPayee{
			ID:   transaction.PayeeID,
			Name: beans.Name(transaction.PayeeName.String()),
		}
	}

	return transactionResponse{
		ID:       transaction.ID.String(),
		Account:  responseFromAccount(transaction.Account),
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

	type response struct {
		ID beans.ID `json:"transaction_id"`
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

		jsonResponse(w, struct {
			Data response `json:"data"`
		}{
			Data: response{ID: transaction.ID},
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

func (s *Server) handleTransactionGetAll() http.HandlerFunc {
	type response struct {
		Data []transactionResponse `json:"data"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		transactions, err := s.transactionContract.GetAll(r.Context(), getBudgetAuth(r))
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
