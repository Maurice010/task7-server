package dto

type CartItemDTO struct {
	ProductID uint `json:"productId"`
	Quantity  int  `json:"quantity"`
}