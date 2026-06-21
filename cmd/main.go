package main

import (
	"fmt"
	"httpserver/server"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var ShutDownTimeout = 5 * time.Second

func main() {

	chn := make(chan os.Signal, 1)

	signal.Notify(chn, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	srv := server.SetupUpRoutes()

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
