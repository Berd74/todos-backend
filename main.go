package main

import (
	"github.com/gin-gonic/gin"
	"todoBackend/routes"
)

func main() {
	router := gin.Default()

	routes.Todo(router.Group("/todo"))
	routes.Collection(router.Group("/collection"))

	router.Run()
}
