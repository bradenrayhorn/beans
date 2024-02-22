package specification

import (
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testPayee(t *testing.T, interactor Interactor) {

	t.Run("create", func(t *testing.T) {

		t.Run("name is required", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			_, err := interactor.PayeeCreate(t, c.ctx, beans.Name(""))
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("can create and get", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			// create payee
			id, err := interactor.PayeeCreate(t, c.ctx, beans.Name("Work"))
			require.NoError(t, err)

			// try to get payee
			payee, err := interactor.PayeeGet(t, c.ctx, id)
			require.NoError(t, err)

			assert.False(t, payee.ID.Empty())
			assert.Equal(t, beans.Name("Work"), payee.Name)
		})
	})

	t.Run("get", func(t *testing.T) {

		t.Run("cannot get non-existent payee", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			_, err := interactor.PayeeGet(t, c.ctx, beans.NewID())
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})

		t.Run("cannot payee from another budget", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)
			c2 := makeUserAndBudget(t, interactor)

			payee := c2.Payee(PayeeOpts{})

			_, err := interactor.PayeeGet(t, c.ctx, payee.ID)
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})
	})

	t.Run("get all", func(t *testing.T) {

		t.Run("can get all", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			// create payees
			payee1 := c.Payee(PayeeOpts{})
			payee2 := c.Payee(PayeeOpts{})

			// get payees
			res, err := interactor.PayeeGetAll(t, c.ctx)
			require.NoError(t, err)

			// check payees are proper
			assert.Equal(t, 2, len(res))

			findPayee(t, res, payee1.ID, func(it beans.Payee) {
				assert.Equal(t, payee1.ID, it.ID)
				assert.Equal(t, payee1.Name, it.Name)
			})
			findPayee(t, res, payee2.ID, func(it beans.Payee) {
				assert.Equal(t, payee2.ID, it.ID)
				assert.Equal(t, payee2.Name, it.Name)
			})
		})

		t.Run("filters by budget", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)
			c2 := makeUserAndBudget(t, interactor)

			// make a payee in budget 2
			c2.Payee(PayeeOpts{})

			// get payees
			res, err := interactor.PayeeGetAll(t, c.ctx)
			require.NoError(t, err)

			// there should be no payees
			assert.Equal(t, 0, len(res))
		})
	})

}
