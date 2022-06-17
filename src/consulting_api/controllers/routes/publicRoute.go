package routes

import (
	"consulting_api/controllers"
	"github.com/gofiber/fiber/v2"
)

func PublicRoutesElectoralAuth(a *fiber.App, controller *controllers.ConsultingElectoralAuthorityController) {
	route := a.Group("/api/v1")
	route.Get("/consulting/vote/:electionId/:voterId", controller.RequestVote)
	route.Get("/consulting/election/:electionId", controller.RequestElectionResult)
}

func PublicRoutesElectionConfig(a *fiber.App, controller *controllers.ConsultingElectoralConfigController) {
	route := a.Group("/api/v1")
	route.Get("/consulting/electionConfig/:electionId", controller.RequestElectionConfiguration)
}
