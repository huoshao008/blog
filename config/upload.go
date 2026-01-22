package config

import (
	"blog/storage"
	"log"
	"os"
)

var LocalStorage *storage.LocalStorage
var OSSStorage *storage.OSSStorage

// 初始化文件上传
func InitUpload() {
	var err error
	//从变量获取文件参数
	STORAGE := getEnv("STORAGE", "local")
	FAILEPATH := getEnv("FAILEPATH", "./uploads")
	OSS_URL := getEnv("OSS_URL", "oss-cn-hangzhou.aliyuncs.com")
	OSS_ACCESS_KEY := getEnv("OSS_ACCESS_KEY", "")
	OSS_KEY_SECRET := getEnv("OSS_KEY_SECRET", "")
	OSS_BUCKET_NAME := getEnv("OSS_BUCKET_NAME", "")

	switch STORAGE {
	case "local":
		LocalStorage, err = storage.NewLocalStorage(FAILEPATH)
		if err != nil {
			log.Fatal(err)
		}
	case "oss":
		OSSStorage, err = storage.NewOSSStorage(OSS_URL, OSS_ACCESS_KEY, OSS_KEY_SECRET, OSS_BUCKET_NAME)
		if err != nil {
			log.Fatal(err)
		}
	default:
		LocalStorage, err = storage.NewLocalStorage("./uploads")
		if err != nil {
			log.Fatal(err)
		}
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getenv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
