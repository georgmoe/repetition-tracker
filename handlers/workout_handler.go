package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/georgmoe/repetition-tracker/configs"
	"github.com/georgmoe/repetition-tracker/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	// retrieve all the users workouts
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

func GetSingleWorkout(c *fiber.Ctx) error {
	var workout models.Workout

	// get primitive user id
	userIdFromLocals := c.Locals(USER_ID)
	userId, err := GetPrimitiveObjectIDFromInterface(userIdFromLocals)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": err.Error()})
	}

	// get workout id from path
	workoutIdStr := c.Params("workoutId")
	workoutId, err := primitive.ObjectIDFromHex(workoutIdStr)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": err.Error()})
	}

	// retrieve workout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	// userId important to not get other users workouts!
	// filter := bson.D{{Key: "_id", Value: workoutId}, {Key: "userId", Value: userId}}
	filter := bson.M{"_id": workoutId, "userId": userId}
	// filter := bson.D{
	// 	{Key: "$and",
	// 		Value: bson.A{
	// 			bson.D{{Key: "_id", Value: workoutId}},
	// 			bson.D{{Key: "userId", Value: userId}},
	// 		},
	// 	},
	// }

	err = workoutCollection.FindOne(ctx, filter).Decode(&workout)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(
		fiber.Map{"message": "success", "data": workout},
	)
}
