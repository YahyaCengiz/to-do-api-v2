package store

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"

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

	for i := range s.todoLists {
		if s.todoLists[i].ID == id {
			return &s.todoLists[i], nil
		}
	}
	return nil, fmt.Errorf("todo list not found")
}

func (s *Store) GetAllTodoLists() ([]*models.TodoList, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	lists := make([]*models.TodoList, len(s.todoLists))
	for i := range s.todoLists {
		lists[i] = &s.todoLists[i]
	}
	return lists, nil
}

func (s *Store) UpdateTodoList(todoList *models.TodoList) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.todoLists {
		if s.todoLists[i].ID == todoList.ID {
			s.todoLists[i] = *todoList
			return s.saveToFile()
		}
	}
	return fmt.Errorf("todo list not found")
}

func (s *Store) DeleteTodoList(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.todoLists {
		if s.todoLists[i].ID == id {
			s.todoLists[i].DeletedAt = time.Now()
			return s.saveToFile()
		}
	}
	return fmt.Errorf("todo list not found")
}

// TodoItem operations
func (s *Store) CreateTodoItem(todoItem *models.TodoItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.todoLists {
		if s.todoLists[i].ID == todoItem.TodoListID {
			// Generate new ID
			if len(s.todoLists[i].TodoItems) == 0 {
				todoItem.ID = 1
			} else {
				todoItem.ID = s.todoLists[i].TodoItems[len(s.todoLists[i].TodoItems)-1].ID + 1
			}
			// Do NOT overwrite user_id here, just append
			s.todoLists[i].TodoItems = append(s.todoLists[i].TodoItems, *todoItem)
			return s.saveToFile()
		}
	}
	return fmt.Errorf("todo list not found")
}

func (s *Store) GetTodoItem(listID, itemID int) (*models.TodoItem, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for i := range s.todoLists {
		if s.todoLists[i].ID == listID {
			for j := range s.todoLists[i].TodoItems {
				if s.todoLists[i].TodoItems[j].ID == itemID {
					return &s.todoLists[i].TodoItems[j], nil
				}
			}
		}
	}
	return nil, fmt.Errorf("todo item not found")
}

func (s *Store) UpdateTodoItem(listID int, todoItem *models.TodoItem) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.todoLists {
		if s.todoLists[i].ID == listID {
			for j := range s.todoLists[i].TodoItems {
				if s.todoLists[i].TodoItems[j].ID == todoItem.ID {
					s.todoLists[i].TodoItems[j] = *todoItem
					return s.saveToFile()
				}
			}
		}
	}
	return fmt.Errorf("todo item not found")
}

func (s *Store) DeleteTodoItem(listID, itemID int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.todoLists {
		if s.todoLists[i].ID == listID {
			for j := range s.todoLists[i].TodoItems {
				if s.todoLists[i].TodoItems[j].ID == itemID {
					s.todoLists[i].TodoItems[j].DeletedAt = time.Now()
					return s.saveToFile()
				}
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

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(data); err != nil {
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