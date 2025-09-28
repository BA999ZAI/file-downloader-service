package storage

import (
	"os"
)

type FileStorage interface {
	EnsureDownloadDir() error
	GetDownloadDir() string
}

type LocalFileStorage struct {
	downloadDir string
}

func NewLocalFileStorage(downloadDir string) *LocalFileStorage {
	return &LocalFileStorage{downloadDir: downloadDir}
}

func (s *LocalFileStorage) EnsureDownloadDir() error {
	return os.MkdirAll(s.downloadDir, 0755)
}

func (s *LocalFileStorage) GetDownloadDir() string {
	return s.downloadDir
}
