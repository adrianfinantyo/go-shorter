package main

import (
	"context"
	"fmt"
	"net/http"

	"adrianfinantyo.com/adrianfinantyo/go-shorter/router"
	"adrianfinantyo.com/adrianfinantyo/go-shorter/util"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	config := util.LoadConfig()
	db, err := util.InitMongoDB(config)
	if err != nil {
		panic("Cannot connect to MongoDB")
	}

	defer runOnExit(db)

	routes := router.InitRouter(config, db)

	server := &http.Server{
		Addr: config.AppHost + ":" + config.AppPort,
		Handler: routes,
	}
	server.ListenAndServe()
}

func runOnExit(db *mongo.Client) {
	err := db.Disconnect(context.Background())
	if err != nil {
		fmt.Printf("Error disconnecting from MongoDB: %v\n", err)
	}
}