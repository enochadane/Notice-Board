package model

import "NoticeBoard/entity"

type UserRepository interface {
	Users() ([]entity.User, error)
	User(id int) (entity.User, error)
	UpdateUser(user entity.User) error
	DeleteUser(id int) error
	StoreUser(user entity.User) error
}

type CompanyRepository interface {
	Companies() ([]entity.Company, error)
	Company(id int) (entity.Company, error)
	UpdateCompany(company entity.Company) error
	DeleteCompany(id int) error
	StoreCompany(company entity.Company) error
}

type PostRepository interface {
	Posts() ([]entity.Post, error)
	Post(id int) (entity.Post, error)
	UpdatePost(post entity.Post) error
	DeletePost(id int) error
	StorePost(post entity.Post) error
}
