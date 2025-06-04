package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/iamtaufik/golang-vercel-deployment/internals/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	FindAll(context context.Context) ([]models.Product, error)
	FindByID(context context.Context, id uuid.UUID) (*models.Product, error)
	Create(context context.Context, product *models.Product) error
}

type productRepository struct {
	DB *gorm.DB
}

func NewProductRepository(db *gorm.DB) *productRepository {
	return &productRepository{DB: db}
}

func (r *productRepository) FindAll(ctx context.Context) ([]models.Product, error) {
	var products []models.Product
	err := r.DB.WithContext(ctx).Find(&products).Error
	return products, err
}

func (r *productRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Product, error)  {
	var product models.Product

	if err := r.DB.WithContext(ctx).Find(&product, id).Error; err != nil {
		return nil, err
	}

	return &product, nil
}

func (r *productRepository) Create(ctx context.Context, product *models.Product) error {
	return r.DB.WithContext(ctx).Create(product).Error
}