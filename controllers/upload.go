package controllers

import (
	"blog/config"
	"blog/storage"
	"blog/utils"
	"context"
	"errors"
	"io"
	"log"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UploadController struct{}

func (up *UploadController) UploadFile(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}
	var filename string
	// 保存文件
	if config.LocalStorage != nil {
		filename, err = upload(config.LocalStorage, file)
		if err != nil {
			utils.BadRequest(c, err.Error())
			return
		}
	}
	if config.OSSStorage != nil {
		filename, err = upload(config.OSSStorage, file)
		if err != nil {
			utils.BadRequest(c, err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "File uploaded successfully",
		"filename": filename,
		"size":     file.Size,
	})
}

func (up *UploadController) DownloadFile(c *gin.Context) {
	filename := c.Param("filename")
	// 设置响应头
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", "application/octet-stream")

	if config.LocalStorage != nil {
		file, err := download(config.LocalStorage, filename)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}
		// 将 io.ReadCloser 内容流式传输到响应
		_, errio := io.Copy(c.Writer, file)
		if errio != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "File download error"})
			return
		}
	}
	if config.OSSStorage != nil {
		file, err := download(config.OSSStorage, filename)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}
		// 将 io.ReadCloser 内容流式传输到响应
		_, errio := io.Copy(c.Writer, file)
		if errio != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "File download error"})
			return
		}
	}
}

func (up *UploadController) DeleteFile(c *gin.Context) {
	filename := c.Param("filename")
	if config.LocalStorage != nil {
		err := delete(config.LocalStorage, filename)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "File delete found"})
			return
		}
	}
	if config.OSSStorage != nil {
		err := delete(config.OSSStorage, filename)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "File delete found"})
			return
		}
	}
	utils.Success(c, gin.H{
		"code":    200,
		"message": "success"})
}

// 统一上传函数：兼容所有Storage实现
func upload(store storage.Storage, file *multipart.FileHeader) (string, error) {
	// 打开本地文件
	rc, err := file.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer rc.Close()

	// 上传到存储后端
	ctx := context.Background()
	url, err := store.Upload(ctx, file.Filename, rc)
	if err != nil {
		return "", err
	}
	log.Printf("Upload success! URL: %s", url)

	// 检查文件是否存在
	exists, err := store.Exists(ctx, file.Filename)
	if err != nil {
		return "", err
	}
	log.Printf("File exists: %v", exists)
	return file.Filename, nil
}

// 下载函数：兼容所有Storage实现
func download(store storage.Storage, filename string) (io.ReadCloser, error) {
	ctx := context.Background()
	// 检查文件是否存在
	exists, err := store.Exists(ctx, filename)
	if err != nil {
		return nil, err
	}
	if exists {
		return store.Download(ctx, filename)
	}
	return nil, errors.New("File is not exists!")
}

// 删除函数：兼容所有Storage实现
func delete(store storage.Storage, filename string) error {
	ctx := context.Background()
	return store.Delete(ctx, filename)
}
