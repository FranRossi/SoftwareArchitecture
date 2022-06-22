package routes

import (
	"certificate_api/controllers"

	"github.com/gofiber/fiber/v2"
)

func PublicRoutes(a *fiber.App, controller *controllers.CertificateController) {
	route := a.Group("/api/v1")

	route.Get("/certificate/request", controller.RequestCertificate)

}
