package request

import "github.com/amthesonofGod/Notice-Board/entity"

// RequestRepository specifies user join request related database operations
type RequestRepository interface {
	Requests() ([]entity.Request, []error)
	Request(id uint) (*entity.Request, []error)
	UpdateRequest(request *entity.Request) (*entity.Request, []error)
	DeleteRequest(id uint) (*entity.Request, []error)
	StoreRequest(request *entity.Request) (*entity.Request, []error)
}
