package controllers

import (
	"auth/jwt"
	"consulting_api/models"
	"consulting_api/repositories"
	"github.com/gofiber/fiber/v2"
	l "own_logger"
	"time"
)

type ConsultingElectionInfoController struct {
	repo       *repositories.ElectionRepo
	jwtManager *jwt.Manager
}

func NewConsultingElectionConfigController(repo *repositories.ElectionRepo, manager *jwt.Manager) *ConsultingElectionInfoController {
	return &ConsultingElectionInfoController{
		repo:       repo,
		jwtManager: manager,
	}
}

func (controller *ConsultingElectionInfoController) RequestElectionConfiguration(c *fiber.Ctx) error {
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
	anotherRoleForApi := jwt.GetRoles().Consulter
	validRole := ValidateRole(claims.Role, roleForApi, anotherRoleForApi)
	if !validRole {
		l.LogInfo("User " + claims.TokenInfo.Id + " tried to request election configuration but is not authorized")
		return c.Status(fiber.ErrForbidden.Code).JSON(fiber.Map{
			"error":   true,
			"msg":     "User is not authorized to request election configuration",
			"request": nil,
		})
	}
	electionId := c.Params("electionId")
	electionConfig, err := controller.repo.RequestElectionConfig(electionId)
	if err != nil {
		l.LogError(err.Error())
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
	l.LogInfo("Election configuration requested")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"msg":     "Config of election" + electionId,
		"request": electionResponse,
	})
}
