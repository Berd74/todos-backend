package database

import (
	"cloud.google.com/go/spanner"
	"context"
	"github.com/google/uuid"
	"time"
	"todoBackend/model"
)

func CreateTodo(args model.CreateTodoArgs) (*model.Todo, error) {
	todoId := uuid.New().String() // Generate a new UUID.

	var isDone bool
	if args.Done == nil || *args.Done == false {
		isDone = false
	} else {
		isDone = true
	}

	newTodo := &model.Todo{
		TodoId:       todoId,
		Name:         args.Name,
		CollectionId: args.CollectionId,
		Description:  args.Description,
		Done:         isDone,
		DueDate:      args.DueDate,
		CreatedAt:    time.Now(),
	}

	m, err := spanner.InsertOrUpdateStruct("todo", newTodo)
	if err != nil {
		return nil, err
	}
	_, err = GetDatabase().Apply(context.Background(), []*spanner.Mutation{m})
	if err != nil {
		return nil, err
	}

	return newTodo, nil
}
