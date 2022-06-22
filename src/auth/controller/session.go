package controller

import (
	"auth/jwt"
	"auth/models"
	"auth/repository"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	l "own_logger"
)

type SessionController struct {
	repo    *repository.UsersRepo
	manager *jwt.Manager
}

func NewSessionController(repo *repository.UsersRepo, manager *jwt.Manager) *SessionController {
	return &SessionController{
		repo:    repo,
		manager: manager,
	}
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
	message := "Invalid credentials"
	user, err := controller.repo.FindUser(login.Id)
	if err != nil {
		l.LogError(message)
		return c.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"error": true,
			"msg":   message,
			"login": nil,
		})
	}
	err2 := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(login.Password))
	if err2 != nil {
		l.LogWarning(message)
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"error": true,
			"msg":   message,
			"login": nil,
		})
	}
	generator := models.TokenInfo{
		Id:   user.Id,
		Role: user.Role,
	}

	token, err := controller.manager.Generate(generator)
	if err != nil {
		l.LogError("Token cannot be generated")
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"error": true,
			"msg":   err.Error() + " cannot generate token",
			"login": nil,
		})
	}
	l.LogInfo("Login successful")
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
