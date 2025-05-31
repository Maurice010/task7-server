package models

import "gorm.io/gorm"

type Cart struct {
	gorm.Model
	UserId    uint       `json:"user_id"`
	CartItems []CartItem `json:"cart_items" gorm:"foreignKey:CartId"`
}