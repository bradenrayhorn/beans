package specification

import (
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testTransaction(t *testing.T, interactor Interactor) {

	t.Run("get", func(t *testing.T) {

		t.Run("cannot get non-existent", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			_, err := interactor.TransactionGet(t, c.ctx, beans.NewID())
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})

		t.Run("cannot get from other budget", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)
			c2 := makeUserAndBudget(t, interactor)

			transaction := c.Transaction(TransactionOpts{})

			_, err := interactor.TransactionGet(t, c2.ctx, transaction.ID)
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})

		t.Run("can get off-budget variant", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			account := c.Account(AccountOpts{OffBudget: true})
			transaction := c.Transaction(TransactionOpts{Account: account})

			res, err := interactor.TransactionGet(t, c.ctx, transaction.ID)
			require.NoError(t, err)
			assert.Equal(t, beans.TransactionOffBudget, res.Variant)
		})
	})

	t.Run("create", func(t *testing.T) {

		basicParams := beans.TransactionCreateParams{
			TransactionParams: beans.TransactionParams{
				Amount: beans.NewAmount(1, 2),
				Date:   testutils.NewDate(t, "2022-06-07"),
			},
		}

		// validation

		t.Run("does validation", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			_, err := interactor.TransactionCreate(t, c.ctx, beans.TransactionCreateParams{})
			testutils.AssertError(t, err, "Account ID is required. Amount is required. Date is required.")
		})

		// actually creating

		t.Run("with all parameters", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			account := c.Account(AccountOpts{})
			category := c.Category(CategoryOpts{})
			payee := c.Payee(PayeeOpts{})

			// create transaction
			params := beans.TransactionCreateParams{
				TransactionParams: beans.TransactionParams{
					AccountID:  account.ID,
					CategoryID: category.ID,
					PayeeID:    payee.ID,
					Amount:     beans.NewAmount(7, 0),
					Date:       testutils.NewDate(t, "2022-06-07"),
					Notes:      beans.NewTransactionNotes("My Notes"),
				},
			}

			id, err := interactor.TransactionCreate(t, c.ctx, params)
			require.NoError(t, err)

			// get transaction and verify
			transaction, err := interactor.TransactionGet(t, c.ctx, id)
			require.NoError(t, err)

			assert.False(t, transaction.ID.Empty())
			assert.Equal(t, account.ID, transaction.AccountID)
			assert.Equal(t, category.ID, transaction.CategoryID)
			assert.Equal(t, payee.ID, transaction.PayeeID)
			assert.Equal(t, beans.NewAmount(7, 0), transaction.Amount)
			assert.Equal(t, testutils.NewDate(t, "2022-06-07"), transaction.Date)
			assert.Equal(t, beans.NewTransactionNotes("My Notes"), transaction.Notes)

			assert.Equal(t, beans.TransactionStandard, transaction.Variant)
			assert.Equal(t, beans.RelatedAccount{ID: account.ID, Name: account.Name}, transaction.Account)
			assert.Equal(t, beans.OptionalWrap(beans.RelatedCategory{ID: category.ID, Name: category.Name}), transaction.Category)
			assert.Equal(t, beans.OptionalWrap(beans.RelatedPayee{ID: payee.ID, Name: payee.Name}), transaction.Payee)
		})

		t.Run("with minimal parameters", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			account := c.Account(AccountOpts{})

			// create transaction
			params := beans.TransactionCreateParams{
				TransactionParams: beans.TransactionParams{
					AccountID: account.ID,
					Amount:    beans.NewAmount(7, 0),
					Date:      testutils.NewDate(t, "2022-06-07"),
				},
			}

			id, err := interactor.TransactionCreate(t, c.ctx, params)
			require.NoError(t, err)

			// get transaction and verify
			transaction, err := interactor.TransactionGet(t, c.ctx, id)
			require.NoError(t, err)

			assert.False(t, transaction.ID.Empty())
			assert.Equal(t, account.ID, transaction.AccountID)
			assert.Equal(t, beans.EmptyID(), transaction.CategoryID)
			assert.Equal(t, beans.EmptyID(), transaction.PayeeID)
			assert.Equal(t, beans.NewAmount(7, 0), transaction.Amount)
			assert.Equal(t, testutils.NewDate(t, "2022-06-07"), transaction.Date)
			assert.Equal(t, beans.TransactionNotes{}, transaction.Notes)

			assert.Equal(t, beans.TransactionStandard, transaction.Variant)
			assert.Equal(t, beans.RelatedAccount{ID: account.ID, Name: account.Name}, transaction.Account)
			assert.True(t, transaction.Category.Empty())
			assert.True(t, transaction.Payee.Empty())
		})

		// account validation

		t.Run("account validation", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)
			c2 := makeUserAndBudget(t, interactor)

			t.Run("cannot use non-existent account", func(t *testing.T) {
				params := basicParams
				params.AccountID = beans.NewID()

				_, err := interactor.TransactionCreate(t, c.ctx, params)
				testutils.AssertErrorAndCode(t, err, beans.EINVALID, "Invalid Account ID")
			})

			t.Run("cannot use account from another budget", func(t *testing.T) {
				params := basicParams
				params.AccountID = c2.Account(AccountOpts{}).ID

				_, err := interactor.TransactionCreate(t, c.ctx, params)
				testutils.AssertErrorAndCode(t, err, beans.EINVALID, "Invalid Account ID")
			})
		})

		t.Run("category validation", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)
			c2 := makeUserAndBudget(t, interactor)

			account := c.Account(AccountOpts{})
			params := basicParams
			params.AccountID = account.ID

			t.Run("cannot use non-existent category", func(t *testing.T) {
				params.CategoryID = beans.NewID()

				_, err := interactor.TransactionCreate(t, c.ctx, params)
				testutils.AssertErrorAndCode(t, err, beans.EINVALID, "Invalid Category ID")
			})

			t.Run("cannot use category from another budget", func(t *testing.T) {
				params.CategoryID = c2.Category(CategoryOpts{}).ID

				_, err := interactor.TransactionCreate(t, c.ctx, params)
				testutils.AssertErrorAndCode(t, err, beans.EINVALID, "Invalid Category ID")
			})

			t.Run("cannot assign category with off-budget account", func(t *testing.T) {
				params.AccountID = c.Account(AccountOpts{OffBudget: true}).ID
				params.CategoryID = c.Category(CategoryOpts{}).ID

				_, err := interactor.TransactionCreate(t, c.ctx, params)
				testutils.AssertErrorAndCode(t, err, beans.EINVALID, "Cannot assign category with off-budget account")
			})
		})

		t.Run("payee validation", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)
			c2 := makeUserAndBudget(t, interactor)

			account := c.Account(AccountOpts{})
			params := basicParams
			params.AccountID = account.ID

			t.Run("cannot use non-existent payee", func(t *testing.T) {
				params.PayeeID = beans.NewID()

				_, err := interactor.TransactionCreate(t, c.ctx, params)
				testutils.AssertErrorAndCode(t, err, beans.EINVALID, "Invalid Payee ID")
			})

			t.Run("cannot use payee from another budget", func(t *testing.T) {
				params.PayeeID = c2.Payee(PayeeOpts{}).ID

				_, err := interactor.TransactionCreate(t, c.ctx, params)
				testutils.AssertErrorAndCode(t, err, beans.EINVALID, "Invalid Payee ID")
			})
		})

	})

	t.Run("update", func(t *testing.T) {

		t.Run("does validation", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			err := interactor.TransactionUpdate(t, c.ctx, beans.TransactionUpdateParams{})
			testutils.AssertError(t, err, "Account ID is required. Amount is required. Date is required.")
		})

		t.Run("validates existing transaction", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)
			account := c.Account(AccountOpts{})
			category := c.Category(CategoryOpts{})
			transaction := c.Transaction(TransactionOpts{
				Account:  account,
				Category: category,
				Amount:   "10",
			})

			baseParams := beans.TransactionUpdateParams{
				TransactionParams: beans.TransactionParams{
					AccountID:  account.ID,
					CategoryID: category.ID,
					Amount:     beans.NewAmount(1, 1),
					Date:       testutils.NewDate(t, "2022-05-01"),
				},
			}

			t.Run("cannot update non-existent transaction", func(t *testing.T) {
				params := baseParams
				params.ID = beans.NewID()

				err := interactor.TransactionUpdate(t, c.ctx, params)
				testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
			})

			t.Run("cannot update transaction from other budget", func(t *testing.T) {
				c2 := makeUserAndBudget(t, interactor)

				params := baseParams
				params.ID = transaction.ID

				err := interactor.TransactionUpdate(t, c2.ctx, params)
				testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
			})
		})

		t.Run("can update", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)
			account := c.Account(AccountOpts{})
			account2 := c.Account(AccountOpts{})
			category := c.Category(CategoryOpts{})
			category2 := c.Category(CategoryOpts{})
			payee := c.Payee(PayeeOpts{})
			payee2 := c.Payee(PayeeOpts{})

			t.Run("with all parameters", func(t *testing.T) {

				// create transaction
				transaction := c.Transaction(TransactionOpts{
					Account:  account,
					Category: category,
					Payee:    payee,
					Amount:   "10",
					Date:     "2022-04-05",
					Notes:    "Notes",
				})

				// update transaction
				params := beans.TransactionUpdateParams{
					ID: transaction.ID,
					TransactionParams: beans.TransactionParams{
						AccountID:  account2.ID,
						CategoryID: category2.ID,
						PayeeID:    payee2.ID,
						Amount:     beans.NewAmount(6, 0),
						Date:       testutils.NewDate(t, "2022-06-07"),
						Notes:      beans.NewTransactionNotes("My Notes"),
					},
				}

				err := interactor.TransactionUpdate(t, c.ctx, params)
				require.NoError(t, err)

				// get and verify transaction
				res, err := interactor.TransactionGet(t, c.ctx, transaction.ID)
				require.NoError(t, err)

				assert.Equal(t, transaction.ID, res.ID)
				assert.Equal(t, account2.ID, res.AccountID)
				assert.Equal(t, category2.ID, res.CategoryID)
				assert.Equal(t, payee2.ID, res.PayeeID)
				assert.Equal(t, beans.NewAmount(6, 0), res.Amount)
				assert.Equal(t, testutils.NewDate(t, "2022-06-07"), res.Date)
				assert.Equal(t, beans.NewTransactionNotes("My Notes"), res.Notes)

				assert.Equal(t, beans.TransactionStandard, res.Variant)
				assert.Equal(t, beans.RelatedAccount{ID: account2.ID, Name: account2.Name}, res.Account)
				assert.Equal(t, beans.OptionalWrap(beans.RelatedCategory{ID: category2.ID, Name: category2.Name}), res.Category)
				assert.Equal(t, beans.OptionalWrap(beans.RelatedPayee{ID: payee2.ID, Name: payee2.Name}), res.Payee)
			})

			t.Run("with minimal parameters", func(t *testing.T) {

				// create transaction
				transaction := c.Transaction(TransactionOpts{
					Account: account,
					Amount:  "10",
					Date:    "2022-04-05",
				})

				// update transaction
				params := beans.TransactionUpdateParams{
					ID: transaction.ID,
					TransactionParams: beans.TransactionParams{
						AccountID: account2.ID,
						Amount:    beans.NewAmount(6, 0),
						Date:      testutils.NewDate(t, "2022-06-07"),
					},
				}

				require.NoError(t, interactor.TransactionUpdate(t, c.ctx, params))

				// get and verify transaction
				res, err := interactor.TransactionGet(t, c.ctx, transaction.ID)
				require.NoError(t, err)

				assert.Equal(t, transaction.ID, res.ID)
				assert.Equal(t, account2.ID, res.AccountID)
				assert.Equal(t, beans.EmptyID(), res.CategoryID)
				assert.Equal(t, beans.EmptyID(), res.PayeeID)
				assert.Equal(t, beans.NewAmount(6, 0), res.Amount)
				assert.Equal(t, testutils.NewDate(t, "2022-06-07"), res.Date)
				assert.Equal(t, beans.NewTransactionNotes(""), res.Notes)

				assert.Equal(t, beans.TransactionStandard, res.Variant)
				assert.Equal(t, beans.RelatedAccount{ID: account2.ID, Name: account2.Name}, res.Account)
				assert.Equal(t, beans.Optional[beans.RelatedCategory]{}, res.Category)
				assert.Equal(t, beans.Optional[beans.RelatedPayee]{}, res.Payee)
			})
		})

		t.Run("validates related models", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)
			c2 := makeUserAndBudget(t, interactor)

			account := c.Account(AccountOpts{})

			transaction := c.Transaction(TransactionOpts{
				Account: account,
				Amount:  "10",
				Date:    "2022-04-05",
			})
			basicParams := beans.TransactionUpdateParams{
				ID: transaction.ID,
				TransactionParams: beans.TransactionParams{
					Amount: beans.NewAmount(1, 2),
					Date:   testutils.NewDate(t, "2022-06-07"),
				},
			}

			t.Run("account", func(t *testing.T) {

				t.Run("cannot use non-existent", func(t *testing.T) {
					params := basicParams
					params.AccountID = beans.NewID()

					err := interactor.TransactionUpdate(t, c.ctx, params)
					testutils.AssertErrorAndCode(t, err, beans.EINVALID, "Invalid Account ID")
				})

				t.Run("cannot use from another budget", func(t *testing.T) {
					params := basicParams
					params.AccountID = c2.Account(AccountOpts{}).ID

					err := interactor.TransactionUpdate(t, c.ctx, params)
					testutils.AssertErrorAndCode(t, err, beans.EINVALID, "Invalid Account ID")
				})
			})

			t.Run("category", func(t *testing.T) {
				params := basicParams
				params.AccountID = account.ID

				t.Run("cannot use non-existent", func(t *testing.T) {
					params.CategoryID = beans.NewID()

					err := interactor.TransactionUpdate(t, c.ctx, params)
					testutils.AssertErrorAndCode(t, err, beans.EINVALID, "Invalid Category ID")
				})

				t.Run("cannot use from another budget", func(t *testing.T) {
					params.CategoryID = c2.Category(CategoryOpts{}).ID

					err := interactor.TransactionUpdate(t, c.ctx, params)
					testutils.AssertErrorAndCode(t, err, beans.EINVALID, "Invalid Category ID")
				})

				t.Run("cannot assign category with off-budget account", func(t *testing.T) {
					params.AccountID = c.Account(AccountOpts{OffBudget: true}).ID
					params.CategoryID = c.Category(CategoryOpts{}).ID

					err := interactor.TransactionUpdate(t, c.ctx, params)
					testutils.AssertErrorAndCode(t, err, beans.EINVALID, "Cannot assign category with off-budget account")
				})
			})

			t.Run("payee", func(t *testing.T) {
				params := basicParams
				params.AccountID = account.ID

				t.Run("cannot use non-existent", func(t *testing.T) {
					params.PayeeID = beans.NewID()

					err := interactor.TransactionUpdate(t, c.ctx, params)
					testutils.AssertErrorAndCode(t, err, beans.EINVALID, "Invalid Payee ID")
				})

				t.Run("cannot use from another budget", func(t *testing.T) {
					params.PayeeID = c2.Payee(PayeeOpts{}).ID

					err := interactor.TransactionUpdate(t, c.ctx, params)
					testutils.AssertErrorAndCode(t, err, beans.EINVALID, "Invalid Payee ID")
				})
			})

		})

	})

	t.Run("delete", func(t *testing.T) {

		t.Run("can delete", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			account := c.Account(AccountOpts{})

			// make two transaction
			transaction1 := c.Transaction(TransactionOpts{
				Account: account,
				Amount:  "10",
				Date:    "2022-04-05",
			})
			transaction2 := c.Transaction(TransactionOpts{
				Account: account,
				Amount:  "10",
				Date:    "2022-04-05",
			})

			// try to delete one of them
			err := interactor.TransactionDelete(t, c.ctx, []beans.ID{transaction1.ID})
			require.NoError(t, err)

			// check that transaction1 is deleted but transaction2 is not
			_, err = interactor.TransactionGet(t, c.ctx, transaction1.ID)
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)

			_, err = interactor.TransactionGet(t, c.ctx, transaction2.ID)
			require.NoError(t, err)
		})

		t.Run("cannot delete from other budget", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)
			c2 := makeUserAndBudget(t, interactor)

			account := c.Account(AccountOpts{})

			// make transaction
			transaction := c.Transaction(TransactionOpts{
				Account: account,
				Amount:  "10",
				Date:    "2022-04-05",
			})

			// try to delete the transaction as budget 2
			err := interactor.TransactionDelete(t, c2.ctx, []beans.ID{transaction.ID})
			require.NoError(t, err)

			// check that transaction is not deleted
			_, err = interactor.TransactionGet(t, c.ctx, transaction.ID)
			require.NoError(t, err)
		})
	})

	t.Run("get all", func(t *testing.T) {

		t.Run("can get all", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			account := c.Account(AccountOpts{})
			category := c.Category(CategoryOpts{})
			payee := c.Payee(PayeeOpts{})

			// make transaction
			transaction := c.Transaction(TransactionOpts{
				Account:  account,
				Category: category,
				Payee:    payee,
				Amount:   "7",
				Date:     "2022-04-05",
				Notes:    "hey",
			})

			// get transactions and verify
			res, err := interactor.TransactionGetAll(t, c.ctx)
			require.NoError(t, err)

			assert.Equal(t, 1, len(res))
			findTransaction(t, res, transaction.ID, func(it beans.TransactionWithRelations) {
				assert.Equal(t, transaction.ID, it.ID)
				assert.Equal(t, account.ID, it.AccountID)
				assert.Equal(t, category.ID, it.CategoryID)
				assert.Equal(t, payee.ID, it.PayeeID)
				assert.Equal(t, beans.NewAmount(7, 0), it.Amount)
				assert.Equal(t, testutils.NewDate(t, "2022-04-05"), it.Date)
				assert.Equal(t, beans.NewTransactionNotes("hey"), it.Notes)

				assert.Equal(t, beans.TransactionStandard, it.Variant)
				assert.Equal(t, beans.RelatedAccount{ID: account.ID, Name: account.Name}, it.Account)
				assert.Equal(t, beans.OptionalWrap(beans.RelatedCategory{ID: category.ID, Name: category.Name}), it.Category)
				assert.Equal(t, beans.OptionalWrap(beans.RelatedPayee{ID: payee.ID, Name: payee.Name}), it.Payee)
			})
		})

		t.Run("can get off-budget variant", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			account := c.Account(AccountOpts{OffBudget: true})

			// make transaction
			transaction := c.Transaction(TransactionOpts{Account: account})

			// get transactions and verify
			res, err := interactor.TransactionGetAll(t, c.ctx)
			require.NoError(t, err)

			assert.Equal(t, 1, len(res))
			findTransaction(t, res, transaction.ID, func(it beans.TransactionWithRelations) {
				assert.Equal(t, beans.TransactionOffBudget, it.Variant)
			})
		})

		t.Run("filters by budget", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)
			c2 := makeUserAndBudget(t, interactor)

			// make transaction
			c.Transaction(TransactionOpts{})

			// budget 2 should have zero transactions
			res, err := interactor.TransactionGetAll(t, c2.ctx)
			require.NoError(t, err)

			assert.Equal(t, 0, len(res))
		})
	})
}
