package database

import (
	"cloud.google.com/go/spanner"
	"context"
	"fmt"
	"net/http"
	"todoBackend/response"
)

func DeleteCollection(collectionId string, userId string) error {

	ctx := context.Background()

	var affectedRowsCount int64
	_, err := GetDatabase().ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {

		delTodoStmt := spanner.Statement{
			SQL: `DELETE FROM todo WHERE collection_id IN (SELECT collection_id FROM collection WHERE collection_id = @collectionId AND user_id = @userId)`,
			Params: map[string]any{
				"collectionId": collectionId,
				"userId":       userId,
			},
		}
		_, err := txn.Update(ctx, delTodoStmt)
		if err != nil {
			return err
		}

		delCollectionStmt := spanner.Statement{
			SQL: `DELETE FROM collection WHERE collection_id = @collectionId AND user_id = @userId`,
			Params: map[string]any{
				"collectionId": collectionId,
				"userId":       userId,
			},
		}
		rowCount, err := txn.Update(ctx, delCollectionStmt)
		if err != nil {
			return err
		}
		affectedRowsCount = rowCount
		return nil
	})

	if err != nil {
		return err
	}

	if affectedRowsCount == 0 {
		return response.ErrorResponse{Code: http.StatusNotFound, Message: fmt.Sprintf("Item with this id not found %v", collectionId)}
	}

	return nil
}
