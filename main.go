package main

import (
	"context"
	"time"

	"github.com/georgmoe/repetition-tracker/configs"
	"github.com/georgmoe/repetition-tracker/models"
	"github.com/georgmoe/repetition-tracker/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

func main() {
	app := fiber.New()

	// registration route before auth middleware
	routes.RegistrationRouter(app)

	app.Use(basicauth.New(basicauth.Config{
		Authorizer: func(username, password string) bool {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			var user models.User
			filter := bson.M{"name": username, "password": password}
			err := userCollection.FindOne(ctx, filter).Decode(&user)
			if err != nil {
				if err == mongo.ErrNoDocuments {
					// This error means your query did not match any documents.
					return false
				}
				panic(err)
			}
			return true
		},
	}))

	//run database
	// configs.ConnectDB()

	//routes
	routes.UserRouter(app)
	routes.WorkoutRouter(app)

	app.Listen(":8080")
}
