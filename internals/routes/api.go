package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iamtaufik/golang-vercel-deployment/internals/handlers"
)

type RouteConfig struct {
	ProductHandler *handlers.ProductHandler
	AuthHandler    *handlers.AuthHandler
}

func RegisterRoutes(app *fiber.App, cfg *RouteConfig)  {
	api := app.Group("/api")

	RegisterAuthRoutes(api, cfg.AuthHandler)
	RegisterProductRoutes(api, cfg.ProductHandler)
}