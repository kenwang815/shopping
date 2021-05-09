package commodity

import "github.com/gin-gonic/gin"

func MakeHandler(r *gin.RouterGroup) {
	g := r.Group("/commodity")
	{
		g.GET("", FindCommodity)
		g.POST("", RegisterCommodity)
		g.DELETE("/:id", DeleteCommodity)
		g.PUT("", UpdateCommodity)
	}
}
