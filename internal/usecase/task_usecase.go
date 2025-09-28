package usecase

import (
	"file-downloader-service/internal/entity"
	"file-downloader-service/internal/repository"
	"file-downloader-service/pkg/downloader"
	"file-downloader-service/pkg/storage"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type TaskUseCase interface {
	CreateTask(urls []string) (*entity.Task, error)
	GetTaskStatus(id string) (*entity.Task, error)
	ProcessPendingTasks()
	StopProcessing()
	RecoverPendingTasks()
}

type taskUseCase struct {
	taskRepo     repository.TaskRepository
	downloader   downloader.Downloader
	fileStorage  storage.FileStorage
	stopChan     chan struct{}
	wg           sync.WaitGroup
	isProcessing bool
	mutex        sync.RWMutex
}

func NewTaskUseCase(taskRepo repository.TaskRepository, downloader downloader.Downloader, fileStorage storage.FileStorage) TaskUseCase {
	return &taskUseCase{
		taskRepo:    taskRepo,
		downloader:  downloader,
		fileStorage: fileStorage,
		stopChan:    make(chan struct{}),
	}
}

func (uc *taskUseCase) CreateTask(urls []string) (*entity.Task, error) {
	task := &entity.Task{
		ID:        uuid.New().String(),
		URLs:      urls,
		Status:    entity.StatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Results:   make([]entity.FileResult, 0, len(urls)),
	}

	if err := uc.taskRepo.Create(task); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return task, nil
}

func (uc *taskUseCase) GetTaskStatus(id string) (*entity.Task, error) {
	task, err := uc.taskRepo.FindByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}
	if task == nil {
		return nil, fmt.Errorf("task not found")
	}
	return task, nil
}

func (uc *taskUseCase) ProcessPendingTasks() {
	uc.mutex.Lock()
	uc.isProcessing = true
	uc.mutex.Unlock()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-uc.stopChan:
			uc.mutex.Lock()
			uc.isProcessing = false
			uc.mutex.Unlock()
			return
		case <-ticker.C:
			uc.processBatch()
		}
	}
}

func (uc *taskUseCase) processBatch() {
	pendingTasks, err := uc.taskRepo.FindByStatus(entity.StatusPending)
	if err != nil || len(pendingTasks) == 0 {
		return
	}

	for _, task := range pendingTasks {
		select {
		case <-uc.stopChan:
			return
		default:
			uc.processTask(task)
		}
	}
}

func (uc *taskUseCase) processTask(task *entity.Task) {
	task.Status = entity.StatusProcessing
	uc.taskRepo.Update(task)

	results := make([]entity.FileResult, 0, len(task.URLs))

	for _, url := range task.URLs {
		select {
		case <-uc.stopChan:
			task.Status = entity.StatusPending
			uc.taskRepo.Update(task)
			return
		default:
			fileName, err := uc.downloader.DownloadFile(url, uc.fileStorage.GetDownloadDir())
			result := entity.FileResult{URL: url}
			if err != nil {
				result.Error = err.Error()
			} else {
				result.FileName = fileName
			}
			results = append(results, result)
		}
	}

	task.Results = results
	if uc.allDownloadsSuccessful(results) {
		task.Status = entity.StatusCompleted
	} else {
		task.Status = entity.StatusFailed
	}

	uc.taskRepo.Update(task)
}

func (uc *taskUseCase) allDownloadsSuccessful(results []entity.FileResult) bool {
	for _, result := range results {
		if result.Error != "" {
			return false
		}
	}
	return true
}

func (uc *taskUseCase) StopProcessing() {
	close(uc.stopChan)
	uc.wg.Wait()
}

func (uc *taskUseCase) RecoverPendingTasks() {
	processingTasks, err := uc.taskRepo.FindByStatus(entity.StatusProcessing)
	if err != nil {
		return
	}

	for _, task := range processingTasks {
		task.Status = entity.StatusPending
		uc.taskRepo.Update(task)
	}
}

func (uc *taskUseCase) IsProcessing() bool {
	uc.mutex.RLock()
	defer uc.mutex.RUnlock()
	return uc.isProcessing
}
