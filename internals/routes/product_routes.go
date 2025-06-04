package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iamtaufik/golang-vercel-deployment/internals/handlers"
)

func RegisterProductRoutes(router fiber.Router, h *handlers.ProductHandler) {
	user := router.Group("/products")

	user.Get("/", h.GetProducts)
}