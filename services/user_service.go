package services

import (
	"blog/models"
	"blog/repositories"
)

type UserService interface {
	CreateUser(user *models.User) error
	GetUserByID(id int) (*models.User, error)
	CreateVisitor(visitor *models.Visitor) (*models.Visitor, error)
}

type userService struct {
	userRepo repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) UserService {
	return &userService{userRepo: *userRepo}
}

func (us *userService) CreateUser(user *models.User) error {
	return us.userRepo.CreateUser(user)
}

func (us *userService) GetUserByID(id int) (*models.User, error) {
	return us.userRepo.GetUserById(id)
}

func (us *userService) CreateVisitor(visitor *models.Visitor) (*models.Visitor, error) {
	return us.userRepo.CreateVisitor(visitor)
}
