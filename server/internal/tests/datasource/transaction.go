package datasource

import (
	"context"
	"testing"
	"time"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testTransaction(t *testing.T, ds beans.DataSource) {
	factory := testutils.NewFactory(t, ds)

	transactionRepository := ds.TransactionRepository()
	ctx := context.Background()

	t.Run("create and get", func(t *testing.T) {
		t.Run("can create", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()
			account := factory.Account(beans.Account{BudgetID: budget.ID})
			payee := factory.Payee(beans.Payee{BudgetID: budget.ID})
			category := factory.Category(beans.Category{BudgetID: budget.ID})

			err := transactionRepository.Create(
				ctx,
				[]beans.Transaction{{
					ID:         beans.NewID(),
					AccountID:  account.ID,
					CategoryID: category.ID,
					PayeeID:    payee.ID,
					Amount:     beans.NewAmount(5, 0),
					Date:       beans.NewDate(time.Now()),
					Notes:      beans.NewTransactionNotes("notes"),
				}},
			)
			require.Nil(t, err)
		})

		t.Run("can create with empty optional fields", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()
			account := factory.Account(beans.Account{BudgetID: budget.ID})

			transaction := beans.Transaction{
				ID:        beans.NewID(),
				AccountID: account.ID,
				Amount:    beans.NewAmount(5, 0),
				Date:      testutils.NewDate(t, "2022-08-28"),
			}
			require.Nil(t, transactionRepository.Create(ctx, []beans.Transaction{transaction}))
		})

		t.Run("can create multiple transactions, with a transfer_id", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()
			account1 := factory.Account(beans.Account{BudgetID: budget.ID})
			account2 := factory.Account(beans.Account{BudgetID: budget.ID})

			id1 := beans.NewID()
			id2 := beans.NewID()
			account1Transaction := beans.Transaction{
				ID:         id1,
				AccountID:  account1.ID,
				Amount:     beans.NewAmount(5, 0),
				Date:       beans.NewDate(time.Now()),
				TransferID: id2,
			}
			account2Transaction := beans.Transaction{
				ID:         id2,
				AccountID:  account2.ID,
				Amount:     beans.NewAmount(-5, 0),
				Date:       beans.NewDate(time.Now()),
				TransferID: id1,
			}

			err := transactionRepository.Create(
				ctx,
				[]beans.Transaction{account1Transaction, account2Transaction},
			)
			require.NoError(t, err)
		})

		t.Run("can create", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()
			account := factory.Account(beans.Account{BudgetID: budget.ID})

			parent := beans.Transaction{
				ID:        beans.NewID(),
				AccountID: account.ID,
				Amount:    beans.NewAmount(5, 0),
				Date:      beans.NewDate(time.Now()),
				IsSplit:   true,
			}
			split1 := beans.Transaction{
				ID:         beans.NewID(),
				AccountID:  account.ID,
				CategoryID: factory.Category(beans.Category{}).ID,
				Amount:     beans.NewAmount(2, 0),
				Date:       beans.NewDate(time.Now()),
				SplitID:    parent.ID,
			}
			split2 := beans.Transaction{
				ID:         beans.NewID(),
				AccountID:  account.ID,
				CategoryID: factory.Category(beans.Category{}).ID,
				Amount:     beans.NewAmount(3, 0),
				Date:       beans.NewDate(time.Now()),
				SplitID:    parent.ID,
			}

			err := transactionRepository.Create(
				ctx,
				[]beans.Transaction{parent, split1, split2},
			)
			require.NoError(t, err)
		})
	})

	t.Run("cannot get nonexistant", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		_, err := transactionRepository.Get(ctx, budget.ID, beans.NewID())
		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
	})

	t.Run("cannot get for other budget", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		budget2, _ := factory.MakeBudgetAndUser()
		transaction := factory.Transaction(budget2.ID, beans.Transaction{})

		_, err := transactionRepository.Get(ctx, budget.ID, transaction.ID)
		testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
	})

	t.Run("can get", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		account := factory.Account(beans.Account{BudgetID: budget.ID})
		payee := factory.Payee(beans.Payee{BudgetID: budget.ID})
		category := factory.Category(beans.Category{BudgetID: budget.ID})

		transaction := beans.Transaction{
			ID:         beans.NewID(),
			AccountID:  account.ID,
			CategoryID: category.ID,
			PayeeID:    payee.ID,
			Amount:     beans.NewAmount(5, 0),
			Date:       testutils.NewDate(t, "2022-08-28"),
			Notes:      beans.NewTransactionNotes("notes"),
		}
		require.Nil(t, transactionRepository.Create(ctx, []beans.Transaction{transaction}))

		res, err := transactionRepository.Get(ctx, budget.ID, transaction.ID)
		require.Nil(t, err)

		assert.Equal(t, transaction, res)
	})

	t.Run("can update", func(t *testing.T) {
		budget, _ := factory.MakeBudgetAndUser()
		account1 := factory.Account(beans.Account{BudgetID: budget.ID})
		account2 := factory.Account(beans.Account{BudgetID: budget.ID})
		payee1 := factory.Payee(beans.Payee{BudgetID: budget.ID})
		payee2 := factory.Payee(beans.Payee{BudgetID: budget.ID})
		category1 := factory.Category(beans.Category{BudgetID: budget.ID})
		category2 := factory.Category(beans.Category{BudgetID: budget.ID})

		transaction := beans.Transaction{
			ID:         beans.NewID(),
			AccountID:  account1.ID,
			CategoryID: category1.ID,
			PayeeID:    payee1.ID,
			Amount:     beans.NewAmount(5, 0),
			Date:       testutils.NewDate(t, "2022-08-28"),
			Notes:      beans.NewTransactionNotes("notes"),
		}
		require.NoError(t, transactionRepository.Create(ctx, []beans.Transaction{transaction}))

		transaction.AccountID = account2.ID
		transaction.CategoryID = category2.ID
		transaction.PayeeID = payee2.ID
		transaction.Amount = beans.NewAmount(6, 0)
		transaction.Date = testutils.NewDate(t, "2022-08-30")
		transaction.Notes = beans.NewTransactionNotes("notes 5")

		require.NoError(t, transactionRepository.Update(ctx, []beans.Transaction{transaction}))

		res, err := transactionRepository.Get(ctx, budget.ID, transaction.ID)
		require.NoError(t, err)

		assert.Equal(t, transaction, res)
	})

	t.Run("delete", func(t *testing.T) {

		t.Run("can delete", func(t *testing.T) {
			budget1, _ := factory.MakeBudgetAndUser()
			budget2, _ := factory.MakeBudgetAndUser()

			transaction1 := factory.Transaction(budget1.ID, beans.Transaction{})
			transaction2 := factory.Transaction(budget1.ID, beans.Transaction{})
			transaction3 := factory.Transaction(budget1.ID, beans.Transaction{})
			transaction4 := factory.Transaction(budget2.ID, beans.Transaction{})

			err := transactionRepository.Delete(ctx, budget1.ID, []beans.ID{transaction1.ID, transaction2.ID, transaction4.ID})
			require.NoError(t, err)

			// transaction1 and transaction2 should be deleted, they are passed in and part of budget 1.
			// transaction3 should not be deleted, it is not passed in.
			// transaction4 should not be deleted, it is passed in but not part of budget 1.
			_, err = transactionRepository.Get(ctx, budget1.ID, transaction1.ID)
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)

			_, err = transactionRepository.Get(ctx, budget1.ID, transaction2.ID)
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)

			_, err = transactionRepository.Get(ctx, budget1.ID, transaction3.ID)
			assert.NoError(t, err)

			_, err = transactionRepository.Get(ctx, budget2.ID, transaction4.ID)
			assert.NoError(t, err)
		})

		t.Run("deleting transfer deletes both ends", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()

			accountA := factory.Account(beans.Account{BudgetID: budget.ID})
			accountB := factory.Account(beans.Account{BudgetID: budget.ID})
			transactions := factory.Transfer(budget.ID, accountA, accountB, beans.NewAmount(5, 0))

			err := transactionRepository.Delete(ctx, budget.ID, []beans.ID{transactions[0].ID})
			require.NoError(t, err)

			// both transactions should be deleted
			_, err = transactionRepository.Get(ctx, budget.ID, transactions[0].ID)
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
			_, err = transactionRepository.Get(ctx, budget.ID, transactions[1].ID)
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})

		t.Run("deleting parent deletes split children", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()

			parent := factory.Transaction(budget.ID, beans.Transaction{
				IsSplit: true,
			})
			child := factory.Transaction(budget.ID, beans.Transaction{
				SplitID: parent.ID,
			})

			err := transactionRepository.Delete(ctx, budget.ID, []beans.ID{parent.ID})
			require.NoError(t, err)

			// both transactions should be deleted
			_, err = transactionRepository.Get(ctx, budget.ID, parent.ID)
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
			_, err = transactionRepository.Get(ctx, budget.ID, child.ID)
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})

		t.Run("cannot delete split directly", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()

			parent := factory.Transaction(budget.ID, beans.Transaction{
				IsSplit: true,
			})
			child := factory.Transaction(budget.ID, beans.Transaction{
				SplitID: parent.ID,
			})

			err := transactionRepository.Delete(ctx, budget.ID, []beans.ID{child.ID})
			require.NoError(t, err)

			// nothing should be deleted
			_, err = transactionRepository.Get(ctx, budget.ID, parent.ID)
			assert.NoError(t, err)
			_, err = transactionRepository.Get(ctx, budget.ID, child.ID)
			assert.NoError(t, err)
		})
	})

	t.Run("get all for budget", func(t *testing.T) {

		t.Run("can get all", func(t *testing.T) {
			budget1, _ := factory.MakeBudgetAndUser()
			budget2, _ := factory.MakeBudgetAndUser()

			account := factory.Account(beans.Account{BudgetID: budget1.ID})
			payee := factory.Payee(beans.Payee{BudgetID: budget1.ID})
			category := factory.Category(beans.Category{BudgetID: budget1.ID})

			transaction1 := factory.Transaction(budget1.ID, beans.Transaction{
				AccountID:  account.ID,
				PayeeID:    payee.ID,
				CategoryID: category.ID,
			})
			// this transaction should not be included
			factory.Transaction(budget2.ID, beans.Transaction{})

			transactions, err := transactionRepository.GetForBudget(ctx, budget1.ID)
			require.NoError(t, err)
			assert.Len(t, transactions, 1)

			assert.Equal(t, beans.TransactionWithRelations{
				ID:       transaction1.ID,
				Date:     transaction1.Date,
				Amount:   transaction1.Amount,
				Notes:    transaction1.Notes,
				Variant:  beans.TransactionStandard,
				Account:  beans.RelatedAccount{ID: account.ID, Name: account.Name, OffBudget: false},
				Category: beans.OptionalWrap(beans.RelatedCategory{ID: category.ID, Name: category.Name}),
				Payee:    beans.OptionalWrap(beans.RelatedPayee{ID: payee.ID, Name: payee.Name}),
			}, transactions[0])
		})

		t.Run("maps off budget variant", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()

			account := factory.Account(beans.Account{BudgetID: budget.ID, OffBudget: true})

			transaction := factory.Transaction(budget.ID, beans.Transaction{AccountID: account.ID})

			res, err := transactionRepository.GetForBudget(ctx, budget.ID)
			require.NoError(t, err)
			require.Equal(t, 1, len(res))

			assert.Equal(t, beans.TransactionWithRelations{
				ID:       transaction.ID,
				Date:     transaction.Date,
				Amount:   transaction.Amount,
				Notes:    transaction.Notes,
				Variant:  beans.TransactionOffBudget,
				Account:  beans.RelatedAccount{ID: account.ID, Name: account.Name, OffBudget: true},
				Category: beans.Optional[beans.RelatedCategory]{},
				Payee:    beans.Optional[beans.RelatedPayee]{},
			}, res[0])
		})

		t.Run("maps split variant", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()

			account := factory.Account(beans.Account{BudgetID: budget.ID})
			transaction := factory.Transaction(budget.ID, beans.Transaction{IsSplit: true, AccountID: account.ID})

			res, err := transactionRepository.GetForBudget(ctx, budget.ID)
			require.NoError(t, err)
			require.Equal(t, 1, len(res))

			assert.Equal(t, beans.TransactionWithRelations{
				ID:      transaction.ID,
				Date:    transaction.Date,
				Amount:  transaction.Amount,
				Notes:   transaction.Notes,
				Variant: beans.TransactionSplit,
				Account: account.ToRelated(),
			}, res[0])
		})

		t.Run("maps transfer", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()

			accountA := factory.Account(beans.Account{BudgetID: budget.ID})
			accountB := factory.Account(beans.Account{BudgetID: budget.ID})

			transactions := factory.Transfer(budget.ID, accountA, accountB, beans.NewAmount(5, 0))

			res, err := transactionRepository.GetForBudget(ctx, budget.ID)
			require.NoError(t, err)
			require.Equal(t, 2, len(res))

			assert.ElementsMatch(t, []beans.TransactionWithRelations{
				{
					ID:              transactions[0].ID,
					Date:            transactions[0].Date,
					Amount:          transactions[0].Amount,
					Notes:           transactions[0].Notes,
					Variant:         beans.TransactionTransfer,
					Account:         beans.RelatedAccount{ID: accountA.ID, Name: accountA.Name, OffBudget: false},
					TransferAccount: beans.OptionalWrap(beans.RelatedAccount{ID: accountB.ID, Name: accountB.Name, OffBudget: false}),
				},
				{
					ID:              transactions[1].ID,
					Date:            transactions[1].Date,
					Amount:          transactions[1].Amount,
					Notes:           transactions[1].Notes,
					Variant:         beans.TransactionTransfer,
					Account:         beans.RelatedAccount{ID: accountB.ID, Name: accountB.Name, OffBudget: false},
					TransferAccount: beans.OptionalWrap(beans.RelatedAccount{ID: accountA.ID, Name: accountA.Name, OffBudget: false}),
				},
			}, res)
		})

		t.Run("excludes splits", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()

			parent := factory.Transaction(budget.ID, beans.Transaction{IsSplit: true})
			factory.Transaction(budget.ID, beans.Transaction{SplitID: parent.ID})

			res, err := transactionRepository.GetForBudget(ctx, budget.ID)
			require.NoError(t, err)
			require.Equal(t, 1, len(res))

			assert.Equal(t, parent.ID, res[0].ID)
		})

		t.Run("maps off-budget transfer", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()

			accountA := factory.Account(beans.Account{BudgetID: budget.ID, OffBudget: true})
			accountB := factory.Account(beans.Account{BudgetID: budget.ID, OffBudget: true})

			transactions := factory.Transfer(budget.ID, accountA, accountB, beans.NewAmount(5, 0))

			res, err := transactionRepository.GetForBudget(ctx, budget.ID)
			require.NoError(t, err)
			require.Equal(t, 2, len(res))

			assert.ElementsMatch(t, []beans.TransactionWithRelations{
				{
					ID:              transactions[0].ID,
					Date:            transactions[0].Date,
					Amount:          transactions[0].Amount,
					Notes:           transactions[0].Notes,
					Variant:         beans.TransactionTransfer,
					Account:         beans.RelatedAccount{ID: accountA.ID, Name: accountA.Name, OffBudget: true},
					TransferAccount: beans.OptionalWrap(beans.RelatedAccount{ID: accountB.ID, Name: accountB.Name, OffBudget: true}),
				},
				{
					ID:              transactions[1].ID,
					Date:            transactions[1].Date,
					Amount:          transactions[1].Amount,
					Notes:           transactions[1].Notes,
					Variant:         beans.TransactionTransfer,
					Account:         beans.RelatedAccount{ID: accountB.ID, Name: accountB.Name, OffBudget: true},
					TransferAccount: beans.OptionalWrap(beans.RelatedAccount{ID: accountA.ID, Name: accountA.Name, OffBudget: true}),
				},
			}, res)
		})
	})

	t.Run("get with relations", func(t *testing.T) {

		t.Run("can get", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()

			account := factory.Account(beans.Account{BudgetID: budget.ID})
			payee := factory.Payee(beans.Payee{BudgetID: budget.ID})
			category := factory.Category(beans.Category{BudgetID: budget.ID})

			transaction := factory.Transaction(budget.ID, beans.Transaction{
				AccountID:  account.ID,
				PayeeID:    payee.ID,
				CategoryID: category.ID,
			})

			res, err := transactionRepository.GetWithRelations(ctx, budget.ID, transaction.ID)
			require.NoError(t, err)

			assert.Equal(t, beans.TransactionWithRelations{
				ID:       transaction.ID,
				Date:     transaction.Date,
				Amount:   transaction.Amount,
				Notes:    transaction.Notes,
				Variant:  beans.TransactionStandard,
				Account:  beans.RelatedAccount{ID: account.ID, Name: account.Name, OffBudget: false},
				Category: beans.OptionalWrap(beans.RelatedCategory{ID: category.ID, Name: category.Name}),
				Payee:    beans.OptionalWrap(beans.RelatedPayee{ID: payee.ID, Name: payee.Name}),
			}, res)
		})

		t.Run("maps off-budget variant", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()

			account := factory.Account(beans.Account{BudgetID: budget.ID, OffBudget: true})
			transaction := factory.Transaction(budget.ID, beans.Transaction{AccountID: account.ID})

			res, err := transactionRepository.GetWithRelations(ctx, budget.ID, transaction.ID)
			require.NoError(t, err)

			assert.Equal(t, beans.TransactionWithRelations{
				ID:      transaction.ID,
				Date:    transaction.Date,
				Amount:  transaction.Amount,
				Notes:   transaction.Notes,
				Variant: beans.TransactionOffBudget,
				Account: beans.RelatedAccount{ID: account.ID, Name: account.Name, OffBudget: true},
			}, res)
		})

		t.Run("maps split variant", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()

			account := factory.Account(beans.Account{BudgetID: budget.ID})
			transaction := factory.Transaction(budget.ID, beans.Transaction{IsSplit: true, AccountID: account.ID})

			res, err := transactionRepository.GetWithRelations(ctx, budget.ID, transaction.ID)
			require.NoError(t, err)

			assert.Equal(t, beans.TransactionWithRelations{
				ID:      transaction.ID,
				Date:    transaction.Date,
				Amount:  transaction.Amount,
				Notes:   transaction.Notes,
				Variant: beans.TransactionSplit,
				Account: account.ToRelated(),
			}, res)
		})

		t.Run("maps transfer", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()

			accountA := factory.Account(beans.Account{BudgetID: budget.ID})
			accountB := factory.Account(beans.Account{BudgetID: budget.ID})
			transactions := factory.Transfer(budget.ID, accountA, accountB, beans.NewAmount(7, 0))

			res, err := transactionRepository.GetWithRelations(ctx, budget.ID, transactions[0].ID)
			require.NoError(t, err)

			assert.Equal(t, beans.TransactionWithRelations{
				ID:              transactions[0].ID,
				Date:            transactions[0].Date,
				Amount:          transactions[0].Amount,
				Notes:           transactions[0].Notes,
				Variant:         beans.TransactionTransfer,
				Account:         beans.RelatedAccount{ID: accountA.ID, Name: accountA.Name, OffBudget: false},
				TransferAccount: beans.OptionalWrap(beans.RelatedAccount{ID: accountB.ID, Name: accountB.Name, OffBudget: false}),
			}, res)
		})

		t.Run("maps off-budget transfer", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()

			accountA := factory.Account(beans.Account{BudgetID: budget.ID, OffBudget: true})
			accountB := factory.Account(beans.Account{BudgetID: budget.ID, OffBudget: true})
			transactions := factory.Transfer(budget.ID, accountA, accountB, beans.NewAmount(7, 0))

			res, err := transactionRepository.GetWithRelations(ctx, budget.ID, transactions[0].ID)
			require.NoError(t, err)

			assert.Equal(t, beans.TransactionWithRelations{
				ID:              transactions[0].ID,
				Date:            transactions[0].Date,
				Amount:          transactions[0].Amount,
				Notes:           transactions[0].Notes,
				Variant:         beans.TransactionTransfer,
				Account:         beans.RelatedAccount{ID: accountA.ID, Name: accountA.Name, OffBudget: true},
				TransferAccount: beans.OptionalWrap(beans.RelatedAccount{ID: accountB.ID, Name: accountB.Name, OffBudget: true}),
			}, res)
		})

		t.Run("filters by budget", func(t *testing.T) {
			budget1, _ := factory.MakeBudgetAndUser()
			budget2, _ := factory.MakeBudgetAndUser()

			transaction := factory.Transaction(budget2.ID, beans.Transaction{})

			_, err := transactionRepository.GetWithRelations(ctx, budget1.ID, transaction.ID)
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})

		t.Run("cannot get non existent", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()

			_, err := transactionRepository.GetWithRelations(ctx, budget.ID, beans.NewID())
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})
	})

	t.Run("get splits", func(t *testing.T) {

		t.Run("filters by budget", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()
			budget2, _ := factory.MakeBudgetAndUser()

			parent := factory.Transaction(budget.ID, beans.Transaction{
				IsSplit: true,
			})
			factory.Transaction(budget.ID, beans.Transaction{
				SplitID: parent.ID,
			})

			res, err := transactionRepository.GetSplits(ctx, budget2.ID, parent.ID)
			require.NoError(t, err)
			assert.Equal(t, 0, len(res))
		})

		t.Run("filters by transaction", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()

			parent := factory.Transaction(budget.ID, beans.Transaction{
				IsSplit: true,
			})
			realParent := factory.Transaction(budget.ID, beans.Transaction{
				IsSplit: true,
			})
			factory.Transaction(budget.ID, beans.Transaction{
				SplitID: realParent.ID,
			})

			res, err := transactionRepository.GetSplits(ctx, budget.ID, parent.ID)
			require.NoError(t, err)
			assert.Equal(t, 0, len(res))
		})

		t.Run("can get splits", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()

			category := factory.Category(beans.Category{BudgetID: budget.ID})

			parent := factory.Transaction(budget.ID, beans.Transaction{
				IsSplit: true,
			})

			child := factory.Transaction(budget.ID, beans.Transaction{
				SplitID:    parent.ID,
				CategoryID: category.ID,
				Amount:     beans.NewAmount(5, 0),
				Notes:      beans.NewTransactionNotes("hi"),
			})

			res, err := transactionRepository.GetSplits(ctx, budget.ID, parent.ID)
			require.NoError(t, err)
			assert.Equal(t, 1, len(res))

			assert.Equal(t, beans.TransactionAsSplit{
				Transaction: child,
				Split: beans.Split{
					ID:       child.ID,
					Amount:   beans.NewAmount(5, 0),
					Notes:    beans.NewTransactionNotes("hi"),
					Category: category.ToRelated(),
				},
			}, res[0])
		})

		t.Run("errors if missing category", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()

			parent := factory.Transaction(budget.ID, beans.Transaction{
				IsSplit: true,
			})
			factory.Transaction(budget.ID, beans.Transaction{
				SplitID: parent.ID,
			})

			_, err := transactionRepository.GetSplits(ctx, budget.ID, parent.ID)
			assert.ErrorContains(t, err, "category null")
		})
	})

	t.Run("can get activity by category", func(t *testing.T) {

		t.Run("groups and sums", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()

			category1 := factory.Category(beans.Category{BudgetID: budget.ID})
			category2 := factory.Category(beans.Category{BudgetID: budget.ID})

			// setup 3 transactions - two in category1 and one in category2
			factory.Transaction(budget.ID, beans.Transaction{
				Amount:     beans.NewAmount(3, 0),
				CategoryID: category1.ID,
			})
			factory.Transaction(budget.ID, beans.Transaction{
				Amount:     beans.NewAmount(2, 0),
				CategoryID: category1.ID,
			})
			factory.Transaction(budget.ID, beans.Transaction{
				Amount:     beans.NewAmount(1, 0),
				CategoryID: category2.ID,
			})

			// make sure they are grouped and summed properly
			res, err := transactionRepository.GetActivityByCategory(ctx, budget.ID, beans.Date{}, beans.Date{})
			require.NoError(t, err)

			assert.Equal(t, 2, len(res))
			assert.Equal(t, beans.NewAmount(5, 0), res[category1.ID])
			assert.Equal(t, beans.NewAmount(1, 0), res[category2.ID])
		})

		t.Run("filters by date", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()

			// setup 4 transactions with varying dates
			category := factory.Category(beans.Category{BudgetID: budget.ID})
			factory.Transaction(budget.ID, beans.Transaction{
				Amount:     beans.NewAmount(3, 0),
				CategoryID: category.ID,
				Date:       testutils.NewDate(t, "2099-09-01"),
			})
			factory.Transaction(budget.ID, beans.Transaction{
				Amount:     beans.NewAmount(2, 0),
				CategoryID: category.ID,
				Date:       testutils.NewDate(t, "2022-09-01"),
			})
			factory.Transaction(budget.ID, beans.Transaction{
				Amount:     beans.NewAmount(1, 0),
				CategoryID: category.ID,
				Date:       testutils.NewDate(t, "2022-08-31"),
			})
			factory.Transaction(budget.ID, beans.Transaction{
				Amount:     beans.NewAmount(8, 0),
				CategoryID: category.ID,
				Date:       testutils.NewDate(t, "1900-08-31"),
			})

			// try to filter the transactions

			t.Run("filters by only from date", func(t *testing.T) {
				res, err := transactionRepository.GetActivityByCategory(ctx, budget.ID, testutils.NewDate(t, "2022-09-01"), beans.Date{})
				require.NoError(t, err)

				assert.Equal(t, 1, len(res))
				assert.Equal(t, beans.NewAmount(5, 0), res[category.ID])
			})

			t.Run("filters by only to date", func(t *testing.T) {
				res, err := transactionRepository.GetActivityByCategory(ctx, budget.ID, beans.Date{}, testutils.NewDate(t, "2022-08-31"))
				require.NoError(t, err)

				assert.Equal(t, 1, len(res))
				assert.Equal(t, beans.NewAmount(9, 0), res[category.ID])
			})

			t.Run("filters by both dates", func(t *testing.T) {
				res, err := transactionRepository.GetActivityByCategory(ctx, budget.ID, testutils.NewDate(t, "2022-08-01"), testutils.NewDate(t, "2022-09-30"))
				require.NoError(t, err)

				assert.Equal(t, 1, len(res))
				assert.Equal(t, beans.NewAmount(3, 0), res[category.ID])
			})

			t.Run("applies no date filter", func(t *testing.T) {
				res, err := transactionRepository.GetActivityByCategory(ctx, budget.ID, beans.Date{}, beans.Date{})
				require.NoError(t, err)

				assert.Equal(t, 1, len(res))
				assert.Equal(t, beans.NewAmount(14, 0), res[category.ID])
			})
		})

		t.Run("filters by budget", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()
			budget2, _ := factory.MakeBudgetAndUser()

			factory.Transaction(budget2.ID, beans.Transaction{
				Amount: beans.NewAmount(1, 0),
				Date:   testutils.NewDate(t, "2022-09-01"),
			})

			res, err := transactionRepository.GetActivityByCategory(ctx, budget.ID, beans.Date{}, beans.Date{})
			require.NoError(t, err)

			assert.Equal(t, 0, len(res))
		})
	})

	t.Run("can get income", func(t *testing.T) {

		t.Run("can get", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()
			incomeGroup := factory.CategoryGroup(beans.CategoryGroup{BudgetID: budget.ID, IsIncome: true})
			incomeCategory := factory.Category(beans.Category{BudgetID: budget.ID, GroupID: incomeGroup.ID})
			otherCategory := factory.Category(beans.Category{BudgetID: budget.ID})

			budget2, _ := factory.MakeBudgetAndUser()
			budget2IncomeGroup := factory.CategoryGroup(beans.CategoryGroup{BudgetID: budget2.ID, IsIncome: true})
			budget2IncomeCategory := factory.Category(beans.Category{BudgetID: budget2.ID, GroupID: budget2IncomeGroup.ID})

			// Earned $1 in September
			factory.Transaction(budget.ID, beans.Transaction{
				Amount:     beans.NewAmount(1, 0),
				Date:       testutils.NewDate(t, "2022-09-01"),
				CategoryID: incomeCategory.ID,
			})
			// Earned $2 in August
			factory.Transaction(budget.ID, beans.Transaction{
				Amount:     beans.NewAmount(2, 0),
				Date:       testutils.NewDate(t, "2022-08-31"),
				CategoryID: incomeCategory.ID,
			})
			// Earned $3 in August
			factory.Transaction(budget.ID, beans.Transaction{
				Amount:     beans.NewAmount(3, 0),
				Date:       testutils.NewDate(t, "2022-08-01"),
				CategoryID: incomeCategory.ID,
			})
			// Earned $3 in July
			factory.Transaction(budget.ID, beans.Transaction{
				Amount:     beans.NewAmount(3, 0),
				Date:       testutils.NewDate(t, "2022-07-31"),
				CategoryID: incomeCategory.ID,
			})

			// Spent $99 in August
			factory.Transaction(budget.ID, beans.Transaction{
				Amount:     beans.NewAmount(99, 0),
				Date:       testutils.NewDate(t, "2022-08-31"),
				CategoryID: otherCategory.ID,
			})
			// Spent $99 in July
			factory.Transaction(budget.ID, beans.Transaction{
				Amount:     beans.NewAmount(99, 0),
				Date:       testutils.NewDate(t, "2022-07-29"),
				CategoryID: otherCategory.ID,
			})

			// Budget 2, earned $99 in August
			factory.Transaction(budget2.ID, beans.Transaction{
				Amount:     beans.NewAmount(99, 0),
				Date:       testutils.NewDate(t, "2022-08-15"),
				CategoryID: budget2IncomeCategory.ID,
			})

			amount, err := transactionRepository.GetIncomeBetween(ctx, budget.ID, testutils.NewDate(t, "2022-08-01"), testutils.NewDate(t, "2022-08-31"))
			require.Nil(t, err)

			// August earnings for budget 1 = $5
			require.Equal(t, beans.NewAmount(5, 0), amount)
		})

		t.Run("can get with no income", func(t *testing.T) {
			budget, _ := factory.MakeBudgetAndUser()
			amount, err := transactionRepository.GetIncomeBetween(ctx, budget.ID, testutils.NewDate(t, "2022-08-01"), testutils.NewDate(t, "2022-08-31"))

			require.Nil(t, err)

			require.Equal(t, beans.NewAmount(0, 0), amount)
		})
	})
}
