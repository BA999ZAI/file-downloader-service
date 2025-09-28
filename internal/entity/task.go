package entity

import (
	"time"
)

type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusProcessing TaskStatus = "processing"
	StatusCompleted  TaskStatus = "completed"
	StatusFailed     TaskStatus = "failed"
)

type Task struct {
	ID        string       `json:"id"`
	URLs      []string     `json:"urls"`
	Status    TaskStatus   `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	Results   []FileResult `json:"results,omitempty"`
}

type FileResult struct {
	URL      string `json:"url"`
	FileName string `json:"file_name,omitempty"`
	Error    string `json:"error,omitempty"`
}

type CreateTaskRequest struct {
	URLs []string `json:"urls" binding:"required,min=1"`
}
