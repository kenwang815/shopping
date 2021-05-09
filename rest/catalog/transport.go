package catalog

import "github.com/gin-gonic/gin"

func MakeHandler(r *gin.RouterGroup) {
	g := r.Group("/catalog")
	{
		g.GET("", FindCatalog)
		g.POST("", RegisterCatalog)
		g.DELETE("/:id", DeleteCatalog)
		g.PUT("", UpdateCatalog)
	}
}
