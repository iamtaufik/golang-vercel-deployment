package dto

import (
	"time"

	"github.com/google/uuid"
)

type ProductRequest struct {
	Name   string  		`json:"name"`
	Price  float64 		`json:"price"`
	UserID uuid.UUID	`json:"userId"`
}
type ProductResponse struct {
	ID			uuid.UUID		`json:"id"`
	Name   		string  		`json:"name"`
	Price  		float64 		`json:"price"`
	CreatedAt	time.Time		`json:"createdAt"`
}

