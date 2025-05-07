package main

import (
	"fmt"
	"net/http"

	"github.com/YahyaCengiz/todo-v2/controllers"
	"github.com/YahyaCengiz/todo-v2/middleware"
	"github.com/YahyaCengiz/todo-v2/services"
	"github.com/YahyaCengiz/todo-v2/store"
)

func main() {

	store := store.NewStore()
	todoService := services.NewTodoService(store)
	userService := services.NewUserService(store)


	todoController := controllers.NewTodoController(todoService)
	authController := controllers.NewAuthController(userService)

	http.HandleFunc("/api/login", authController.Login)


	todoListMux := http.NewServeMux()
	todoListMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			todoController.GetTodoList(w, r)
		case http.MethodPost:
			todoController.CreateTodoList(w, r)
		case http.MethodPut:
			todoController.UpdateTodoList(w, r)
		case http.MethodDelete:
			todoController.DeleteTodoList(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})


	todoItemMux := http.NewServeMux()
	todoItemMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			todoController.CreateTodoItem(w, r)
		case http.MethodPut:
			todoController.UpdateTodoItem(w, r)
		case http.MethodDelete:
			todoController.DeleteTodoItem(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.Handle("/api/todo-lists", middleware.AuthMiddleware(todoListMux))
	http.Handle("/api/todo-items", middleware.AuthMiddleware(todoItemMux))

	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}
