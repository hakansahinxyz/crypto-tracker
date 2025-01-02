package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	request "github.com/hakansahinxyz/crypto-tracker-backend/internal/dto/request"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/mappers"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/services"
)

type WalletBalanceHandler struct {
	service  *services.WalletBalanceService
	validate *validator.Validate
}

func NewWalletBalanceHandler(service *services.WalletBalanceService) *WalletBalanceHandler {
	return &WalletBalanceHandler{
		service:  service,
		validate: validator.New(),
	}
}

func (h *WalletBalanceHandler) GetAllBalances(c *gin.Context) {
	var req request.WalletBalanceRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters"})
		return
	}

	if err := h.validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
