package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID 			uuid.UUID 	`gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name		string		`json:"name"`
	Email		string		`gorm:"unique" json:"email"`
	Password	string		`json:"password"`

	Products	[]Product	`gorm:"foreignKey:UserID" json:"reviews"`

	CreatedAt	time.Time	`json:"createdAt"`
	UpdatedAt	time.Time	`json:"updatedAt"`
}