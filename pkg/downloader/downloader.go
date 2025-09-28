package downloader

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Downloader interface {
	DownloadFile(url string, downloadDir string) (string, error)
}

type HTTPDownloader struct {
	client *http.Client
}

func NewHTTPDownloader() *HTTPDownloader {
	return &HTTPDownloader{
		client: &http.Client{
			Timeout: 30 * time.Minute,
		},
	}
}

func (d *HTTPDownloader) DownloadFile(url string, downloadDir string) (string, error) {
	resp, err := d.client.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s", resp.Status)
	}

	fileName := getFileNameFromURL(url)
	if fileName == "" {
		fileName = fmt.Sprintf("download_%d", time.Now().UnixNano())
	}

	filePath := filepath.Join(downloadDir, fileName)

	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		os.Remove(filePath)
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return fileName, nil
}

func getFileNameFromURL(url string) string {
	parts := strings.Split(url, "/")
	fileName := parts[len(parts)-1]

	if fileName == "" || strings.Contains(fileName, "?") {
		return ""
	}

	return fileName
}
