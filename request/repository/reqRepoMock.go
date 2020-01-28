package repository

import (
	"errors"

	"github.com/amthesonofGod/Notice-Board/entity"
	"github.com/amthesonofGod/Notice-Board/request"
	"github.com/jinzhu/gorm"
)

// RequestGormRepoMock implements request.RequestRepository interface
type RequestGormRepoMock struct {
	conn *gorm.DB
}

// NewRequestGormRepoMock returns new object of RequestGormRepo
func NewRequestGormRepoMock(db *gorm.DB) request.RequestRepository {
	return &RequestGormRepoMock{conn: db}
}

// Requests returns all user event join requests stored in the database
func (reqRepo *RequestGormRepoMock) Requests() ([]entity.Request, []error) {
	reqs := []entity.Request{entity.ReqMock}

	return reqs, nil

}

// Request retrieves a user event join request from the database by its id
func (reqRepo *RequestGormRepoMock) Request(id uint) (*entity.Request, []error) {
	req := entity.Request{}

	if id != 1 {
		return nil, []error{errors.New("request never found")}
	}
	return &req, nil
}

// UpdateRequest updates a given user event join request in the database
func (reqRepo *RequestGormRepoMock) UpdateRequest(request *entity.Request) (*entity.Request, []error) {
	req := request

	return req, nil
}

// DeleteRequest deletes a given user event join request from the database
func (reqRepo *RequestGormRepoMock) DeleteRequest(id uint) (*entity.Request, []error) {
	req := entity.ReqMock

	if id != 1 {
		return nil, []error{errors.New("post not found")}
	}

	return &req, nil
}

// StoreRequest stores a given user event join request in the database
func (reqRepo *RequestGormRepoMock) StoreRequest(request *entity.Request) (*entity.Request, []error) {
	req := request

	return req, nil
}
