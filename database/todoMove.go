package database

import (
	"cloud.google.com/go/spanner"
	"context"
	"fmt"
	"google.golang.org/api/iterator"
)

func MoveTodo(clientId string, todoToMoveId string, todoTargetId string, isAfter bool) error {

	_, err := GetDatabase().ReadWriteTransaction(context.Background(), func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {

		// get rank of target
		iterRank := txn.Query(ctx, spanner.Statement{
			SQL: `SELECT t.rank 
				FROM todo t INNER JOIN collection c ON t.collection_id = c.collection_id
                WHERE t.todo_id = @todoTargetId AND c.user_id = @clientId`,
			Params: map[string]interface{}{
				"todoTargetId": todoTargetId,
				"clientId":     clientId,
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
			SQL: `SELECT t.rank, t.todo_id FROM todo t
                     INNER JOIN collection c ON t.collection_id = c.collection_id
                     WHERE t.rank ` + operator + ` @rankOfTarget AND c.user_id = @clientId AND t.todo_id != @todoToMoveId ORDER BY rank ASC`,
			Params: map[string]interface{}{
				"rankOfTarget": rankOfTarget,
				"clientId":     clientId,
				"todoToMoveId": todoToMoveId,
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
			var todo struct {
				TodoId string `spanner:"todo_id" json:"todoId"`
				Rank   int64  `spanner:"rank" json:"rank"`
			}
			if err := row.ToStruct(&todo); err != nil {
				return err
			}

			mutation := spanner.Update("todo", []string{"todo_id", "rank"}, []interface{}{todo.TodoId, movingRank})
			mutations = append(mutations, mutation)
			movingRank++
		}

		var newRank = rankOfTarget
		if isAfter {
			newRank = newRank + 1
		}

		// move selected item
		mutation := spanner.Update("todo", []string{"todo_id", "rank"}, []interface{}{todoToMoveId, newRank})
		mutations = append(mutations, mutation)

		return txn.BufferWrite(mutations)
	})

	return err
}
