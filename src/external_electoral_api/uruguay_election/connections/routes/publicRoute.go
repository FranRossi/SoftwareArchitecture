package routes

import (
	"external_electoral_api/uruguay_election/controllers"
	"github.com/gofiber/fiber/v2"
)

func PublicRoutes(a *fiber.App, controller *controllers.ElectionController) {
	// Create routes group.
	route := a.Group("/api/v1")

	// Routes for Get method:
	route.Get("/election/uruguay", controller.GetElection)
}
