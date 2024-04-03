package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
	"todoBackend/database"
	"todoBackend/model"
	"todoBackend/response"
	"todoBackend/utils"
)

func Todo(rg *gin.RouterGroup) {

	rg.GET("/", utils.VerifyToken, func(c *gin.Context) {
		userIdString := c.GetString("userId")

		todoIds := utils.SplitOrNil(utils.StringOrNil(c.Query("todoIds")))
		userIds := utils.SplitOrNil(utils.StringOrNil(c.Query("userIds")))
		collectionIds := utils.SplitOrNil(utils.StringOrNil(c.Query("collectionIds")))
		done := utils.StringToBoolOrNil(c.Query("done"))

		todos, err := database.SelectTodos(userIdString, todoIds, userIds, collectionIds, done)

		if err != nil {
			response.SendError(c, err)
			return
		}

		response.SendOk(c, todos)
	})

	rg.DELETE("/", utils.VerifyToken, func(c *gin.Context) {
		userIdString := c.GetString("userId")
		_todoIds := c.Query("todoIds")
		todoIds := strings.Split(_todoIds, ",")

		test, err := database.AreUserTodos(userIdString, todoIds)

		if err != nil {
			response.SendError(c, err)
			return
		}
		if !test {
			response.ErrorResponse{Code: http.StatusForbidden, Message: "One of the todo ID that doesn't belong to you or does not exist."}.Send(c)
			return
		}

		amount, err := database.DeleteTodo(todoIds)
		if err != nil {
			response.SendError(c, err)
			return
		}

		response.SendOk(c, map[string]any{
			"removedItems": amount,
		})
	})

	rg.POST("/", utils.VerifyToken, func(c *gin.Context) {
		clientId := c.GetString("userId")

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

		// check if fields have correct types
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
		if _, errSelect := database.AreUserCollections(clientId, []string{body.CollectionId}); errSelect != nil {
			response.SendError(c, errSelect)
			return
		}

		// create
		todo, err := database.CreateTodo(clientId, body)
		if err != nil {
			response.ErrorResponse{http.StatusInternalServerError, err.Error()}.Send(c)
			return
		}
		response.SendOk(c, todo)
	})

	rg.PUT("/:id", utils.VerifyToken, func(c *gin.Context) {
		userIdString := c.GetString("userId")

		todoId := c.Param("id")
		test, errSelect := database.AreUserTodos(userIdString, []string{todoId})
		if errSelect != nil {
			response.SendError(c, errSelect)
			return
		}

		if !test {
			response.ErrorResponse{Code: http.StatusForbidden, Message: "Provided todo ID that doesn't belong to you or does not exist."}.Send(c)
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

		updateErr := database.UpdateTodo(userIdString, todoId, bodyFull)
		if updateErr != nil {
			response.SendError(c, updateErr)
			return
		}

		collection, errSelect := database.SelectTodos(userIdString, &[]string{todoId}, nil, nil, nil)
		if errSelect != nil {
			response.SendError(c, errSelect)
			return
		}

		response.SendOk(c, collection[0])
	})

	rg.PUT("/move/:id", utils.VerifyToken, func(c *gin.Context) {
		clientId := c.GetString("userId")
		todoId := c.Param("id")

		// checking body
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
		var lookFor = []string{"after", "before"}
		if len(bodyFull) == 0 {
			response.ErrorResponse{http.StatusBadRequest, `you need to provide "after" or "before" value to move item`}.Send(c)
			return
		}
		if unexpectedKeys := utils.GetUnexpectedKeys(&bodyFull, lookFor); unexpectedKeys != nil {
			response.ErrorResponse{Code: http.StatusBadRequest, Message: fmt.Sprintf("unexpected keys in body:%v", unexpectedKeys)}.Send(c)
			return
		}
		if len(bodyFull) > 1 {
			response.ErrorResponse{http.StatusBadRequest, `only one parameter can be set: "after" or "before"`}.Send(c)
			return
		}

		// getting parameters
		var isAfter = bodyFull["after"] != nil
		var targetId string
		if isAfter {
			targetId = bodyFull["after"].(string)
		} else {
			targetId = bodyFull["before"].(string)
		}

		if targetId == todoId {
			response.ErrorResponse{http.StatusBadRequest, `ids of the items must be different`}.Send(c)
			return
		}

		// checking collections ids
		test, errSelect := database.AreUserTodos(clientId, []string{todoId, targetId})
		if errSelect != nil {
			response.SendError(c, errSelect)
			return
		}
		if !test {
			response.ErrorResponse{Code: http.StatusForbidden, Message: "Provided todo ID that doesn't belong to you or does not exist."}.Send(c)
			return
		}

		moveErr := database.MoveTodo(clientId, todoId, targetId, isAfter)
		if moveErr != nil {
			response.SendError(c, moveErr)
			return
		}

		response.SendOk(c, "ok")
	})

}
