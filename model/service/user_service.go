package service

import (
	"github.com/motikingo/Notice-Board/entity"
	"github.com/motikingo/Notice-Board/model"
)

// UserServiceImpl ...
type UserServiceImpl struct {
	userRepo model.UserRepository
}

// NewUserServiceImpl ...
func NewUserServiceImpl(UserRepos model.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{userRepo: UserRepos}
}

// Users ...
func (us *UserServiceImpl) Users() ([]entity.User, error) {

	users, err := us.userRepo.Users()

	if err != nil {
		return nil, err
	}

	return users, nil
}

// StoreUser ...
func (us *UserServiceImpl) StoreUser(user entity.User) error {

	err := us.userRepo.StoreUser(user)

	if err != nil {
		return err
	}

	return nil
}

// User ...
func (us *UserServiceImpl) User(id int) (entity.User, error) {

	u, err := us.userRepo.User(id)

	if err != nil {
		return u, err
	}

	return u, nil
}

// UpdateUser ...
func (us *UserServiceImpl) UpdateUser(user entity.User) error {

	err := us.userRepo.UpdateUser(user)

	if err != nil {
		return err
	}

	return nil
}

// DeleteUser ...
func (us *UserServiceImpl) DeleteUser(id int) error {

	err := us.userRepo.DeleteUser(id)

	if err != nil {
		return err
	}

	return nil
}
