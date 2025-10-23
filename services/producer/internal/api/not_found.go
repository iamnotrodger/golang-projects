package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iamnotrodger/golang-kafka/services/producer/internal/model"
)

func NotFound() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusNotFound, &model.Error{
			Error:       http.StatusText(http.StatusNotFound),
			Code:        http.StatusNotFound,
			Description: "route not found",
		})
	}
}
