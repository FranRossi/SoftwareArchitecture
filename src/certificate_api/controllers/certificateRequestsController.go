package controllers

import (
	"certificate_api/repositories"
	"github.com/gofiber/fiber/v2"
)

type CertificateRequestsController struct {
	repo *repositories.CertificateRequestsRepo
}

func CertificateRequestsController(repo *repositories.CertificateRequestsRepo) *CertificateRequestsController {
	return &CertificateController{repo: repo}
}

func (controller *CertificateRequestsController) RequestCertificate(c *fiber.Ctx) error {
	var request models.CertificateRequestsModel

	request.Timestamp = time.Now().Format(time.RFC3339)
	err := json.Unmarshal(c.Body(), &request)
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
			"request":  nil,
		})
	}

	err = controller.repo.AddRequest(request)
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
			"request":  nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error": false,
		"msg":   "Request created succesfully!",
		"request":  request,
	})
}