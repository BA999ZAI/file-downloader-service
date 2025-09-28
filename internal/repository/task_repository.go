package repository

import (
	"file-downloader-service/internal/entity"
	"sync"
	"time"
)

type TaskRepository interface {
	Create(task *entity.Task) error
	FindByID(id string) (*entity.Task, error)
	Update(task *entity.Task) error
	FindAll() ([]*entity.Task, error)
	FindByStatus(status entity.TaskStatus) ([]*entity.Task, error)
}

type inMemoryTaskRepository struct {
	tasks map[string]*entity.Task
	mutex sync.RWMutex
}

func NewInMemoryTaskRepository() TaskRepository {
	return &inMemoryTaskRepository{
		tasks: make(map[string]*entity.Task),
	}
}

func (r *inMemoryTaskRepository) Create(task *entity.Task) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.tasks[task.ID] = task
	return nil
}

func (r *inMemoryTaskRepository) FindByID(id string) (*entity.Task, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	task, exists := r.tasks[id]
	if !exists {
		return nil, nil
	}
	return task, nil
}

func (r *inMemoryTaskRepository) Update(task *entity.Task) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	task.UpdatedAt = time.Now()
	r.tasks[task.ID] = task
	return nil
}

func (r *inMemoryTaskRepository) FindAll() ([]*entity.Task, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	tasks := make([]*entity.Task, 0, len(r.tasks))
	for _, task := range r.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

func (r *inMemoryTaskRepository) FindByStatus(status entity.TaskStatus) ([]*entity.Task, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	var tasks []*entity.Task
	for _, task := range r.tasks {
		if task.Status == status {
			tasks = append(tasks, task)
		}
	}
	return tasks, nil
}
