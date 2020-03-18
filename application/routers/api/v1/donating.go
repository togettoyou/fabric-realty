/**
 * @Author: 夜央 Oh oh oh oh oh oh (https://github.com/togettoyou)
 * @Email: zoujh99@qq.com
 * @Date: 2020/3/19 1:54 上午
 * @Description: 捐赠相关接口
 */
package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	bc "github.com/togettoyou/blockchain-real-estate/application/blockchain"
	"github.com/togettoyou/blockchain-real-estate/application/pkg/app"
	"net/http"
)

type DonatingRequestBody struct {
	ObjectOfDonating string `json:"objectOfDonating"` //捐赠对象
	Donor            string `json:"donor"`            //捐赠人
	Grantee          string `json:"grantee"`          //受赠人
}

type DonatingListQueryRequestBody struct {
	Donor string `json:"donor"`
}

type DonatingListQueryByGranteeRequestBody struct {
	Grantee string `json:"grantee"`
}

type UpdateDonatingRequestBody struct {
	ObjectOfDonating string `json:"objectOfDonating"` //捐赠对象
	Donor            string `json:"donor"`            //捐赠人
	Grantee          string `json:"grantee"`          //受赠人
	Status           string `json:"status"`           //需要更改的状态
}

// @Summary 发起捐赠
// @Param donating body DonatingRequestBody true "donating"
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/createDonating [post]
func CreateDonating(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(DonatingRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	if body.ObjectOfDonating == "" || body.Donor == "" || body.Grantee == "" {
		appG.Response(http.StatusBadRequest, "失败", "ObjectOfDonating捐赠对象和Donor捐赠人和Grantee受赠人不能为空")
		return
	}
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.ObjectOfDonating))
	bodyBytes = append(bodyBytes, []byte(body.Donor))
	bodyBytes = append(bodyBytes, []byte(body.Grantee))
	//调用智能合约
	resp, err := bc.ChannelExecute("createDonating", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	var data map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	appG.Response(http.StatusOK, "成功", data)
}

// @Summary 查询捐赠列表(可查询所有，也可根据发起捐赠人查询)
// @Param donatingListQuery body DonatingListQueryRequestBody true "donatingListQuery"
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/queryDonatingList [post]
func QueryDonatingList(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(DonatingListQueryRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	var bodyBytes [][]byte
	if body.Donor != "" {
		bodyBytes = append(bodyBytes, []byte(body.Donor))
	}
	//调用智能合约
	resp, err := bc.ChannelQuery("queryDonatingList", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	// 反序列化json
	var data []map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	appG.Response(http.StatusOK, "成功", data)
}

// @Summary 根据受赠人(受赠人AccountId)查询捐赠(受赠的)(供受赠人查询)
// @Param donatingListQueryByGrantee body DonatingListQueryByGranteeRequestBody true "donatingListQueryByGrantee"
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/queryDonatingListByGrantee [post]
func QueryDonatingListByGrantee(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(DonatingListQueryByGranteeRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	if body.Grantee == "" {
		appG.Response(http.StatusBadRequest, "失败", "必须指定AccountId查询")
		return
	}
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.Grantee))
	//调用智能合约
	resp, err := bc.ChannelQuery("queryDonatingListByGrantee", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	// 反序列化json
	var data []map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	appG.Response(http.StatusOK, "成功", data)
}

// @Summary 更新捐赠状态（确认受赠、取消）
// @Param updateDonating body UpdateDonatingRequestBody true "updateDonating"
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/updateDonating [post]
func UpdateDonating(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(UpdateDonatingRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	if body.ObjectOfDonating == "" || body.Donor == "" || body.Grantee == "" || body.Status == "" {
		appG.Response(http.StatusBadRequest, "失败", "参数不能为空")
		return
	}
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.ObjectOfDonating))
	bodyBytes = append(bodyBytes, []byte(body.Donor))
	bodyBytes = append(bodyBytes, []byte(body.Grantee))
	bodyBytes = append(bodyBytes, []byte(body.Status))
	//调用智能合约
	resp, err := bc.ChannelExecute("updateDonating", bodyBytes)
	if err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	var data map[string]interface{}
	if err = json.Unmarshal(bytes.NewBuffer(resp.Payload).Bytes(), &data); err != nil {
		appG.Response(http.StatusInternalServerError, "失败", err.Error())
		return
	}
	appG.Response(http.StatusOK, "成功", data)
}
