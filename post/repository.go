package post

import "github.com/amthesonofGod/Notice-Board/entity"

// PostRepository specifies post database operations
<<<<<<< HEAD
=======
//PostRepository ...
>>>>>>> 56480e1450127de4cec062eea25b723b5216035f
type PostRepository interface {
	Posts() ([]entity.Post, []error)
	Post(id uint) (*entity.Post, []error)
	UpdatePost(post *entity.Post) (*entity.Post, []error)
	DeletePost(id uint) (*entity.Post, []error)
	StorePost(post *entity.Post) (*entity.Post, []error)
}
