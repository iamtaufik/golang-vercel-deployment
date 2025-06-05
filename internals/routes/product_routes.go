package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iamtaufik/golang-vercel-deployment/internals/handlers"
	"github.com/iamtaufik/golang-vercel-deployment/internals/middlewares"
)

func RegisterProductRoutes(router fiber.Router, h *handlers.ProductHandler) {
	user := router.Group("/products")

	user.Get("/", middlewares.JWTProtected(), h.GetProducts)
}