package post

import "github.com/amthesonofGod/Notice-Board/entity"

// PostService specifies post services
type PostService interface {
	Posts() ([]entity.Post, []error)
	Post(id uint) (*entity.Post, []error)
	UpdatePost(post *entity.Post) (*entity.Post, []error)
	DeletePost(id uint) (*entity.Post, []error)
	StorePost(post *entity.Post) (*entity.Post, []error)
	// StoreSession(session *entity.PostSession) (*entity.PostSession, []error)
	// Session(uuid string) (*entity.PostSession, []error)
	// DeleteSession(uuid string) (*entity.PostSession, []error)
}