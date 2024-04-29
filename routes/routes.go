package routes

import (
	"short-url/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	server.GET("/health", services.GetHealth)
	server.POST("/add", services.AddUrl)
	server.GET("/:url", services.ResolveUrl)
}
