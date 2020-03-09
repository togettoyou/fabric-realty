/**
 * @Author: 夜央 Oh oh oh oh oh oh (https://github.com/togettoyou)
 * @Email: zoujh99@qq.com
 * @Date: 2020/3/5 12:55 上午
 * @Description: 账户信息相关接口
 */
package v1

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	bc "github.com/togettoyou/blockchain-real-estate/application/blockchain"
	"github.com/togettoyou/blockchain-real-estate/application/pkg/app"
	"net/http"
)

// @Summary 获取账户信息
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/queryAccountList [get]
func QueryAccountList(c *gin.Context) {
	appG := app.Gin{C: c}
	resp, err := bc.ChannelQuery("queryAccountList", [][]byte{})
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	// 反序列化json
	var data []map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	appG.Response(http.StatusOK, "成功", data)
}
