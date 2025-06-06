package services

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/iamtaufik/golang-vercel-deployment/internals/dto"
	"github.com/iamtaufik/golang-vercel-deployment/internals/models"
	"github.com/iamtaufik/golang-vercel-deployment/internals/repository"
)

type ProductService interface {
	GetProducts(ctx context.Context)([]dto.ProductResponse, error)
	GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error)
	CreateProduct(ctx context.Context, product *models.Product) error
}

type productService struct {
	Repository 		repository.ProductRepository
	UserRepository 	repository.UserRepository
}

func NewProductService(repository repository.ProductRepository, userRepository repository.UserRepository) *productService {
	return &productService{
		Repository: repository,
		UserRepository: userRepository,
	}
}

func (s *productService) CreateProduct(ctx context.Context, product *models.Product ) error {
	if _, err := uuid.Parse(product.UserID.String()); err != nil {
		return errors.New("invalid user id")
	}

	_, err := s.UserRepository.FindByID(ctx, product.UserID); 
	if err != nil {
		return errors.New("user not found")
	}

	return s.Repository.Create(ctx, product)
}

func (s *productService) GetProducts(ctx context.Context) ([]dto.ProductResponse, error) {
	var productResp []dto.ProductResponse

	products, err := s.Repository.FindAll(ctx)
	if err != nil {
		return productResp, err
	}

	for _, product := range products {
		productResp = append(productResp, dto.ProductResponse{
			ID: product.ID,
			Name: product.Name,
			Price: product.Price,
			CreatedAt: product.CreatedAt,
		})
	}

	return productResp, nil
}

func (s *productService) GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error){
	return s.Repository.FindByID(ctx, id)
}