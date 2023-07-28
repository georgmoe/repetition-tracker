package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/georgmoe/repetition-tracker/configs"
	"github.com/georgmoe/repetition-tracker/models"
	"github.com/georgmoe/repetition-tracker/responses"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

var workoutCollection *mongo.Collection = configs.GetCollection(configs.DB, "workout")

func CreateWorkout(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var workout models.Workout
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&workout); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Message: "error", Data: fiber.Map{"data": err.Error()}})
	}

	//use the validator library to validate required fields
	// if validationErr := validate.Struct(&user); validationErr != nil {
	// 	return c.Status(http.StatusBadRequest).JSON(responses.Response{Message: "error", Data: fiber.Map{"data": validationErr.Error()}})
	// }

	result, err := workoutCollection.InsertOne(ctx, workout)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Message: "error", Data: fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.Response{Message: "success", Data: fiber.Map{"data": result}})
}
