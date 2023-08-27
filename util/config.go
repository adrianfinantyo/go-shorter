package util

import (
	"context"

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

func InitMongoDB(config *model.Config) *mongo.Client {
	clientOptions := options.Client().ApplyURI(config.MongoDBURI)

	conn, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		panic(err)
	}
	return conn
}