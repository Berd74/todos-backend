package routes

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
	"todoBackend/database"
	"todoBackend/response"
	"todoBackend/utils"
)

func Collection(rg *gin.RouterGroup) {

	rg.GET("/", utils.VerifyToken, func(c *gin.Context) {
		userIdString := c.GetString("userId")

		userIds := utils.SplitOrNil(utils.StringOrNil(c.Query("userIds")))
		collectionIds := utils.SplitOrNil(utils.StringOrNil(c.Query("collectionIds")))

		collection, err := database.SelectCollection(userIdString, userIds, collectionIds)

		if err != nil {
			response.SendError(c, err)
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
		userIdString := c.GetString("userId")
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
		userIdString := c.GetString("userId")

		collectionId := c.Param("id")
		test, errSelect := database.AreUserCollections(userIdString, []string{collectionId})
		if errSelect != nil {
			response.SendError(c, errSelect)
			return
		}
		if !test {
			response.ErrorResponse{Code: http.StatusForbidden, Message: "Provided collection ID that doesn't belong to you or does not exist."}.Send(c)
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

		collection, errSelect := database.SelectCollection(userIdString, &[]string{collectionId}, nil)
		if errSelect != nil {
			response.SendError(c, errSelect)
			return
		}

		response.SendOk(c, collection)
	})

}
