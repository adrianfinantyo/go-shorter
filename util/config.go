package util

import (
	"context"
	"fmt"
	"time"

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

	if config.MongoDBUser != "" && config.MongoDBPassword != "" {
		config.MongoDBURI = fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", config.MongoDBUser, config.MongoDBPassword, config.MongoDBHost, config.MongoDBPort, config.MongoDBName)
	} else {
		config.MongoDBURI = fmt.Sprintf("mongodb://%s:%s/%s", config.MongoDBHost, config.MongoDBPort, config.MongoDBName)
	}

	return &config
}

func createMongoDBConnection(config *model.Config) *mongo.Client {
	clientOptions := options.Client().ApplyURI(config.MongoDBURI)

	conn, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}
	return conn
}

func InitMongoDB(config *model.Config) (*mongo.Client, error) {
	var conn *mongo.Client
	var err error
	var maxRetry int = config.MongoDBRetry

	for i := 0; i < maxRetry; i++ {
		fmt.Printf("[%d] Create establish connection with MongoDB...\n", i+1)
		conn = createMongoDBConnection(config)
		err = conn.Ping(context.Background(), nil)
		if err != nil {
			fmt.Printf("Error connecting to MongoDB: %v\n", err)
			if i < maxRetry-1 {
				sleepDuration := time.Duration(config.MongoDBRetryInterval * i+1)
				fmt.Printf("Retrying in %ds...\n", sleepDuration/time.Second)
				time.Sleep(sleepDuration)
			}
		} else {
			break
		}
	}

	return conn, err
}