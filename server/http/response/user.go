package response

import "github.com/bradenrayhorn/beans/server/beans"

type User struct {
	ID       beans.ID `json:"id"`
	Username string   `json:"username"`
}

type GetMeResponse User
