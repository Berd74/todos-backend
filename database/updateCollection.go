package database

import (
	"cloud.google.com/go/spanner"
	"context"
	"fmt"
	"net/http"
	"todoBackend/response"
)

func UpdateCollection(collectionId string, name *string, description *string) error {

	columns := []string{"collection_id"}
	values := []any{collectionId}
	if name != nil {
		columns = append(columns, "name")
		values = append(values, name)
	}
	if description != nil {
		columns = append(columns, "description")
		values = append(values, description)
	}

	if len(columns) <= 1 {
		fmt.Println("no fields provided for update")
		return response.ErrorResponse{http.StatusBadRequest, "no fields provided for update"}
	}

	_, err := GetDatabase().ReadWriteTransaction(context.Background(), func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		fmt.Println(columns)
		fmt.Println(values)
		err := txn.BufferWrite([]*spanner.Mutation{
			spanner.Update("collection", columns, values),
		})
		return err
	})
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
