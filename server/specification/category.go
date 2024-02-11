package specification

import (
	"testing"

	"github.com/bradenrayhorn/beans/server/beans"
	"github.com/bradenrayhorn/beans/server/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func testCategory(t *testing.T, interactor Interactor) {

	t.Run("create and get group", func(t *testing.T) {
		t.Run("cannot create with invalid name", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			_, err := interactor.CategoryGroupCreate(t, c.ctx, "")
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("can create and get", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			// create group
			groupID, err := interactor.CategoryGroupCreate(t, c.ctx, "General Expenses")
			require.NoError(t, err)

			// get group
			group, err := interactor.CategoryGroupGet(t, c.ctx, groupID)
			require.NoError(t, err)

			assert.Equal(t, groupID, group.ID)
			assert.Equal(t, beans.Name("General Expenses"), group.Name)
			assert.Equal(t, false, group.IsIncome)
			assert.Empty(t, group.Categories)
		})

		t.Run("can get with categories", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)
			group := c.CategoryGroup(CategoryGroupOpts{})
			category := c.Category(CategoryOpts{Group: group})

			// get group
			res, err := interactor.CategoryGroupGet(t, c.ctx, group.ID)
			require.NoError(t, err)

			assert.Equal(t, res.Categories, []beans.RelatedCategory{
				{ID: category.ID, Name: category.Name},
			})
		})

		t.Run("cannot get non-existent group", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			_, err := interactor.CategoryGroupGet(t, c.ctx, beans.NewBeansID())
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})

		t.Run("cannot group of another user", func(t *testing.T) {
			c1 := makeUserAndBudget(t, interactor)
			c2 := makeUserAndBudget(t, interactor)
			group := c1.CategoryGroup(CategoryGroupOpts{})

			_, err := interactor.CategoryGroupGet(t, c2.ctx, group.ID)
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})
	})

	t.Run("create and get category", func(t *testing.T) {
		t.Run("cannot create with invalid name", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)
			group := c.CategoryGroup(CategoryGroupOpts{})

			_, err := interactor.CategoryCreate(t, c.ctx, group.ID, "")
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("cannot create with no group id", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			_, err := interactor.CategoryCreate(t, c.ctx, beans.EmptyID(), "Electric")
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("cannot create with group from other budget", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)
			c2 := makeUserAndBudget(t, interactor)
			group := c2.CategoryGroup(CategoryGroupOpts{})

			_, err := interactor.CategoryCreate(t, c.ctx, group.ID, "Electric")
			testutils.AssertErrorCode(t, err, beans.EINVALID)
		})

		t.Run("can create and get", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)
			group := c.CategoryGroup(CategoryGroupOpts{})

			// create category
			id, err := interactor.CategoryCreate(t, c.ctx, group.ID, "Electric")
			require.NoError(t, err)

			// get category and check values
			category, err := interactor.CategoryGet(t, c.ctx, id)
			require.NoError(t, err)

			assert.Equal(t, id, category.ID)
			assert.Equal(t, beans.Name("Electric"), category.Name)
			assert.Equal(t, group.ID, category.GroupID)
		})

		t.Run("cannot get with non-existent id", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)

			_, err := interactor.CategoryGet(t, c.ctx, beans.NewBeansID())
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})

		t.Run("cannot get with category from other budget", func(t *testing.T) {
			c1 := makeUserAndBudget(t, interactor)
			c2 := makeUserAndBudget(t, interactor)
			group := c2.CategoryGroup(CategoryGroupOpts{})
			category := c2.Category(CategoryOpts{Group: group})

			_, err := interactor.CategoryGet(t, c1.ctx, category.ID)
			testutils.AssertErrorCode(t, err, beans.ENOTFOUND)
		})
	})

	t.Run("get all", func(t *testing.T) {
		t.Run("can get all", func(t *testing.T) {
			c := makeUserAndBudget(t, interactor)
			group1 := c.CategoryGroup(CategoryGroupOpts{})
			category1 := c.Category(CategoryOpts{Group: group1})
			group2 := c.CategoryGroup(CategoryGroupOpts{})
			category2 := c.Category(CategoryOpts{Group: group2})
			group3 := c.CategoryGroup(CategoryGroupOpts{})

			// get all categories in the budget
			res, err := interactor.CategoryGetAll(t, c.ctx)
			require.NoError(t, err)

			// check that all the groups exist in the response
			findCategoryGroup(t, res, group1.ID, func(it beans.CategoryGroupWithCategories) {
				assert.Equal(t, group1.ID, it.ID)
				assert.Equal(t, group1.Name, it.Name)
				assert.Equal(t, false, it.IsIncome)
				assert.Equal(t, []beans.RelatedCategory{{ID: category1.ID, Name: category1.Name}}, it.Categories)
			})

			findCategoryGroup(t, res, group2.ID, func(it beans.CategoryGroupWithCategories) {
				assert.Equal(t, group2.ID, it.ID)
				assert.Equal(t, group2.Name, it.Name)
				assert.Equal(t, false, it.IsIncome)
				assert.Equal(t, []beans.RelatedCategory{{ID: category2.ID, Name: category2.Name}}, it.Categories)
			})

			findCategoryGroup(t, res, group3.ID, func(it beans.CategoryGroupWithCategories) {
				assert.Equal(t, group3.ID, it.ID)
				assert.Equal(t, group3.Name, it.Name)
				assert.Equal(t, false, it.IsIncome)
				assert.Equal(t, []beans.RelatedCategory{}, it.Categories)
			})
		})

		t.Run("filters by budget", func(t *testing.T) {
			c1 := makeUserAndBudget(t, interactor)
			c2 := makeUserAndBudget(t, interactor)
			c2.CategoryGroup(CategoryGroupOpts{})

			res, err := interactor.CategoryGetAll(t, c1.ctx)
			require.NoError(t, err)

			// only the income category should exist
			assert.Equal(t, 1, len(res))
			assert.Equal(t, true, res[0].IsIncome)
		})
	})
}
