package User

import "github.com/amthesonofGod/Notice-Board/entity"

// UserRepository specifies user database operations
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
