package controller

import (
	"context"
	"net/http"

	"adrianfinantyo.com/adrianfinantyo/go-shorter/model"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ShorterController struct{
	config *model.Config
	db *mongo.Client
}

func NewShorterController(config *model.Config, db *mongo.Client) *ShorterController {
	return &ShorterController{
		config: config,
		db: db,
	}
}

func (controller *ShorterController) CreateShortLink(c *gin.Context) {
	collection := controller.db.Database("go-shorter").Collection("links")
	
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

	c.JSON(http.StatusOK, response)
}

func (controller *ShorterController) GetShortLink(c *gin.Context) {
	collection := controller.db.Database("go-shorter").Collection("links")

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

func (controller *ShorterController) RedirectShortLink(c *gin.Context) {
	collection := controller.db.Database("go-shorter").Collection("links")

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

	originalURL := result["original_url"].(string)
	c.Redirect(http.StatusMovedPermanently, originalURL)
}