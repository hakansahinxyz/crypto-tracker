package handlers

/*
import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/models"
	"github.com/hakansahinxyz/crypto-tracker-backend/internal/services"
)

func CreateCoin(c *gin.Context) {
	var coin models.Coin
	if err := c.ShouldBindJSON(&coin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.CreateCoin(&coin); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create coin"})
		return
	}

	c.JSON(http.StatusOK, coin)
}

func GetAllCoins(c *gin.Context) {
	coins, err := services.GetAllCoins()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch coins"})
		return
	}

	c.JSON(http.StatusOK, coins)
}

func GetCoinByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid coin ID"})
		return
	}

	coin, err := services.GetCoinByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch coin"})
		return
	}

	c.JSON(http.StatusOK, coin)
}

func UpdateCoin(c *gin.Context) {
	var coin models.Coin
	if err := c.ShouldBindJSON(&coin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := services.UpdateCoin(&coin); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update coin"})
		return
	}

	c.JSON(http.StatusOK, coin)
}

func DeleteCoin(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid coin ID"})
		return
	}

	if err := services.DeleteCoin(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete coin"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Coin deleted successfully"})
}
*/
