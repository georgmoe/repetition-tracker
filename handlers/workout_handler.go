package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/georgmoe/repetition-tracker/configs"
	"github.com/georgmoe/repetition-tracker/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var workoutCollection *mongo.Collection = configs.GetCollection(configs.DB, "workout")

func CreateWorkout(c *fiber.Ctx) error {
	var workout models.Workout

	// get primitive user id
	userIdFromLocals := c.Locals(USER_ID)
	userId, err := GetPrimitiveObjectIDFromInterface(userIdFromLocals)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": err.Error()})
	}

	// parse request body
	if err := c.BodyParser(&workout); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "error", "data": err.Error()})
	}
	workout.UserId = userId

	// insert new workout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := workoutCollection.InsertOne(ctx, workout)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "success", "data": result})
}

func GetAllWorkouts(c *fiber.Ctx) error {
	var workouts []models.Workout

	// get primitive user id
	userIdFromLocals := c.Locals(USER_ID)
	userId, err := GetPrimitiveObjectIDFromInterface(userIdFromLocals)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": err.Error()})
	}

	// retrieve all workouts for the user
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{Key: "userId", Value: userId}}
	cursor, err := workoutCollection.Find(ctx, filter)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": err.Error()})
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var singleWorkout models.Workout
		if err = cursor.Decode(&singleWorkout); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": err.Error()})
		}

		workouts = append(workouts, singleWorkout)
	}

	return c.Status(http.StatusOK).JSON(
		fiber.Map{"message": "success", "data": workouts},
	)
}
