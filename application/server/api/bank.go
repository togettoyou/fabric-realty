package api

import (
	"application/service"
	"application/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BankHandler struct {
	bankService *service.BankService
}

func NewBankHandler() *BankHandler {
	return &BankHandler{
		bankService: &service.BankService{},
	}
}

// CompleteTransaction 完成交易（仅银行组织可以调用）
func (h *BankHandler) CompleteTransaction(c *gin.Context) {
	txID := c.Param("txId")
	err := h.bankService.CompleteTransaction(txID)
	if err != nil {
		utils.ServerError(c, "完成交易失败："+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "交易完成", nil)
}

// QueryTransaction 查询交易信息
func (h *BankHandler) QueryTransaction(c *gin.Context) {
	txID := c.Param("txId")
	transaction, err := h.bankService.QueryTransaction(txID)
	if err != nil {
		utils.ServerError(c, "查询交易信息失败："+err.Error())
		return
	}

	utils.Success(c, transaction)
}

// QueryTransactionList 分页查询交易列表
func (h *BankHandler) QueryTransactionList(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	bookmark := c.DefaultQuery("bookmark", "")
	status := c.DefaultQuery("status", "")

	result, err := h.bankService.QueryTransactionList(int32(pageSize), bookmark, status)
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}

	utils.Success(c, result)
}

// QueryBlockList 分页查询区块列表
func (h *BankHandler) QueryBlockList(c *gin.Context) {
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))
	pageNum, _ := strconv.Atoi(c.DefaultQuery("pageNum", "1"))

	result, err := h.bankService.QueryBlockList(pageSize, pageNum)
	if err != nil {
		utils.ServerError(c, err.Error())
		return
	}

	utils.Success(c, result)
}
