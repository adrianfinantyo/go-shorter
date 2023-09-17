package util

import (
	"context"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"adrianfinantyo.com/adrianfinantyo/go-shorter/model"
)

func LoadConfig() *model.Config {
	viper.SetConfigType("env")
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	// Set default value
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("APP_HOST", "localhost")

	// Bind config to struct
	var config model.Config
	err = viper.Unmarshal(&config)
	if err != nil {
		panic(err)
	}

	return &config
}

func createMongoDBInstance(config *model.Config) (*mongo.Client, error) {
	var conn *mongo.Client
	var err error
	var maxRetry int = config.MongoDBRetry
	clientOptions := options.Client().ApplyURI(config.MongoDBURI)

	for i := 0; i <= maxRetry; i++ {
		log.Info("Create establish connection with MongoDB...")
		conn, _ = mongo.Connect(context.Background(), clientOptions)
		err = conn.Ping(context.Background(), nil)
		if err != nil {
			log.Error(err)
			if i < maxRetry-1 {
				sleepDuration := time.Duration(config.MongoDBRetryInterval * (i+1)) * time.Millisecond
				log.Warnf("Retrying in %ds...\n", sleepDuration/time.Second)
				os.Stdout.Sync()
				time.Sleep(sleepDuration)
			}
		} else {
			break
		}
	}

	return conn, err
}

func createRedisInstance(config *model.Config) (*redis.Client, error) {
	var conn *redis.Client
	var err error
	redisOptions := &redis.Options{
		Addr: config.RedisHost + ":" + config.RedisPort,
		Password: config.RedisPassword,
		DB: 0,
	}

	log.Info("Create establish connection with Redis...")
	conn = redis.NewClient(redisOptions)
	_, err = conn.Ping(context.Background()).Result()
	if err != nil {
		log.Error(err)
	}

	return conn, err
}

func InitDBConnection(config *model.Config) * model.DatabaseConnection {
	var dbConn model.DatabaseConnection
	var err error

	dbConn.MongoDB, err = createMongoDBInstance(config)
	if err != nil {
		log.Error(err)
		panic("Failed to connect to MongoDB")
	}

	dbConn.Redis, err = createRedisInstance(config)
	if err != nil {
		log.Error(err)
		panic("Failed to connect to Redis")
	}

	return &dbConn
}

func RunOnExit(dbConn *model.DatabaseConnection) {
	dbConn.MongoDB.Disconnect(context.Background())
	dbConn.Redis.Close()
}