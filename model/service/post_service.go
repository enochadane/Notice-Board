package service

import (
	"github.com/motikingo/Notice-Board/entity"
	"github.com/motikingo/Notice-Board/model"
)

//PostServiceImpl ...
type PostServiceImpl struct {
	postRepo model.PostRepository
}

// NewPostServiceImpl ...
func NewPostServiceImpl(PostRepos model.PostRepository) *PostServiceImpl {
	return &PostServiceImpl{postRepo: PostRepos}
}

//Posts ...
func (ps *PostServiceImpl) Posts() ([]entity.Post, error) {

	posts, err := ps.postRepo.Posts()

	if err != nil {
		return nil, err
	}

	return posts, nil
}

// StorePost ...
func (ps *PostServiceImpl) StorePost(post entity.Post) error {

	err := ps.postRepo.StorePost(post)

	if err != nil {
		return err
	}

	return nil
}

//Post ...
func (ps *PostServiceImpl) Post(id int) (entity.Post, error) {

	post, err := ps.postRepo.Post(id)

	if err != nil {
		return post, err
	}

	return post, nil
}

//UpdatePost ...
func (ps *PostServiceImpl) UpdatePost(post entity.Post) error {

	err := ps.postRepo.UpdatePost(post)

	if err != nil {
		return err
	}

	return nil
}

// DeletePost ...
func (ps *PostServiceImpl) DeletePost(id int) error {

	err := ps.postRepo.DeletePost(id)

	if err != nil {
		return err
	}

	return nil
}
