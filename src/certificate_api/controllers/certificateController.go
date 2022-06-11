package controllers

import (
	"certificate_api/repositories"
	"github.com/gofiber/fiber/v2"
)

type CertificateController struct {
	repo *repositories.CertificatesRepo
}

func CertificatesController(repo *repositories.CertificatesRepo) *CertificateController {
	return &CertificateController{repo: repo}
}

