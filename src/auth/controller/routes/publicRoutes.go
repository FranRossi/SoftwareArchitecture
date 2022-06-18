package routes

import (
	"auth/controller"

	"github.com/gofiber/fiber/v2"
)

func PublicRoutes(a *fiber.App, controller *controller.SessionController) {
	route := a.Group("/api/v1")

	route.Post("/users/login", controller.Login)
	route.Post("/users/", controller.RegisterUser)

}
