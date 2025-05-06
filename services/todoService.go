package services

import (
	"errors"
	"time"

	"github.com/YahyaCengiz/todo-v2/models"
	"github.com/YahyaCengiz/todo-v2/store"
)

type TodoService struct {
	store *store.Store
}

func NewTodoService(store *store.Store) *TodoService {
	return &TodoService{store: store}
}

// TodoList operations
func (s *TodoService) CreateTodoList(name string) (*models.TodoList, error) {
	todoList := &models.TodoList{
		Name:                 name,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
		CompletionPercentage: 0,
		TodoItems:           []models.TodoItem{},
	}
	
	if err := s.store.CreateTodoList(todoList); err != nil {
		return nil, err
	}
	return todoList, nil
}

func (s *TodoService) GetTodoList(id int) (*models.TodoList, error) {
	return s.store.GetTodoList(id)
}

func (s *TodoService) GetAllTodoLists() ([]models.TodoList, error) {
	return s.store.GetAllTodoLists()
}

func (s *TodoService) UpdateTodoList(id int, name string) (*models.TodoList, error) {
	todoList, err := s.store.GetTodoList(id)
	if err != nil {
		return nil, err
	}
	
	todoList.Name = name
	todoList.UpdatedAt = time.Now()
	
	if err := s.store.UpdateTodoList(todoList); err != nil {
		return nil, err
	}
	return todoList, nil
}

func (s *TodoService) DeleteTodoList(id int) error {
	return s.store.DeleteTodoList(id)
}

// TodoItem operations
func (s *TodoService) CreateTodoItem(listID int, content string) (*models.TodoItem, error) {
	todoList, err := s.store.GetTodoList(listID)
	if err != nil {
		return nil, err
	}

	todoItem := &models.TodoItem{
		TodoListID:  listID,
		Content:     content,
		IsCompleted: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.store.CreateTodoItem(todoItem); err != nil {
		return nil, err
	}

	// Update completion percentage
	s.updateCompletionPercentage(todoList)
	return todoItem, nil
}

func (s *TodoService) UpdateTodoItem(listID, itemID int, content string, isCompleted bool) (*models.TodoItem, error) {
	todoList, err := s.store.GetTodoList(listID)
	if err != nil {
		return nil, err
	}

	todoItem, err := s.store.GetTodoItem(itemID)
	if err != nil {
		return nil, err
	}

	if todoItem.TodoListID != listID {
		return nil, errors.New("todo item does not belong to the specified list")
	}

	todoItem.Content = content
	todoItem.IsCompleted = isCompleted
	todoItem.UpdatedAt = time.Now()

	if err := s.store.UpdateTodoItem(todoItem); err != nil {
		return nil, err
	}

	// Update completion percentage
	s.updateCompletionPercentage(todoList)
	return todoItem, nil
}

func (s *TodoService) DeleteTodoItem(listID, itemID int) error {
	todoList, err := s.store.GetTodoList(listID)
	if err != nil {
		return err
	}

	if err := s.store.DeleteTodoItem(itemID); err != nil {
		return err
	}

	// Update completion percentage
	s.updateCompletionPercentage(todoList)
	return nil
}

func (s *TodoService) updateCompletionPercentage(todoList *models.TodoList) {
	if len(todoList.TodoItems) == 0 {
		todoList.CompletionPercentage = 0
		return
	}

	completed := 0
	for _, item := range todoList.TodoItems {
		if item.IsCompleted {
			completed++
		}
	}

	todoList.CompletionPercentage = (completed * 100) / len(todoList.TodoItems)
	todoList.UpdatedAt = time.Now()
	s.store.UpdateTodoList(todoList)
} 