package database

import (
	"cloud.google.com/go/spanner"
	"context"
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