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

