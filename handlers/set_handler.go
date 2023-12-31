package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/georgmoe/repetition-tracker/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PostSet(c *fiber.Ctx) error {
	var set models.Set

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
	if err := c.BodyParser(&set); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "error", "data": err.Error()})
	}

	// push exercise to array in workout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": workoutId, "userId": userId}
	update := bson.M{
		"$push": bson.M{fmt.Sprintf("exercises.%v.sets", exerciseIdxStr): set},
	}
	result, err := workoutCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "success", "data": result})
}

func PutSet(c *fiber.Ctx) error {
	var set models.Set

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

	// get set index from path and check type integer
	setIdxStr := c.Params("setIdx")
	_, err = strconv.Atoi(setIdxStr)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": "Exercise index must be of type integer!"})
	}

	// parse request body
	if err := c.BodyParser(&set); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "error", "data": err.Error()})
	}

	// update exercise at exerciseIdx
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": workoutId, "userId": userId}
	update := bson.M{
		"$set": bson.M{fmt.Sprintf("exercises.%v.sets.%v", exerciseIdxStr, setIdxStr): set},
	}
	result, err := workoutCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "success", "data": result})
}

func GetSet(c *fiber.Ctx) error {
	var set models.Set

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

	// get set index from path and check type integer
	setIdxStr := c.Params("setIdx")
	setIdx, err := strconv.Atoi(setIdxStr)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": "Exercise index must be of type integer!"})
	}

	// pipeline for an aggregation operation to find the exercise at the specified index
	pipeline := bson.A{
		// 1st stage: find the user's workout document
		bson.M{
			"$match": bson.M{
				"_id":    workoutId,
				"userId": userId,
			},
		},
		// 2nd stage: use projection to get the specified exercise in a new field "exercise"
		bson.M{
			"$project": bson.M{
				"exercise": bson.M{
					"$arrayElemAt": bson.A{"$exercises", exerciseIdx},
				},
			},
		},
		// 3rd stage: use projection to get the specified set in a new field "set"
		bson.M{
			"$project": bson.M{
				"set": bson.M{
					"$arrayElemAt": bson.A{"$exercise.sets", setIdx},
				},
			},
		},
		// 4th stage: make set the root to not have the unnecessary "set" key
		bson.M{
			"$replaceRoot": bson.M{
				"newRoot": "$set",
			},
		},
	}

	// execute aggregation operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := workoutCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": err.Error()})
	}
	defer cursor.Close(ctx)

	// decode exercise
	if cursor.Next(ctx) {
		if err := cursor.Decode(&set); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": err.Error()})
		}
	}

	// check if really only one result come from the aggregation pipeline
	if cursor.TryNext(ctx) {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": "query resulted in multiple sets"})
	}

	return c.Status(http.StatusOK).JSON(
		fiber.Map{"message": "success", "data": set},
	)
}
