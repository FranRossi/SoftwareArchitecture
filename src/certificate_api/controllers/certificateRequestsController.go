package controllers

import (
	"certificate_api/models"
	"certificate_api/providers"
	"certificate_api/repositories"
	electoral_service_models "electoral_service/models"
	"encoding/json"
	"encrypt"
	"fmt"
	"os"
	l "own_logger"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type CertificateController struct {
	repo *repositories.CertificatesRepo
}

func CertificateRequestsController(repo *repositories.CertificatesRepo) *CertificateController {
	return &CertificateController{repo: repo}
}

func (controller *CertificateController) RequestCertificate(c *fiber.Ctx) error {
	var request models.CertificateRequestModel

	request.Timestamp = time.Now().Format(time.RFC3339)
	err := json.Unmarshal(c.Body(), &request)
	if err != nil {
		l.LogError(err.Error())
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"error":   true,
			"msg":     err.Error(),
			"request": nil,
		})
	}
	err = controller.repo.StoreRequest(request)
	if err != nil {
		l.LogError(err.Error())
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"error":   true,
			"msg":     err.Error(),
			"request": nil,
		})
	}
	go controller.checkMaxCertificateLimit(request.VoterId)
	l.LogInfo("Certificate request stored successfully")
	go controller.sendCertificateRequestedToEmail(request)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"msg":     "Certificate requested successfully, check your email!",
		"request": request,
	})
}

func (controller *CertificateController) checkMaxCertificateLimit(voterId string) {
	amountOfCertificates, err := controller.repo.FindAmountOfRequest(voterId)
	if err != nil {
		l.LogError(err.Error())
	}
	limit := os.Getenv("MAX_REQUEST")
	limitCertificates, err := strconv.Atoi(limit)
	if err != nil {
		l.LogError(err.Error())
	}
	if amountOfCertificates >= limitCertificates {
		controller.sendCertificateAlert(voterId)
	}
}

func (controller *CertificateController) sendCertificateAlert(voterId string) {
	go providers.SendAlert(voterId)
}

func (controller *CertificateController) sendCertificateRequestedToEmail(request models.CertificateRequestModel) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	var certificate models.CertificateModel
	var err error
	go func() {
		certificate, err = controller.repo.FindCertificate(request.VoterId, request.VoteIdentification)
		if err != nil {
			l.LogError(err.Error())
			fmt.Println(err.Error() + "cannot find certificate")
			return
		}
		wg.Done()
	}()
	var voter electoral_service_models.VoterModel
	go func() {
		voter, err = controller.repo.FindVoter(request.VoterId)
		if err != nil {
			l.LogError(err.Error())
			fmt.Println("Error finding voter with id: " + request.VoterId)
			return
		}
		wg.Done()
	}()
	wg.Wait()
	encrypt.DecryptVoter(&voter)
	go providers.SendEmail(certificate, voter)
}
