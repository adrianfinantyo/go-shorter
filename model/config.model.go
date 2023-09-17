package model

import (
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type Config struct {
	AppPort string `mapstructure:"APP_PORT"`
	AppHost string `mapstructure:"APP_HOST"`
	AppPrefix string `mapstructure:"APP_PREFIX"`
	AppSecret string `mapstructure:"APP_SECRET"`
	AppEnv string `mapstructure:"APP_ENV"`
	MongoDBURI string `mapstructure:"MONGODB_URI"`
	MongoDBRetry int `mapstructure:"MONGODB_RETRY"`
	MongoDBRetryInterval int `mapstructure:"MONGODB_RETRY_INTERVAL"`
	RedisHost string `mapstructure:"REDIS_HOST"`
	RedisPort string `mapstructure:"REDIS_PORT"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	CacheTTL int `mapstructure:"CACHE_TTL"`
	CacheKeyPrefix string `mapstructure:"CACHE_KEY_PREFIX"`
	CacheKeyAppPrefix string `mapstructure:"CACHE_KEY_APP_PREFIX"`
}

type BaseController struct {
	Config *Config
	MongoDB *mongo.Client
	Redis *redis.Client
}

type Response struct {
	Message string `json:"message"`
	Status int `json:"status"`
	Data interface{} `json:"data"`
}

type DatabaseConnection struct {
	MongoDB *mongo.Client
	Redis *redis.Client
}