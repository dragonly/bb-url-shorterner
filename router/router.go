package router

import (
	"shurl/dao"

	"github.com/gin-gonic/gin"
)

func SetupRouter(db dao.Database) *gin.Engine {
	r := gin.Default()

	r.Static("/web", "./web")

	api := r.Group("/api")
	{
		registerShortUrlAPIs(api, db)
	}

	return r
}
