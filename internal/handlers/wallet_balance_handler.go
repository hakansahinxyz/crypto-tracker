package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	request "github.com/hakansahinxyz/crypto-tracker-backend/internal/dto/request"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/mappers"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/services"
)

type WalletBalanceHandler struct {
	service *services.WalletBalanceService
}

func NewWalletBalanceHandler(service *services.WalletBalanceService) *WalletBalanceHandler {
	return &WalletBalanceHandler{
		service: service,
	}
}

func (h *WalletBalanceHandler) GetAllBalances(c *gin.Context) {
	var req request.WalletBalanceRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	balances, err := h.service.GetAllBalances(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := mappers.ToWalletBalancesResponse(balances)

	c.JSON(http.StatusOK, response)
}
