package v1

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Status(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, "")
}
