package controller

import (
	"context"
	"net/http"

	"adrianfinantyo.com/adrianfinantyo/go-shorter/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type ShorterController struct{
	config *model.Config
	db *mongo.Client
}

type CreateShortLinkRequest struct {
	OriginalURL string `json:"original_url" binding:"required"`
	ShortURL string `json:"short_url" binding:"required"`
}

func NewShorterController(config *model.Config, db *mongo.Client) *ShorterController {
	return &ShorterController{
		config: config,
		db: db,
	}
}

func (controller *ShorterController) CreateShortLink(c *gin.Context) {
	var request CreateShortLinkRequest

	if c.Request.ContentLength == 0 {
		response := model.Response{
			Status: http.StatusBadRequest,
			Message: "Error, bad request",
			Data: nil,
		}
		c.JSON(http.StatusBadRequest, response)
	 }

	validationErr := c.ShouldBindJSON(&request)
	if validationErr != nil {
		var validationErrors []string
		for _, err := range validationErr.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, err.Tag() + " " + err.Field())
		}

		response := model.Response{
			Status: http.StatusBadRequest,
			Message: "Error, bad request",
			Data: validationErrors,
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	collection := controller.db.Database("go-shorter").Collection("links")
	
	// find the link in the database
	result := bson.M{}
	findErr := collection.FindOne(context.Background(), bson.M{
		"short_url": "google",
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
		"original_url": "https://www.google.com",
		"short_url": "google",
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