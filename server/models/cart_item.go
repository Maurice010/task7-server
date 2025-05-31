package models

import "gorm.io/gorm"

type CartItem struct {
    gorm.Model
    CartId     uint     `json:"cart_id"`
    ProductId  uint     `json:"product_id"`
    Product    Product  `gorm:"foreignKey:ProductId"`
    Quantity   int      `json:"quantity"`
}
