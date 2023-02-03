package auth

import (
	"time"

	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nesarptr/chat-app-go/config"
	"github.com/nesarptr/chat-app-go/models"
	"github.com/nesarptr/chat-app-go/utils"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *fiber.Ctx) error {
	user := new(models.Client)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errors := utils.ValidateStruct(*user)

	if errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	db := config.GetDB()
	foundUser := new(models.Client)

	db.Where("user_name = ?", user.UserName).First(foundUser)

	if foundUser.ID != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "user already exist with this username",
		})
	}

	if err := user.Create(config.GetDB()); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":       user.ID,
		"username": user.UserName,
	})
}

func SignIn(c *fiber.Ctx) error {
	type loginInput struct {
		UserName string `json:"username" validate:"required,min=3,max=32"`
		Password string `json:"password" validate:"required,min=6"`
	}

	input := new(loginInput)

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errors := utils.ValidateStruct(*input)

	if errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	user := new(models.Client)

	config.GetDB().Where("user_name = ?", input.UserName).First(user)

	if user.ID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "user does not exist with this username",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "password is not valid",
		})
	}

	claims := jwt.MapClaims{
		"id":       user.ID,
		"username": user.UserName,
		"exp":      time.Now().Add(time.Minute * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := config.GetEnv("TOKEN_SECRET")

	t, err := token.SignedString([]byte(secret))

	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Success login", "token": t})
}

func Jwt(c *fiber.Ctx) error {
	userId := c.Locals("userId")
	userName := c.Locals("username")
	claims := jwt.MapClaims{
		"id":       userId,
		"username": userName,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := config.GetEnv("TOKEN_SECRET")

	t, err := token.SignedString([]byte(secret))

	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Success jwt renew", "token": t})
}

func Protected() []fiber.Handler {
	secret := config.GetEnv("TOKEN_SECRET")
	return []fiber.Handler{jwtware.New(jwtware.Config{
		SigningKey:   []byte(secret),
		ErrorHandler: jwtError,
	}), func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		foundUser := new(models.Client)
		db := config.GetDB()
		db.First(foundUser, claims["id"].(float64))
		if foundUser.ID == 0 {
			return fiber.ErrUnauthorized
		}
		c.Locals("username", claims["username"].(string))
		c.Locals("userId", claims["id"].(float64))
		return c.Next()
	}}
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}
