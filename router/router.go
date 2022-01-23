package router

import (
	"net/http"
	"shurl/dao"

	"github.com/gin-gonic/gin"
)

func SetupRouter(db dao.Database) *gin.Engine {
	r := gin.Default()

	r.Static("/web", "./web")
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/web")
	})

	api := r.Group("/api")
	{
		registerShortUrlAPIs(api, db)
	}

	return r
}
