package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/georgmoe/repetition-tracker/configs"
	"github.com/georgmoe/repetition-tracker/models"
	"github.com/georgmoe/repetition-tracker/responses"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var workoutCollection *mongo.Collection = configs.GetCollection(configs.DB, "workout")

func CreateWorkout(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	var workout models.Workout
	defer cancel()

	//validate the request body
	if err := c.BodyParser(&workout); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.Response{Message: "error", Data: err.Error()})
	}

	//use the validator library to validate required fields
	// if validationErr := validate.Struct(&user); validationErr != nil {
	// 	return c.Status(http.StatusBadRequest).JSON(responses.Response{Message: "error", Data: validationEr.Error()}})
	// }

	result, err := workoutCollection.InsertOne(ctx, workout)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Message: "error", Data: err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(responses.Response{Message: "success", Data: result})
}

func GetAllWorkouts(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var workouts []models.Workout

	// empty map because no filter
	cursor, err := workoutCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Message: "error", Data: err.Error()})
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var singleWorkout models.Workout
		if err = cursor.Decode(&singleWorkout); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(responses.Response{Message: "error", Data: err.Error()})
		}

		workouts = append(workouts, singleWorkout)
	}

	return c.Status(http.StatusOK).JSON(
		responses.Response{Message: "success", Data: workouts},
	)
}
