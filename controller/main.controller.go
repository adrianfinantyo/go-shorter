package controller

import (
	"net/http"

	"adrianfinantyo.com/adrianfinantyo/go-shorter/model"
	"github.com/gin-gonic/gin"
)

type MainController struct{
	config *model.Config
}

func NewMainController(config *model.Config) *MainController {
	return &MainController{
		config: config,
	}
}

// Ping is a function to handle GET request at /ping endpoint
func (controller *MainController) Ping(c *gin.Context) {
	clientData := c.MustGet("clientData").(model.ClientData)

	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
		"status": http.StatusOK,
		"data": clientData,
	})
}