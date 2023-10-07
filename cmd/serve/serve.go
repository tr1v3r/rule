package main

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	"github.com/tr1v3r/rule"
	_ "github.com/tr1v3r/rule/docs"
	"github.com/tr1v3r/rule/web"
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

var timeout, _ = time.ParseDuration(os.Getenv("SHUTDOWN_TIMEOUT"))

func main() {
	var rules []rule.Rule
	web.InitForest(web.DefaultBuilder(rules...))

	if timeout == 0 {
		timeout = 3 * time.Second
	}
	web.Serve(timeout, register(gin.Default()))
}

func register(r *gin.Engine) *gin.Engine {
	apiV1 := r.Group("api/v1")
	{
		web.RegisterApi(apiV1)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
