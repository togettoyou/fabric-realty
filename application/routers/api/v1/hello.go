/**
 * @Author: 夜央 Oh oh oh oh oh oh (https://github.com/togettoyou)
 * @Email: zoujh99@qq.com
 * @Date: 2020/3/3 11:24 下午
 * @Description: 测试api
 */
package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/togettoyou/blockchain-real-estate/application/pkg/app"
	"net/http"
)

// @Summary 测试输出Hello
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/hello [get]
func Hello(c *gin.Context) {
	appG := app.Gin{C: c}
	appG.Response(http.StatusOK, "成功", map[string]interface{}{
		"msg": "Hello",
	})
}
