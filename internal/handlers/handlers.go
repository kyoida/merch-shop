// internal/handlers/handlers.go
package handlers

import (
	"net/http"

	"merch-shop/internal/middleware"
	"merch-shop/internal/model"
	"merch-shop/internal/repository"
	"merch-shop/internal/service"

	"github.com/gin-gonic/gin"
)

// AuthHandler обрабатывает аутентификацию и регистрацию.
// Он вызывает сервисную функцию AuthenticateUser и генерирует JWT.
func AuthHandler(repo repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.AuthRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": "Invalid request"})
			return
		}

		user, err := service.AuthenticateUser(repo, req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
			return
		}

		token, err := middleware.GenerateToken(user.ID, user.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"errors": "Failed to generate token"})
			return
		}
		c.JSON(http.StatusOK, model.AuthResponse{Token: token})
	}
}

// InfoHandler возвращает информацию о монетах, инвентаре и истории транзакций.
// Сервисный слой формирует объект InfoResponse.
func InfoHandler(repo repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDVal, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"errors": "User not found in token"})
			return
		}
		userID := userIDVal.(int64)

		info, err := service.GetInfo(repo, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
			return
		}
		c.JSON(http.StatusOK, info)
	}
}

// SendCoinHandler обрабатывает перевод монет между сотрудниками.
// Он использует функцию TransferCoins из сервисного слоя.
func SendCoinHandler(repo repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Из контекста получаем имя пользователя (установлено JWT-мидлваром)
		senderUsername := c.GetString("username")
		var req model.SendCoinRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": "Invalid request payload"})
			return
		}

		err := service.TransferCoins(repo, senderUsername, req.ToUser, req.Amount)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Coins transferred successfully"})
	}
}

// BuyHandler обрабатывает покупку мерча.
// Вызывает сервисную функцию PurchaseItem, которая проверяет наличие товара, баланс пользователя и записывает покупку.
func BuyHandler(repo repository.Repository) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.GetString("username")
		item := c.Param("item")
		if item == "" {
			c.JSON(http.StatusBadRequest, gin.H{"errors": "Item parameter is required"})
			return
		}
		err := service.PurchaseItem(repo, username, item)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Purchase successful"})
	}
}
