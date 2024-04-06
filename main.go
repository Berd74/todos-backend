package main

import (
	"github.com/gin-gonic/gin"
	"todoBackend/database"
	"todoBackend/firebase"
	"todoBackend/routes"
)

func main() {
	firebase.InitFirebase()
	database.InitDatabase()
	router := gin.Default()

	routes.Todo(router.Group("/todo"))
	routes.Collection(router.Group("/collection"))

	router.Run()
}
