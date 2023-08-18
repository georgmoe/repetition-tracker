package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func ReturnSessionAuth() func(c *fiber.Ctx) error {
	return SessionAuth
}

func SessionAuth(c *fiber.Ctx) error {
	sess, err := store.Get(c)

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "not authorized",
		})
	}

	if sess.Get(AUTH_KEY) != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "not authorized",
		})
	}

	return c.Next()
}
