package v1

import (
	bc "application/blockchain"
	"application/pkg/app"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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
