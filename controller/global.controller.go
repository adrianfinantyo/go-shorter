package controller

import (
	"context"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"

	"adrianfinantyo.com/adrianfinantyo/go-shorter/model"
	"github.com/gin-gonic/gin"
)

type GlobalController struct{
	*model.BaseController
}

func NewGlobalController(config *model.Config, dbConn *model.DatabaseConnection) *GlobalController {
	return &GlobalController{
		BaseController: &model.BaseController{
			Config: config,
			MongoDB: dbConn.MongoDB,
			Redis: dbConn.Redis,
		},
	}
}

// Ping is a function to handle GET request at /ping endpoint
func (controller *GlobalController) Ping(c *gin.Context) {
	clientData := c.MustGet("clientData").(model.ClientData)

	// test redis
	key := "test"
	value := "test"
	
	redisErr := controller.Redis.Set(context.Background(), key, value, 0).Err()
	if redisErr != nil {
		log.Error(redisErr)
	}

	redisResult, redisErr := controller.Redis.Get(context.Background(), key).Result()
	if redisErr != nil {
		log.Error(redisErr)
	}

	response := model.Response{
		Status: http.StatusOK,
		Message: "Success, ping",
		Data: gin.H{
			"clientData": clientData,
			"redisResult": redisResult,
		},
	}

	c.JSON(response.Status, response)
}

func (controller *GlobalController) ToggleCache(c *gin.Context) {
	cacheKey := fmt.Sprintf("%s:status", controller.Config.CacheKeyPrefix)
	cacheStatus, err := controller.Redis.Get(context.Background(), cacheKey).Result()
	if err != nil {
		log.Error(err)
	}

	if cacheStatus == "on" {
		controller.Redis.Set(context.Background(), cacheKey, "off", 0)
	} else {
		controller.Redis.Set(context.Background(), cacheKey, "on", 0)
	}

	response := model.Response{
		Status: http.StatusOK,
		Message: "Success, toggle cache",
		Data: gin.H{
			"cacheStatus": cacheStatus,
		},
	}

	c.JSON(response.Status, response)
}

func (controller *GlobalController) GetCacheStatus(c *gin.Context) {
	cacheStatus := c.MustGet("cacheStatus").(string)

	response := model.Response{
		Status: http.StatusOK,
		Message: "Success, get cache status",
		Data: gin.H{
			"cacheStatus": cacheStatus,
		},
	}

	c.JSON(response.Status, response)
}