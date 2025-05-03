package models

import "time"

type TodoList struct {
	ID                   int        `json:"id"`
	Name                 string     `json:"name"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
	DeletedAt            time.Time  `json:"deleted_at"`
	CompletionPercentage int        `json:"completion_percentage"`
	TodoItems            []TodoItem `json:"todo_items"`
}

type TodoItem struct {
	ID          int       `json:"id"`
	TodoListID  int       `json:"todo_list_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   time.Time `json:"deleted_at"`
	Content     string    `json:"content"`
	IsCompleted bool      `json:"is_completed"`
}
