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

	repo := &repositories.CertificatesRepo{}
	repo2 := &repositories.CertificateRequestsRepo{}

	// Creo una instancia de mis controladores con mi instancia de repo
	controller := controllers.CertificatesController(repo)
	controller2 := controllers.CertificatesRequestController(repo2)

	// Routes.
	// Aqui defino cuales van a ser las rutas accesibles
	routes.PublicRoutes(app, controller) // Register a public routes for app.
	routes.PublicRoutes(app, controller2) // Register a public routes for app.

	// Aqui inicializamos el servidor en el puerto 8081
	app.Listen(":8081")
}
