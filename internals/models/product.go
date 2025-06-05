package models

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID			uuid.UUID	`gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name		string		`json:"name"`
	Price		float64		`json:"price"`

	UserID		uuid.UUID	`gorm:"type:uuid" json:"userId"`
	User		User		`gorm:"foreignKey:UserID" json:"user"`

	CreatedAt	time.Time	`json:"createdAt"`
	UpdatedAt	time.Time	`json:"updatedAt"`
}