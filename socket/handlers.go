package socket

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nesarptr/chat-app-go/config"
	"github.com/nesarptr/chat-app-go/models"
	"github.com/nesarptr/chat-app-go/utils"
)

func Upgrade(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		go RunHub()
		return c.Next()
	}
	return c.SendStatus(fiber.StatusUpgradeRequired)
}

func WSHandler(c *websocket.Conn) {
	jwtoken := c.Query("token")
	if jwtoken == "" {
		websocket.FormatCloseMessage(fiber.StatusUnauthorized, "Missing or malformed JWT")
		UnRegister <- c
	}
	token, err := jwt.Parse(jwtoken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.GetEnv("TOKEN_SECRET")), nil
	})
	if err != nil {
		websocket.FormatCloseMessage(fiber.StatusUnauthorized, "Invalid or expired JWT")
		c.Close()
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	user := new(models.Client)
	config.GetDB().First(user, claims["id"].(float64))
	if user.ID == 0 {
		websocket.FormatCloseMessage(fiber.StatusUnauthorized, "Invalid JWT, user does not exist")
		c.Close()
		return
	}
	Clients[c] = user
	defer func() {
		UnRegister <- c
		c.Close()
	}()
	Register <- c

	for {
		message := new(models.Text)
		message.SenderName = user.UserName
		if err := c.ReadJSON(message); err != nil {
			fmt.Println(err.Error())
			websocket.FormatCloseMessage(fiber.StatusBadRequest, "message is not supported")
			return
		}
		if errors := utils.ValidateStruct(message); errors != nil {
			websocket.FormatCloseMessage(fiber.StatusBadRequest, "message is not supported")
			return
		}
		db := config.GetDB()
		to := new(models.Client)
		db.Where("user_name = ?", message.ReceiverName).First(to)
		if to.ID == 0 {
			websocket.FormatCloseMessage(fiber.StatusBadRequest, "receiver is not a valid user")
			return
		}
		if err := message.Create(db); err != nil {
			websocket.FormatCloseMessage(fiber.StatusInternalServerError, "could not proccess the message")
			return
		}
		Broadcast <- message
	}

}
