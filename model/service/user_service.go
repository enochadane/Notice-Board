package service

import (
	"NoticeBoard/entity"
	"NoticeBoard/model"
)

type UserServiceImpl struct {
	userRepo model.UserRepository
}

func NewUserServiceImpl(UserRepos model.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{userRepo: UserRepos}
}

func (us *UserServiceImpl) Users() ([]entity.User, error)  {
	
	users, err := us.userRepo.Users()

	if err != nil {
		return nil, err
	}

	return users, nil
}

func (us *UserServiceImpl) StoreUser(user entity.User) error {
	
	err := us.userRepo.StoreUser(user)

	if err != nil {
		return err
	}

	return nil
}

func (us *UserServiceImpl) User(id int) (entity.User, error) {

	u, err := us.userRepo.User(id)

	if err != nil {
		return u, err
	}

	return u, nil
}

func (us *UserServiceImpl) UpdateUser(user entity.User) error {
	
	err := us.userRepo.UpdateUser(user)

	if err != nil {
		return err
	}

	return nil
}

func (us *UserServiceImpl) DeleteUser(id int) error {
	
	err := us.userRepo.DeleteUser(id)

	if err != nil {
		return err
	}

	return nil
}

