package service

import (
	"github.com/amthesonofGod/Notice-Board/entity"
	"github.com/amthesonofGod/Notice-Board/rest-api/application"
)

// ApplicationService implements application.ApplicationService interface
type ApplicationService struct {
	applicationRepo application.ApplicationRepository
}

// NewApplicationService returns a new ApplicationService object
func NewApplicationService(appRepo application.ApplicationRepository) application.ApplicationService {
	return &ApplicationService{applicationRepo: appRepo}
}

// Applications returns all stored applications
func (as *ApplicationService) Applications() ([]entity.Application, []error) {
	apps, errs := as.applicationRepo.Applications()
	if len(errs) > 0 {
		return nil, errs
	}
	return apps, errs
}

// Application retrieves stored application by its id
func (as *ApplicationService) Application(id uint) (*entity.Application, []error) {
	app, errs := as.applicationRepo.Application(id)
	if len(errs) > 0 {
		return nil, errs
	}
	return app, errs
}

// UpdateApplication updates a given application
func (as *ApplicationService) UpdateApplication(application *entity.Application) (*entity.Application, []error) {
	app, errs := as.applicationRepo.UpdateApplication(application)
	if len(errs) > 0 {
		return nil, errs
	}
	return app, errs
}

// DeleteApplication deletes a given application
func (as *ApplicationService) DeleteApplication(id uint) (*entity.Application, []error) {
	app, errs := as.applicationRepo.DeleteApplication(id)
	if len(errs) > 0 {
		return nil, errs
	}
	return app, errs
}

// StoreApplication stores a given application
func (as *ApplicationService) StoreApplication(application *entity.Application) (*entity.Application, []error) {
	app, errs := as.applicationRepo.StoreApplication(application)
	if len(errs) > 0 {
		return nil, errs
	}
	return app, errs
}
