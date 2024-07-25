package model

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserCreatedResponse struct {
	UserID int `json:"user_id"`
}

type TokenCreatedResponse struct {
	Token string `json:"token"`
}
