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

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

func GetAllUsers(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var users []models.User

	cursor, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": err.Error()})
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var singleUser models.User
		if err = cursor.Decode(&singleUser); err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": err.Error()})
		}

		users = append(users, singleUser)
	}

	return c.Status(http.StatusOK).JSON(
		fiber.Map{"message": "success", "data": users},
	)
}

// func GetAUser(c *fiber.Ctx) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	userId := c.Params("userId")
// 	var user models.User
// 	defer cancel()

// 	objId, _ := primitive.ObjectIDFromHex(userId)

// 	err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": err.Error()})
// 	}

// 	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "success", "data": user})
// }

// func CreateUser(c *fiber.Ctx) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	var user models.User
// 	defer cancel()

// 	if err := c.BodyParser(&user); err != nil {
// 		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "error", "data": err.Error()})
// 	}

// 	user.CreatedAt = time.Now()
// 	result, err := userCollection.InsertOne(ctx, user)
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": err.Error()})
// 	}

// 	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "success", "data": result})
// }

// func EditAUser(c *fiber.Ctx) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	userId := c.Params("userId")
// 	var user models.User
// 	defer cancel()

// 	objId, _ := primitive.ObjectIDFromHex(userId)

// 	//validate the request body
// 	if err := c.BodyParser(&user); err != nil {
// 		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "error", "data": err.Error()})
// 	}

// 	update := bson.M{"name": user.Name, "location": user.Location, "title": user.Title}

// 	result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": err.Error()})
// 	}
// 	//get updated user details
// 	var updatedUser models.User
// 	if result.MatchedCount == 1 {
// 		err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)

// 		if err != nil {
// 			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": err.Error()})
// 		}
// 	}

// 	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "success", "data": updatedUser})
// }

// func DeleteAUser(c *fiber.Ctx) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	userId := c.Params("userId")
// 	defer cancel()

// 	objId, _ := primitive.ObjectIDFromHex(userId)

// 	result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})
// 	if err != nil {
// 		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": err.Error()})
// 	}

// 	if result.DeletedCount < 1 {
// 		return c.Status(http.StatusNotFound).JSON(
// 			fiber.Map{"message": "error", "data": "User with specified ID not found!"},
// 		)
// 	}

// 	return c.Status(http.StatusOK).JSON(
// 		fiber.Map{"message": "success", "data": "User successfully deleted!"},
// 	)
// }
