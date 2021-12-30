package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"

	handlers "github.com/abhinavmsra/go-api/internal/api/v1"
)

type routeHandler func(c *gin.Context)

func AuthorizationMiddleware(c *gin.Context) {
	token := c.Request.Header.Get("Authorization")
	parts := strings.Split(token, " ")

	if len(parts) != 2 || parts[0] != "Bearer" {
		c.AbortWithStatusJSON(http.StatusForbidden, "Invalid Bearer Token Format")
		return
	}

	if parts[1] == "" {
		c.AbortWithStatusJSON(http.StatusForbidden, "Missing Token")
		return
	}

	c.Next()
}

func (s *Server) Routes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.Use(cors.Default())
	router.Use(AuthorizationMiddleware)

	v1 := router.Group("/api/v1")
	{
		v1.GET("/status", handlers.Status)

		v1.GET("/merchants", handlers.IndexMerchant(s.Storage))
		v1.POST("/merchants", handlers.CreateMerchant(s.Storage))
		v1.GET("/merchants/:id", handlers.ShowMerchant(s.Storage))
		v1.PATCH("/merchants/:id", handlers.UpdateMerchant(s.Storage))
		v1.DELETE("/merchants/:id", handlers.DeleteMerchant(s.Storage))

		v1.GET("/members", handlers.IndexMember(s.Storage))
		v1.POST("/members", handlers.CreateMember(s.Storage))
		v1.GET("/members/:id", handlers.ShowMember(s.Storage))
		v1.PATCH("/members/:id", handlers.UpdateMember(s.Storage))
		v1.DELETE("/members/:id", handlers.DeleteMember(s.Storage))
	}

	return router
}
