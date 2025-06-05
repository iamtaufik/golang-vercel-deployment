package main

import (
	"fmt"
	"log"
	"os"

	// Import your Vercel function's package
	// The path here depends on your project structure.
	// If api/index.go is at the root, it might be "your-module-name/api"
	// If your go.mod is at the root and api/index.go is there, it's just "api"
	// For this example, let's assume your go.mod is at the root and
	// your Vercel function is in `api/index.go` with package `handler`.
	// The import path will be like this:
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/iamtaufik/golang-vercel-deployment/internals/db"
	"github.com/iamtaufik/golang-vercel-deployment/internals/handlers"
	"github.com/iamtaufik/golang-vercel-deployment/internals/repository"
	"github.com/iamtaufik/golang-vercel-deployment/internals/routes"
	"github.com/iamtaufik/golang-vercel-deployment/internals/services"
	"github.com/joho/godotenv"
)

func init(){
	err := godotenv.Load(".env") 
	if err != nil {
		fmt.Println("Failed to load .env file")
	}
}

func main() {
	app := fiber.New()
	db := db.ConnectDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	pRepository := repository.NewProductRepository(db)
	pService 	:= services.NewProductService(pRepository)
	pHandler	:= handlers.NewProductHandler(pService)
	
	uRepository := repository.NewUserRepository(db)
	aService 	:= services.NewAuthService(uRepository)
	aHandler	:= handlers.NewAuthService(aService)

	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("CLIENT_URL"), // Your Vue development server origin
		AllowCredentials: true,
		// AllowHeaders: "Origin, Content-Type, Accept, Authorization", // Good to explicitly allow headers if you send them
	}))

	routes.RegisterRoutes(app, &routes.RouteConfig{
		ProductHandler: pHandler,
		AuthHandler: aHandler,
	})

	log.Fatal(app.Listen(":" + port))
}