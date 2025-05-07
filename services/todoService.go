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
	return s.store.GetTodoList(todoList.ID)
}

func (s *TodoService) GetTodoList(id int) (*models.TodoList, error) {
	todoList, err := s.store.GetTodoList(id)
	if err != nil {
		return nil, err
	}
	if !todoList.DeletedAt.IsZero() {
		return nil, errors.New("todo list not found")
	}
	// Filter out deleted items
	filteredItems := make([]models.TodoItem, 0)
	for _, item := range todoList.TodoItems {
		if item.DeletedAt.IsZero() {
			filteredItems = append(filteredItems, item)
		}
	}
	todoList.TodoItems = filteredItems
	return todoList, nil
}

func (s *TodoService) GetAllTodoLists() ([]*models.TodoList, error) {
	lists, err := s.store.GetAllTodoLists()
	if err != nil {
		return nil, err
	}
	filteredLists := make([]*models.TodoList, 0)
	for _, list := range lists {
		if list.DeletedAt.IsZero() {
			// Filter out deleted items
			filteredItems := make([]models.TodoItem, 0)
			for _, item := range list.TodoItems {
				if item.DeletedAt.IsZero() {
					filteredItems = append(filteredItems, item)
				}
			}
			list.TodoItems = filteredItems
			filteredLists = append(filteredLists, list)
		}
	}
	return filteredLists, nil
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
	todoList, err := s.store.GetTodoList(id)
	if err != nil {
		return err
	}
	todoList.DeletedAt = time.Now()
	return s.store.UpdateTodoList(todoList)
}

// TodoItem operations
func (s *TodoService) CreateTodoItem(listID int, content string) (*models.TodoItem, error) {
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

	// Fetch the pointer from store
	todoList, err := s.store.GetTodoList(listID)
	if err != nil {
		return nil, err
	}
	if len(todoList.TodoItems) == 0 {
		return nil, errors.New("failed to add todo item")
	}
	return &todoList.TodoItems[len(todoList.TodoItems)-1], nil
}

func (s *TodoService) UpdateTodoItem(listID, itemID int, content string, isCompleted bool) (*models.TodoItem, error) {
	todoList, err := s.store.GetTodoList(listID)
	if err != nil {
		return nil, err
	}

	todoItem, err := s.store.GetTodoItem(listID, itemID)
	if err != nil {
		return nil, err
	}

	if todoItem.TodoListID != listID {
		return nil, errors.New("todo item does not belong to the specified list")
	}

	todoItem.Content = content
	todoItem.IsCompleted = isCompleted
	todoItem.UpdatedAt = time.Now()

	if err := s.store.UpdateTodoItem(listID, todoItem); err != nil {
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
	todoItem, err := s.store.GetTodoItem(listID, itemID)
	if err != nil {
		return err
	}
	todoItem.DeletedAt = time.Now()
	if err := s.store.UpdateTodoItem(listID, todoItem); err != nil {
		return err
	}
	// Update completion percentage
	s.updateCompletionPercentage(todoList)
	return nil
}

func (s *TodoService) updateCompletionPercentage(todoList *models.TodoList) {
	total := 0
	completed := 0
	for _, item := range todoList.TodoItems {
		if item.DeletedAt.IsZero() {
			total++
			if item.IsCompleted {
				completed++
			}
		}
	}
	if total == 0 {
		todoList.CompletionPercentage = 0
	} else {
		todoList.CompletionPercentage = (completed * 100) / total
	}
	todoList.UpdatedAt = time.Now()
	s.store.UpdateTodoList(todoList)
} 