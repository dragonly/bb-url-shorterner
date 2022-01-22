package router

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Static("/web", "./web")

	api := r.Group("/api")
	{
		registerShortUrlAPIs(api)
	}

	return r
}
