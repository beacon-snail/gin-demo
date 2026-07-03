package main

import (
	"fmt"
	"log"

	"gin-mysql-demo/config"
	"gin-mysql-demo/database"
	"gin-mysql-demo/migrations"
	"gin-mysql-demo/routes"
)

func main() {
	// 1. 加载配置
	config.LoadConfig()
	fmt.Println("✅ Config loaded successfully!")

	// 2. 初始化数据库
	database.InitDB()
	defer database.CloseDB()

	// 3. 数据库迁移
	migrations.AutoMigrate()

	// 4. 设置路由
	r := routes.SetupRouter()

	// 5. 启动服务
	port := config.AppConfig.AppPort
	fmt.Printf("🚀 Server is running on port %s\n", port)
	fmt.Printf("📝 API Documentation: http://localhost:%s/ping\n", port)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
