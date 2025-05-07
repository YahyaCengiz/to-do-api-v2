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
func (s *TodoService) CreateTodoList(name string, userID int) (*models.TodoList, error) {
	todoList := &models.TodoList{
		Name:                 name,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
		CompletionPercentage: 0,
		TodoItems:           []models.TodoItem{},
		UserID:              userID,
	}
	
	if err := s.store.CreateTodoList(todoList); err != nil {
		return nil, err
	}
	return s.store.GetTodoList(todoList.ID)
}

func (s *TodoService) GetTodoList(id int, userID int, role string) (*models.TodoList, error) {
	todoList, err := s.store.GetTodoList(id)
	if err != nil {
		return nil, err
	}
	if !todoList.DeletedAt.IsZero() {
		return nil, errors.New("todo list not found")
	}
	if role != "admin" && todoList.UserID != userID {
		return nil, errors.New("forbidden")
	}
	// Filter out deleted items and items not owned by user (unless admin)
	filteredItems := make([]models.TodoItem, 0)
	for _, item := range todoList.TodoItems {
		if item.DeletedAt.IsZero() && (role == "admin" || item.UserID == userID) {
			filteredItems = append(filteredItems, item)
		}
	}
	todoList.TodoItems = filteredItems
	return todoList, nil
}

func (s *TodoService) GetAllTodoLists(userID int, role string) ([]*models.TodoList, error) {
	lists, err := s.store.GetAllTodoLists()
	if err != nil {
		return nil, err
	}
	filteredLists := make([]*models.TodoList, 0)
	for _, list := range lists {
		if list.DeletedAt.IsZero() && (role == "admin" || list.UserID == userID) {
			// Filter out deleted items and items not owned by user (unless admin)
			filteredItems := make([]models.TodoItem, 0)
			for _, item := range list.TodoItems {
				if item.DeletedAt.IsZero() && (role == "admin" || item.UserID == userID) {
					filteredItems = append(filteredItems, item)
				}
			}
			list.TodoItems = filteredItems
			filteredLists = append(filteredLists, list)
		}
	}
	return filteredLists, nil
}

func (s *TodoService) UpdateTodoList(id int, name string, userID int, role string) (*models.TodoList, error) {
	todoList, err := s.store.GetTodoList(id)
	if err != nil {
		return nil, err
	}
	if role != "admin" && todoList.UserID != userID {
		return nil, errors.New("forbidden")
	}
	todoList.Name = name
	todoList.UpdatedAt = time.Now()
	
	if err := s.store.UpdateTodoList(todoList); err != nil {
		return nil, err
	}
	return todoList, nil
}

func (s *TodoService) DeleteTodoList(id int, userID int, role string) error {
	todoList, err := s.store.GetTodoList(id)
	if err != nil {
		return err
	}
	if role != "admin" && todoList.UserID != userID {
		return errors.New("forbidden")
	}
	todoList.DeletedAt = time.Now()
	return s.store.UpdateTodoList(todoList)
}

// TodoItem operations
func (s *TodoService) CreateTodoItem(listID int, content string, userID int, role string) (*models.TodoItem, error) {
	todoList, err := s.store.GetTodoList(listID)
	if err != nil {
		return nil, err
	}
	if role != "admin" && todoList.UserID != userID {
		return nil, errors.New("forbidden")
	}
	todoItem := &models.TodoItem{
		TodoListID:  listID,
		Content:     content,
		IsCompleted: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		UserID:      userID,
	}

	if err := s.store.CreateTodoItem(todoItem); err != nil {
		return nil, err
	}

	// Fetch the pointer from store
	todoList, err = s.store.GetTodoList(listID)
	if err != nil {
		return nil, err
	}
	// Find the item with the highest ID and matching user_id
	var createdItem *models.TodoItem
	maxID := -1
	for i := range todoList.TodoItems {
		item := &todoList.TodoItems[i]
		if item.UserID == userID && item.ID > maxID {
			maxID = item.ID
			createdItem = item
		}
	}
	if createdItem == nil {
		return nil, errors.New("failed to add todo item")
	}
	return createdItem, nil
}

func (s *TodoService) UpdateTodoItem(listID, itemID int, content string, isCompleted bool, userID int, role string) (*models.TodoItem, error) {
	todoList, err := s.store.GetTodoList(listID)
	if err != nil {
		return nil, err
	}
	todoItem, err := s.store.GetTodoItem(listID, itemID)
	if err != nil {
		return nil, err
	}
	if role != "admin" && todoItem.UserID != userID {
		return nil, errors.New("forbidden")
	}
	todoItem.Content = content
	todoItem.IsCompleted = isCompleted
	todoItem.UpdatedAt = time.Now()
	if err := s.store.UpdateTodoItem(listID, todoItem); err != nil {
		return nil, err
	}
	s.updateCompletionPercentage(todoList)
	return todoItem, nil
}

func (s *TodoService) DeleteTodoItem(listID, itemID int, userID int, role string) error {
	todoList, err := s.store.GetTodoList(listID)
	if err != nil {
		return err
	}
	todoItem, err := s.store.GetTodoItem(listID, itemID)
	if err != nil {
		return err
	}
	if role != "admin" && todoItem.UserID != userID {
		return errors.New("forbidden")
	}
	todoItem.DeletedAt = time.Now()
	if err := s.store.UpdateTodoItem(listID, todoItem); err != nil {
		return err
	}
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