package api

import (
	"application/service"
	"application/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TradingPlatformHandler struct {
	tradingService *service.TradingPlatformService
}

func NewTradingPlatformHandler() *TradingPlatformHandler {
	return &TradingPlatformHandler{
		tradingService: &service.TradingPlatformService{},
	}
}

// CreateTransaction 生成交易（仅交易平台组织可以调用）
func (h *TradingPlatformHandler) CreateTransaction(c *gin.Context) {
	var req struct {
		TxID         string  `json:"txId"`
		RealEstateID string  `json:"realEstateId"`
		Seller       string  `json:"seller"`
		Buyer        string  `json:"buyer"`
		Price        float64 `json:"price"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "交易信息格式错误")
		return
	}

	err := h.tradingService.CreateTransaction(req.TxID, req.RealEstateID, req.Seller, req.Buyer, req.Price)
	if err != nil {
		utils.ServerError(c, "生成交易失败："+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "交易创建成功", nil)
}

// QueryRealEstate 查询房产信息
func (h *TradingPlatformHandler) QueryRealEstate(c *gin.Context) {
	id := c.Param("id")
	realEstate, err := h.tradingService.QueryRealEstate(id)
	if err != nil {
		utils.ServerError(c, "查询房产信息失败："+err.Error())
		return
	}

	utils.Success(c, realEstate)
}

// QueryTransaction 查询交易信息
func (h *TradingPlatformHandler) QueryTransaction(c *gin.Context) {
	txID := c.Param("txId")
	transaction, err := h.tradingService.QueryTransaction(txID)
	if err != nil {
		utils.ServerError(c, "查询交易信息失败："+err.Error())
		return
	}

	utils.Success(c, transaction)
}

// QueryTransactionList 分页查询交易列表
func (h *TradingPlatformHandler) QueryTransactionList(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	bookmark := c.DefaultQuery("bookmark", "")
	status := c.DefaultQuery("status", "")

	result, err := h.tradingService.QueryTransactionList(int32(pageSize), bookmark, status)
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}

	utils.Success(c, result)
}

// QueryBlockList 分页查询区块列表
func (h *TradingPlatformHandler) QueryBlockList(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))

	result, err := h.tradingService.QueryBlockList(pageSize, pageNum)
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}

	utils.Success(c, result)
}
