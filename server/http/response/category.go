package response

import "github.com/bradenrayhorn/beans/server/beans"

type AssociatedCategory struct {
	ID   beans.ID   `json:"id"`
	Name beans.Name `json:"name"`
}

type Category struct {
	ID   beans.ID   `json:"id"`
	Name beans.Name `json:"name"`
}

type CategoryGroup struct {
	ID         beans.ID   `json:"id"`
	Name       beans.Name `json:"name"`
	IsIncome   bool       `json:"is_income"`
	Categories []Category `json:"categories"`
}

type CreateCategoryResponse Data[ID]
type CreateCategoryGroupResponse Data[ID]
type GetCategoriesResponse Data[[]CategoryGroup]
type GetCategoryResponse Data[Category]
type GetCategoryGroupResponse Data[CategoryGroup]
