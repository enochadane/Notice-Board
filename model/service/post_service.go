package service

import (
	"NoticeBoard/entity"
	"NoticeBoard/model"
)

type PostServiceImpl struct {
	postRepo model.PostRepository
}

func NewPostServiceImpl(PostRepos model.PostRepository) *PostServiceImpl {
	return &PostServiceImpl{postRepo: PostRepos}
}

func (ps *PostServiceImpl) Posts() ([]entity.Post, error) {

	posts, err := ps.postRepo.Posts()

	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (ps *PostServiceImpl) StorePost(post entity.Post) error {

	err := ps.postRepo.StorePost(post)

	if err != nil {
		return err
	}

	return nil
}

func (ps *PostServiceImpl) Post(id int) (entity.Post, error) {

	post, err := ps.postRepo.Post(id)

	if err != nil {
		return post, err
	}

	return post, nil
}

func (ps *PostServiceImpl) UpdatePost(post entity.Post) error {

	err := ps.postRepo.UpdatePost(post)

	if err != nil {
		return err
	}

	return nil
}

func (ps *PostServiceImpl) DeletePost(id int) error {

	err := ps.postRepo.DeletePost(id)

	if err != nil {
		return err
	}

	return nil
}