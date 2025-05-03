package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/YahyaCengiz/todo-v2/controllers"
	"github.com/YahyaCengiz/todo-v2/middleware"
	"github.com/YahyaCengiz/todo-v2/store"
)

func withStore(store *store.Store, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "store", store)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func main() {
	store := store.NewStore()

	http.HandleFunc("/login", withStore(store, controllers.LoginHandler))

	http.Handle("/todos", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Protected todos endpoint"))
	})))

	fmt.Println("Server is running on port 8080...")
	http.ListenAndServe(":8080", nil)
}
