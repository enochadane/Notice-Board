package User

import (
	"github.com/amthesonofGod/Notice-Board/entity"
)

var b = [2]entity.User{
	entity.User{Name: "moti", Email: "king@zmountain", Password: "12"},
	entity.User{Name: "king", Email: "king@zvalley", Password: "123"},
}
var a []entity.User = b[0:len(b)]

//UserRepositoryMock ...
type UserRepositoryMock interface {
	UsersMock() ([]entity.User, []error)
	UserMock(id uint) (*entity.User, []error)
	UpdateUserMock(user *entity.User) (*entity.User, []error)
	DeleteUserMock(id uint) (*entity.User, []error)
	StoreUserMock(user *entity.User) (*entity.User, []error)
	StoreSession(session *entity.UserSession) (*entity.UserSession, []error)
	Session(uuid string) (*entity.UserSession, []error)
	DeleteSession(uuid string) (*entity.UserSession, []error)
}

//UsersMock ...
func UsersMock() ([]entity.User, error) {
	return a, nil
}

// UserMock ...
func UserMock(id uint) (*entity.User, error) {

	for _, s := range a {
		if s.ID == id {
			return &s, nil
		}
	}

	return nil, nil
}

// UpdateUserMock ...
func UpdateUserMock(user *entity.User) (*entity.User, error) {
	a[1] = *user

	return user, nil

}

// DeleteUserMock ...
func DeleteUserMock(id uint) (*entity.User, error) {
	for i, q := range a {
		if id == q.ID && i == 0 {

			c := a[i+1 : len(a)-1]
			for j := i; j < len(c); j++ {
				a[j] = c[j-i]
			}
			return &a[i], nil
		}
	}
	return nil, nil

}

//StoreUserMock ...
func StoreUserMock(user *entity.User) (*entity.User, error) {
	a[len(a)] = *user

	return user, nil
}
