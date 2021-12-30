package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"

	handlers "github.com/abhinavmsra/go-api/internal/api/v1"
)

type routeHandler func(c *gin.Context)

func DummyMiddleware(c *gin.Context) {
	c.JSON(http.StatusForbidden, "")
}

func (s *Server) Routes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.Use(cors.Default())

	v1 := router.Group("/api/v1")
	{
		v1.GET("/status", handlers.Status)

		v1.GET("/merchants", handlers.IndexMerchant(s.Storage))
		v1.POST("/merchants", handlers.CreateMerchant(s.Storage))
		v1.GET("/merchants/:id", handlers.ShowMerchant(s.Storage))
		v1.PATCH("/merchants/:id", handlers.UpdateMerchant(s.Storage))
		v1.DELETE("/merchants/:id", handlers.DeleteMerchant(s.Storage))
	}

	return router
}
