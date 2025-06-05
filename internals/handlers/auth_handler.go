package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/iamtaufik/golang-vercel-deployment/internals/dto"
	"github.com/iamtaufik/golang-vercel-deployment/internals/models"
	"github.com/iamtaufik/golang-vercel-deployment/internals/services"
)

type AuthHandler struct {
	Service services.AuthService
}

func NewAuthService(service services.AuthService) *AuthHandler {
	return &AuthHandler{Service: service}
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var request dto.LoginRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request!"})
	}

	accessToken, refreshToken, err := h.Service.Login(c.Context(), request.Email, request.Password)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	c.Cookie(&fiber.Cookie{
		Name: "refreshToken",
		Value: refreshToken,
		Path: "/",
		HTTPOnly: true,
		Secure: true,
		SameSite: "None",
		Expires:  time.Now().AddDate(0, 0, 7),
	})

	loginResp := dto.LoginResponse{
		AccessToken: accessToken,
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": loginResp})
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var request dto.RegisterRequest

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request!"})
	}
	
	if request.Password != request.ConfPassword {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "password and confirm password not match"})
	}

	body := models.User{
		Name: request.Name,
		Email: request.Email,
		Password: request.Password,
	}

	err := h.Service.Register(c.Context(), &body)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	resp := dto.RegisterResponse{
		Name: request.Name,
		Email: request.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": resp})
}