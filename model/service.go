package model

import "NoticeBoard/entity"

type UserService interface{
	Users() ([]entity.User, error)
	User(id int) (entity.User, error)
	UpdateUser(user entity.User) error
	DeleteUser(id int) error
	StoreUser(user entity.User) error
}

type CompanyService interface {
	Companies() ([]entity.Company, error)
	Company(id int) (entity.Company, error)
	UpdateCompany(company entity.Company) error
	DeleteCompany(id int) error
	StoreCompany(company entity.Company) error
}