package web

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Ping status check
//
//	@Summary		Check status
//	@Description	return pong
//	@Tags			status
//	@Accept			plain
//	@Produce		json
//	@Success		200	{object}	map[string]any
//	@Router			/ping [get]
func Ping(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"message": "pong"}) }

// GetRule get target node's rule
//
//	@Summary		Get rule
//	@Description	return node's rule
//	@Tags			rule
//	@Accept			plain
//	@Produce		json
//	@Success		200	{object}	map[string]any
//	@Router			/rule [get]
func GetRule(c *gin.Context) {
	name := c.Query("name")
	path := c.Query("path")

	rule, err := f.Get(name).Get(path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": fmt.Sprintf("query %s on %s fail: %s", path, name, err),
		})

	}
	c.JSON(http.StatusOK, rule)
}
