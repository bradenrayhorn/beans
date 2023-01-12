package beans

type AuthContext struct {
	userID ID
}

func NewAuthContext(userID ID) *AuthContext {
	return &AuthContext{userID}
}

func (c *AuthContext) UserID() ID {
	return c.userID
}

func NewBudgetAuthContext(auth *AuthContext, budget *Budget) (*BudgetAuthContext, error) {
	if !budget.UserHasAccess(auth.UserID()) {
		return nil, NewError(EFORBIDDEN, "No access to budget")
	}

	return &BudgetAuthContext{budget, auth}, nil
}

type BudgetAuthContext struct {
	budget *Budget

	*AuthContext
}

func (c *BudgetAuthContext) BudgetID() ID {
	return c.budget.ID
}
