package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/georgmoe/repetition-tracker/models"
	"github.com/georgmoe/repetition-tracker/responses"
	"github.com/gofiber/fiber/v2"
)

func Register(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var user models.User
	defer cancel()

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Message: "error", Data: err.Error()})
	}

	user.CreatedAt = time.Now()
	result, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Message: "error", Data: err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(responses.Response{Message: "success", Data: result})
}
