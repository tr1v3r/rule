package main

import (
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	"github.com/riverchu/rule/biz/handler"
	"github.com/riverchu/rule/biz/handler/rule"
	_ "github.com/riverchu/rule/docs"
)

//	@title			R1v3r's rule engine
//	@version		1.0
//	@description	This is a rule engine server.
//	@termsOfService	http://localhost/terms/

//	@contact.name	R1v3r
//	@contact.url	http://localhost/support
//	@contact.email	churiver@outlook.com

//	@license.name	MIT License
//	@license.url	https://www.mit.edu/~amini/LICENSE.md

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.basic	BasicAuth

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/

func register(r *gin.Engine) *gin.Engine {
	r.GET("ping", handler.Ping)

	apiV1 := r.Group("api/v1")
	{
		rule.RegisterApi(apiV1)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
