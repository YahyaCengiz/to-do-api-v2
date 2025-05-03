package store

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/YahyaCengiz/todo-v2/models"
)

func LoadUsersFromFile(filePath string) ([]models.User, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var users []models.User
	if err := json.NewDecoder(file).Decode(&users); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return users, nil
}

func SaveUsersToFile(filePath string, users []models.User) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	if err := json.NewEncoder(file).Encode(users); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}
