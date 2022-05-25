package controllers

import (
	"ElectoralService/UruguayanElection/models"
	"ElectoralService/UruguayanElection/repositories"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

type ElectionController struct {
	repo *repositories.ElectionRepo
}

func NewElectionController(repo *repositories.ElectionRepo) *ElectionController {
	return &ElectionController{repo: repo}
}

func (controller *ElectionController) SendElectionSettings(c *fiber.Ctx) error {
	var election models.Election
	// Tengo que pasar la direccion de memoria asi va a pone las variables
	// Que vienen en el contexto declaradas segun el json del modelo
	fmt.Println(string(c.Body()))

	//election.CreatedAt = time.Now()
	err := json.Unmarshal(c.Body(), &election)
	if err != nil {
		// Return, if book has invalid fields
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"error":    true,
			"msg":      err.Error(),
			"election": nil,
		})
	}

	err = controller.repo.SendElectionSettings(election)
	if err != nil {
		// Return, if book not found.
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
			"book":  nil,
		})
	}

	// Return status 201 Created.
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":    false,
		"msg":      "Election settings send successfully!",
		"election": election,
	})
}
