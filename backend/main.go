package main

import (
	"backend/config"
	"backend/db"
	"backend/middleware"
	"backend/router"
	"fmt"
	"github.com/hibiken/asynq"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var asynqServer *asynq.Server
	var asynqServerStop func()

	if config.Config.Server.Environment == "server" {
		server := asynq.NewServer(
			middleware.RedisClient,
			asynq.Config{
				Concurrency: 10,
				Queues: map[string]int{
					"critical": 6,
					"default":  3,
					"low":      1,
				},
			},
		)
		mux := asynq.NewServeMux()
		middleware.RegisterTaskHandlers(mux)

		asynqServer = server
		go func() {
			if err := server.Run(mux); err != nil {
				log.Printf("Asynq server stopped: %v", err)
			}
		}()

		asynqServerStop = func() {
			middleware.AsynqClient.Close()
		}
	}

	r := router.Router
	go func() {
		r.Run(":" + config.Config.Server.Port)
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")

	if asynqServer != nil {
		fmt.Println("Shutting down Asynq server...")
		asynqServer.Stop()
	}

	if asynqServerStop != nil {
		asynqServerStop()
	}

	defer db.DB.Close()

	fmt.Println("Server exiting...")
}
