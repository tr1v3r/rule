package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping return pong in json format
// used for status check
func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
