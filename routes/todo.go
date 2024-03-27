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
		userId, _ := c.Get("userId")
		userIdString := userId.(string)

		userIds := utils.SplitOrNil(utils.StringOrNil(c.Query("userIds")))
		collectionIds := utils.SplitOrNil(utils.StringOrNil(c.Query("collectionIds")))
		done := utils.StringToBoolOrNil(c.Query("done"))

		todos, err := database.SelectTodos(userIdString, userIds, collectionIds, done)

		if err != nil {
			response.SendError(c, err)
			return
		}

		response.SendOk(c, todos)
	})

	rg.DELETE("/", utils.VerifyToken, func(c *gin.Context) {
		userId, _ := c.Get("userId")
		userIdString := userId.(string)
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

}
