package controller

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"

	"adrianfinantyo.com/adrianfinantyo/go-shorter/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type ShorterController struct{
	*model.BaseController
}

func NewShorterController(config *model.Config, dbConn *model.DatabaseConnection) *ShorterController {
	return &ShorterController{
		BaseController: &model.BaseController{
			Config: config,
			MongoDB: dbConn.MongoDB,
			Redis: dbConn.Redis,
		},
	}
}

func (controller *ShorterController) CreateShortLink(c *gin.Context) {
	collection := controller.MongoDB.Database("go-shorter").Collection("links")
	
	// find the link in the database
	request := c.MustGet("request").(*model.CreateShortLinkRequest)
	result := bson.M{}
	findErr := collection.FindOne(context.Background(), bson.M{
		"short_url": request.ShortURL,
	}).Decode(&result)
	if findErr == nil {
		response := model.Response{
			Status: http.StatusConflict,
			Message: "Error, short link already exist",
			Data: nil,
		}
		c.JSON(http.StatusConflict, response)
		return
	}

	_, insertErr := collection.InsertOne(context.Background(), bson.M{
		"original_url": request.OriginalURL,
		"short_url": request.ShortURL,
	})
	if insertErr != nil {
		panic(insertErr)
	}

	response := model.Response{
		Status: http.StatusCreated,
		Message: "Success, create short link",
		Data: nil,
	}

	c.JSON(response.Status, response)
}

func (controller *ShorterController) GetAllShortLink(c *gin.Context) {
	collection := controller.MongoDB.Database("go-shorter").Collection("links")

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		panic(err)
	}
	
	data := []bson.M{}
	if err = cursor.All(context.Background(), &data); err != nil {
		panic(err)
	}
	
	response := model.Response{
		Status: http.StatusOK,
		Message: "Success, get all data",
		Data: data,
	}
	c.JSON(http.StatusOK, response)
}

func (controller *ShorterController) GetShortLink(c *gin.Context) {
	collection := controller.MongoDB.Database("go-shorter").Collection("links")

	shortURL := c.Param("shortURL")
	result := bson.M{}
	findErr := collection.FindOne(context.Background(), bson.M{
		"short_url": shortURL,
	}).Decode(&result)
	if findErr != nil {
		response := model.Response{
			Status: http.StatusNotFound,
			Message: "Error, short link not found",
			Data: nil,
		}
		c.JSON(http.StatusNotFound, response)
		return
	}

	response := model.Response{
		Status: http.StatusOK,
		Message: "Success, get short link",
		Data: result,
	}
	c.JSON(http.StatusOK, response)
}

func (controller *ShorterController) DeleteShortLink(c *gin.Context) {
	collection := controller.MongoDB.Database("go-shorter").Collection("links")

	shortURL := c.Param("shortURL")
	result := bson.M{}
	findErr := collection.FindOne(context.Background(), bson.M{
		"short_url": shortURL,
	}).Decode(&result)
	if findErr != nil {
		response := model.Response{
			Status: http.StatusNotFound,
			Message: "Error, short link not found",
			Data: nil,
		}
		c.JSON(http.StatusNotFound, response)
		return
	}

	_, deleteErr := collection.DeleteOne(context.Background(), bson.M{
		"short_url": shortURL,
	})
	if deleteErr != nil {
		panic(deleteErr)
	}

	response := model.Response{
		Status: http.StatusOK,
		Message: "Success, delete short link",
		Data: nil,
	}
	c.JSON(http.StatusOK, response)
}

func getDataFromDB(controller *ShorterController, key string) (interface{}, error) {
	collection := controller.MongoDB.Database("go-shorter").Collection("links")
	data := bson.M{}
	findErr := collection.FindOne(context.Background(), bson.M{
		"short_url": key,
	}).Decode(&data)
	if findErr != nil {
		return nil, findErr
	}

	return data, nil
}

func (controller *ShorterController) RedirectShortLink(c *gin.Context) {
	var originalURL string

	cacheStatus := c.MustGet("cacheStatus").(string)
	shortURL := c.Param("shortURL")
	cacheKey := fmt.Sprintf("%s:%s", controller.Config.CacheKeyPrefix, shortURL)

	cachedURL, redisErr := controller.Redis.Get(context.Background(), cacheKey).Result()
	if redisErr != nil {
		log.Error(redisErr)
	}

	if cachedURL != "" && cacheStatus == "on" {
		log.Info("Get data from cache")
		originalURL = cachedURL
	} else if cacheStatus == "on" && redisErr == redis.Nil {
		log.Info("Get data from database")
		data, dbErr := getDataFromDB(controller, shortURL)
		if dbErr != nil {
			response := model.Response{
				Status: http.StatusNotFound,
				Message: "Error, short link not found",
				Data: nil,
			}
			c.JSON(http.StatusNotFound, response)
			return
		}

		originalURL = data.(bson.M)["original_url"].(string)
		controller.Redis.Set(context.Background(), cacheKey, originalURL, time.Duration(controller.Config.CacheTTL) * time.Millisecond)
	} else {
		data, dbErr := getDataFromDB(controller, shortURL)
		if dbErr != nil {
			response := model.Response{
				Status: http.StatusNotFound,
				Message: "Error, short link not found",
				Data: nil,
			}
			c.JSON(http.StatusNotFound, response)
			return
		}

		originalURL = data.(bson.M)["original_url"].(string)
	}
	
	c.Redirect(http.StatusMovedPermanently, originalURL)
}