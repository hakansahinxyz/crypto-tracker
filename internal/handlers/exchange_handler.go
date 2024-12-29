// controllers/coinController.go
package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/models"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/services"
)

func CreateExchange(c *gin.Context) {
	var exchange models.Exchange
	if err := c.ShouldBindJSON(&exchange); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.CreateExchange(&exchange); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create coin"})
		return
	}

	c.JSON(http.StatusOK, exchange)
}
