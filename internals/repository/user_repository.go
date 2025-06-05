package repository

import (
	"context"

	"github.com/iamtaufik/golang-vercel-deployment/internals/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	FindByEmail(context context.Context, email string) (*models.User, error)
	Create(context context.Context, user *models.User) error
}

type userRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{DB: db}
}

func (r *userRepository) FindByEmail(context context.Context, email string) (*models.User, error){
	var user models.User

	if err := r.DB.WithContext(context).First(&user, "email = ?", email).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) Create(context context.Context, user *models.User) error {
	return r.DB.WithContext(context).Create(user).Error
}