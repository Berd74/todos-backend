package database

import (
	"cloud.google.com/go/spanner"
	"context"
	"github.com/google/uuid"
	"todoBackend/model"
)

func CreateCollection(name string, description *string, clientId string) (*model.Collection, error) {
	collectionId := uuid.New().String()

	num, err := GetNewestCollectionRank(clientId)

	if err != nil {
		return nil, err
	}

	newCollection := &model.Collection{
		CollectionId: collectionId,
		Name:         name,
		UserId: clientId,
		Description:  description,
		Rank:   num + 1,
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
