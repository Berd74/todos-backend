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
	"todoBackend/database"
	"todoBackend/response"
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

	router.GET("/allCollections", verifyToken, func(c *gin.Context) {
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

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			response.SendError(c, err)
			return
		}

		response.SendOk(c, collections)
	})

	router.GET("/collection/:id", verifyToken, func(c *gin.Context) {
		userId, _ := c.Get("userId")
		role, _ := c.Get("role")
		collectionId := c.Param("id")

		collection, err := database.SelectCollection(collectionId)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			response.SendError(c, err)
			return
		}

		if userId != collection.UserId && role != Admin {
			response.ErrorResponse{Code: http.StatusForbidden, Message: "This collection doesn't belong to you"}.Send(c)
			return
		}

		response.SendOk(c, collection)
	})

	type CollectionInput struct {
		Description string `json:"description"`
		Name        string `json:"name"`
	}

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
		role, _ := c.Get("role")
		userIdString := userId.(string)
		collectionId := c.Param("id")

		collection, err := database.SelectCollection(collectionId)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			response.SendError(c, err)
			return
		}

		if userId != collection.UserId && role != Admin {
			response.ErrorResponse{Code: http.StatusForbidden, Message: "This collection doesn't belong to you"}.Send(c)
			return
		}

		if err != nil {
			response.SendError(c, err)
		}

		err = database.DeleteCollection(collectionId, userIdString)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			response.SendError(c, err)
			return
		}

		response.SendOk(c, collection)
	})

	router.DELETE("/allCollections", verifyToken, func(c *gin.Context) {
		userId, _ := c.Get("userId")
		userIdString := userId.(string)

		num, err := database.DeleteAllCollections(userIdString)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			response.SendError(c, err)
			return
		}

		response.SendOk(c, map[string]any{
			"removedItems": num,
		})
	})

	router.Run()
}
