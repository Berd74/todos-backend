package collection

import (
	"cloud.google.com/go/spanner"
	"context"
	"google.golang.org/api/iterator"
	"log"
	"todoBackend/api2/database"
)

type Collection struct {
	CollectionId string  `spanner:"collection_id" json:"collectionId"`
	Name         string  `spanner:"name" json:"name"`
	Description  *string `spanner:"description" json:"description"`
	UserId       string  `spanner:"user_id" json:"userId"`
}

func SelectCollections(userIds []string) ([]Collection, error) {
	stmt := spanner.Statement{
		SQL: `SELECT c.collection_id, c.name, c.description, c.user_id 
              FROM collection c 
              WHERE c.user_id IN UNNEST(@user_ids)`,
		Params: map[string]interface{}{
			"user_ids": userIds,
		},
	}

	ctx := context.Background()
	iter := database.GetDatabase().Single().Query(ctx, stmt)
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
