package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/georgmoe/repetition-tracker/models"
	"github.com/georgmoe/repetition-tracker/responses"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(c *fiber.Ctx) error {
	var user user
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Message: "error", Data: err.Error()})
	}

	password, bcErr := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if bcErr != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Message: "error", Data: bcErr.Error()})
	}

	newUser := models.User{
		Name:      user.Name,
		Email:     user.Email,
		Password:  string(password),
		CreatedAt: time.Now(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Message: "error", Data: err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(responses.Response{Message: "success", Data: result})
}

func Login(c *fiber.Ctx) error {
	var user user
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Message: "error", Data: err.Error()})
	}

	var loggedInUser models.User

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.D{{Key: "email", Value: user.Email}}
	err := userCollection.FindOne(ctx, filter).Decode(&loggedInUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return c.Status(http.StatusUnauthorized).JSON(responses.Response{Message: "error", Data: err.Error()})
		}
		panic(err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(loggedInUser.Password), []byte(user.Password))
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(responses.Response{Message: "error", Data: err.Error()})
	}

	// set session
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Message: "error", Data: err.Error()})
	}

	sess.Set(AUTH_KEY, true)
	sess.Set(USER_ID, loggedInUser.ID)

	sessErr := sess.Save()
	if sessErr != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Message: "error", Data: err.Error()})
	}

	return c.Status(http.StatusOK).JSON(responses.Response{Message: "successfully logged in"})
}

func Logout(c *fiber.Ctx) error {
	sess, err := store.Get(c)
	if err != nil {
		return c.Status(http.StatusOK).JSON(responses.Response{Message: "No session", Data: err.Error()})
	}

	err = sess.Destroy()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.Response{Message: "error", Data: err.Error()})
	}

	return c.Status(http.StatusOK).JSON(responses.Response{Message: "logged out successfully"})

}
