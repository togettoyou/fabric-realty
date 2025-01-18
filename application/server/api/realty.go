package api

import (
	"application/service"
	"application/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RealtyHandler struct {
	realtyService *service.RealtyService
}

func NewRealtyHandler() *RealtyHandler {
	return &RealtyHandler{
		realtyService: &service.RealtyService{},
	}
}

// CreateRealEstate 创建房产信息
func (h *RealtyHandler) CreateRealEstate(c *gin.Context) {
	var req struct {
		ID      string  `json:"id"`
		Address string  `json:"address"`
		Area    float64 `json:"area"`
		Owner   string  `json:"owner"`
		Price   float64 `json:"price"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "房产信息格式错误")
		return
	}

	err := h.realtyService.CreateRealEstate(req.ID, req.Address, req.Area, req.Owner, req.Price)
	if err != nil {
		utils.ServerError(c, "创建房产信息失败："+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "房产信息创建成功", nil)
}

// QueryRealEstate 查询房产信息
func (h *RealtyHandler) QueryRealEstate(c *gin.Context) {
	id := c.Param("id")
	realEstate, err := h.realtyService.QueryRealEstate(id)
	if err != nil {
		utils.ServerError(c, "查询房产信息失败："+err.Error())
		return
	}

	utils.Success(c, realEstate)
}

// CreateTransaction 创建交易
func (h *RealtyHandler) CreateTransaction(c *gin.Context) {
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

	err := h.realtyService.CreateTransaction(req.TxID, req.RealEstateID, req.Seller, req.Buyer, req.Price)
	if err != nil {
		utils.ServerError(c, "创建交易失败："+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "交易创建成功", nil)
}

// ConfirmEscrow 确认资金托管
func (h *RealtyHandler) ConfirmEscrow(c *gin.Context) {
	txID := c.Param("txId")
	err := h.realtyService.ConfirmEscrow(txID)
	if err != nil {
		utils.ServerError(c, "确认资金托管失败："+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "资金托管确认成功", nil)
}

// CompleteTransaction 完成交易
func (h *RealtyHandler) CompleteTransaction(c *gin.Context) {
	txID := c.Param("txId")
	err := h.realtyService.CompleteTransaction(txID)
	if err != nil {
		utils.ServerError(c, "完成交易失败："+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "交易完成", nil)
}

// QueryTransaction 查询交易信息
func (h *RealtyHandler) QueryTransaction(c *gin.Context) {
	txID := c.Param("txId")
	transaction, err := h.realtyService.QueryTransaction(txID)
	if err != nil {
		utils.ServerError(c, "查询交易信息失败："+err.Error())
		return
	}

	utils.Success(c, transaction)
}

// QueryRealEstateList 分页查询房产列表
func (h *RealtyHandler) QueryRealEstateList(c *gin.Context) {
	// 获取分页参数
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	bookmark := c.DefaultQuery("bookmark", "")

	result, err := h.realtyService.QueryRealEstateList(int32(pageSize), bookmark)
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}

	utils.Success(c, result)
}

// QueryTransactionList 分页查询交易列表
func (h *RealtyHandler) QueryTransactionList(c *gin.Context) {
	// 获取分页参数
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	bookmark := c.DefaultQuery("bookmark", "")

	result, err := h.realtyService.QueryTransactionList(int32(pageSize), bookmark)
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}

	utils.Success(c, result)
}

// QueryRealEstateByFilter 按条件查询房产列表
func (h *RealtyHandler) QueryRealEstateByFilter(c *gin.Context) {
	// 获取查询参数
	owner := c.Query("owner")
	status := c.Query("status")
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	bookmark := c.DefaultQuery("bookmark", "")

	result, err := h.realtyService.QueryRealEstateByFilter(owner, status, int32(pageSize), bookmark)
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}

	utils.Success(c, result)
}

// QueryTransactionByFilter 按条件查询交易列表
func (h *RealtyHandler) QueryTransactionByFilter(c *gin.Context) {
	// 获取查询参数
	seller := c.Query("seller")
	buyer := c.Query("buyer")
	status := c.Query("status")
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	bookmark := c.DefaultQuery("bookmark", "")

	result, err := h.realtyService.QueryTransactionByFilter(seller, buyer, status, int32(pageSize), bookmark)
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}

	utils.Success(c, result)
}
