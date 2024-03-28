package database

import (
	"cloud.google.com/go/spanner"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"net/http"
	"todoBackend/model"
	"todoBackend/response"
)

func SelectTodos(clientId string, todoIds *[]string, userIds *[]string, collectionIds *[]string, done *bool) ([]model.Todo, error) {
	var params = make(map[string]any)
	var query = ""
	var conditions = 0
	var addConditionPrefix = func() {
		if conditions == 0 {
			query += "WHERE "
		} else {
			query += "AND "
		}
		conditions++
	}

	query += `SELECT t.todo_id, t.name, t.description, t.done, t.collection_id, t.due_date, t.created_at, c.user_id FROM todo t 
    INNER JOIN collection c ON t.collection_id = c.collection_id `

	if todoIds != nil {
		addConditionPrefix()
		query += "t.todo_id IN UNNEST(@todoIds) "
		params["todoIds"] = *todoIds
	}
	if userIds != nil {
		addConditionPrefix()
		query += "c.user_id IN UNNEST(@userIds) "
		params["userIds"] = *userIds
	}
	if collectionIds != nil {
		addConditionPrefix()
		query += "t.collection_id IN UNNEST(@collectionIds) "
		params["collectionIds"] = *collectionIds
	}
	if done != nil {
		addConditionPrefix()
		query += "t.done = @done "
		params["done"] = *done
	}
	query += ";"

	stmt := spanner.Statement{
		SQL:    query,
		Params: params,
	}

	ctx := context.Background()
	iter := GetDatabase().Single().Query(ctx, stmt)
	defer iter.Stop()

	todos := make([]model.Todo, 0)

	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var val model.SelectTodoDbResponse

		if err := row.ToStruct(&val); err != nil {
			return nil, err
		}
		if clientId != val.UserId {
			return nil, response.ErrorResponse{http.StatusUnauthorized, "You do not have access to one of the selected todo."}
		}
		todos = append(todos, model.Todo{
			TodoId:       val.TodoId,
			Name:         val.Name,
			CollectionId: val.CollectionId,
			Description:  val.Description,
			Done:         val.Done,
			DueDate:      val.DueDate,
			CreatedAt:    val.CreatedAt,
		})
	}

	return todos, nil
}

func AreUserTodos(userId string, todoIds []string) (bool, error) {

	stmt := spanner.Statement{
		SQL: `SELECT COUNT(1) 
				FROM todo t INNER JOIN collection c ON t.collection_id = c.collection_id
                WHERE user_id = @userId AND todo_id IN UNNEST(@todoIds)`,
		Params: map[string]any{
			"userId":  userId,
			"todoIds": todoIds,
		},
	}

	ctx := context.Background()
	iter := GetDatabase().Single().Query(ctx, stmt)
	defer iter.Stop()

	row, err := iter.Next()
	if err != nil {
		// Correct handling of Next() error.
		return false, fmt.Errorf("query failed: %v", err)
	}

	var _amount int64
	var amount int
	if err := row.Column(0, &_amount); err != nil {
		// Error handling for reading the result.
		return false, fmt.Errorf("failed to read result: %v", err)
	}

	amount = int(_amount)

	return amount == len(todoIds), nil
}
