package store

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/YahyaCengiz/todo-v2/models"
)

type Store struct {
	mu        sync.RWMutex
	todoLists []models.TodoList
	users     []models.User
	filePath  string
}

func NewStore() *Store {
	s := &Store{
		todoLists: make([]models.TodoList, 0),
		users:     make([]models.User, 0),
		filePath:  "data/store.json",
	}
	if err := s.loadFromFile(); err != nil {
		panic(fmt.Sprintf("Failed to load store.json: %v", err))
	}
	return s
}

// TodoList operations
func (s *Store) CreateTodoList(todoList *models.TodoList) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Generate new ID
	if len(s.todoLists) == 0 {
		todoList.ID = 1
	} else {
		todoList.ID = s.todoLists[len(s.todoLists)-1].ID + 1
	}

	s.todoLists = append(s.todoLists, *todoList)
	return s.saveToFile()
}

func (s *Store) GetTodoList(id int) (*models.TodoList, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, list := range s.todoLists {
		if list.ID == id {
			return &list, nil
		}
	}
	return nil, fmt.Errorf("todo list not found")
}

func (s *Store) GetAllTodoLists() ([]models.TodoList, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.todoLists, nil
}

func (s *Store) UpdateTodoList(todoList *models.TodoList) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, list := range s.todoLists {
		if list.ID == todoList.ID {
			s.todoLists[i] = *todoList
			return s.saveToFile()
		}
	}
	return fmt.Errorf("todo list not found")
}

func (s *Store) DeleteTodoList(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, list := range s.todoLists {
		if list.ID == id {
			s.todoLists = append(s.todoLists[:i], s.todoLists[i+1:]...)
			return s.saveToFile()
		}
	}
	return fmt.Errorf("todo list not found")
}

// TodoItem operations
func (s *Store) CreateTodoItem(todoItem *models.TodoItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Find the todo list
	for i, list := range s.todoLists {
		if list.ID == todoItem.TodoListID {
			// Generate new ID
			if len(list.TodoItems) == 0 {
				todoItem.ID = 1
			} else {
				todoItem.ID = list.TodoItems[len(list.TodoItems)-1].ID + 1
			}

			s.todoLists[i].TodoItems = append(s.todoLists[i].TodoItems, *todoItem)
			return s.saveToFile()
		}
	}
	return fmt.Errorf("todo list not found")
}

func (s *Store) GetTodoItem(id int) (*models.TodoItem, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, list := range s.todoLists {
		for _, item := range list.TodoItems {
			if item.ID == id {
				return &item, nil
			}
		}
	}
	return nil, fmt.Errorf("todo item not found")
}

func (s *Store) UpdateTodoItem(todoItem *models.TodoItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, list := range s.todoLists {
		if list.ID == todoItem.TodoListID {
			for j, item := range list.TodoItems {
				if item.ID == todoItem.ID {
					s.todoLists[i].TodoItems[j] = *todoItem
					return s.saveToFile()
				}
			}
		}
	}
	return fmt.Errorf("todo item not found")
}

func (s *Store) DeleteTodoItem(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, list := range s.todoLists {
		for j, item := range list.TodoItems {
			if item.ID == id {
				s.todoLists[i].TodoItems = append(list.TodoItems[:j], list.TodoItems[j+1:]...)
				return s.saveToFile()
			}
		}
	}
	return fmt.Errorf("todo item not found")
}

// User operations
func (s *Store) GetUsers() []models.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.users
}

func (s *Store) AddUser(user models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users = append(s.users, user)
	return s.saveToFile()
}

// File operations
func (s *Store) saveToFile() error {
	data := struct {
		TodoLists []models.TodoList `json:"todo_lists"`
		Users     []models.User     `json:"users"`
	}{
		TodoLists: s.todoLists,
		Users:     s.users,
	}

	file, err := os.Create(s.filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

func (s *Store) loadFromFile() error {
	file, err := os.Open(s.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // File doesn't exist yet, that's okay
		}
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var data struct {
		TodoLists []models.TodoList `json:"todo_lists"`
		Users     []models.User     `json:"users"`
	}

	if err := json.NewDecoder(file).Decode(&data); err != nil {
		return fmt.Errorf("failed to decode JSON: %w", err)
	}

	s.todoLists = data.TodoLists
	s.users = data.Users
	return nil
} 