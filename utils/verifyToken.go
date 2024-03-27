package utils

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"todoBackend/firebase"
	"todoBackend/types"
)

// Middleware to verify Firebase token and add custom fields to the context
func VerifyToken(c *gin.Context) {
	token := c.GetHeader("Authorization")

	if strings.HasPrefix(token, "Bearer ") {
		token = strings.TrimPrefix(token, "Bearer ")
	}

	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		c.Abort()
		return
	}

	decodedToken, err := firebase.AuthClient.VerifyIDToken(context.Background(), token)

	if err != nil {
		log.Printf("error verifying ID token: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "bad token"})
		c.Abort()
		return
	}

	// Assuming googleIdToUuid is a function you've implemented that converts a Google ID to your UUID format
	userId := googleIdToUuid(decodedToken.UID) // Replace with actual function call

	// Set custom request fields
	c.Set("userId", userId)
	c.Set("role", types.User)

	c.Next()
}
