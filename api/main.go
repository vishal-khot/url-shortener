package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/vishal-khot/url-shortener/routes"
)

func main() {
	app := fiber.New()
	app.Use(logger.New())
	app.Static("/", "./static")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("./static/index.html")
	})
	app.Get("/:url", routes.ResolveURL)
	app.Post("/shorten", routes.ShortenURL)

	log.Fatal(app.Listen(":8080"))
}
