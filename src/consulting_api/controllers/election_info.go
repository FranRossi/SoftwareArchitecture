package controllers

import (
	"consulting_api/models"
	"consulting_api/repositories"
	"github.com/gofiber/fiber/v2"
	"time"
)

type ConsultingElectionInfoController struct {
	repo *repositories.ElectionRepo
}

func NewConsultingElectionConfigController(repo *repositories.ElectionRepo) *ConsultingElectionInfoController {
	return &ConsultingElectionInfoController{repo: repo}
}

func (controller *ConsultingElectionInfoController) RequestElectionConfiguration(c *fiber.Ctx) error {
	timeQueryRequest := time.Now()
	electionId := c.Params("electionId")
	electionConfig, err := controller.repo.RequestElectionConfig(electionId)
	if err != nil {
		return c.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"error":   true,
			"msg":     err.Error(),
			"request": nil,
		})
	}
	electionResponse := models.ElectionConfig{
		ElectionId:       electionId,
		Emails:           electionConfig.Emails,
		MaxVotes:         electionConfig.MaxVotes,
		MaxCertificates:  electionConfig.MaxCertificates,
		QueryRequestTime: timeQueryRequest.Format(time.RFC3339),
	}
	timeQueryResponse := time.Now()
	electionResponse.QueryRequestTime = timeQueryRequest.Format(time.RFC3339)
	electionResponse.QueryResponseTime = timeQueryResponse.Format(time.RFC3339)
	electionResponse.QueryProcessingTime = timeQueryResponse.Sub(timeQueryRequest).String()
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"msg":     "Config of election" + electionId,
		"request": electionResponse,
	})
}
