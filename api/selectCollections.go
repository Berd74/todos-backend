package main

import (
	"cloud.google.com/go/spanner"
	"context"
	"google.golang.org/api/iterator"
	"log"
)

type Collection struct {
	CollectionId string  `spanner:"collection_id"`
	Name         string  `spanner:"name"`
	Description  *string `spanner:"description"`
	UserId       string  `spanner:"user_id"`
}

func SelectCollections(userIds []string, limit int, offset int) ([]Collection, error) {
	stmt := spanner.Statement{
		SQL: `SELECT c.collection_id, c.name, c.description, c.user_id 
              FROM database c 
              WHERE c.user_id IN UNNEST(@user_ids)`,
		Params: map[string]interface{}{
			"user_ids": userIds,
		},
	}

	ctx := context.Background()
	iter := GetDatabase().Single().Query(ctx, stmt)
	defer iter.Stop()

	var collections []Collection

	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var col Collection

		if err := row.ToStruct(&col); err != nil {
			log.Fatalf("Failed to parse row: %v", err)
		}
		collections = append(collections, col)
	}
	return collections, nil
}
