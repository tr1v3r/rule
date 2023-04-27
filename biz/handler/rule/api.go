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
		rule.GET("", handler.Ping)
		// return root template
		rule.GET("template", handler.Ping)
		// return info about node
		rule.GET("info", handler.Ping)

		// modify rule
		m := rule.Group("mod")
		{
			// update node rule
			m.POST("", handler.Ping)
			// create or delete rule node
			m.POST("node", handler.Ping)
			// replace template
			m.POST("template", handler.Ping)
			// check operates on node
			m.POST("check", handler.Ping)
		}
	}

	// list nodes
	list := r.Group("list")
	{
		// list all rule tree name
		list.GET("name", handler.Ping)
		// list all nodes
		list.GET("node", handler.Ping)
	}

	// configure about rule enginee
	config := r.Group("config")
	{
		// return config
		config.GET("", handler.Ping)
	}
}
