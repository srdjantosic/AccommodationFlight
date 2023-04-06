package service

import (
	"log"
	"user-service/model"
	"user-service/repository"
)

type UserService struct {
	Logger *log.Logger
	Repo   *repository.UserRepository
}

func NewUserService(l *log.Logger, r *repository.UserRepository) *UserService {
	return &UserService{l, r}
}
func (us *UserService) Insert(newUser *model.User) (*model.User, error) {

	user, err := us.Repo.GetByEmail(newUser.Email)

	if user != nil {
		us.Logger.Println("User with this email already exists!")
		return nil, err
	}

	return us.Repo.Insert(newUser) //newUser
}

func (us *UserService) GetUserByEmailAndPassword(email string, password string) (*model.User, error) {
	return us.Repo.GetUserByEmailAndPassword(email, password)
}
