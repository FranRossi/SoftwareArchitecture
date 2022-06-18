package controllers

import (
	"consulting_api/models"
	"consulting_api/repositories"
	"github.com/gofiber/fiber/v2"
	l "own_logger"
	"time"
)

type ConsultingElectionVotesController struct {
	repo *repositories.VotesRepo
}

func NewConsultingController(repo *repositories.VotesRepo) *ConsultingElectionVotesController {
	return &ConsultingElectionVotesController{repo: repo}
}

func (controller *ConsultingElectionVotesController) RequestVote(c *fiber.Ctx) error {
	timeQueryRequest := time.Now()
	voterId := c.Params("voterId")
	electionId := c.Params("electionId")
	vote, err := controller.repo.RequestVote(voterId, electionId)
	if err != nil {
		l.LogError(err.Error())
		return c.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"error":   true,
			"msg":     err.Error(),
			"request": nil,
		})
	}
	timeQueryResponse := time.Now()
	vote.ElectionId = electionId
	vote.QueryRequestTime = timeQueryRequest.Format(time.RFC3339)
	vote.QueryResponseTime = timeQueryResponse.Format(time.RFC3339)
	vote.QueryProcessingTime = timeQueryResponse.Sub(timeQueryRequest).String()
	l.LogInfo("Requested voted by " + voterId + " for election " + electionId)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"msg":     "Vote found",
		"request": vote,
	})
}

func (controller *ConsultingElectionVotesController) RequestElectionResult(c *fiber.Ctx) error {
	timeQueryRequest := time.Now()
	electionId := c.Params("electionId")
	electionResult, err := controller.repo.RequestElectionResult(electionId)
	if err != nil {
		l.LogError(err.Error())
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
	l.LogInfo("Requested election result for election " + electionId)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"msg":     "Result of election" + electionId,
		"request": electionResponse,
	})
}

func (controller *ConsultingElectionVotesController) RequestAverageVotingTime(c *fiber.Ctx) error {
	timeQueryRequest := time.Now()
	electionId := c.Params("electionId")
	votesPerHours, err := controller.repo.RequestPopularVotingTimes(electionId)
	if err != nil {
		l.LogError(err.Error())
		return c.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"error":   true,
			"msg":     err.Error(),
			"request": nil,
		})
	}
	response := models.VotesPerHours{
		ElectionId:            electionId,
		AmountOfVotesPerHours: votesPerHours,
		QueryRequestTime:      timeQueryRequest.Format(time.RFC3339),
	}
	timeQueryResponse := time.Now()
	response.QueryResponseTime = timeQueryResponse.Format(time.RFC3339)
	response.QueryProcessingTime = timeQueryResponse.Sub(timeQueryRequest).String()
	l.LogInfo("Requested popular voting times for election " + electionId)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"msg":     "Votes Per Hours",
		"request": response,
	})
}
