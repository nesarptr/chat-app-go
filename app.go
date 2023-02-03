package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/helmet/v2"
	"github.com/nesarptr/chat-app-go/config"
)

func main() {

	port := config.GetEnv("PORT")

	if port == "" {
		port = "4000"
	}

	app := fiber.New()

	app.Use(cors.New())
	app.Use(helmet.New())
	app.Use(compress.New())

	app.Use("/", func(c *fiber.Ctx) error {
		return fiber.ErrNotFound
	})
	fmt.Println(app.Listen(":" + port))
}

func init() {
	err := config.Connect()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Database successfully connected!")
	}
}