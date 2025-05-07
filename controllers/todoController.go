package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/YahyaCengiz/todo-v2/middleware"
	"github.com/YahyaCengiz/todo-v2/services"
)

type TodoController struct {
	todoService *services.TodoService
}

func NewTodoController(todoService *services.TodoService) *TodoController {
	return &TodoController{todoService: todoService}
}

func (c *TodoController) CreateTodoList(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*middleware.Claims)
	var request struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	todoList, err := c.todoService.CreateTodoList(request.Name, claims.UserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todoList)
}

func (c *TodoController) GetTodoList(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*middleware.Claims)
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		todoLists, err := c.todoService.GetAllTodoLists(claims.UserID, claims.Role)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(todoLists)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	todoList, err := c.todoService.GetTodoList(id, claims.UserID, claims.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todoList)
}

func (c *TodoController) UpdateTodoList(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*middleware.Claims)
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var request struct {
		Name string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	todoList, err := c.todoService.UpdateTodoList(id, request.Name, claims.UserID, claims.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todoList)
}

func (c *TodoController) DeleteTodoList(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*middleware.Claims)
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := c.todoService.DeleteTodoList(id, claims.UserID, claims.Role); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *TodoController) CreateTodoItem(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*middleware.Claims)
	listIDStr := r.URL.Query().Get("list_id")
	if listIDStr == "" {
		http.Error(w, "List ID is required", http.StatusBadRequest)
		return
	}

	listID, err := strconv.Atoi(listIDStr)
	if err != nil {
		http.Error(w, "Invalid List ID", http.StatusBadRequest)
		return
	}

	var request struct {
		Content string `json:"content"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	todoItem, err := c.todoService.CreateTodoItem(listID, request.Content, claims.UserID, claims.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todoItem)
}

func (c *TodoController) UpdateTodoItem(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*middleware.Claims)
	listIDStr := r.URL.Query().Get("list_id")
	itemIDStr := r.URL.Query().Get("item_id")
	if listIDStr == "" || itemIDStr == "" {
		http.Error(w, "List ID and Item ID are required", http.StatusBadRequest)
		return
	}

	listID, err := strconv.Atoi(listIDStr)
	if err != nil {
		http.Error(w, "Invalid List ID", http.StatusBadRequest)
		return
	}

	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		http.Error(w, "Invalid Item ID", http.StatusBadRequest)
		return
	}

	var request struct {
		Content     string `json:"content"`
		IsCompleted bool   `json:"is_completed"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	todoItem, err := c.todoService.UpdateTodoItem(listID, itemID, request.Content, request.IsCompleted, claims.UserID, claims.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todoItem)
}

func (c *TodoController) DeleteTodoItem(w http.ResponseWriter, r *http.Request) {
	claims := r.Context().Value("claims").(*middleware.Claims)
	listIDStr := r.URL.Query().Get("list_id")
	itemIDStr := r.URL.Query().Get("item_id")
	if listIDStr == "" || itemIDStr == "" {
		http.Error(w, "List ID and Item ID are required", http.StatusBadRequest)
		return
	}

	listID, err := strconv.Atoi(listIDStr)
	if err != nil {
		http.Error(w, "Invalid List ID", http.StatusBadRequest)
		return
	}

	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		http.Error(w, "Invalid Item ID", http.StatusBadRequest)
		return
	}

	if err := c.todoService.DeleteTodoItem(listID, itemID, claims.UserID, claims.Role); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
} 