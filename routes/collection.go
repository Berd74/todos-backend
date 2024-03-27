package routes

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
	"todoBackend/database"
	"todoBackend/response"
	"todoBackend/types"
	"todoBackend/utils"
)

func Collection(rg *gin.RouterGroup) {

	rg.GET("/:id", utils.VerifyToken, func(c *gin.Context) {
		userId, _ := c.Get("userId")
		role, _ := c.Get("role")
		collectionId := c.Param("id")

		collection, err := database.SelectCollection(collectionId)

		if err != nil {
			response.SendError(c, err)
			return
		}

		if userId != collection.UserId && role != types.Admin {
			response.ErrorResponse{Code: http.StatusForbidden, Message: "This collection doesn't belong to you"}.Send(c)
			return
		}

		response.SendOk(c, collection)
	})

	rg.POST("/", utils.VerifyToken, func(c *gin.Context) {
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

	rg.DELETE("/", utils.VerifyToken, func(c *gin.Context) {
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
			response.ErrorResponse{Code: http.StatusForbidden, Message: "Provided collection ID that doesn't belong to you or does not exist."}.Send(c)
			return
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

	rg.PUT("/:id", utils.VerifyToken, func(c *gin.Context) {
		userId, _ := c.Get("userId")
		role, _ := c.Get("role")
		collectionId := c.Param("id")
		collection, errSelect := database.SelectCollection(collectionId)

		if errSelect != nil {
			response.SendError(c, errSelect)
			return
		}

		if userId != collection.UserId && role != types.Admin {
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

}
