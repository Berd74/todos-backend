package database

import (
	"cloud.google.com/go/spanner"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"todoBackend/model"
)

func SelectCollections(userIds []string) ([]model.Collection, error) {
	stmt := spanner.Statement{
		SQL: `SELECT c.collection_id, c.name, c.description, c.user_id 
              FROM collection c 
              WHERE c.user_id IN UNNEST(@user_ids)`,
		Params: map[string]interface{}{
			"user_ids": userIds,
		},
	}

	ctx := context.Background()
	iter := GetDatabase().Single().Query(ctx, stmt)
	defer iter.Stop()

	var collections []model.Collection

	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var col model.Collection

		if err := row.ToStruct(&col); err != nil {
			return nil, err
		}
		collections = append(collections, col)
	}

	return collections, nil
}

func AreUserCollections(userId string, collectionIds []string) (bool, error) {
	fmt.Println(collectionIds)
	stmt := spanner.Statement{
		SQL: `SELECT COUNT(1) FROM collection WHERE user_id = @userId AND collection_id IN UNNEST(@collectionIds)`,
		Params: map[string]any{
			"userId":        userId,
			"collectionIds": collectionIds,
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

	return amount == len(collectionIds), nil
}
