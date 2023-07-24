package responses

import "github.com/gofiber/fiber/v2"

type UserResponse struct {
	Message string     `json:"message"`
	Data    *fiber.Map `json:"data"`
}
