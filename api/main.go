package main

import (
	"fmt"
	"log"
)

func main() {

	userIds := []string{"fcfec921-d938-49dc-9d07-412fa5c2d256"} // Example user IDs
	limit := 10
	offset := 0

	fmt.Println("server started")

	collections, err := SelectCollections(userIds, limit, offset)

	fmt.Println(collections)
	if err != nil {
		log.Fatalf("Failed to select collections: %v", err)
	}

	for _, collection := range collections {
		fmt.Printf("%+v\n", collection)
	}

}
