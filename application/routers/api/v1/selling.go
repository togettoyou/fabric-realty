/**
 * @Author: 夜央 Oh oh oh oh oh oh (https://github.com/togettoyou)
 * @Email: zoujh99@qq.com
 * @Date: 2020/3/12 12:09 下午
 * @Description: 销售相关接口
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
	"strconv"
)

type SellingRequestBody struct {
	ObjectOfSale string  `json:"objectOfSale"` //销售对象(正在出售的房地产RealEstateID)
	Seller       string  `json:"seller"`       //发起销售人、卖家(卖家AccountId)
	Price        float64 `json:"price"`        //价格
	SalePeriod   int     `json:"salePeriod"`   //智能合约的有效期(单位为天)
}

type SellingByBuyRequestBody struct {
	ObjectOfSale string `json:"objectOfSale"` //销售对象(正在出售的房地产RealEstateID)
	Seller       string `json:"seller"`       //发起销售人、卖家(卖家AccountId)
	Buyer        string `json:"buyer"`        //买家(买家AccountId)
}
type SellingListQueryRequestBody struct {
	Seller string `json:"seller"` //发起销售人、卖家(卖家AccountId)
}
type SellingListQueryByBuyRequestBody struct {
	Buyer string `json:"buyer"` //买家(买家AccountId)
}
type UpdateSellingRequestBody struct {
	ObjectOfSale string `json:"objectOfSale"` //销售对象(正在出售的房地产RealEstateID)
	Seller       string `json:"seller"`       //发起销售人、卖家(卖家AccountId)
	Buyer        string `json:"buyer"`        //买家(买家AccountId)
	Status       string `json:"status"`       //需要更改的状态
}

// @Summary 发起销售
// @Param selling body SellingRequestBody true "selling"
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/createSelling [post]
func CreateSelling(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(SellingRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	if body.ObjectOfSale == "" || body.Seller == "" {
		appG.Response(http.StatusBadRequest, "失败", "ObjectOfSale销售对象和Seller发起销售人不能为空")
		return
	}
	if body.Price <= 0 || body.SalePeriod <= 0 {
		appG.Response(http.StatusBadRequest, "失败", "Price价格和SalePeriod智能合约的有效期(单位为天)必须大于0")
		return
	}
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.ObjectOfSale))
	bodyBytes = append(bodyBytes, []byte(body.Seller))
	bodyBytes = append(bodyBytes, []byte(strconv.FormatFloat(body.Price, 'E', -1, 64)))
	bodyBytes = append(bodyBytes, []byte(strconv.Itoa(body.SalePeriod)))
	//调用智能合约
	resp, err := bc.ChannelExecute("createSelling", bodyBytes)
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

// @Summary 买家购买
// @Param sellingByBuy body SellingByBuyRequestBody true "sellingByBuy"
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/createSellingByBuy [post]
func CreateSellingByBuy(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(SellingByBuyRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	if body.ObjectOfSale == "" || body.Seller == "" || body.Buyer == "" {
		appG.Response(http.StatusBadRequest, "失败", "参数不能为空")
		return
	}
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.ObjectOfSale))
	bodyBytes = append(bodyBytes, []byte(body.Seller))
	bodyBytes = append(bodyBytes, []byte(body.Buyer))
	//调用智能合约
	resp, err := bc.ChannelExecute("createSellingByBuy", bodyBytes)
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

// @Summary 查询销售(可查询所有，也可根据发起销售人查询)(发起的)
// @Param sellingListQuery body SellingListQueryRequestBody true "sellingListQuery"
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/querySellingList [post]
func QuerySellingList(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(SellingListQueryRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	var bodyBytes [][]byte
	if body.Seller != "" {
		bodyBytes = append(bodyBytes, []byte(body.Seller))
	}
	//调用智能合约
	resp, err := bc.ChannelQuery("querySellingList", bodyBytes)
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

// @Summary 根据参与销售人、买家(买家AccountId)查询销售(参与的)
// @Param sellingListQueryByBuy body SellingListQueryByBuyRequestBody true "sellingListQueryByBuy"
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/querySellingListByBuyer [post]
func QuerySellingListByBuyer(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(SellingListQueryByBuyRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	if body.Buyer == "" {
		appG.Response(http.StatusBadRequest, "失败", "必须指定买家AccountId查询")
		return
	}
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.Buyer))
	//调用智能合约
	resp, err := bc.ChannelQuery("querySellingListByBuyer", bodyBytes)
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

// @Summary 更新销售状态（买家确认、买卖家取消）Status取值为 完成"done"、取消"cancelled" 当处于销售中状态，卖家要取消时，buyer为""空
// @Param updateSelling body UpdateSellingRequestBody true "updateSelling"
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/updateSelling [post]
func UpdateSelling(c *gin.Context) {
	appG := app.Gin{C: c}
	body := new(UpdateSellingRequestBody)
	//解析Body参数
	if err := c.ShouldBind(body); err != nil {
		appG.Response(http.StatusBadRequest, "失败", fmt.Sprintf("参数出错%s", err.Error()))
		return
	}
	if body.ObjectOfSale == "" || body.Seller == "" || body.Status == "" {
		appG.Response(http.StatusBadRequest, "失败", "参数不能为空")
		return
	}
	var bodyBytes [][]byte
	bodyBytes = append(bodyBytes, []byte(body.ObjectOfSale))
	bodyBytes = append(bodyBytes, []byte(body.Seller))
	bodyBytes = append(bodyBytes, []byte(body.Buyer))
	bodyBytes = append(bodyBytes, []byte(body.Status))
	//调用智能合约
	resp, err := bc.ChannelExecute("updateSelling", bodyBytes)
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
