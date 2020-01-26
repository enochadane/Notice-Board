package application

import (
	"github.com/amthesonofGod/Notice-Board/entity"
)

var b = [2]entity.Application{
	entity.Application{FullName: "moti", Email: "king@zmountain", Phone: "12765", Letter: "info le hulum", Resume: "resume", PostID: 1, UserID: 1},
	entity.Application{FullName: "king", Email: "king@zvalley", Phone: "1254233", Letter: "information", Resume: "resume again", PostID: 2, UserID: 2},
}
var a []entity.Application = b[0:len(b)]

//ApplicationRepositoryMock ...
type ApplicationRepositoryMock interface {
	ApplicationsMock() ([]entity.Application, []error)
	ApplicationMock(id uint) (*entity.Application, []error)
	UpdateApplicationMock(comment *entity.Application) (*entity.Application, []error)
	DeleteApplicationMock(id uint) (*entity.Application, []error)
	StoreApplicationMock(comment *entity.Application) (*entity.Application, []error)
}

//ApplicationsMock ...
func ApplicationsMock() ([]entity.Application, error) {
	return a, nil
}

// ApplicationMock ...
func ApplicationMock(id uint) (*entity.Application, error) {

	for _, s := range a {
		if s.ID == id {
			return &s, nil
		}
	}

	return nil, nil
}

// UpdateApplicationMock ...
func UpdateApplicationMock(user *entity.Application) (*entity.Application, error) {
	a[1] = *user

	return user, nil

}

// DeleteApplicationMock ...
func DeleteApplicationMock(id uint) (*entity.Application, error) {
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

//StoreApplicationMock ...
func StoreApplicationMock(user *entity.Application) (*entity.Application, error) {
	a[len(a)] = *user

	return user, nil
}
