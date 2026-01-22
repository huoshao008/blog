package storage

import (
	"context"
	"io"
)

// Storage 通用存储接口
type Storage interface {
	// Upload 上传文件：返回文件访问路径
	Upload(ctx context.Context, key string, reader io.Reader) (string, error)
	// Download 下载文件：返回文件读取流
	Download(ctx context.Context, key string) (io.ReadCloser, error)
	// Delete 删除文件
	Delete(ctx context.Context, key string) error
	// Exists 检查文件是否存在
	Exists(ctx context.Context, key string) (bool, error)
}
