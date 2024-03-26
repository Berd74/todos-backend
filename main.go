package main

import (
	"bytes"
	"context"
	"encoding/json"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"fmt"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	"io"
	"log"
	"net/http"
	"strings"
	"todoBackend/database"
	"todoBackend/model"
	"todoBackend/response"
	"todoBackend/utils"
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

	// COLLECTIONS

	router.GET("/ownCollections", verifyToken, func(c *gin.Context) {
		userId, _ := c.Get("userId")
		userIdString := userId.(string)
		var userIds = []string{userIdString}

		collections, err := database.SelectCollections(userIds)

		if err != nil {
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
			response.SendError(c, err)
			return
		}

		if userId != collection.UserId && role != Admin {
			response.ErrorResponse{Code: http.StatusForbidden, Message: "This collection doesn't belong to you"}.Send(c)
			return
		}

		response.SendOk(c, collection)
	})

	router.POST("/collection", verifyToken, func(c *gin.Context) {
		var body struct {
			Description *string `json:"description,omitempty"`
			Name        string  `json:"name"`
		}

		if err := c.BindJSON(&body); err != nil {
			response.SendError(c, err)
			return
		}

		userId := c.GetString("userId")

		collection, err := database.CreateCollection(body.Name, body.Description, userId)
		if err != nil {
			response.ErrorResponse{Code: http.StatusInternalServerError, Message: "Internal error"}.Send(c)
			return
		}
		response.SendOk(c, collection)
	})

	router.DELETE("/collection/", verifyToken, func(c *gin.Context) {
		userId, _ := c.Get("userId")
		//role, _ := c.Get("role")
		userIdString := userId.(string)
		ids := c.Query("ids")
		collectionIds := strings.Split(ids, ",")

		test, err := database.AreUserCollections(userIdString, collectionIds)

		if err != nil {
			response.SendError(c, err)
			return
		}
		if !test {
			response.ErrorResponse{Code: http.StatusForbidden, Message: "Provided collection that doesn't belong to you"}.Send(c)
		}

		amount, err := database.DeleteCollection(collectionIds, userIdString)
		if err != nil {
			response.SendError(c, err)
			return
		}

		response.SendOk(c, map[string]any{
			"removedItems": amount,
		})
	})

	router.PUT("/collection/:id", verifyToken, func(c *gin.Context) {
		userId, _ := c.Get("userId")
		role, _ := c.Get("role")
		collectionId := c.Param("id")
		collection, errSelect := database.SelectCollection(collectionId)

		if errSelect != nil {
			response.SendError(c, errSelect)
			return
		}

		if userId != collection.UserId && role != Admin {
			response.ErrorResponse{Code: http.StatusForbidden, Message: "This collection doesn't belong to you"}.Send(c)
			return
		}

		var bodyBytes, bodyBytesErr = io.ReadAll(c.Request.Body)
		if bodyBytesErr != nil {
			response.ErrorResponse{Code: http.StatusInternalServerError, Message: bodyBytesErr.Error()}.Send(c)
			return
		}

		var bodyFull map[string]interface{}

		if err := json.Unmarshal(bodyBytes, &bodyFull); err != nil {
			response.ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()}.Send(c)
			return
		}

		updateErr := database.UpdateCollection(collectionId, bodyFull)
		if updateErr != nil {
			response.SendError(c, updateErr)
			return
		}

		collection, errSelect = database.SelectCollection(collectionId)
		if errSelect != nil {
			response.SendError(c, errSelect)
			return
		}

		response.SendOk(c, collection)
	})

	// TO-DO

	router.DELETE("/todo", verifyToken, func(c *gin.Context) {
		userId, _ := c.Get("userId")
		//role, _ := c.Get("role")
		//userIdString := userId.(string)
		//collectionId := c.Param("id")
		ids := c.Param("ids")
		ids2 := c.Query("ids")
		ids3 := strings.Split(ids2, ",")
		fmt.Println(userId)
		fmt.Println(ids)
		fmt.Println(ids2)
		fmt.Println(ids3)
		//
		//collection, err := database.SelectCollection(collectionId)
		//
		//if err != nil {
		//	response.SendError(c, err)
		//	return
		//}
		//
		//if userId != collection.UserId && role != Admin {
		//	response.ErrorResponse{Code: http.StatusForbidden, Message: "This collection doesn't belong to you"}.Send(c)
		//	return
		//}
		//
		//err = database.DeleteCollection(collectionId, userIdString)
		//
		//if err != nil {
		//	response.SendError(c, err)
		//	return
		//}

		response.SendOk(c, "ok")
	})

	router.POST("/todo", verifyToken, func(c *gin.Context) {

		// check if unexpected keys in json
		if bodyBytes, bodyBytesErr := io.ReadAll(c.Request.Body); bodyBytesErr != nil {
			response.ErrorResponse{Code: http.StatusInternalServerError, Message: bodyBytesErr.Error()}.Send(c)
			return
		} else {
			var bodyFull map[string]interface{}
			if err := json.Unmarshal(bodyBytes, &bodyFull); err != nil {
				response.ErrorResponse{Code: http.StatusInternalServerError, Message: err.Error()}.Send(c)
				return
			}
			expectedKeysInJSON := []string{"name", "collectionId", "description", "done", "dueDate"}
			if unexpectedKeys := utils.GetUnexpectedKeys(&bodyFull, expectedKeysInJSON); unexpectedKeys != nil {
				response.ErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("unexpected keys in body:%v", unexpectedKeys)}.Send(c)
				return
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // Restore the body for future reads
		}

		// check if fields have correct type
		var body model.CreateTodoArgs
		if err := c.ShouldBindJSON(&body); err != nil {
			response.ErrorResponse{http.StatusBadRequest, err.Error()}.Send(c)
			return
		}

		// custom validations
		if len(body.Name) < 3 {
			response.ErrorResponse{http.StatusBadRequest, "'name' must be set and have at least 3 characters"}.Send(c)
			return
		}
		if body.CollectionId == "" {
			response.ErrorResponse{http.StatusBadRequest, "'collectionId' must be set"}.Send(c)
			return
		}

		// [DB] check if user owns this collection
		if _, errSelect := database.SelectCollection(body.CollectionId); errSelect != nil {
			response.SendError(c, errSelect)
			return
		}

		// create
		todo, err := database.CreateTodo(body)
		if err != nil {
			response.ErrorResponse{http.StatusInternalServerError, err.Error()}.Send(c)
			return
		}
		response.SendOk(c, todo)
	})

	router.Run()
}
