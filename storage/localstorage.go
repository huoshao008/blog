package storage

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

// LocalStorage 本地存储实现
type LocalStorage struct {
	RootDir string // 本地存储根目录（如 "./uploads"）
}

// NewLocalStorage 创建本地存储实例
func NewLocalStorage(rootDir string) (*LocalStorage, error) {
	// 确保根目录存在
	if err := os.MkdirAll(rootDir, 0755); err != nil {
		return nil, err
	}
	return &LocalStorage{RootDir: rootDir}, nil
}

// Upload 上传文件到本地
func (s *LocalStorage) Upload(ctx context.Context, key string, reader io.Reader) (string, error) {
	filePath := filepath.Join(s.RootDir, key)
	// 创建父目录
	if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
		return "", err
	}
	// 写入文件
	file, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	if _, err := io.Copy(file, reader); err != nil {
		return "", err
	}
	return filePath, nil // 返回本地文件路径
}

// Download 从本地下载文件
func (s *LocalStorage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	filePath := filepath.Join(s.RootDir, key)
	return os.Open(filePath)
}

// Delete 从本地删除文件
func (s *LocalStorage) Delete(ctx context.Context, key string) (err error) {
	filePath := filepath.Join(s.RootDir, key)
	return os.Remove(filePath)
}

// Exists 检查本地文件是否存在
func (s *LocalStorage) Exists(ctx context.Context, key string) (bool, error) {
	filePath := filepath.Join(s.RootDir, key)
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false, nil
	}
	return err == nil, err
}
