package rule

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// QueryRule return rule
func QueryRule(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
