package database

import (
	"cloud.google.com/go/spanner"
	"context"
)

func DeleteTodo(todoIds []string) (int64, error) {

	ctx := context.Background()

	var affectedRowsCount int64
	_, err := GetDatabase().ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {

		delTodoStmt := spanner.Statement{
			SQL: `DELETE FROM todo WHERE todo_id IN UNNEST(@todoIds)`,
			Params: map[string]any{
				"todoIds": todoIds,
			},
		}
		rowCount, err := txn.Update(ctx, delTodoStmt)
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
