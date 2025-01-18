package main

import (
	"application/api"
	"application/config"
	"application/utils"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化配置
	if err := config.InitConfig(); err != nil {
		log.Fatalf("初始化配置失败：%v", err)
	}

	// 初始化 Fabric 客户端
	if err := utils.InitFabric(); err != nil {
		log.Fatalf("初始化Fabric客户端失败：%v", err)
	}

	// 创建 Gin 路由
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	apiGroup := r.Group("/api")

	// 注册路由
	realtyHandler := api.NewRealtyHandler()

	// 不动产登记机构的接口
	realty := apiGroup.Group("/realty-agency")
	{
		// 创建房产信息
		realty.POST("/realty/create", realtyHandler.CreateRealEstate)
	}

	// 交易平台的接口
	trading := apiGroup.Group("/trading-platform")
	{
		// 创建交易
		trading.POST("/transaction/create", realtyHandler.CreateTransaction)
	}

	// 银行的接口
	bank := apiGroup.Group("/bank")
	{
		// 完成交易
		bank.POST("/transaction/complete/:txId", realtyHandler.CompleteTransaction)
	}

	// 公共查询接口（所有组织都可以访问）
	query := apiGroup.Group("/query")
	{
		// 房产相关查询
		realty := query.Group("/realty")
		{
			// 查询房产信息
			realty.GET("/:id", realtyHandler.QueryRealEstate)
			// 分页查询房产列表
			realty.GET("/list", realtyHandler.QueryRealEstateList)
		}

		// 交易相关查询
		transaction := query.Group("/transaction")
		{
			// 查询交易信息
			transaction.GET("/:txId", realtyHandler.QueryTransaction)
			// 分页查询交易列表
			transaction.GET("/list", realtyHandler.QueryTransactionList)
		}
	}

	// 启动服务器
	addr := fmt.Sprintf(":%d", config.GlobalConfig.Server.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("启动服务器失败：%v", err)
	}
}
