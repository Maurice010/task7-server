package dto

type RegisterDTO struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

type LoginDTO struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}
