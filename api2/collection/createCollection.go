package collection

import (
	"cloud.google.com/go/spanner"
	"context"
	"github.com/google/uuid"
	"todoBackend/api2/database"
)

func CreateCollection(name string, description string, userId string) (*Collection, error) {
	collectionId := uuid.New().String() // Generate a new UUID.

	newCollection := &Collection{
		CollectionId: collectionId,
		Name:         name,
		Description:  &description,
		UserId:       userId,
	}

	m, err := spanner.InsertOrUpdateStruct("collection", newCollection)
	if err != nil {
		return nil, err
	}
	_, err = database.GetDatabase().Apply(context.Background(), []*spanner.Mutation{m})
	if err != nil {
		return nil, err
	}

	return newCollection, nil
}
