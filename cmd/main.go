package main

import (
	"blog/config"
	"blog/routes"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	//加载环境变量
	if err := godotenv.Load("/home/lsp/work/blog/.env"); err != nil {
		logrus.Warn("No .env file found, using system environment variables")
	}

	//初始化日志
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.Info("Starting blog application...")

	//初始化数据库
	config.InitDatabase()

	//初始化文件
	config.InitUpload()

	//设置路由
	r := routes.SetupRoutes()

	// 生成256位（32字节）的密钥
	// key, err := utils.GenerateJWTKey()
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "生成密钥时出错: %v\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Println("生成的JWT密钥:")
	// fmt.Println(key)

	//启动服务器
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	port = ":" + port
	logrus.WithField("port", port).Info("Server starting")

	if err := r.Run(port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}

}
