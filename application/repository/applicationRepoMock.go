package repository

import (
	"errors"

	"github.com/amthesonofGod/Notice-Board/application"
	"github.com/amthesonofGod/Notice-Board/entity"
	"github.com/jinzhu/gorm"
)

// ApplicationGormRepoMock implements application.ApplicationRepository interface
type ApplicationGormRepoMock struct {
	conn *gorm.DB
}

// NewApplicationGormRepoMock returns new object of ApplicationGormRepo
func NewApplicationGormRepoMock(db *gorm.DB) application.ApplicationRepository {
	return &ApplicationGormRepoMock{conn: db}
}

// Applications returns all user Applications stored in the database
func (appRepo *ApplicationGormRepoMock) Applications() ([]entity.Application, []error) {
	ap := []entity.Application{entity.AplMock}

	return ap, nil
}

// Application retrieves a user application from the database by its id
func (appRepo *ApplicationGormRepoMock) Application(id uint) (*entity.Application, []error) {

	if id == 1 {
		return &entity.AplMock, nil
	}

	return nil, nil
}

// UpdateApplication updates a given user application in the database
func (appRepo *ApplicationGormRepoMock) UpdateApplication(application *entity.Application) (*entity.Application, []error) {
	apl := entity.AplMock

	return &apl, nil
}

// DeleteApplication deletes a given user application from the database
func (appRepo *ApplicationGormRepoMock) DeleteApplication(id uint) (*entity.Application, []error) {

	ap := entity.AplMock
	if id != 1 {
		return nil, []error{errors.New("post not found")}
	}
	return &ap, nil
}

// StoreApplication stores a given user application in the database
func (appRepo *ApplicationGormRepoMock) StoreApplication(application *entity.Application) (*entity.Application, []error) {
	ap := application

	return ap, nil
}
