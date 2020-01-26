package service

import (
	"github.com/amthesonofGod/Notice-Board/user"
	"github.com/amthesonofGod/Notice-Board/entity"
)

// UserService implements model.UserRepository interface
type UserService struct {
	userRepo user.UserRepository
}

// NewUserService will create new UserService object
func NewUserService(UserRepos user.UserRepository) user.UserService {
	return &UserService{userRepo: UserRepos}
}

// Users returns list of users
func (us *UserService) Users() ([]entity.User, []error) {

	users, errs := us.userRepo.Users()

	if len(errs) > 0 {
		return nil, errs
	}

	return users, nil
}

// StoreUser persists new user information
func (us *UserService) StoreUser(user *entity.User) (*entity.User, []error) {

	usr, errs := us.userRepo.StoreUser(user)

	if len(errs) > 0 {
		return nil, errs
	}

	return usr, nil
}

// User returns a user object with a given id
func (us *UserService) User(id uint) (*entity.User, []error) {

	user, errs := us.userRepo.User(id)

	if len(errs) > 0 {
		return nil, errs
	}

	return user, nil
}

// UpdateUser updates a user with new data
func (us *UserService) UpdateUser(user *entity.User) (*entity.User, []error) {

	usr, errs := us.userRepo.UpdateUser(user)

	if len(errs) > 0 {
		return nil, errs
	}

	return usr, nil
}

// DeleteUser deletes a user by its id
func (us *UserService) DeleteUser(id uint) (*entity.User, []error) {

	usr, errs := us.userRepo.DeleteUser(id)

	if len(errs) > 0 {
		return nil, errs
	}

	return usr, nil
}

// UserByEmail retrieves an application user by its email address
func (us *UserService) UserByEmail(email string) (*entity.User, []error) {
	usr, errs := us.userRepo.UserByEmail(email)
	if len(errs) > 0 {
		return nil, errs
	}
	return usr, errs
}

// PhoneExists check if there is a user with a given phone number
func (us *UserService) PhoneExists(phone string) bool {
	exists := us.userRepo.PhoneExists(phone)
	return exists
}

// EmailExists checks if there exist a user with a given email address
func (us *UserService) EmailExists(email string) bool {
	exists := us.userRepo.EmailExists(email)
	return exists
}

// UserRoles returns list of roles a user has
// func (us *UserService) UserRoles(user *entity.User) ([]entity.Role, []error) {
// 	userRoles, errs := us.userRepo.UserRoles(user)
// 	if len(errs) > 0 {
// 		return nil, errs
// 	}
// 	return userRoles, errs
// }