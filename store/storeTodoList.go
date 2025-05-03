package store

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/YahyaCengiz/todo-v2/models"
)

func LoadTodosFromFile(filePath string) ([]models.TodoList, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var todoLists []models.TodoList
	if err := json.NewDecoder(file).Decode(&todoLists); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return todoLists, nil
}

func SaveTodosToFile(filePath string, todoLists []models.TodoList) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(todoLists); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}
