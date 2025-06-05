package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/iamtaufik/golang-vercel-deployment/internals/db"
	"github.com/iamtaufik/golang-vercel-deployment/internals/handlers"
	"github.com/iamtaufik/golang-vercel-deployment/internals/repository"
	"github.com/iamtaufik/golang-vercel-deployment/internals/routes"
	"github.com/iamtaufik/golang-vercel-deployment/internals/services"
	"github.com/joho/godotenv"
)
 
var app *fiber.App

func init() {
  	err := godotenv.Load(".env") 
	if err != nil {
		fmt.Println("Failed to load .env file")
	}

	db := db.ConnectDB()

	pRepository := repository.NewProductRepository(db)
	pService 	:= services.NewProductService(pRepository)
	pHandler	:= handlers.NewProductHandler(pService)

	uRepository := repository.NewUserRepository(db)
	aService 	:= services.NewAuthService(uRepository)
	aHandler	:= handlers.NewAuthService(aService)

	// Initialize Fiber app once
	app = fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("CLIENT_URL"), // Your Vue development server origin
		AllowCredentials: true,
		// AllowHeaders: "Origin, Content-Type, Accept, Authorization", // Good to explicitly allow headers if you send them
	}))

	routes.RegisterRoutes(app, &routes.RouteConfig{
		ProductHandler: pHandler,
		AuthHandler: aHandler,
	})

	// Define your Fiber routes here
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Hello from Fiber on Vercel!",
			"path":    c.Path(),
			"query":   c.Query("name"),
		})
	})

	// Add more routes as needed
}

func Handler(w http.ResponseWriter, r *http.Request) {
  adaptor.FiberApp(app)(w, r)
}