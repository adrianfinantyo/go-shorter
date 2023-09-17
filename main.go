package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"adrianfinantyo.com/adrianfinantyo/go-shorter/model"
	"adrianfinantyo.com/adrianfinantyo/go-shorter/router"
	"adrianfinantyo.com/adrianfinantyo/go-shorter/util"
)

func main() {
	config := util.LoadConfig()
	dbConn := util.InitDBConnection(config)

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	
	// Handle system interrupt
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt, syscall.SIGTERM)
	go systemInterruptHandler(channel, dbConn)

	// Run on exit
	defer runOnExit(dbConn)

	// Initialize router and server
	routes := router.InitRouter(config, dbConn)
	server := &http.Server{
		Addr: config.AppHost + ":" + config.AppPort,
		Handler: routes,
	}

	server.ListenAndServe()
}

func systemInterruptHandler(channel chan os.Signal, dbConn *model.DatabaseConnection) {
	<-channel
	fmt.Println("Interrupt signal received.")
	runOnExit(dbConn)
	os.Exit(0)
}

func runOnExit(dbConn *model.DatabaseConnection) {
	util.RunOnExit(dbConn)
	fmt.Println("Server is shutting down...")
}