package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iamtaufik/golang-vercel-deployment/internals/handlers"
	"github.com/iamtaufik/golang-vercel-deployment/internals/middlewares"
)

func RegisterAuthRoutes(router fiber.Router, h *handlers.AuthHandler) {
	auth := router.Group("/auth")

	auth.Post("/login", h.Login)
	auth.Post("/register", h.Register)
	auth.Get("/me", middlewares.JWTProtected(), h.Me)
}
