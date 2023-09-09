package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/georgmoe/repetition-tracker/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PostExercise(c *fiber.Ctx) error {
	var exercise models.Exercise

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

	// parse request body
	if err := c.BodyParser(&exercise); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "error", "data": err.Error()})
	}

	// push exercise to array in workout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": workoutId, "userId": userId}
	update := bson.M{
		"$push": bson.M{"exercises": exercise},
	}
	result, err := workoutCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "success", "data": result})
}

func PutExercise(c *fiber.Ctx) error {
	var exercise models.Exercise

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

	// get exercise index from path and check type integer
	exerciseIdxStr := c.Params("exerciseIdx")
	_, err = strconv.Atoi(exerciseIdxStr)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": "Exercise index must be of type integer!"})
	}

	// parse request body
	if err := c.BodyParser(&exercise); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "error", "data": err.Error()})
	}

	// update exercise at exerciseIdx
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": workoutId, "userId": userId}
	update := bson.M{
		"$set": bson.M{"exercises." + exerciseIdxStr: exercise},
	}
	result, err := workoutCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "success", "data": result})
}

func GetExercise(c *fiber.Ctx) error {
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

	// get exercise index from path and check type integer
	exerciseIdxStr := c.Params("exerciseIdx")
	exerciseIdx, err := strconv.Atoi(exerciseIdxStr)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": "Exercise index must be of type integer!"})
	}

	// find workout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": workoutId, "userId": userId}
	err = workoutCollection.FindOne(ctx, filter).Decode(&workout)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": err.Error()})
	}

	// check valid exercise index and get exercise
	var exercise models.Exercise

	if exerciseIdx >= 0 && exerciseIdx < len(workout.Exercises) {
		exercise = workout.Exercises[exerciseIdx]
	} else {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": "exercise index out of range"})
	}

	return c.Status(http.StatusOK).JSON(
		fiber.Map{"message": "success", "data": exercise},
	)
}
