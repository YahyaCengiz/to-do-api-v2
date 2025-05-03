package store

import (
	"github.com/YahyaCengiz/todo-v2/models"
)

type Store struct {
	Todos []models.TodoList
	Users []models.User
}

func NewStore() *Store {
	todoLists, err := LoadTodosFromFile("data/todos.json")
	if err != nil {
		panic(err)
	}

	users, err := LoadUsersFromFile("data/users.json")
	if err != nil {
		panic(err)
	}

	return &Store{
		Todos: todoLists,
		Users: users,
	}
} 