package dto

import (
	"time"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
}

type RegisterRequest struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	ConfPassword string `json:"confPassword"`
}

type RegisterResponse struct {
	Name         string 	`json:"name"`
	Email        string 	`json:"email"`
	CreatedAt	 time.Time	`json:"createdAt"`
	UpdatedAt	 time.Time	`json:"updatedAt"`
}