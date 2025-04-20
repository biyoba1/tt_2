package main

import (
	"context"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	mdk "test_task"
	initializer "test_task/initializators"
	"test_task/internal/handler"
	"test_task/internal/repository"
	"test_task/internal/service"
)

func init() {
	initializer.LoadEnvVariables()
	initializer.PingDatabase()
}

/*
	TODO
	Можно дополнительно покрыть код тестами,
	если сервер будет разрастаться логикой,
	на данном моменте не стал
*/

func main() {
	postgres, err := repository.NewPostgresDB(repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		DBName:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		log.Fatalf("Failed to initialize db: %s", err.Error())
	}
	defer postgres.Close()
	repos := repository.NewRepository(postgres)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	srv := new(mdk.Server)
	port := os.Getenv("PORT")

	go func() {
		log.Printf("Server is running on port %s", port)
		if err := srv.Run(port, handlers.RegisterRoutes()); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error running server: %s", err.Error())
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("Shutting down server gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	services.CancelAllTasks()
	log.Println("All active tasks have been canceled.")

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("HTTP server forced to shutdown: %v", err)
	}
	if err := repos.UpdateTasksOnShutdown(ctx); err != nil {
		log.Printf("Failed to update task statuses on shutdown: %v", err)
	}
	log.Println("Server exited gracefully")
}
