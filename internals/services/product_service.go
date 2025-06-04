package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/iamtaufik/golang-vercel-deployment/internals/models"
	"github.com/iamtaufik/golang-vercel-deployment/internals/repository"
)

type ProductService interface {
	GetProducts(ctx context.Context)([]models.Product, error)
	GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error)
	CreateProduct(ctx context.Context, product *models.Product) error
}

type productService struct {
	Repository repository.ProductRepository
}

func NewProductService(repository repository.ProductRepository) *productService {
	return &productService{Repository: repository}
}

func (s *productService) CreateProduct(ctx context.Context, product *models.Product) error {
	return s.Repository.Create(ctx, product)
}

func (s *productService) GetProducts(ctx context.Context) ([]models.Product, error) {
	return s.Repository.FindAll(ctx)
}

func (s *productService) GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error){
	return s.Repository.FindByID(ctx, id)
}