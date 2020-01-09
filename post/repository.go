package post

import "github.com/amthesonofGod/Notice-Board/entity"

// PostRepository specifies post database operations
type PostRepository interface {
	Posts() ([]entity.Post, []error)
	Post(id uint) (*entity.Post, []error)
	UpdatePost(post *entity.Post) (*entity.Post, []error)
	DeletePost(id uint) (*entity.Post, []error)
	StorePost(post *entity.Post) (*entity.Post, []error)
}
