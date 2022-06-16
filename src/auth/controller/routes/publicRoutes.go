package routes

import (
	"auth/controller"

	"github.com/gofiber/fiber/v2"
)

func PublicRoutes(a *fiber.App, controller *controllers.SessionController) {
	route := a.Group("/api/v1")

	route.Post("/login", controller.Login)

}
