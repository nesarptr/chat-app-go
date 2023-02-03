package main

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/helmet/v2"
	"github.com/nesarptr/chat-app-go/auth"
	"github.com/nesarptr/chat-app-go/config"
	"github.com/nesarptr/chat-app-go/controllers"
	"github.com/nesarptr/chat-app-go/models"
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

	app.Post("/signup", auth.SignUp)
	app.Post("/signin", auth.SignIn)
	protected := app.Group("/", auth.Protected()...)
	protected.Get("/jwt", auth.Jwt)
	protected.Get("/users", controllers.GetUsers)
	protected.Get("/:from", controllers.GetMessages)

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
