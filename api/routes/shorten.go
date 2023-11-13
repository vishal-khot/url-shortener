package routes

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"net"
	"net/url"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/vishal-khot/url-shortener/database"
)

type reqBody struct {
	UrlToShorten string `json:"urlToShorten"`
}

func isValidURL(tocheck string) (bool, error) {
	uri, err := url.ParseRequestURI(tocheck)
	if err != nil {
		return false, err
	}

	switch uri.Scheme {
	case "http":
	case "https":
	default:
		return false, errors.New("invalid scheme")
	}

	_, err = net.LookupHost(uri.Host)
	if err != nil {
		return false, err
	}

	return true, nil
}

func ShortenURL(ctx *fiber.Ctx) error {
	redisClient, error := database.CreateRedisClient(0)
	if error != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"error":    error,
			"shortUrl": "",
		})
	}
	defer redisClient.Close()

	baseUrl := ctx.BaseURL()
	body := reqBody{}
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error":    "could not parse url",
			"shortUrl": "",
		})
	}

	valid, _ := isValidURL(body.UrlToShorten)
	if valid {
		hash := sha1.New()
		hash.Write([]byte(body.UrlToShorten))
		shortenedUrl := baseUrl + "/" + base64.URLEncoding.EncodeToString(hash.Sum(nil))[:8]

		value, err := redisClient.Get(database.Ctx, shortenedUrl).Result()
		if err == redis.Nil {
			_, err := redisClient.Set(database.Ctx, shortenedUrl, body.UrlToShorten, 0).Result()
			if err != nil {
				return ctx.Status(fiber.StatusFound).JSON(&fiber.Map{
					"error":    err,
					"shortUrl": "",
				})
			} else {
				return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
					"error":       "",
					"description": "URL shortened",
					"shortUrl":    shortenedUrl,
				})
			}
		} else if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
				"error":    err,
				"shortUrl": value,
			})
		} else {
			return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
				"error":       "",
				"description": "URL already shortened",
				"shortUrl":    shortenedUrl,
			})
		}
	} else {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"error":    "not a valid url",
			"shortUrl": "",
		})
	}
}
