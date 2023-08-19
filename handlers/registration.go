package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/georgmoe/repetition-tracker/models"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// User data from request
type user struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *fiber.Ctx) error {
	// parse request data
	var user user
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": err.Error()})
	}

	// hash password
	password, bcErr := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if bcErr != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": bcErr.Error()})
	}

	// create new model user
	newUser := models.User{
		ID:        primitive.NewObjectID(),
		Name:      user.Name,
		Email:     user.Email,
		Password:  string(password),
		CreatedAt: time.Now(),
	}

	// insert new model user
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"message": "success", "data": result})
}

func Login(c *fiber.Ctx) error {
	// parse request data
	var user user
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": err.Error()})
	}

	// get existing model user from db
	var loggedInUser models.User

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{Key: "email", Value: user.Email}}
	err := userCollection.FindOne(ctx, filter).Decode(&loggedInUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "error", "data": err.Error()})
		}
		panic(err)
	}

	// compare password and hash
	err = bcrypt.CompareHashAndPassword([]byte(loggedInUser.Password), []byte(user.Password))
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"message": "error", "data": err.Error()})
	}

	// get session from context
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": err.Error()})
	}

	// set session id
	userId, err := loggedInUser.ID.MarshalJSON()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": err.Error()})
	}
	sess.Set(USER_ID, userId)

	// save session
	sessErr := sess.Save()
	if sessErr != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": sessErr.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "logged in"})
}

func Logout(c *fiber.Ctx) error {
	// get session
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(http.StatusOK).JSON(fiber.Map{"message": "No session", "data": err.Error()})
	}

	// destroy session
	err = sess.Destroy()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "error", "data": err.Error()})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "logged out"})
}
