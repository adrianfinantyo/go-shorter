package router

import (
	"adrianfinantyo.com/adrianfinantyo/go-shorter/controller"
	"adrianfinantyo.com/adrianfinantyo/go-shorter/middleware"
	"adrianfinantyo.com/adrianfinantyo/go-shorter/model"
	"github.com/gin-gonic/gin"
)

// InitRouter is a function to initialize the router
func InitRouter(config *model.Config, conn *model.DatabaseConnection) *gin.Engine {
	// Create a default gin router
	r := gin.Default()

	// Initialize controller
	globalController := controller.NewGlobalController(config, conn)
	shorterController := controller.NewShorterController(config, conn)
	
	// Define the route
	r.Use(middleware.GetClientData(), middleware.RequestLogs())

	// Define the route group
	mainRoute := r.Group("/")
	apiRoute := r.Group(config.AppPrefix)

	// Main route handler
	apiRoute.GET("/ping", globalController.Ping)
	apiRoute.POST("/url", middleware.ValidateRequestBody(model.CreateShortLinkRequest{}), shorterController.CreateShortLink)
	apiRoute.GET("/url", shorterController.GetAllShortLink)
	apiRoute.GET("/url/:shortURL", shorterController.GetShortLink)
	apiRoute.DELETE("/url/:shortURL", shorterController.DeleteShortLink)
	apiRoute.PUT("/cache", globalController.ToggleCache)
	apiRoute.GET("/cache", middleware.GetCacheStatus(config, conn), globalController.GetCacheStatus)

	// Shorter route handler
	mainRoute.GET("/:shortURL", middleware.GetCacheStatus(config, conn), shorterController.RedirectShortLink)

	return r
}