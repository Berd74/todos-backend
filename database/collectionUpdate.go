package database

import (
	"cloud.google.com/go/spanner"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"
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

	if _, err := AddToUpdate(lookFor[0], nil, "string", changes, &columns, &values); err != nil {
		return response.ErrorResponse{Code: http.StatusBadRequest, Message: err.Error()}
	}
	if _, err := AddToUpdate(lookFor[1], nil, "string", changes, &columns, &values); err != nil {
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

func AddToUpdate(keyName string, DbName *string, targetType string, changes map[string]interface{}, DBColumns *[]string, DBValues *[]any) (any, error) {
	var castedVal any
	var ok bool = false
	value, exists := changes[keyName]
	// ignore if does not exist in changes
	if !exists {
		return nil, nil
	}
	switch targetType {
	case "string":
		castedVal, ok = value.(string)
	case "int":
		if floatVal, okVal := value.(float64); okVal {
			ok = true
			castedVal = int64(floatVal)
		}
	case "timestamp":
		if floatVal, okVal := value.(float64); okVal {
			ok = true
			intVal := int64(floatVal)
			castedVal = time.Unix(intVal, 0)
		}
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
			return nil, errors.New(fmt.Sprintf("parameter: \"%v\" must be %v", keyName, targetType))
		}
	}
	if DbName == nil {
		*DBColumns = append(*DBColumns, keyName)
	} else {
		*DBColumns = append(*DBColumns, *DbName)
	}
	*DBValues = append(*DBValues, castedVal)
	return castedVal, nil
}
