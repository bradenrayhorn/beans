package http

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bradenrayhorn/beans/beans"
	"github.com/bradenrayhorn/beans/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateTransaction(t *testing.T) {
	transactionService := new(mocks.TransactionService)
	sv := &Server{transactionService: transactionService}

	date, err := time.Parse("2006-01-02", "2022-08-29")
	require.Nil(t, err)
	transaction := beans.Transaction{
		ID:        beans.NewBeansID(),
		AccountID: beans.NewBeansID(),
		Amount:    beans.NewAmount(1456, -2),
		Date:      beans.NewDate(date),
		Notes:     beans.TransactionNotes("My Notes"),
	}
	transactionService.On("Create", mock.Anything, mock.Anything, mock.Anything).Return(&transaction, nil)

	budget := &beans.Budget{ID: beans.NewBeansID(), Name: "Budget1"}

	req := httptest.NewRequest("POST", "/transactions", bytes.NewReader([]byte(fmt.Sprintf(`
    {
      "account_id": "%s",
      "amount": 14.56,
      "date": "2022-08-29"
    }
    `, beans.NewBeansID()))))
	req = req.WithContext(context.WithValue(req.Context(), "budget", budget))
	w := httptest.NewRecorder()
	sv.handleTransactionCreate().ServeHTTP(w, req)
	res := w.Result()
	require.Equal(t, http.StatusOK, res.StatusCode)
	data, err := ioutil.ReadAll(res.Body)
	require.Nil(t, err)
	assert.JSONEq(t, string(data), fmt.Sprintf(`{"data":{
    "id": "%s",
    "account_id": "%s",
    "amount": {
      "coefficient": 1456,
      "exponent": -2
    },
    "date": "2022-08-29",
    "notes": "My Notes"
    }}`, transaction.ID, transaction.AccountID))
}
