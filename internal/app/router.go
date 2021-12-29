package app

import (
	"net/http"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	handlers "github.com/abhinavmsra/go-api/internal/api/v1"
)

func DummyMiddleware(c *gin.Context) {
	c.JSON(http.StatusForbidden, "")
  }

func Router() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	
	router := gin.Default()
	router.Use(cors.Default())
	router.Use(DummyMiddleware)

	v1 := router.Group("/api/v1")
	{
		v1.GET("/status", handlers.Status)
	}

	return router
}
