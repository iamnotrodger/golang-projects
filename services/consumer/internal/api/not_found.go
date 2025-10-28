package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error":       http.StatusText(http.StatusNotFound),
			"code":        http.StatusNotFound,
			"description": "route not found",
		})
	}
}
