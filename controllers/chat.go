package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nesarptr/chat-app-go/config"
	"github.com/nesarptr/chat-app-go/models"
)

func GetUsers(c *fiber.Ctx) error {
	userId := c.Locals("userId")
	username := c.Locals("username")
	db := config.GetDB()
	users := make([]models.Client, 0)
	db.Not("id = ?", userId).Find(&users)
	messages := make([]models.Text, 0)
	db.Where("sender_name = ?", username).Or("receiver_name = ?", username).Order("created_at desc").Find(&messages)
	for i, user := range users {
		for _, message := range messages {
			if message.SenderName == user.UserName || message.ReceiverName == user.UserName {
				users[i].LatestText = message.Body
				break
			}
		}
	}
	return c.Status(fiber.StatusOK).JSON(users)
}

func GetMessages(c *fiber.Ctx) error {
	username := c.Locals("username").(string)
	from := c.Params("from")
	user := new(models.Client)
	db := config.GetDB()
	db.Where("user_name = ?", from).First(user)
	if user.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "invalid username"})
	}
	usernames := []string{username, from}
	messages := make([]models.Text, 0)
	db.Where("sender_name IN ? AND receiver_name IN ?", usernames, usernames).Order("created_at desc").Find(&messages)
	return c.Status(fiber.StatusOK).JSON(messages)
}
