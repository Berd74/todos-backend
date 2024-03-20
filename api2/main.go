package main

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"log"
	"net/http"
	"strings"
	"todoBackend/api2/collection"
	"todoBackend/api2/errorResponse"
)

type Role int

const (
	User Role = iota
	Admin
)

// Initialize Firebase Auth
var app *firebase.App
var authClient *auth.Client

func init() {
	opt := option.WithCredentialsFile("todos-645f4-firebase-adminsdk-u8h0h-6f635bea6e.json")

	var err error
	app, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	authClient, err = app.Auth(context.Background())
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}
}

func main() {
	router := gin.Default()

	// Use the middleware
	router.GET("/test", verifyToken, func(c *gin.Context) {
		userId, _ := c.Get("userId")
		role, _ := c.Get("role")

		// Type assertion
		userRole, _ := role.(Role)
		// Print the result of the type assertion
		fmt.Printf("%v", userRole)

		c.JSON(http.StatusOK, gin.H{
			"userId": userId,
			"role":   role,
		})
	})

	// Use the middleware
	router.GET("/collections", verifyToken, func(c *gin.Context) {
		userId, _ := c.Get("userId")
		userIdString := userId.(string)
		role, _ := c.Get("role")
		userIdsString := c.Query("userIds")
		var userIds []string

		if userIdsString != "" {
			if role != Admin {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "Only Admin cannot use userIds",
				})
				return
			}
			userIds = strings.Split(userIdsString, " ")
		} else {
			userIds = []string{userIdString}
		}

		collections, err := collection.SelectCollections(userIds)

		fmt.Println(collections)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error": fmt.Sprintf("Failed to select collections: %v", err),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"collections": collections,
		})
	})

	router.POST("/collection", verifyToken, func(c *gin.Context) {
		//convert body to object
		var body CollectionInput
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//convert empty string to nil
		//if body.Description != nil && *body.Description == "" {
		//	body.Description = nil
		//}
		//other vars
		userId := c.GetString("userId")

		fmt.Println(userId)
		fmt.Println(body.Name)
		fmt.Println(body.Description)
		collection.CreateCollection(body.Name, body.Description, userId)

	})

	router.DELETE("/collection/:id", verifyToken, func(c *gin.Context) {
		userId, _ := c.Get("userId")
		userIdString := userId.(string)

		collectionId := c.Param("id")

		err := collection.DeleteCollection(collectionId, userIdString)

		if err != nil {
			errorResponse.SendErrorResponse(c, err)
			return
		}

		c.Status(http.StatusOK)
	})

	router.Run()
}

type CollectionInput struct {
	Description string `json:"description"`
	Name        string `json:"name"`
}
