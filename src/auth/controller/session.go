package controller

import(
	"auth/models"
	"github.com/gofiber/fiber/v2"
	"encoding/json"
	jwt "auth/jwt"

)

type SessionController struct {
	repo *repositories.UsersRepo 
}

func SessionController(repo *repositories.UsersRepo) *SessionController {
	return &SessionController{repo: repo}
}

func Login(c *fiber.Ctx) error{
	var login models.Login
	err := json.Unmarshal(c.Body(), &login)
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"error":   true,
			"msg":     err.Error(),
			"login": nil,
		})
	}

	user:= //get from db
	if user {
		token := jwt.Generate(user)
		result := models.ResponseModel{token}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"error":   false,
			"msg":     "Login succesful",
			"token": 	result,
		})
	}


}