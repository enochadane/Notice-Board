package repository

import (
	"errors"

	"github.com/amthesonofGod/Notice-Board/entity"
	"github.com/amthesonofGod/Notice-Board/post"
	"github.com/jinzhu/gorm"
)

// PostGormRepoMock implements the post.PostRepository interface
type PostGormRepoMock struct {
	conn *gorm.DB
}

// NewPostGormRepoMock will create a new object of PostGormRepoMock
func NewPostGormRepoMock(db *gorm.DB) post.PostRepository {
	return &PostGormRepoMock{conn: db}
}

// Posts returns all posts stored in the database
func (pRepo *PostGormRepoMock) Posts() ([]entity.Post, []error) {

	posts := []entity.Post{entity.PostMock}

	return posts, nil
}

// Post retrieve a post from the database by its id
func (pRepo *PostGormRepoMock) Post(id uint) (*entity.Post, []error) {

	if id == 1 {
		return &entity.PostMock, nil
	}

	return nil, nil

}

// UpdatePost updates a given post in the database
func (pRepo *PostGormRepoMock) UpdatePost(post *entity.Post) (*entity.Post, []error) {

	pst := entity.PostMock

	return &pst, nil
}

// DeletePost deletes a given post from the database
func (pRepo *PostGormRepoMock) DeletePost(id uint) (*entity.Post, []error) {
	post := entity.PostMock
	if id != 1 {
		return nil, []error{errors.New("post not found")}
	}

	return &post, nil
}

// StorePost stores a given post in the database
func (pRepo *PostGormRepoMock) StorePost(post *entity.Post) (*entity.Post, []error) {
	pst := post

	return pst, nil
}
