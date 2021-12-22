package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	handlers "github.com/abhinavmsra/go-api/internal/api/v1"
)

func Router() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	v1 := router.Group("/api/v1")
	{
		v1.GET("/status", handlers.Status)
	}

	return router
}
