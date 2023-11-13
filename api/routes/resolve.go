package routes

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/vishal-khot/url-shortener/database"
)

func ResolveURL(ctx *fiber.Ctx) error {
	redisClient, error := database.CreateRedisClient(0)
	if error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error": "could not connect to redis",
		})
	} else {
		defer redisClient.Close()
		url := ctx.BaseURL() + string(ctx.Request().RequestURI())
		value, err := redisClient.Get(database.Ctx, url).Result()
		if err == redis.Nil {
			return ctx.Status(fiber.StatusNotFound).JSON(&fiber.Map{
				"error": "url not found",
			})
		} else if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"error": "internal server error",
			})
		} else {
			return ctx.Redirect(value)
		}
	}
}
