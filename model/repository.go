package model

import "github.com/amthesonofGod/Notice-Board/entity"

// UserRepository specifies user database operations
type UserRepository interface {
	Users() ([]entity.User, []error)
	User(id uint) (*entity.User, []error)
	UpdateUser(user *entity.User) (*entity.User, []error)
	DeleteUser(id uint) (*entity.User, []error)
	StoreUser(user *entity.User) (*entity.User, []error)
}

// CompanyRepository specifies company database operations
type CompanyRepository interface {
	Companies() ([]entity.Company, []error)
	Company(id uint) (*entity.Company, []error)
	UpdateCompany(company *entity.Company) (*entity.Company, []error)
	DeleteCompany(id uint) (*entity.Company, []error)
	StoreCompany(company *entity.Company) (*entity.Company ,[]error)
	StoreSession(session *entity.Session) (*entity.Session, []error)
	Session(uuid string) (*entity.Session, []error)
	DeleteSession(uuid string) (*entity.Session, []error)
}

