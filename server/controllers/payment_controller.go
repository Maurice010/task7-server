package controllers

import (
	"net/http"
	"server/database"
	"server/dto"
	"server/models"

	"github.com/labstack/echo/v4"
)

func HandlePayment(c echo.Context) error {
	var items []dto.CartItemDTO
	if err := c.Bind(&items); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid input"})
	}

	if len(items) == 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "empty cart"})
	}

	var total float64
	for _, item := range items {
		var product models.Product
		if err := database.DB.First(&product, item.ProductID).Error; err != nil {
			continue
		}
		total += product.Price * float64(item.Quantity)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status": "ok",
		"total":  total,
	})
}