package user

import "github.com/amthesonofGod/Notice-Board/entity"

// UserService specifies user services
type UserService interface {
	Users() ([]entity.User, []error)
	User(id uint) (*entity.User, []error)
	UserByEmail(email string) (*entity.User, []error)
	UpdateUser(user *entity.User) (*entity.User, []error)
	DeleteUser(id uint) (*entity.User, []error)
	StoreUser(user *entity.User) (*entity.User, []error)
	PhoneExists(phone string) bool
	EmailExists(email string) bool
	// UserRoles(*entity.User) ([]entity.Role, []error)
}

//SessionService ...
type SessionService interface {
	Session(sessionID string) (*entity.UserSession, []error)
	StoreSession(session *entity.UserSession) (*entity.UserSession, []error)
	DeleteSession(sessionID string) (*entity.UserSession, []error)
}
