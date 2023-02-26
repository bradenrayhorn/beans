package httpcontext

type key string

func (k key) String() string {
	return string(k)
}

var Auth = key("auth")
var Budget = key("budget")
var BudgetAuth = key("budget_auth")
var UserID = key("userID")
