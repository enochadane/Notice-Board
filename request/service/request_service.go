package service

import (
	"github.com/amthesonofGod/Notice-Board/entity"
	"github.com/amthesonofGod/Notice-Board/request"
)

// RequestService implements request.RequestService interface
type RequestService struct {
	requestRepo request.RequestRepository
}

// NewRequestService returns a new RequestService object
func NewRequestService(reqRepo request.RequestRepository) request.RequestService {
	return &RequestService{requestRepo: reqRepo}
}

// Requests returns all stored requests
func (rs *RequestService) Requests() ([]entity.Request, []error) {
	reqs, errs := rs.requestRepo.Requests()
	if len(errs) > 0 {
		return nil, errs
	}
	return reqs, errs
}

// Request retrieves stored request by its id
func (rs *RequestService) Request(id uint) (*entity.Request, []error) {
	req, errs := rs.requestRepo.Request(id)
	if len(errs) > 0 {
		return nil, errs
	}
	return req, errs
}

// UpdateRequest updates a given request
func (rs *RequestService) UpdateRequest(request *entity.Request) (*entity.Request, []error) {
	req, errs := rs.requestRepo.UpdateRequest(request)
	if len(errs) > 0 {
		return nil, errs
	}
	return req, errs
}

// DeleteRequest deletes a given request
func (rs *RequestService) DeleteRequest(id uint) (*entity.Request, []error) {
	req, errs := rs.requestRepo.DeleteRequest(id)
	if len(errs) > 0 {
		return nil, errs
	}
	return req, errs
}

// StoreRequest stores a given request
func (rs *RequestService) StoreRequest(request *entity.Request) (*entity.Request, []error) {
	req, errs := rs.requestRepo.StoreRequest(request)
	if len(errs) > 0 {
		return nil, errs
	}
	return req, errs
}
