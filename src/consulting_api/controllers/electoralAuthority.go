package controllers

import (
	"consulting_api/models"
	"consulting_api/repositories"
	"github.com/gofiber/fiber/v2"
	"time"
)

type ConsultingElectoralAuthorityController struct {
	repo *repositories.VotesRepo
}

func NewConsultingController(repo *repositories.VotesRepo) *ConsultingElectoralAuthorityController {
	return &ConsultingElectoralAuthorityController{repo: repo}
}

func (controller *ConsultingElectoralAuthorityController) RequestVote(c *fiber.Ctx) error {
	timeQueryRequest := time.Now()
	voterId := c.Params("voterId")
	electionId := c.Params("electionId")
	vote, err := controller.repo.RequestVote(voterId, electionId)
	if err != nil {
		return c.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"error":   true,
			"msg":     err.Error(),
			"request": nil,
		})
	}
	timeQueryResponse := time.Now()
	vote.QueryRequestTime = timeQueryRequest.Format(time.RFC3339)
	vote.QueryResponseTime = timeQueryResponse.Format(time.RFC3339)
	vote.QueryProcessingTime = timeQueryResponse.Sub(timeQueryRequest).String()
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"msg":     "Vote found",
		"request": vote,
	})
}

func (controller *ConsultingElectoralAuthorityController) RequestElectionResult(c *fiber.Ctx) error {
	timeQueryRequest := time.Now()
	electionId := c.Params("electionId")
	electionResult, err := controller.repo.RequestElectionResult(electionId)
	if err != nil {
		return c.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"error":   true,
			"msg":     err.Error(),
			"request": nil,
		})
	}
	electionResponse := models.ElectionModel{
		Result:           electionResult,
		QueryRequestTime: timeQueryRequest.Format(time.RFC3339),
	}
	timeQueryResponse := time.Now()
	electionResponse.QueryRequestTime = timeQueryRequest.Format(time.RFC3339)
	electionResponse.QueryResponseTime = timeQueryResponse.Format(time.RFC3339)
	electionResponse.QueryProcessingTime = timeQueryResponse.Sub(timeQueryRequest).String()
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"msg":     "Result of election" + electionId,
		"request": electionResponse,
	})
}
