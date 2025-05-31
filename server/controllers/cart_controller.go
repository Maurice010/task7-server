package controllers

import (
	"net/http"
	"server/database"
	"server/dto"
	"server/models"

	"github.com/labstack/echo/v4"
)

func SaveCart(c echo.Context) error {
	var items []dto.CartItemDTO
	if err := c.Bind(&items); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid input")
	}

	if len(items) == 0 {
		return c.JSON(http.StatusBadRequest, "Empty cart")
	}

	userId := uint(1)

	var cart models.Cart
	if err := database.DB.Where("user_id = ?", userId).First(&cart).Error; err != nil {
		cart = models.Cart{UserId: userId}
		if err := database.DB.Create(&cart).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, "Failed to create cart")
		}
	}

	database.DB.Unscoped().Where("cart_id = ?", cart.ID).Delete(&models.CartItem{})

	for _, item := range items {
		cartItem := models.CartItem{
			CartId:    cart.ID,
			ProductId: item.ProductID,
			Quantity:  item.Quantity,
		}
		database.DB.Create(&cartItem)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":  "saved",
		"cart_id": cart.ID,
	})
}
