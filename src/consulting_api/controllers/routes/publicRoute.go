package routes

import (
	"consulting_api/controllers"
	"github.com/gofiber/fiber/v2"
)

func PublicRoutes(a *fiber.App, controller *controllers.ConsultingController) {
	route := a.Group("/api/v1")
	route.Get("/consulting/vote/:electionId/:voterId", controller.RequestVote)
}
