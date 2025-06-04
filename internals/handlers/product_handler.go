package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/iamtaufik/golang-vercel-deployment/internals/services"
)

type ProductHandler struct {
	Services services.ProductService
}

func NewProductHandler(service services.ProductService) *ProductHandler {
	return &ProductHandler{Services: service}
}

func (h *ProductHandler) GetProducts(c *fiber.Ctx) error {
	products, err := h.Services.GetProducts(c.Context())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": products})
}