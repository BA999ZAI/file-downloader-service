package rest

import (
	"file-downloader-service/internal/service"

	"github.com/gin-gonic/gin"
)

func SetupRouter(taskService service.TaskService) *gin.Engine {
	router := gin.Default()

	taskController := NewTaskController(taskService)

	api := router.Group("/api/v1")
	{
		tasks := api.Group("/tasks")
		{
			tasks.POST("", taskController.CreateTask)
			tasks.GET("/:id", taskController.GetTaskStatus)
		}
	}

	return router
}
