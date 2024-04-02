package database

import (
	"cloud.google.com/go/spanner"
	"context"
	"fmt"
	"net/http"
	"todoBackend/response"
	"todoBackend/utils"
)

func UpdateTodo(userIdString string, todoId string, changes map[string]interface{}) error {

	var lookFor = []string{"name", "description", "collectionId", "done", "dueDate"}

	if len(changes) == 0 {
		fmt.Println("no fields provided for update")
		return response.ErrorResponse{http.StatusBadRequest, "no fields provided for update"}
	}

	if unexpectedKeys := utils.GetUnexpectedKeys(&changes, lookFor); unexpectedKeys != nil {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("unexpected keys in body:%v", unexpectedKeys)}
	}

	columns := []string{"todo_id"}
	values := []any{todoId}

	if _, err := AddToUpdate(lookFor[0], nil, "string", changes, &columns, &values); err != nil {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()}
	}
	if _, err := AddToUpdate(lookFor[1], nil, "string", changes, &columns, &values); err != nil {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()}
	}
	var a = "collection_id"
	if collectionId, err := AddToUpdate(lookFor[2], &a, "string", changes, &columns, &values); err != nil {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()}
	} else if collectionId != nil {
		collectionIdString, ok := collectionId.(string)
		if !ok {
			return response.ErrorResponse{Code: http.StatusInternalServerError, Message: "Something wrong with parsing collectionId."}
		}
		test, errSelect := AreUserCollections(userIdString, []string{collectionIdString})
		if errSelect != nil {
			return errSelect
		}
		if !test {
			return response.ErrorResponse{Code: http.StatusForbidden, Message: "Provided collection ID that doesn't belong to you or does not exist."}
		}
	}
	if _, err := AddToUpdate(lookFor[3], nil, "bool", changes, &columns, &values); err != nil {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()}
	}
	var b = "due_date"
	if _, err := AddToUpdate(lookFor[4], &b, "timestamp", changes, &columns, &values); err != nil {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()}
	}

	fmt.Println(" === columns")
	fmt.Println(columns)
	fmt.Println(values)

	_, err := GetDatabase().ReadWriteTransaction(context.Background(), func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		err := txn.BufferWrite([]*spanner.Mutation{
			spanner.Update("todo", columns, values),
		})
		return err
	})
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
