package connections

import (
	"certificate_api/configs"
	"certificate_api/controllers"
	"certificate_api/repositories"
	"certificate_api/datasources"
	"certificate_api/connections/routes"
	"github.com/gofiber/fiber/v2"
	"context"
)

func Connection() {
	config := configs.FiberConfig()
	app := fiber.New(config)

	mongoClient, err := datasources.NewMongoDataSource("mongodb://docker:mongopw@localhost:55000")

	if err != nil {
		panic(err)
	}

	defer mongoClient.Disconnect(context.TODO())

	repo := repositories.NewRequestsRepo(mongoClient, "certificates")

	controller := controllers.CertificateRequestsController(repo)

	routes.PublicRoutes(app, controller) 

	app.Listen(":8081")
}
