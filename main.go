package main

import (
	"fmt"

	"github.com/YahyaCengiz/todo-v2/controllers"
	"github.com/YahyaCengiz/todo-v2/models"
	"github.com/YahyaCengiz/todo-v2/store"

	"net/http"
)

type Store struct {
	Todos []models.TodoList
	Users []models.User
}

func NewStore() *Store {
	todoLists, err := store.LoadTodosFromFile("data/todos.json")
	if err != nil {
		panic(err)
	}

	users, err := store.LoadUsersFromFile("data/users.json")
	if err != nil {
		panic(err)
	}

	return &Store{
		Todos: todoLists,
		Users: users,
	}
}

func main() {
	store := NewStore()

	http.HandleFunc("/login", controllers.LoginHandler)

	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}
