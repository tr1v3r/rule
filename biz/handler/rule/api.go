package rule

import (
	"github.com/gin-gonic/gin"

	"github.com/riverchu/rule/biz/handler"
)

func RegisterApi(r *gin.RouterGroup) {
	// return rules
	rule := r.Group("rule")
	{
		// ping pong
		rule.GET("ping", handler.Ping)

		// return node rule
		rule.GET("")
		// return root template
		rule.GET("template")
		// return info about node
		rule.GET("info")

		// modify rule
		m := rule.Group("mod")
		{
			// update node rule
			m.POST("")
			// create or delete rule node
			m.POST("node")
			// replace template
			m.POST("template")
			// check operates on node
			m.POST("check")
		}
	}

	// list nodes
	list := r.Group("list")
	{
		// list all rule tree name
		list.GET("name")
		// list all nodes
		list.GET("node")
	}

	// configure about rule enginee
	config := r.Group("config")
	{
		// return config
		config.GET("")
	}
}
