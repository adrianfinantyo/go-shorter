package router

import (
	"adrianfinantyo.com/adrianfinantyo/go-shorter/controller"
	"adrianfinantyo.com/adrianfinantyo/go-shorter/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

// InitRouter is a function to initialize the router
func InitRouter(config *model.Config, db *mongo.Client) *gin.Engine {
	// Create a default gin router
	r := gin.Default()

	// Initialize controller
	mainController := controller.NewMainController(config)
	shorterController := controller.NewShorterController(config, db)
	
	// Define the route
	mainRouter := r.Group(config.AppPrefix)
	

	// Main route handler
	mainRouter.GET("/ping", mainController.Ping)
	mainRouter.POST("/shorten", shorterController.CreateShortLink)
	mainRouter.GET("/getAllLinks", shorterController.GetShortLink)

	return r
}