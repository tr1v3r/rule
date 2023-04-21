package main

import (
	"github.com/gin-gonic/gin"

	"github.com/riverchu/rule/biz/handler"
)

func register(r *gin.Engine) *gin.Engine {
	r.GET("ping", handler.Ping)

	return r
}
