package model

import "github.com/amthesonofGod/Notice-Board/entity"

// UserService specifies user services
type UserService interface{
	Users() ([]entity.User, []error)
	User(id uint) (*entity.User, []error)
	UpdateUser(user *entity.User) (*entity.User, []error)
	DeleteUser(id uint) (*entity.User, []error)
	StoreUser(user *entity.User) (*entity.User, []error)
}

// CompanyService specifies compny services
type CompanyService interface {
	Companies() ([]entity.Company, []error)
	Company(id uint) (*entity.Company, []error)
	UpdateCompany(company *entity.Company) (*entity.Company, []error)
	DeleteCompany(id uint) (*entity.Company, []error)
	StoreCompany(company *entity.Company) (*entity.Company ,[]error)
	StoreSession(session *entity.Session) (*entity.Session, []error)
	Session(uuid string) (*entity.Session, []error)
	DeleteSession(uuid string) (*entity.Session, []error)
}

