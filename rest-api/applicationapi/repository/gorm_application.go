package repository

import (
	"github.com/amthesonofGod/Notice-Board/entity"
	"github.com/amthesonofGod/Notice-Board/rest-api/application"
	"github.com/jinzhu/gorm"
)

// ApplicationGormRepo implements application.ApplicationRepository interface
type ApplicationGormRepo struct {
	conn *gorm.DB
}

// type ApplicationRepository interface {
// 	Applications() ([]entity.Application, []error)
// 	Application(id uint) (*entity.Application, []error)
// 	UpdateApplication(comment *entity.Application) (*entity.Application, []error)
// 	DeleteApplication(id uint) (*entity.Application, []error)
// 	StoreApplication(comment *entity.Application) (*entity.Application, []error)
// }

// NewApplicationGormRepo returns new object of ApplicationGormRepo
func NewApplicationGormRepo(db *gorm.DB) application.ApplicationRepository {
	return &ApplicationGormRepo{conn: db}
}

// Applications returns all user Applications stored in the database
func (appRepo *ApplicationGormRepo) Applications() ([]entity.Application, []error) {
	apps := []entity.Application{}
	errs := appRepo.conn.Find(&apps).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return apps, errs
}

// Application retrieves a user application from the database by its id
func (appRepo *ApplicationGormRepo) Application(id uint) (*entity.Application, []error) {
	app := entity.Application{}
	errs := appRepo.conn.First(&app, id).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return &app, errs
}

// UpdateApplication updates a given user application in the database
func (appRepo *ApplicationGormRepo) UpdateApplication(application *entity.Application) (*entity.Application, []error) {
	app := application
	errs := appRepo.conn.Save(app).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return app, errs
}

// DeleteApplication deletes a given user application from the database
func (appRepo *ApplicationGormRepo) DeleteApplication(id uint) (*entity.Application, []error) {
	app, errs := appRepo.Application(id)

	if len(errs) > 0 {
		return nil, errs
	}

	errs = appRepo.conn.Delete(app, id).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return app, errs
}

// StoreApplication stores a given user application in the database
func (appRepo *ApplicationGormRepo) StoreApplication(application *entity.Application) (*entity.Application, []error) {
	app := application
	errs := appRepo.conn.Create(app).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return app, errs
}
