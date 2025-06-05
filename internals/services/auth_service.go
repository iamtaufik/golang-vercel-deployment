package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iamtaufik/golang-vercel-deployment/internals/models"
	"github.com/iamtaufik/golang-vercel-deployment/internals/repository"
	"github.com/iamtaufik/golang-vercel-deployment/internals/utils/crypto"
	"github.com/iamtaufik/golang-vercel-deployment/internals/utils/jwt"
	"gorm.io/gorm"
)

type AuthService interface {
	Login(context context.Context, email, password string) (string, string, error)
	Register(context context.Context, user *models.User) error
}

type authService struct {
	Repository repository.UserRepository
}

func NewAuthService(repository repository.UserRepository) *authService {
	return &authService{Repository: repository}
}

func (s *authService) Login(ctx context.Context, email, password string) (string, string, error){
	user, err := s.Repository.FindByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", errors.New("user not found")
		}
		return "", "", err
	}

	if isMatch := crypto.CheckPasswordHash(password, user.Password); !isMatch {
		return "", "", errors.New("invalid credentials")
	}

	accessToken, err := jwt.GenerateAccessToken(user.ID.String())

	if err != nil {
		return "", "", errors.New("failed create access token")
	}
	
	refreshToken, err := jwt.GenerateRefreshToken(user.ID.String())

	if err != nil {
		return "", "", errors.New("failed create refresh token")
	}
	
	return accessToken, refreshToken, nil
}

func (s *authService) Register(ctx context.Context, input *models.User) error {

	existedUser, err := s.Repository.FindByEmail(ctx, input.Email)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if existedUser != nil {
		return errors.New("email already used")
	}

	hashedPassword, err :=  crypto.HashPassword(input.Password)

	if err != nil {
		return errors.New("failed to hashed password")
	}

	user := models.User{
		ID: uuid.New(),
		Name: input.Name,
		Email: input.Email,
		Password: hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.Repository.Create(ctx, &user); err != nil {
		return errors.New("failed to create user")
	}

	return nil
}