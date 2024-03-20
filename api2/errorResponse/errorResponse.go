package errorResponse

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("ErrorResponse: %v, %v", e.Code, e.Message)
}

func SendErrorResponse(c *gin.Context, err error) {
	if err == nil {
		return // No error, so just return
	}

	// Attempt to type assert 'err' to 'ErrorResponse'
	if errResp, ok := err.(ErrorResponse); ok {
		// 'err' is an ErrorResponse, so use its 'Code' and 'Message'
		c.JSON(errResp.Code, gin.H{"error": errResp.Message})
	} else {
		// 'err' is not an ErrorResponse, so fall back to a default error response
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
