package collection

import (
	"cloud.google.com/go/spanner"
	"context"
	"database/sql"
	"github.com/google/uuid"
	"todoBackend/api2/database"
)

func CreateCollection(name string, description string, userId string) (*Collection, error) {
	collectionId := uuid.New().String() // Generate a new UUID.

	stmt := spanner.Statement{
		SQL: `INSERT INTO collection (collection_id, name, description, user_id) 
              VALUES (@collectionId, @name, @description, @userId)`,
		Params: map[string]any{
			"collectionId": collectionId,
			"name":         name,
			"description":  sql.NullString{String: description, Valid: description != ""},
			"userId":       userId,
		},
	}

	returnStmt := spanner.Statement{
		SQL: `SELECT collection_id, name, description, user_id FROM collection WHERE collection_id = @collectionId`,
		Params: map[string]any{
			"collectionId": collectionId,
		},
	}

	var newCollection *Collection

	_, err := database.GetDatabase().ReadWriteTransaction(context.Background(), func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		_, err := txn.Update(ctx, stmt)
		if err != nil {
			return err // Return if the INSERT fails
		}

		iter := txn.Query(ctx, returnStmt)
		defer iter.Stop()
		row, err := iter.Next()
		if err != nil {
			return err // Return if the SELECT fails
		}
		newCollection = &Collection{}
		if err := row.ToStruct(newCollection); err != nil {
			return err // Return if convert the row into the Collection struct fails
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return newCollection, nil
}
