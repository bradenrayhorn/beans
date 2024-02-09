package beans

type AuthContext struct {
	userID    ID
	sessionID SessionID
}

func NewAuthContext(userID ID, sessionID SessionID) *AuthContext {
	return &AuthContext{userID, sessionID}
}

func (c *AuthContext) UserID() ID {
	return c.userID
}

func (c *AuthContext) SessionID() SessionID {
	return c.sessionID
}

func NewBudgetAuthContext(auth *AuthContext, budget Budget) (*BudgetAuthContext, error) {
	return &BudgetAuthContext{budget, auth}, nil
}

type BudgetAuthContext struct {
	budget Budget

	*AuthContext
}

func (c *BudgetAuthContext) BudgetID() ID {
	return c.budget.ID
}
