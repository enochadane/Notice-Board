package request

import "github.com/amthesonofGod/Notice-Board/entity"

// RequestService specifies user join request related service
type RequestService interface{
	Requests() ([]entity.Request, []error)
	Request(id uint) (*entity.Request, []error)
	UpdateRequest(request *entity.Request) (*entity.Request, []error)
	DeleteRequest(id uint) (*entity.Request, []error)
	StoreRequest(request *entity.Request) (*entity.Request, []error)
}
