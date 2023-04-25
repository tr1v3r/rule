package main

import (
	"github.com/gin-gonic/gin"

	"github.com/riverchu/rule/biz/handler"
	"github.com/riverchu/rule/biz/handler/rule"
)

func register(r *gin.Engine) *gin.Engine {
	r.GET("ping", handler.Ping)

	apiV1 := r.Group("api/v1")
	{
		rule.RegisterApi(apiV1)
	}

	return r
}
