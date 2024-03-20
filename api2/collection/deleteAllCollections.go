package collection

func DeleteAllCollections() error {

	//stmt := spanner.Statement{
	//	SQL: `DELETE FROM collection WHERE collection_id = @collectionId`,
	//	Params: map[string]any{
	//		"collectionId": collectionId,
	//	},
	//}
	//
	//ctx := context.Background()
	//
	//var affectedRowsCount int64
	//_, err := database.GetDatabase().ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
	//	rowCount, err := txn.Update(ctx, stmt)
	//	affectedRowsCount = rowCount
	//	return err
	//})
	//
	//if err != nil {
	//	return err
	//}
	//
	//if affectedRowsCount == 0 {
	//	return errorResponse.ErrorResponse{Code: http.StatusNotFound, Message: fmt.Sprintf("No collection found with ID %v", collectionId)}
	//}

	return nil
}
