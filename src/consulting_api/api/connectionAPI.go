package api

import (
	"consulting_api/configs"
	"consulting_api/controllers"
	"consulting_api/controllers/routes"
	"os"

	"github.com/gofiber/fiber/v2"
)

func ConnectionAPI(controller *controllers.ConsultingElectionVotesController, electionController *controllers.ConsultingElectionInfoController) {
	// Define Fiber config.
	config := configs.FiberConfig()

	// Define a new Fiber app with config.
	app := fiber.New(config)

	// Routes.
	routes.PublicRoutesElectionVotes(app, controller)        // Register a public routes for app.
	routes.PublicRoutesElectionInfo(app, electionController) // Register a public routes for app.

	app.Listen(os.Getenv("PORT"))
}
