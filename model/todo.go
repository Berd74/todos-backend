package model

import "time"

type Todo struct {
	TodoId       string     `spanner:"todo_id" json:"todoId"`
	Name         string     `spanner:"name" json:"name"`
	CollectionId string     `spanner:"collection_id" json:"collectionId"`
	Description  *string    `spanner:"description" json:"description"`
	Done         bool       `spanner:"done" json:"done"`
	DueDate      *time.Time `spanner:"due_date" json:"dueDate"`
	CreatedAt    time.Time  `spanner:"created_at" json:"createdAt"`
}

type SelectTodoDbResponse struct {
	TodoId       string     `spanner:"todo_id" json:"todoId"`
	UserId       string     `spanner:"user_id" json:"userId"`
	Name         string     `spanner:"name" json:"name"`
	CollectionId string     `spanner:"collection_id" json:"collectionId"`
	Description  *string    `spanner:"description" json:"description"`
	Done         bool       `spanner:"done" json:"done"`
	DueDate      *time.Time `spanner:"due_date" json:"dueDate"`
	CreatedAt    time.Time  `spanner:"created_at" json:"createdAt"`
}

type CreateTodoArgs struct {
	Name         string     `json:"name"`
	CollectionId string     `json:"collectionId"`
	Description  *string    `json:"description,omitempty"`
	Done         *bool      `json:"done,omitempty"`
	DueDate      *time.Time `json:"dueDate,omitempty"`
}
