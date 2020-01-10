package service

import (
	"github.com/amthesonofGod/Notice-Board/entity"
	"github.com/amthesonofGod/Notice-Board/post"
)

// PostService implements post.PostRepository interface
type PostService struct {
	postRepo post.PostRepository
}

// NewPostService will create new PostService object
func NewPostService(PostRepos post.PostRepository) post.PostService {
	return &PostService{postRepo: PostRepos}
}

// Posts returns list of posts
func (ps *PostService) Posts() ([]entity.Post, []error)  {
	
	posts, errs := ps.postRepo.Posts()

	if len(errs) > 0 {
		return nil, errs
	}

	return posts, nil
}

// StorePost persists new post information
func (ps *PostService) StorePost(post *entity.Post) (*entity.Post, []error) {
	
	cmp, errs := ps.postRepo.StorePost(post)

	if len(errs) > 0 {
		return nil, errs
	}

	return cmp, nil
}

// Post returns a post object with a given id
func (ps *PostService) Post(id uint) (*entity.Post, []error) {

	post, errs := ps.postRepo.Post(id)

	if len(errs) > 0 {
		return nil, errs
	}

	return post, nil
}

// UpdatePost updates a post with new data
func (ps *PostService) UpdatePost(post *entity.Post) (*entity.Post, []error) {
	
	pst, errs := ps.postRepo.UpdatePost(post)

	if len(errs) > 0 {
		return nil, errs
	}

	return pst, nil
}

// DeletePost deletes a post by its id
func (ps *PostService) DeletePost(id uint) (*entity.Post, []error) {
	
	pst, errs := ps.postRepo.DeletePost(id)

	if len(errs) > 0 {
		return nil, errs
	}

	return pst, nil
}

