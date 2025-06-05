package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iamtaufik/golang-vercel-deployment/internals/handlers"
)

func RegisterAuthRoutes(router fiber.Router, h *handlers.AuthHandler) {
	auth := router.Group("/auth")

	auth.Post("/login", h.Login)
}
