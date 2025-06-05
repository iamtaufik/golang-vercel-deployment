package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/iamtaufik/golang-vercel-deployment/internals/models"
	"github.com/iamtaufik/golang-vercel-deployment/internals/utils/crypto"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type mockAuthRepository struct {
	mockRegister     func(ctx context.Context, user *models.User) error
	mockFindByEmail  func(ctx context.Context, email string) (*models.User, error)
}

func (m *mockAuthRepository) Create(ctx context.Context, user *models.User) error {
	return m.mockRegister(ctx, user)
}

func (m *mockAuthRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	if m.mockFindByEmail != nil {
		return m.mockFindByEmail(ctx, email)
	}
	return nil, nil
}

func TestRegister_Success(t *testing.T) {
	mockRepo := &mockAuthRepository{
		mockRegister: func(context context.Context, user *models.User) error {
			if user.Email == "" || user.Name == "" || user.Password == "" {
				return errors.New("all fields are required")
			}
			return nil
		},
		mockFindByEmail: func(ctx context.Context, email string) (*models.User, error) {
			// Simulasi bahwa email belum digunakan
			return nil, gorm.ErrRecordNotFound
		},
	}

	service := NewAuthService(mockRepo)

	user := models.User{
		ID: uuid.New(),
		Name: "Taufik",
		Email: "taufik@dev.com",
		Password: "1234567890",
	}
	err := service.Register(context.Background(), &user)

	if err != nil {
		t.Fatalf("expected no error, but get %v", err)
	}
}

func TestRegister_Error(t *testing.T) {
	mockRepo := &mockAuthRepository{
		mockRegister: func(context context.Context, user *models.User) error {
			return errors.New("email already used")
		},
		mockFindByEmail: func(ctx context.Context, email string) (*models.User, error) {
			// Simulasi bahwa email sudah digunakan
			return &models.User{
				Name: "Taufik",
				Email: "taufik@dev.com",
				Password: "1234567890",
			}, nil
		},
	}

	service := NewAuthService(mockRepo)

	user := models.User{
		ID: uuid.New(),
		Name: "Taufik",
		Email: "taufik@dev.com",
		Password: "1234567890",
	}

	err := service.Register(context.Background(), &user)

	assert.Equal(t, "email already used", err.Error())
}

func TestLogin_Success(t *testing.T) {
	mockRepo := &mockAuthRepository{
		mockFindByEmail: func(ctx context.Context, email string) (*models.User, error) {

			hashedPassword, err := crypto.HashPassword("1234567890")

			if err != nil {
				return nil, err
			}

			return &models.User{
				Name: "Taufik",
				Email: "taufik@dev.com",
				Password: hashedPassword,
			}, nil
		},
	}

	service := NewAuthService(mockRepo)

	request := &models.User{
		Email: "taufik@dev.com",
		Password: "1234567890",
	}

	accessToken, refreshToken, err := service.Login(context.Background(), request.Email, request.Password)

	if err != nil {
		t.Fatalf("expected no error, but get %v", err)
	}
	fmt.Sprintf("accessToken: %v\n", accessToken)
	assert.NotEmpty(t, accessToken)
	assert.NotEmpty(t, refreshToken)
}

func TestLogin_Error(t *testing.T) {
	mockRepo := &mockAuthRepository{
		mockFindByEmail: func(ctx context.Context, email string) (*models.User, error) {

			hashedPassword, err := crypto.HashPassword("rahasia123")

			if err != nil {
				return nil, err
			}

			return &models.User{
				Name: "Taufik",
				Email: "taufik@dev.com",
				Password: hashedPassword,
			}, nil
		},
	}

	service := NewAuthService(mockRepo)

	request := &models.User{
		Email: "taufik@dev.com",
		Password: "1234567890",
	}

	accessToken, refreshToken, err := service.Login(context.Background(), request.Email, request.Password)

	if err != nil {
		assert.Equal(t, "invalid credentials", err.Error())
	}

	assert.Empty(t, accessToken)
	assert.Empty(t, refreshToken)
}