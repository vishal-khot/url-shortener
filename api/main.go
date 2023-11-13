package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/vishal-khot/url-shortener/routes"
)

func main() {
	app := fiber.New()
	app.Get("/:url", routes.ResolveURL)
	app.Post("/shorten", routes.ShortenURL)
	app.Use(logger.New())
	log.Fatal(app.Listen(":8080"))
}
