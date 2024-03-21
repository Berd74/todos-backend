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
	"os"
	"strings"
	"todoBackend/api2/database"
	"todoBackend/api2/response"
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
	opt := option.WithCredentialsFile("firebase-adminsdk.json")

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

		collections, err := database.SelectCollections(userIds)

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
		var body CollectionInput
		if err := c.BindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		userId := c.GetString("userId")

		collection, err := database.CreateCollection(body.Name, body.Description, userId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			response.ErrorResponse{Code: http.StatusInternalServerError, Message: "Internal error"}.Send(c)
			return
		}
		response.SendOk(c, collection)

	})

	router.DELETE("/collection/:id", verifyToken, func(c *gin.Context) {
		userId, _ := c.Get("userId")
		userIdString := userId.(string)

		collectionId := c.Param("id")

		err := database.DeleteCollection(collectionId, userIdString)

		if err != nil {
			response.SendError(c, err)
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
