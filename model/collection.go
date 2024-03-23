package model

type Collection struct {
	CollectionId string `spanner:"collection_id" json:"collectionId"`
	Name         string `spanner:"name" json:"name"`
	UserId       string `spanner:"user_id" json:"userId"`
	Description  string `spanner:"description" json:"description"`
}
