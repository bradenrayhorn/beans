package response

import "github.com/bradenrayhorn/beans/server/beans"

type User struct {
	ID       beans.ID `json:"id"`
	Username string   `json:"username"`
}

type SessionID struct {
	SessionID beans.SessionID `json:"sessionID"`
}

type GetMe User

type Login Data[SessionID]
