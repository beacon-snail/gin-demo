package routes

import (
	"gin-mysql-demo/controllers"
	"gin-mysql-demo/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// 跨域中间件
	r.Use(middleware.CORS())

	// 健康检查
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	api := r.Group("/api/v1")
	{
		// 认证路由（无需登录）
		auth := api.Group("/auth")
		{
			userController := controllers.NewUserController()
			auth.POST("/register", userController.Register)
			auth.POST("/login", userController.Login)
		}

		// 用户路由（需要登录）
		user := api.Group("/users")
		user.Use(middleware.AuthMiddleware())
		{
			userController := controllers.NewUserController()
			user.GET("", userController.GetUsers) // 支持智能过滤
			user.GET("/:id", userController.GetUser)
			user.PUT("/:id", userController.UpdateUser)
			user.DELETE("/:id", userController.DeleteUser)
		}
	}

	return r
}
