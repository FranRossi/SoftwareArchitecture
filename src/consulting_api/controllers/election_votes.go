package controllers

import (
	"auth/jwt"
	"consulting_api/models"
	"consulting_api/repositories"

	"github.com/gofiber/fiber/v2"

	l "own_logger"
	"time"
)

type ConsultingElectionVotesController struct {
	repo       *repositories.VotesRepo
	jwtManager *jwt.Manager
}

func NewConsultingController(repo *repositories.VotesRepo, manager *jwt.Manager) *ConsultingElectionVotesController {
	return &ConsultingElectionVotesController{
		repo:       repo,
		jwtManager: manager,
	}
}

func (controller *ConsultingElectionVotesController) RequestVote(c *fiber.Ctx) error {
	timeQueryRequest := time.Now()
	token := c.GetReqHeaders()["Authorization"]
	claims, err := controller.jwtManager.Verify(token)
	if err != nil {
		l.LogError(err.Error())
		return c.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"error":   true,
			"msg":     err.Error(),
			"request": nil,
		})
	}
	roleForApi := jwt.GetRoles().Electoral
	validRole := ValidateRole(claims.Role, roleForApi)
	if !validRole {
		l.LogInfo("User " + claims.TokenInfo.Id + " tried to request a vote but is not authorized")
		return c.Status(fiber.ErrForbidden.Code).JSON(fiber.Map{
			"error":   true,
			"msg":     "User is not authorized to request a vote",
			"request": nil,
		})
	}
	voterId := c.Params("voterId")
	electionId := c.Params("electionId")
	vote, err := controller.repo.RequestVote(voterId, electionId)
	if err != nil {
		l.LogError(err.Error())
		return c.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"error":   true,
			"msg":     err.Error() + " election not found",
			"request": nil,
		})
	}
	timeQueryResponse := time.Now()
	vote.ElectionId = electionId
	vote.QueryRequestTime = timeQueryRequest.Format(time.RFC3339)
	vote.QueryResponseTime = timeQueryResponse.Format(time.RFC3339)
	vote.QueryProcessingTime = timeQueryResponse.Sub(timeQueryRequest).String()
	l.LogInfo("Requested voted by " + voterId + " for election " + electionId)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"msg":     "Vote found",
		"request": vote,
	})
}

func (controller *ConsultingElectionVotesController) RequestElectionResult(c *fiber.Ctx) error {
	timeQueryRequest := time.Now()
	token := c.GetReqHeaders()["Authorization"]
	claims, err := controller.jwtManager.Verify(token)
	if err != nil {
		l.LogError(err.Error())
		return c.Status(fiber.ErrUnauthorized.Code).JSON(fiber.Map{
			"error":   true,
			"msg":     err.Error(),
			"request": nil,
		})
	}
	roleForApi := jwt.GetRoles().Electoral
	validRole := ValidateRole(claims.Role, roleForApi)
	if !validRole {
		l.LogInfo("User " + claims.TokenInfo.Id + " tried to request election result but is not authorized")
		return c.Status(fiber.ErrForbidden.Code).JSON(fiber.Map{
			"error":   true,
			"msg":     "User is not authorized to request election result",
			"request": nil,
		})
	}
	electionId := c.Params("electionId")
	electionResult, err := controller.repo.RequestElectionResult(electionId)
	if err != nil {
		l.LogError(err.Error())
		return c.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"error":   true,
			"msg":     err.Error() + " election not found",
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
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"msg":     "Result of election" + electionId,
		"request": electionResponse,
	})
}

func (controller *ConsultingElectionVotesController) RequestPopularVotingTimes(c *fiber.Ctx) error {
	timeQueryRequest := time.Now()
	electionId := c.Params("electionId")
	votesPerHours, err := controller.repo.RequestPopularVotingTimes(electionId)
	if err != nil {
		l.LogError(err.Error())
		return c.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"error":   true,
			"msg":     err.Error() + " election not found",
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
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"msg":     "Votes Per Hours",
		"request": response,
	})
}

func ValidateRole(roleFromRequest string, roleForApi ...string) bool {
	for _, role := range roleForApi {
		if roleFromRequest == role {
			return true
		}
	}
	return false
}

func (controller *ConsultingElectionVotesController) RequestVotesPerCircuits(c *fiber.Ctx) error {
	timeQueryRequest := time.Now()
	electionId := c.Params("electionId")
	circuit := c.Params("circuitId")
	response, err := controller.repo.RequestVotesPerCircuits(electionId, circuit)
	if err != nil {
		l.LogError(err.Error())
		return c.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"error":   true,
			"msg":     err.Error() + " election not found",
			"request": nil,
		})
	}

	timeQueryResponse := time.Now()
	response.QueryRequestTime = timeQueryRequest.Format(time.RFC3339)
	response.QueryResponseTime = timeQueryResponse.Format(time.RFC3339)
	response.QueryProcessingTime = timeQueryResponse.Sub(timeQueryRequest).String()

	l.LogInfo("Requested votes per circuit for election " + electionId)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"msg":     "Votes Per Group",
		"request": response,
	})
}

func (controller *ConsultingElectionVotesController) RequestVotesPerRegions(c *fiber.Ctx) error {
	timeQueryRequest := time.Now()
	electionId := c.Params("electionId")
	region := c.Params("regionId")
	response, err := controller.repo.RequestVotesPerRegions(electionId, region)
	if err != nil {
		l.LogError(err.Error())
		return c.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"error":   true,
			"msg":     err.Error() + " election not found",
			"request": nil,
		})
	}
	timeQueryResponse := time.Now()
	response.QueryRequestTime = timeQueryRequest.Format(time.RFC3339)
	response.QueryResponseTime = timeQueryResponse.Format(time.RFC3339)
	response.QueryProcessingTime = timeQueryResponse.Sub(timeQueryRequest).String()

	l.LogInfo("Requested votes per circuit for election " + electionId)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"msg":     "Votes Per Group",
		"request": response,
	})
}
