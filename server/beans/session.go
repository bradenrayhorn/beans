package beans

type SessionID string

type Session struct {
	ID     SessionID
	UserID ID
}

type SessionRepository interface {
	Create(userID ID) (Session, error)
	Get(id SessionID) (Session, error)
	Delete(id SessionID) error
}
