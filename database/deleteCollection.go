package database

import (
	"cloud.google.com/go/spanner"
	"context"
)

func DeleteCollection(collectionIds []string, userId string) (int64, error) {

	ctx := context.Background()

	var affectedRowsCount int64
	_, err := GetDatabase().ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {

		delTodoStmt := spanner.Statement{
			SQL: `DELETE FROM todo WHERE collection_id IN (SELECT collection_id FROM collection WHERE collection_id IN UNNEST(@collectionIds) AND user_id = @userId)`,
			Params: map[string]any{
				"collectionIds": collectionIds,
				"userId":        userId,
			},
		}
		_, err := txn.Update(ctx, delTodoStmt)
		if err != nil {
			return err
		}

		delCollectionStmt := spanner.Statement{
			SQL: `DELETE FROM collection WHERE collection_id IN UNNEST(@collectionIds) AND user_id = @userId`,
			Params: map[string]any{
				"collectionIds": collectionIds,
				"userId":        userId,
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
		return 0, err
	}

	return affectedRowsCount, nil
}
