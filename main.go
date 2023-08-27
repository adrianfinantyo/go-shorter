package main

import (
	"net/http"

	"adrianfinantyo.com/adrianfinantyo/go-shorter/router"
	"adrianfinantyo.com/adrianfinantyo/go-shorter/util"
)

func main() {
	config := util.LoadConfig()
	db := util.InitMongoDB(config)

	routes := router.InitRouter(config, db)

	server := &http.Server{
		Addr: config.AppHost + ":" + config.AppPort,
		Handler: routes,
	}
	server.ListenAndServe()
}