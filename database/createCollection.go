package database

import (
	"cloud.google.com/go/spanner"
	"context"
	"github.com/google/uuid"
	"todoBackend/model"
)

func CreateCollection(name string, description *string, userId string) (*model.Collection, error) {
	collectionId := uuid.New().String() // Generate a new UUID.

	newCollection := &model.Collection{
		CollectionId: collectionId,
		Name:         name,
		UserId:       userId,
		Description:  description,
	}

	m, err := spanner.InsertOrUpdateStruct("collection", newCollection)
	if err != nil {
		return nil, err
	}
	_, err = GetDatabase().Apply(context.Background(), []*spanner.Mutation{m})
	if err != nil {
		return nil, err
	}

	return newCollection, nil
}
