package controllers

import (
	"external_electoral_api/uruguay_election/repositories"
	"github.com/gofiber/fiber/v2"
)

type ElectionController struct {
	repo *repositories.ElectionRepo
}

func NewElectionController(repo *repositories.ElectionRepo) *ElectionController {
	return &ElectionController{repo: repo}
}

func (controller *ElectionController) GetElection(c *fiber.Ctx) error {

	election, err := controller.repo.GetElection(c.Params("id"))
	if err != nil {

		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":    true,
			"msg":      err.Error(),
			"election": nil,
		})
	}
	// Return status 200 OK.
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":    false,
		"msg":      nil,
		"election": election,
	})
}
