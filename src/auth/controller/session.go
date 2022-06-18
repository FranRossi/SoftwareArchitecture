package controller

import(
	"auth/models"
	"github.com/gofiber/fiber/v2"
	"auth/repository"
	"encoding/json"
	"log"
	"auth/jwt"
	"fmt"
	"golang.org/x/crypto/bcrypt"

)

type SessionController struct {
	repo *repository.UsersRepo 
}

func NewSessionController(repo *repository.UsersRepo) *SessionController {
	return &SessionController{repo: repo}
}

func  (controller *SessionController) Login(c *fiber.Ctx) error{
	var login models.Login
	err := json.Unmarshal(c.Body(), &login)
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"error":   true,
			"msg":     err.Error(),
			"login": nil,
		})
	}

	user, err:= controller.repo.FindUser(login.Id)
	if err != nil {
		return c.Status(fiber.ErrNotFound.Code).JSON(fiber.Map{
			"error":   true,
			"msg":     err.Error(),
			"login": nil,
		})
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(login.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
			"error":   true,
			"msg":     err.Error(),
			"login": nil,
		})
	}
	if (user.HashedPassword == string(hashedPassword)){
		generator := &models.TokenInfo{
			Id: user.Id,
			Role: user.Role,
		}
		manager:= jwt.NewJWTManager(privateKey, time.now() )
		
		token, err  := manager.Generate(*generator)
		if err != nil {
			return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
				"error":   true,
				"msg":     err.Error(),
				"login": nil,
			})
		}
		result := &models.ResponseToken{Token:token,}
		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"error":   false,
			"msg":     "Login succesful",
			"token": 	result,
		})
	}

	return c.Status(fiber.ErrBadRequest.Code).JSON(fiber.Map{
		"error":   true,
		"msg":     err.Error(),
		"login": nil,
	})
}

func (controller *SessionController) RegisterUser(c *fiber.Ctx) (*models.TokenInfo, error) {
	var register models.UserRegister
	user, err := controller.createUser(&register)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return user, err
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
		Id: user.Id,
		Role: user.Role,
	}
	return result, nil
}
