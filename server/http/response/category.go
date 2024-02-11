package response

import "github.com/bradenrayhorn/beans/server/beans"

type AssociatedCategory struct {
	ID   beans.ID   `json:"id"`
	Name beans.Name `json:"name"`
}

type Category struct {
	ID      beans.ID   `json:"id"`
	Name    beans.Name `json:"name"`
	GroupID beans.ID   `json:"group_id"`
}

type CategoryGroup struct {
	ID         beans.ID             `json:"id"`
	Name       beans.Name           `json:"name"`
	IsIncome   bool                 `json:"is_income"`
	Categories []AssociatedCategory `json:"categories"`
}

type CreateCategoryResponse Data[ID]
type CreateCategoryGroupResponse Data[ID]
type GetCategoriesResponse Data[[]CategoryGroup]
type GetCategoryResponse Data[Category]
type GetCategoryGroupResponse Data[CategoryGroup]
