package rest

import (
	"file-downloader-service/internal/entity"
	"file-downloader-service/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	taskService service.TaskService
}

func NewTaskController(taskService service.TaskService) *TaskController {
	return &TaskController{taskService: taskService}
}

type CreateTaskRequest struct {
	URLs []string `json:"urls" binding:"required,min=1"`
}

type CreateTaskResponse struct {
	TaskID string `json:"task_id"`
}

type TaskStatusResponse struct {
	Task *entity.Task `json:"task"`
}

func (c *TaskController) CreateTask(ctx *gin.Context) {
	var req CreateTaskRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := c.taskService.CreateTask(req.URLs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, CreateTaskResponse{TaskID: task.ID})
}

func (c *TaskController) GetTaskStatus(ctx *gin.Context) {
	taskID := ctx.Param("id")

	task, err := c.taskService.GetTaskStatus(taskID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	ctx.JSON(http.StatusOK, TaskStatusResponse{Task: task})
}
