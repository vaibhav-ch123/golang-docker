package main

import (
	"fmt"
	"httpserver/database"
	"httpserver/server"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

var ShutDownTimeout = 5 * time.Second

func main() {

	chn := make(chan os.Signal, 1)

	signal.Notify(chn, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := server.SetupUpRoutes()
	if err := database.ConnectAndMigrate(
		"localhost",
		"5432",
		"todo",
		"postgres",
		"postgres123",
		database.SSLModeDisable); err != nil {
		logrus.Panicf("Failed to initialize and migrate database with error: %+v", err)
	}
	logrus.Print("migration successful!")

	go func() {
		if err := srv.Run(":8080"); err != nil {
			fmt.Println("server is fail to run ", err)
		}
	}()

	fmt.Println("server is running")

	<-chn

	fmt.Println("server is shutdowning..")

	if err := srv.ShutDown(ShutDownTimeout); err != nil {
		fmt.Println("problem in shutdown!")
	}

	fmt.Println("server is shutdown")
}
