package connections

import (
	"certificate_api/connections/configs"
	"certificate_api/connections/routes"
	"certificate_api/controllers"
	"certificate_api/repositories"
	"github.com/gofiber/fiber/v2"
)

func Connection() {
	config := configs.FiberConfig()
	app := fiber.New(config)

	repo := &repositories.CertificateRequestsRepo{}

	controller := controllers.CertificateRequestsController(repo2)

	routes.PublicRoutes(app, controller) 

	app.Listen(":8081")
}
