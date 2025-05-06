package services

import (
	"fmt"

	"github.com/YahyaCengiz/todo-v2/models"
	"github.com/YahyaCengiz/todo-v2/store"
)

type UserService struct {
	store *store.Store
}

func NewUserService(store *store.Store) *UserService {
	return &UserService{
		store: store,
	}
}

func (s *UserService) Login(username, password string) (*models.User, error) {
	if username == "" || password == "" {
		return nil, fmt.Errorf("username and password are required")
	}

	for _, user := range s.store.GetUsers() {
		if user.Username == username && user.Password == password {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("invalid credentials")
}

func (s *UserService) Register(user models.User) error {
	return s.store.AddUser(user)
}

func (s *UserService) GetUserByID(id string) (*models.User, error) {
	// TODO: Implement get user by ID
	return nil, nil
} 