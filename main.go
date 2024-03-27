package main

import (
	"github.com/gin-gonic/gin"
	"todoBackend/routes"
)

func main() {
	router := gin.Default()

	routes.Todo(router.Group("/todo"))
	routes.Collectionn(router.Group("/collection"))

	// COLLECTIONS

	//todo add it to GET
	//router.GET("/ownCollections", utils.verifyToken, func(c *gin.Context) {
	//	userId, _ := c.Get("userId")
	//	userIdString := userId.(string)
	//	var userIds = []string{userIdString}
	//
	//	collections, err := database.SelectCollections(userIds)
	//
	//	if err != nil {
	//		response.SendError(c, err)
	//		return
	//	}
	//
	//	response.SendOk(c, collections)
	//})

	router.Run()
}
