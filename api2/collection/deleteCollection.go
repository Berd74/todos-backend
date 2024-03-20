package collection

import (
	"cloud.google.com/go/spanner"
	"context"
	"fmt"
	"net/http"
	"todoBackend/api2/database"
	"todoBackend/api2/errorResponse"
)

func DeleteCollection(collectionId string, userId string) error {

	stmt := spanner.Statement{
		SQL: `DELETE FROM collection WHERE collection_id = @collectionId AND user_id = @userId`,
		Params: map[string]any{
			"collectionId": collectionId,
			"userId":       userId,
		},
	}

	ctx := context.Background()

	var affectedRowsCount int64
	_, err := database.GetDatabase().ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		rowCount, err := txn.Update(ctx, stmt)
		affectedRowsCount = rowCount
		return err
	})

	if err != nil {
		return err
	}

	if affectedRowsCount == 0 {
		return errorResponse.ErrorResponse{Code: http.StatusNotFound, Message: fmt.Sprintf("No collection found ", collectionId)}
	}

	return nil
}
