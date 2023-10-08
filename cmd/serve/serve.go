package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware

	"github.com/tr1v3r/pkg/log"
	"github.com/tr1v3r/rule"
	_ "github.com/tr1v3r/rule/docs"
	"github.com/tr1v3r/rule/driver"
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
	web.InitForest(web.DefaultBuilder(load))

	go func() {
		for range time.Tick(5 * time.Second) {
			log.Info("refreshing forest...")
			web.RefreshForest()
		}
	}()

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

var defaultFilename = "../../conf/rules.json"

type RuleDataItem struct {
	Path      string `json:"path"`
	Operators []struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	} `json:"operators"`
}

func load() (rules []*rule.Rule) {
	var filename = os.Getenv("RULES_FILE")
	if filename == "" {
		filename = defaultFilename
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Error("read file fail: %s", err)
		return nil
	}

	var items = []RuleDataItem{}
	if err = json.Unmarshal(data, &items); err != nil {
		log.Error("unmarshal data fail: %s", err)
		return nil
	}
	for _, line := range items {
		var ops []driver.Operator
		for _, opData := range line.Operators {
			var op driver.Operator
			switch opData.Type {
			case "json":
				op = new(driver.JSONOperator)
			case "yaml":
				op = new(driver.YAMLOperator)
			case "curl":
				op = new(driver.CURLOperator)
			}
			if op != nil {
				op.Load(opData.Data)
			}
			ops = append(ops, op)
		}
		rules = append(rules, &rule.Rule{Path: line.Path, Operators: ops})
	}
	return rules
}
