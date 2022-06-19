package routes

import (
	"consulting_api/controllers"
	"github.com/gofiber/fiber/v2"
)

func PublicRoutesElectionVotes(a *fiber.App, controller *controllers.ConsultingElectionVotesController) {
	route := a.Group("/api/v1")
	route.Get("/consulting/vote/:electionId/:voterId", controller.RequestVote)
	route.Get("/consulting/election/:electionId", controller.RequestElectionResult)
	route.Get("/consulting/election/votesHours/:electionId", controller.RequestPopularVotingTimes)
}

func PublicRoutesElectionInfo(a *fiber.App, controller *controllers.ConsultingElectionInfoController) {
	route := a.Group("/api/v1")
	route.Get("/consulting/election/config/:electionId", controller.RequestElectionConfiguration)
}
