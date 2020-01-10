package application

import "github.com/amthesonofGod/Notice-Board/entity"

// ApplicationRepository specifies user job application related database operations
type ApplicationRepository interface {
	Applications() ([]entity.Application, []error)
	Application(id uint) (*entity.Application, []error)
	UpdateApplication(comment *entity.Application) (*entity.Application, []error)
	DeleteApplication(id uint) (*entity.Application, []error)
	StoreApplication(comment *entity.Application) (*entity.Application, []error)
}