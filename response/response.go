package response

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OkResponse struct {
	Data any `json:"data"`
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("ErrorResponse: %v, %v", e.Code, e.Message)
}

func (e ErrorResponse) Send(c *gin.Context) {
	SendError(c, e)
}

func SendError(c *gin.Context, err error) {
	if err == nil {
		return // No error, so just return
	}

	// Attempt to types assert 'err' to 'ErrorResponse'
	if errResp, ok := err.(ErrorResponse); ok {
		// 'err' is an ErrorResponse, so use its 'Code' and 'Message'
		c.JSON(errResp.Code, gin.H{"error": errResp.Message})
	} else {
		// 'err' is not an ErrorResponse, so fall back to a default error response
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func SendOk(c *gin.Context, data any) {
	if data == nil {
		c.Status(http.StatusOK)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": data,
	})
}
