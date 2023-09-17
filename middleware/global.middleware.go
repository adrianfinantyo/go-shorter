package middleware

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"

	"adrianfinantyo.com/adrianfinantyo/go-shorter/model"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func GetClientData() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientData := model.ClientData{
			ClientIP: c.ClientIP(),
			RequestOrigin: c.Request.Header.Get("Origin"),
			Agent: c.Request.UserAgent(),
		}

		c.Set("clientData", clientData)
		c.Next()
	}
}

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

func GetCacheStatus(config *model.Config, conn *model.DatabaseConnection) gin.HandlerFunc {
	return func(c *gin.Context) {
		cacheKey := fmt.Sprintf("%s:status", config.CacheKeyPrefix)
		cacheStatus, err := conn.Redis.Get(context.Background(), cacheKey).Result()
		if err != nil {
			log.Error(err)
		}

		if err == redis.Nil {
			cacheStatus = "on"
			conn.Redis.Set(context.Background(), cacheKey, cacheStatus, 0)
		}

		c.Set("cacheStatus", cacheStatus)
		c.Next()
	}
}

func RequestLogs() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		end := time.Now()
        latency := end.Sub(start)
        clientIP := c.ClientIP()
        status := c.Writer.Status()
        method := c.Request.Method
        path := c.Request.URL.Path

        log.WithFields(log.Fields{
            "status":     status,
            "latency":    latency,
            "clientIP":   clientIP,
            "method":     method,
            "path":       path,
        }).Info("Request details")
	}
}