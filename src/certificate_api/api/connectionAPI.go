package api

import (
	"certificate_api/configs"
	"certificate_api/controllers"
	"certificate_api/controllers/routes"
	"os"

	"github.com/gofiber/fiber/v2"
)

func ConnectionAPI(controller *controllers.CertificateController) {
	// Define Fiber config.
	config := configs.FiberConfig()

	// Define a new Fiber app with config.
	app := fiber.New(config)

	// Routes.
	// Define accesable routes
	routes.PublicRoutes(app, controller) // Register a public routes for app.

	// Start server in port 8081 (specify in the .env)
	app.Listen(os.Getenv("API_PORT"))
}
