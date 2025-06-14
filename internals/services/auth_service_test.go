package services

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/iamtaufik/golang-vercel-deployment/internals/models"
	"github.com/iamtaufik/golang-vercel-deployment/internals/utils/crypto"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type mockAuthRepository struct {
	mockRegister    func(ctx context.Context, user *models.User) error
	mockFindByEmail func(ctx context.Context, email string) (*models.User, error)
	mockFindByID	func(ctx context.Context, id uuid.UUID) (*models.User, error)
}

func (m *mockAuthRepository) Create(ctx context.Context, user *models.User) error {
	return m.mockRegister(ctx, user)
}

func (m *mockAuthRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	return m.mockFindByID(ctx, id)
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

func TestMe_Success(t *testing.T) {
	expectedID := uuid.New()

	mockRepo := &mockAuthRepository{
		mockFindByID: func(ctx context.Context, id uuid.UUID) (*models.User, error) {
			if id != expectedID {
				return nil, errors.New("user not found")
			}
			return &models.User{
				ID:    expectedID,
				Name:  "Taufik",
				Email: "taufik@dev.com",
			}, nil
		},
	}

	service := NewAuthService(mockRepo)

	user, err := service.Me(context.Background(), expectedID.String())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if user.ID != expectedID {
		t.Errorf("expected ID %v, got %v", expectedID, user.ID)
	}

	if user.Name != "Taufik" {
		t.Errorf("expected Name Taufik, got %s", user.Name)
	}

}


func TestMe_Error(t *testing.T) {
	expectedID := uuid.New()

	mockRepo := &mockAuthRepository{
		mockFindByID: func(ctx context.Context, id uuid.UUID) (*models.User, error) {
			return nil, errors.New("invalid user id")
		},
	}

	service := NewAuthService(mockRepo)

	user, err := service.Me(context.Background(), expectedID.String())
	
	if err != nil {
		assert.Equal(t, "invalid user id", err.Error())
	}

	assert.Empty(t, user)
}

func TestRefresh_Success(t *testing.T) {
	mockRepo := &mockAuthRepository{
		mockFindByEmail: func(ctx context.Context, email string) (*models.User, error) {

			hashedPassword, err := crypto.HashPassword("1234567890")

			if err != nil {
				return nil, err
			}

			return &models.User{
				ID: uuid.New(),
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

	_, refreshToken, err := service.Login(context.Background(), request.Email, request.Password)

	if err != nil {
		t.Fatalf("expected no error, but get %v", err)
	}

	accessToken, err := service.Refresh(context.Background(), refreshToken)

	if err != nil {
		t.Fatalf("expected no error, but get %v", err)
	}

	assert.NotEmpty(t, accessToken)
}

func TestRefresh_Test(t *testing.T) {
	mockRepo := &mockAuthRepository{
		mockFindByEmail: func(ctx context.Context, email string) (*models.User, error) {

			hashedPassword, err := crypto.HashPassword("1234567890")

			if err != nil {
				return nil, err
			}

			return &models.User{
				ID: uuid.New(),
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

	_, refreshToken, err := service.Login(context.Background(), request.Email, request.Password)

	
	if err != nil {
		t.Fatalf("expected no error, but get %v", err)
	}
	// Invalid refresh token	
	refreshToken = "abcde"

	accessToken, err := service.Refresh(context.Background(), refreshToken)

	if err != nil {
		assert.Equal(t, "invalid or expired refresh token", err.Error())
	}

	assert.Empty(t, accessToken)
}