package users

import "github.com/bryopsida/gofiber-pug-starter/interfaces"

type usersService struct {
	repo interfaces.IUserRepository
}

func NewUsersService(repo interfaces.IUserRepository) *usersService {
	return &usersService{repo: repo}
}

func (s *usersService) CreateUser(user *interfaces.User) error {
	return s.repo.CreateUser(user)
}

func (s *usersService) GetUserByID(id uint) (*interfaces.User, error) {
	return s.repo.GetUserByID(id)
}

func (s *usersService) GetUserByUsername(username string) (*interfaces.User, error) {
	return s.repo.GetUserByUsername(username)
}

func (s *usersService) UpdateUser(user *interfaces.User) error {
	return s.repo.UpdateUser(user)
}

func (s *usersService) DeleteUser(id uint) error {
	return s.repo.DeleteUser(id)
}
