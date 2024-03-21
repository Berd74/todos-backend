package database

import (
	"cloud.google.com/go/spanner"
	"context"
	"google.golang.org/api/iterator"
	"net/http"
	"todoBackend/model"
	"todoBackend/response"
)

func SelectCollection(collectionId string) (*model.Collection, error) {
	stmt := spanner.Statement{
		SQL: `SELECT c.collection_id, c.name, c.description, c.user_id 
              FROM collection c 
              WHERE c.collection_id = @collectionId`,
		Params: map[string]interface{}{
			"collectionId": collectionId,
		},
	}

	ctx := context.Background()
	iter := GetDatabase().Single().Query(ctx, stmt)
	defer iter.Stop()

	var collection *model.Collection

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
		collection = &col
	}

	if collection == nil {
		return nil, response.ErrorResponse{Code: http.StatusNotFound, Message: "item with this id has not been found"}

	}

	return collection, nil
}
