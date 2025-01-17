package api

import (
	"application/service"
	"net/http"

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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.realtyService.CreateRealEstate(req.ID, req.Address, req.Area, req.Owner, req.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Real estate created successfully"})
}

// QueryRealEstate 查询房产信息
func (h *RealtyHandler) QueryRealEstate(c *gin.Context) {
	id := c.Param("id")
	realEstate, err := h.realtyService.QueryRealEstate(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, realEstate)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.realtyService.CreateTransaction(req.TxID, req.RealEstateID, req.Seller, req.Buyer, req.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction created successfully"})
}

// ConfirmEscrow 确认资金托管
func (h *RealtyHandler) ConfirmEscrow(c *gin.Context) {
	txID := c.Param("txId")
	err := h.realtyService.ConfirmEscrow(txID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Escrow confirmed successfully"})
}

// CompleteTransaction 完成交易
func (h *RealtyHandler) CompleteTransaction(c *gin.Context) {
	txID := c.Param("txId")
	err := h.realtyService.CompleteTransaction(txID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Transaction completed successfully"})
}
