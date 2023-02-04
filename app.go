package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/helmet/v2"
	"github.com/nesarptr/chat-app-go/config"
	"github.com/nesarptr/chat-app-go/models"
	"github.com/nesarptr/chat-app-go/routes"
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
	routes.SetUpRoutes(app)

	fmt.Println(app.Listen(":" + port))
}

func init() {
	err := config.Connect()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		db := config.GetDB()
		db.AutoMigrate(&models.Client{}, &models.Text{})
		fmt.Println("Database successfully connected!")
	}
}
