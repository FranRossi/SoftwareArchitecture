package controller

import (
	"auth/jwt"
	"auth/models"
	"auth/repository"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	l "own_logger"
	"time"
)

type SessionController struct {
	repo *repository.UsersRepo
}

func NewSessionController(repo *repository.UsersRepo) *SessionController {
	return &SessionController{repo: repo}
}

func (controller *SessionController) Login(c *fiber.Ctx) error {
	var login models.Login
	err := json.Unmarshal(c.Body(), &login)
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
			"login": nil,
		})
	}

	user, err := controller.repo.FindUser(login.Id)
	if err != nil {
		return c.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error(),
			"login": nil,
		})
	}
	err2 := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(login.Password))
	if err2 != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"error": true,
			"msg":   err2.Error() + " password is incorrect",
			"login": nil,
		})
	}
	generator := models.TokenInfo{
		Id:   user.Id,
		Role: user.Role,
	}
	duration := 30 * time.Minute

	privateKey, err := ioutil.ReadFile("./private.rsa")
	if err != nil {
		l.LogError(err.Error())
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error() + " cannot read private key",
			"login": nil,
		})
	}
	manager := jwt.NewJWTManager(privateKey, duration)

	token, err := manager.Generate(generator)
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error() + " cannot generate token",
			"login": nil,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error": false,
		"msg":   "Login successful",
		"token": token,
	})
}

func (controller *SessionController) RegisterUser(c *fiber.Ctx) error {
	var register models.UserRegister
	err := json.Unmarshal(c.Body(), &register)
	if err != nil {
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"error": true,
			"msg":   "User malformed",
			"user":  nil,
		})
	}
	user, err := controller.createUser(&register)
	if err != nil {
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"error": true,
			"msg":   "Register failed",
			"user":  nil,
		})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error": false,
		"msg":   "Register successful",
		"user":  user,
	})
}

func (controller *SessionController) createUser(user *models.UserRegister) (*models.TokenInfo, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("cannot hash password: %w", err)
	}

	userdb := &models.UserDB{
		Id:             user.Id,
		HashedPassword: string(hashedPassword),
		Role:           user.Role,
	}
	return controller.storeUser(userdb)
}

func (controller *SessionController) storeUser(user *models.UserDB) (*models.TokenInfo, error) {
	err := controller.repo.RegisterUser(user)
	if err != nil {
		return nil, fmt.Errorf("user cannot be created: %w", err)
	}
	result := &models.TokenInfo{
		Id:   user.Id,
		Role: user.Role,
	}
	return result, nil
}
