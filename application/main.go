package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/togettoyou/blockchain-real-estate/application/blockchain"
	_ "github.com/togettoyou/blockchain-real-estate/application/docs"
	"github.com/togettoyou/blockchain-real-estate/application/pkg/setting"
	"github.com/togettoyou/blockchain-real-estate/application/routers"
	"github.com/togettoyou/blockchain-real-estate/application/service"
	"log"
	"net/http"
	"time"
)

func init() {
	setting.Setup()
}

// @title 基于区块链技术的房地产交易系统api文档
// @version 1.0
// @description 基于区块链技术的房地产交易系统api文档
// @contact.name togettoyou
// @contact.email zoujh99@qq.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	timeLocal, err := time.LoadLocation("Asia/Chongqing")
	if err != nil {
		log.Printf("时区设置失败 %s", err)
	}
	time.Local = timeLocal
	blockchain.Init()
	go service.Init()
	gin.SetMode(setting.ServerSetting.RunMode)
	routersInit := routers.InitRouter()
	readTimeout := setting.ServerSetting.ReadTimeout
	writeTimeout := setting.ServerSetting.WriteTimeout
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        routersInit,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Printf("[info] start http server listening %s", endPoint)

	if err := server.ListenAndServe(); err != nil {
		log.Printf("start http server failed %s", err)
	}
}
