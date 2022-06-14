package controllers

import (
	"consulting_api/repositories"
	"github.com/gofiber/fiber/v2"
	"time"
)

type ConsultingController struct {
	repo *repositories.ConsultingRepo
}

func NewConsultingController(repo *repositories.ConsultingRepo) *ConsultingController {
	return &ConsultingController{repo: repo}
}

func (controller *ConsultingController) RequestVote(c *fiber.Ctx) error {
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
