package database

import (
	"cloud.google.com/go/spanner"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
)

func MoveCollection(clientId string, collectionToMoveId string, collectionTargetId string, isAfter bool) error {

	_, err := GetDatabase().ReadWriteTransaction(context.Background(), func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {

		// get rank of target
		iterRank := txn.Query(ctx, spanner.Statement{
			SQL: `SELECT rank FROM collection WHERE collection_id = @collectionTargetId AND user_id = @clientId`,
			Params: map[string]interface{}{
				"collectionTargetId": collectionTargetId,
				"clientId":           clientId,
			},
		})

		defer iterRank.Stop()
		row, err := iterRank.Next()
		if err == iterator.Done {
			return fmt.Errorf("the collectionTargetId cannot be found")
		}
		if err != nil {
			return err // Handle error
		}

		var rankOfTarget int64
		if err := row.Column(0, &rankOfTarget); err != nil {
			return err // Handle error
		}

		// get items depending on rankOfTarget without collectionToMoveId

		var operator string
		if isAfter {
			operator = ">"
		} else {
			operator = ">="
		}
		iterItems := txn.Query(ctx, spanner.Statement{
			SQL: `SELECT rank, collection_id FROM collection WHERE rank ` + operator + ` @rankOfTarget AND user_id = @clientId AND collection_id != @collectionToMoveId ORDER BY rank ASC`,
			Params: map[string]interface{}{
				"rankOfTarget":       rankOfTarget,
				"clientId":           clientId,
				"collectionToMoveId": collectionToMoveId,
			},
		})
		defer iterItems.Stop()

		// move all above items
		var movingRank = rankOfTarget
		if isAfter {
			movingRank = movingRank + 2
		} else {
			movingRank = movingRank + 1
		}
		var mutations []*spanner.Mutation

		for {
			row, err = iterItems.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}
			var collection struct {
				CollectionId string `spanner:"collection_id" json:"collectionId"`
				Rank         int64  `spanner:"rank" json:"rank"`
			}
			if err := row.ToStruct(&collection); err != nil {
				return err
			}

			mutation := spanner.Update("collection", []string{"collection_id", "rank"}, []interface{}{collection.CollectionId, movingRank})
			mutations = append(mutations, mutation)
			movingRank++
		}

		var newRank = rankOfTarget
		if isAfter {
			newRank = newRank + 1
		}

		// move selected item
		mutation := spanner.Update("collection", []string{"collection_id", "rank"}, []interface{}{collectionToMoveId, newRank})
		mutations = append(mutations, mutation)

		return txn.BufferWrite(mutations)
	})

	return err
}
