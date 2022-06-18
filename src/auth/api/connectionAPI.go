package api

import (
	"auth/configs"
	"auth/controller"
	"auth/controller/routes"

	"github.com/gofiber/fiber/v2"
)

func ConnectionAPI(controller *controllers.SessionController) {
	config := configs.FiberConfig()
	app := fiber.New(config)
	routes.PublicRoutes(app, controller) 

	
	app.Listen(":8082")
}
