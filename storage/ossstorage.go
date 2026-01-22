package storage

import (
	"context"
	"io"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

// OSSStorage 阿里云OSS存储实现
type OSSStorage struct {
	bucket *oss.Bucket // OSS Bucket实例
}

// NewOSSStorage 创建OSS存储实例
func NewOSSStorage(endpoint, accessKeyID, accessKeySecret, bucketName string) (*OSSStorage, error) {
	// 创建OSS客户端
	client, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, err
	}
	// 获取Bucket
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return nil, err
	}
	return &OSSStorage{bucket: bucket}, nil
}

// Upload 上传文件到OSS
func (s *OSSStorage) Upload(ctx context.Context, key string, reader io.Reader) (string, error) {
	// 上传文件（支持流式上传）
	if err := s.bucket.PutObject(key, reader); err != nil {
		return "", err
	}
	// 返回OSS文件的公网访问URL（需确保Bucket已开启公开读权限）
	return s.bucket.SignURL(key, oss.HTTPGet, 3600) // 生成1小时有效期的临时URL
}

// Download 从OSS下载文件
func (s *OSSStorage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	return s.bucket.GetObject(key)
}

// Delete 从OSS删除文件
func (s *OSSStorage) Delete(ctx context.Context, key string) error {
	return s.bucket.DeleteObject(key)
}

// Exists 检查OSS文件是否存在
func (s *OSSStorage) Exists(ctx context.Context, key string) (bool, error) {
	return s.bucket.IsObjectExist(key)
}
