package repository

import (
	"github.com/amthesonofGod/Notice-Board/entity"
	"github.com/amthesonofGod/Notice-Board/post"
	"github.com/jinzhu/gorm"
)

// PostGormRepo implements the post.PostRepository interface
type PostGormRepo struct {
	conn *gorm.DB
}

// NewPostGormRepo will create a new object of PostGormRepo
func NewPostGormRepo(db *gorm.DB) post.PostRepository {
	return &PostGormRepo{conn: db}
}

// Posts returns all posts stored in the database
func (pRepo *PostGormRepo) Posts() ([]entity.Post, []error)  {
	posts := []entity.Post{}
	errs := pRepo.conn.Find(&posts).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}

	return posts, errs
}

// Post retrieve a post from the database by its id
func (pRepo *PostGormRepo) Post(id uint) (*entity.Post, []error) {
	
	post := entity.Post{}
	errs := pRepo.conn.First(&post, id).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return &post, errs
}

// UpdatePost updates a given post in the database
func (pRepo *PostGormRepo) UpdatePost(post *entity.Post) (*entity.Post, []error) {
	
	pst := post
	errs := pRepo.conn.Save(pst).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return pst, errs
}

// DeletePost deletes a given post from the database
func (pRepo *PostGormRepo) DeletePost(id uint) (*entity.Post, []error) {
	
	post, errs := pRepo.Post(id)
	if len(errs) > 0 {
		return nil, errs
	}

	errs = pRepo.conn.Delete(post, post.ID).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return post, errs
}

// StorePost stores a given post in the database
func (pRepo *PostGormRepo) StorePost(post *entity.Post) (*entity.Post, []error) {
	
	pst := post
	errs := pRepo.conn.Create(pst).GetErrors()
	if len(errs) > 0 {
		return nil, errs
	}
	return pst, errs
}

// // StoreSession stores a given session in the database
// func (pRepo *PostGormRepo) StoreSession(session *entity.PostSession) (*entity.PostSession, []error) {
// 	s := session
// 	errs := pRepo.conn.Create(s).GetErrors()
// 	if len(errs) > 0 {
// 		return nil, errs
// 	}
// 	return s, errs
// }

// // DeleteSession deletes a given session from the database
// func (pRepo *PostGormRepo) DeleteSession(uuid string) (*entity.PostSession, []error) {
// 	s, errs := pRepo.Session(uuid)
// 	if len(errs) > 0 {
// 		return nil, errs
// 	}

// 	errs = pRepo.conn.Delete(s, s.UUID).GetErrors()
// 	if len(errs) > 0 {
// 		return nil, errs
// 	}
// 	return s, errs
// }

// // Session retrieve a session from the database by its id
// func (pRepo *PostGormRepo) Session(uuid string) (*entity.PostSession, []error) {
// 	s := entity.CompanySession{}
// 	errs := pRepo.conn.Where("UUID = ?", uuid).First(&s).GetErrors()
// 	if len(errs) > 0 {
// 		return nil, errs
// 	}
// 	return &s, errs
// }
