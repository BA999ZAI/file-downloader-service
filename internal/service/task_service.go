package service

import (
	"file-downloader-service/internal/entity"
	"file-downloader-service/internal/usecase"
)

type TaskService interface {
	CreateTask(urls []string) (*entity.Task, error)
	GetTaskStatus(id string) (*entity.Task, error)
}

type taskService struct {
	taskUseCase usecase.TaskUseCase
}

func NewTaskService(taskUseCase usecase.TaskUseCase) TaskService {
	return &taskService{taskUseCase: taskUseCase}
}

func (s *taskService) CreateTask(urls []string) (*entity.Task, error) {
	return s.taskUseCase.CreateTask(urls)
}

func (s *taskService) GetTaskStatus(id string) (*entity.Task, error) {
	return s.taskUseCase.GetTaskStatus(id)
}
