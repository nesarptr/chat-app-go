package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nesarptr/chat-app-go/config"
	"github.com/nesarptr/chat-app-go/models"
)

func GetUsers(c *fiber.Ctx) error {
	userId := c.Locals("userId")
	db := config.GetDB()
	users := make([]models.Client, 0)
	db.Where("id != ?", userId).Find(&users)
	return c.Status(fiber.StatusOK).JSON(users)
}
