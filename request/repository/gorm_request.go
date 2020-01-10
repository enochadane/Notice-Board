package repository

import (
	"github.com/amthesonofGod/Notice-Board/entity"
	"github.com/amthesonofGod/Notice-Board/request"
	"github.com/jinzhu/gorm"
)
// RequestGormRepo implements request.RequestRepository interface
type RequestGormRepo struct {
	conn *gorm.DB
}

// NewRequestGormRepo returns new object of RequestGormRepo
func NewRequestGormRepo(db *gorm.DB) request.RequestRepository {
	return &RequestGormRepo{conn: db}
}

// Requests returns all user event join requests stored in the database
func (reqRepo *RequestGormRepo) Requests() ([]entity.Request, []error) {
	reqs := []entity.Request{}
	errs := reqRepo.conn.Find(&reqs).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return reqs, errs
}

// Request retrieves a user event join request from the database by its id
func (reqRepo *RequestGormRepo) Request(id uint) (*entity.Request, []error) {
	req := entity.Request{}
	errs := reqRepo.conn.First(&req, id).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return &req, errs
}

// UpdateRequest updates a given user event join request in the database
func (reqRepo *RequestGormRepo) UpdateRequest(request *entity.Request) (*entity.Request, []error) {
	req := request
	errs := reqRepo.conn.Save(req).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return req, errs
}

// DeleteRequest deletes a given user event join request from the database
func (reqRepo *RequestGormRepo) DeleteRequest(id uint) (*entity.Request, []error) {
	req, errs := reqRepo.Request(id)

	if len(errs) > 0 {
		return nil, errs
	}

	errs = reqRepo.conn.Delete(req, id).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return req, errs
}

// StoreRequest stores a given user event join request in the database
func (reqRepo *RequestGormRepo) StoreRequest(request *entity.Request) (*entity.Request, []error) {
	req := request
	errs := reqRepo.conn.Create(req).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return req, errs
}
