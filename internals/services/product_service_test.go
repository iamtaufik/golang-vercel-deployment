package services

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/iamtaufik/golang-vercel-deployment/internals/models"
)

// mockProductRepository adalah mock manual yang implement ProductRepository
type mockProductRepository struct {
	mockFindAll func(ctx context.Context) ([]models.Product, error)
	mockGetByID func(ctx context.Context, id uuid.UUID) (*models.Product, error)
	mockCreate func(ctx context.Context, product *models.Product) error
}

func (m *mockProductRepository) FindAll(ctx context.Context) ([]models.Product, error) {
	return m.mockFindAll(ctx)
}

func (m *mockProductRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	return m.mockGetByID(ctx, id)
}
func (m *mockProductRepository) Create(ctx context.Context, product *models.Product) error {
	return m.mockCreate(ctx, product)
}

func TestCreateProduct_Success(t *testing.T) {
	mockRepo := &mockProductRepository{
		mockCreate: func(ctx context.Context, product *models.Product) error {
			if product.Name == "" {
				return errors.New("Name product is required")
			} else if product.Price <= 0 {
				return errors.New("Price must greater than 0")
			}

			return nil
		},
	}

	service := NewProductService(mockRepo)

	product := &models.Product{
		ID: uuid.New(),
		Name: "Product B",
		Price: 2000,
	}

	err := service.CreateProduct(context.Background(), product)

	if err != nil {
		t.Fatalf("expected no error, but get %v", err)
	}
}

func TestGetProducts_Success(t *testing.T) {
	// Arrange
	mockRepo := &mockProductRepository{
		mockFindAll: func(ctx context.Context) ([]models.Product, error) {
			return []models.Product{
				{ID: uuid.New(), Name: "Produk A", Price: 10000},
				{ID: uuid.New(), Name: "Produk B", Price: 20000},
			}, nil
		},
	}

	service := NewProductService(mockRepo)

	// Act
	products, err := service.GetProducts(context.Background())

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(products) != 2 {
		t.Errorf("expected 2 products, got %d", len(products))
	}
}

func TestGetProducts_Error(t *testing.T) {
	// Arrange
	mockRepo := &mockProductRepository{
		mockFindAll: func(ctx context.Context) ([]models.Product, error) {
			return nil, errors.New("database error")
		},
	}

	service := NewProductService(mockRepo)

	// Act
	products, err := service.GetProducts(context.Background())

	// Assert
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if products != nil {
		t.Errorf("expected nil products, got %+v", products)
	}
}

func TestGetProduct_Success(t *testing.T){
	expectedID := uuid.New()

	expectedProduct := &models.Product{
		ID: expectedID,
		Name: "Product A",
		Price: 20000,
	}

	mockRepo := &mockProductRepository{
		mockGetByID: func(ctx context.Context, id uuid.UUID) (*models.Product, error) {
			if id == expectedID {
				return expectedProduct, nil
			}

			return nil, errors.New("product not found")
		},
	}

	service := NewProductService(mockRepo)
	
	product, err := service.GetProduct(context.Background(), expectedID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if product == nil {
		t.Errorf("expected product, got %+v", product)
	}
}

func TestGetProduct_Error(t *testing.T) {
	expectedID := uuid.New()

	expectedProduct := &models.Product{
		ID: expectedID,
		Name: "Product A",
		Price: 20000,
	}

	mockRepo := &mockProductRepository{
		mockGetByID: func(ctx context.Context, id uuid.UUID) (*models.Product, error) {
			randomID := uuid.New()
			if id == randomID {
				return expectedProduct, nil
			}

			return nil, errors.New("product not found")
		},
	}

	service := NewProductService(mockRepo)
	
	product, err := service.GetProduct(context.Background(), expectedID)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if product != nil {
		t.Errorf("expected nil products, got %+v", product)
	}
}
