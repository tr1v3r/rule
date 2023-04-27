package handler

import (
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
