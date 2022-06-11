package routes

import (
	"certificat_api/controllers"
	"github.com/gofiber/fiber/v2"
)

func PublicRoutes(a *fiber.App, controller *controllers.ElectionController) {
	route := a.Group("/api/v1")

}