package user

import "github.com/amthesonofGod/Notice-Board/entity"

// UserRepository specifies user database operations

//UserRepository ...
type UserRepository interface {
	Users() ([]entity.User, []error)
	User(id uint) (*entity.User, []error)
	UpdateUser(user *entity.User) (*entity.User, []error)
	DeleteUser(id uint) (*entity.User, []error)
	StoreUser(user *entity.User) (*entity.User, []error)
	StoreSession(session *entity.UserSession) (*entity.UserSession, []error)
	Session(uuid string) (*entity.UserSession, []error)
	DeleteSession(uuid string) (*entity.UserSession, []error)
}

// SessionRepository ...
type SessionRepository interface {
	Session(sessionID string) (*entity.UserSession, []error)
	StoreSession(session *entity.UserSession) (*entity.UserSession, []error)
	DeleteSession(sessionID string) (*entity.UserSession, []error)
}
