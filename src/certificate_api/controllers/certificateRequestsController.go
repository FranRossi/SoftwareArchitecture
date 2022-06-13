package controllers

import (
	"certificate_api/models"
	"certificate_api/repositories"
	"encoding/json"
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
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"error":   true,
			"msg":     err.Error(),
			"request": nil,
		})
	}

	err = controller.repo.StoreRequest(request)
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"error":   true,
			"msg":     err.Error(),
			"request": nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"msg":     "Request created succesfully!",
		"request": request,
	})
}
