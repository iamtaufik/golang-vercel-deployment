package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/iamtaufik/golang-vercel-deployment/internals/dto"
	"github.com/iamtaufik/golang-vercel-deployment/internals/models"
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

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	userIDRaw := c.Locals("user_id")
	userIDstr, ok := userIDRaw.(string)

	if !ok || userIDstr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized or userID not found in context",
		})
	}
	
	userID, err := uuid.Parse(userIDstr)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user id is not valid uuid",
		})
	}

	var product dto.ProductRequest
	
	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	newProduct := models.Product{
		UserID: userID,
		Name: product.Name,
		Price: product.Price,
	}

	if err := h.Services.CreateProduct(c.Context(), &newProduct); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	productResp := dto.ProductResponse{
		ID: newProduct.ID,
		Name: newProduct.Name,
		Price: newProduct.Price,
		CreatedAt: newProduct.CreatedAt,
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": productResp})
}