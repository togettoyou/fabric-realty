package main

import (
	"application/api"
	"application/config"
	"application/middleware"
	"application/model"
	"application/utils"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	if err := config.InitConfig(); err != nil {
		log.Fatalf("Failed to initialize config: %v", err)
	}

	// 初始化数据库
	if err := model.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// 初始化 Fabric 客户端
	if err := utils.InitFabric(); err != nil {
		log.Fatalf("Failed to initialize Fabric client: %v", err)
	}

	// 创建 Gin 路由
	r := gin.Default()

	// 创建处理器实例
	userHandler := api.NewUserHandler()
	realtyHandler := api.NewRealtyHandler()

	// 用户相关路由
	userGroup := r.Group("/api/user")
	{
		userGroup.POST("/register", userHandler.Register)
		userGroup.POST("/login", userHandler.Login)
	}

	// 需要认证的路由
	authorized := r.Group("/api")
	authorized.Use(middleware.AuthMiddleware())
	{
		// 房产相关路由
		realtyGroup := authorized.Group("/realty")
		{
			// 房管局管理员路由
			realtyAdmin := realtyGroup.Group("/admin")
			realtyAdmin.Use(middleware.RequireRoles(model.RealtyAdmin))
			{
				realtyAdmin.POST("/create", realtyHandler.CreateRealEstate)
			}

			// 银行管理员路由
			bankAdmin := realtyGroup.Group("/bank")
			bankAdmin.Use(middleware.RequireRoles(model.BankAdmin))
			{
				bankAdmin.POST("/escrow/:txId", realtyHandler.ConfirmEscrow)
				bankAdmin.POST("/complete/:txId", realtyHandler.CompleteTransaction)
			}

			// 普通用户路由
			realtyGroup.GET("/:id", realtyHandler.QueryRealEstate)
			realtyGroup.POST("/transaction", realtyHandler.CreateTransaction)
		}
	}

	// 启动服务器
	addr := fmt.Sprintf(":%d", config.GlobalConfig.Server.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
