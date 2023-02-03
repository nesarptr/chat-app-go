package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/helmet/v2"
)

func main() {
	app := fiber.New()

	app.Use(cors.New())
	app.Use(helmet.New())
	app.Use(compress.New())
}