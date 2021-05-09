package rest

import (
	"github.com/gin-gonic/gin"

	"shopping/rest/catalog"
	"shopping/rest/commodity"
)

func Init() *gin.Engine {
	r := gin.New()

	r.Use(gin.Recovery())

	v1 := r.Group("/v1")
	{
		catalog.MakeHandler(v1)
		commodity.MakeHandler(v1)
	}

	return r
}
