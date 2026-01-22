package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getUserID(c *gin.Context) (string, bool) {
	userID, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return "", false
	}
	return userID.(string), true
}
