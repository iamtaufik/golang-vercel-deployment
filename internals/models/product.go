package models

import (
	"github.com/google/uuid"
)

type Product struct {
	ID		uuid.UUID	`gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name	string		`json:"name"`
	Price	float64		`json:"price"`
}