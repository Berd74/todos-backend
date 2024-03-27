package database

import (
	"cloud.google.com/go/spanner"
	"context"
	"errors"
	"fmt"
	"net/http"
	"todoBackend/response"
	"todoBackend/utils"
)

func UpdateCollection(collectionId string, changes map[string]interface{}) error {

	var lookFor = []string{"name", "description"}

	if len(changes) == 0 {
		fmt.Println("no fields provided for update")
		return response.ErrorResponse{http.StatusBadRequest, "no fields provided for update"}
	}

	if unexpectedKeys := utils.GetUnexpectedKeys(&changes, lookFor); unexpectedKeys != nil {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("unexpected keys in body:%v", unexpectedKeys)}
	}

	columns := []string{"collection_id"}
	values := []any{collectionId}

	if err := AddToUpdate(lookFor[0], "string", changes, &columns, &values); err != nil {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()}
	}
	if err := AddToUpdate(lookFor[1], "string", changes, &columns, &values); err != nil {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()}
	}

	_, err := GetDatabase().ReadWriteTransaction(context.Background(), func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
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

func AddToUpdate(keyName string, targetType string, changes map[string]interface{}, DBColumns *[]string, DBValues *[]any) error {
	var castedVal any
	var ok bool
	value, exists := changes[keyName]
	// ignore if does not exist in changes
	if !exists {
		return nil
	}
	switch targetType {
	case "string":
		castedVal, ok = value.(string)
	case "int":
		castedVal, ok = value.(int)
	case "bool":
		castedVal, ok = value.(bool)
	default:
		fmt.Println("Unsupported target types")
	}
	if !ok {
		// if nil we want to add nil to DB
		if value == nil {
			castedVal = nil
		} else {
			return errors.New(fmt.Sprintf("parameter: \"%v\" must be %v", keyName, targetType))
		}
	}
	*DBColumns = append(*DBColumns, keyName)
	*DBValues = append(*DBValues, castedVal)
	return nil
}
