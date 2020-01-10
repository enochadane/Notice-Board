package post

import "github.com/amthesonofGod/Notice-Board/entity"

// PostService specifies post services
type PostService interface {
	Posts() ([]entity.Post, []error)
	Post(id uint) (*entity.Post, []error)
	UpdatePost(post *entity.Post) (*entity.Post, []error)
	DeletePost(id uint) (*entity.Post, []error)
	StorePost(post *entity.Post) (*entity.Post, []error)
<<<<<<< HEAD
}
=======
}
>>>>>>> 56480e1450127de4cec062eea25b723b5216035f
