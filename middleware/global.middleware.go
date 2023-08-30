package middleware

import (
	"net/http"
	"reflect"

	"adrianfinantyo.com/adrianfinantyo/go-shorter/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidateRequestBody(requestModel interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestPointer := reflect.New(reflect.TypeOf(requestModel)).Interface()

		if(c.Request.ContentLength == 0) {
			var errors []string

			// Get all fields from the request model
			requestType := reflect.TypeOf(requestModel)
			for i := 0; i < requestType.NumField(); i++ {
				field := requestType.Field(i)
				errors = append(errors, field.Tag.Get("json") + " is required")
			}

			response := model.Response{
				Status: http.StatusBadRequest,
				Message: "Error, bad request",
				Data: errors,
			}
			c.JSON(http.StatusBadRequest, response)
			c.Abort()
			return
		}

		if err := c.ShouldBindJSON(requestPointer); err != nil {
			var errors []string
			for _, err := range err.(validator.ValidationErrors) {
				field, _ := reflect.TypeOf(requestModel).FieldByName(err.Field())
				errors = append(errors, field.Tag.Get("json") + " is " + err.Tag())
			}
			response := model.Response{
				Status: http.StatusBadRequest,
				Message: "Error, bad request",
				Data: errors,
			}
			c.JSON(http.StatusBadRequest, response)
			c.Abort()
			return
		}

		// Proceed with the next handler
		c.Set("request", requestPointer)
		c.Next()
	}
}