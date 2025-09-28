package main

import (
	"file-downloader-service/internal/controller/rest"
	"file-downloader-service/internal/repository"
	"file-downloader-service/internal/service"
	"file-downloader-service/internal/usecase"
	"file-downloader-service/pkg/downloader"
	"file-downloader-service/pkg/storage"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	downloadDir = "./downloads"
	port        = ":8080"
)

func main() {
	fileStorage := storage.NewLocalFileStorage(downloadDir)
	if err := fileStorage.EnsureDownloadDir(); err != nil {
		log.Fatalf("Failed to create download directory: %v", err)
	}

	taskRepo := repository.NewInMemoryTaskRepository()
	downloader := downloader.NewHTTPDownloader()
	taskUseCase := usecase.NewTaskUseCase(taskRepo, downloader, fileStorage)
	taskService := service.NewTaskService(taskUseCase)

	taskUseCase.RecoverPendingTasks()

	go taskUseCase.ProcessPendingTasks()

	router := rest.SetupRouter(taskService)

	go func() {
		fmt.Printf("Server starting on port %s\n", port)
		if err := router.Run(port); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")

	taskUseCase.StopProcessing()

	fmt.Println("Server stopped")
}
