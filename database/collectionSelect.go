package database

import (
	"cloud.google.com/go/spanner"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
	"net/http"
	"todoBackend/model"
	"todoBackend/response"
)

func SelectCollection(clientId string, userIds *[]string, collectionIds *[]string) ([]model.Collection, error) {
	var params = make(map[string]any)
	var query = ""
	var conditions = 0
	var addConditionPrefix = func() {
		if conditions == 0 {
			query += "WHERE "
		} else {
			query += "AND "
		}
		conditions++
	}

	query += `SELECT c.collection_id, c.name, c.description, c.user_id FROM collection c `

	if userIds != nil {
		addConditionPrefix()
		query += "c.user_id IN UNNEST(@userIds) "
		params["userIds"] = *userIds
	}
	if collectionIds != nil {
		addConditionPrefix()
		query += "c.collection_id IN UNNEST(@collectionIds) "
		params["collectionIds"] = *collectionIds
	}
	query += ";"
	fmt.Println(query)

	stmt := spanner.Statement{
		SQL:    query,
		Params: params,
	}

	ctx := context.Background()
	iter := GetDatabase().Single().Query(ctx, stmt)
	defer iter.Stop()

	collections := make([]model.Collection, 0)

	for {
		row, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var val model.Collection

		if err := row.ToStruct(&val); err != nil {
			return nil, err
		}
		fmt.Println("xxx")
		fmt.Println(clientId)
		fmt.Println(val.UserId)
		if clientId != val.UserId {
			return nil, response.ErrorResponse{http.StatusUnauthorized, "You do not have access to one of the selected collection."}
		}
		collections = append(collections, val)
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

func GetNewestCollectionRank(userId string) (int64, error) {
	stmt := spanner.Statement{
		SQL: `SELECT c.rank FROM collection c WHERE user_id = @userId ORDER BY c.rank DESC LIMIT 1`,
		Params: map[string]any{
			"userId": userId,
		},
	}

	ctx := context.Background()
	iter := GetDatabase().Single().Query(ctx, stmt)
	defer iter.Stop()

	row, err := iter.Next()
	if err == iterator.Done {
		return 0, nil
	}
	if err != nil {
		// Correct handling of Next() error.
		return 0, fmt.Errorf("query failed: %v", err)
	}

	var amount *int64
	if err := row.Column(0, &amount); err != nil {
		// Error handling for reading the result.
		return 0, fmt.Errorf("failed to read result: %v", err)
	}

	return *amount, nil
}
