package handlers

import (
	"github.com/gofiber/fiber/v2"
)

func sessionAuth(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "not authorized",
		})
	}

	userId := sess.Get(USER_ID)
	if userId != nil {
		c.Locals(USER_ID, userId)
		return c.Next()
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message": "not authorized",
	})
}
