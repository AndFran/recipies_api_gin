package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
	"recipies_api_gin/handlers"
)

var recipesHandler *handlers.RecipesHandler
var authHandler *handlers.AuthHandler

func init() {
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))
	if err != nil {
		log.Fatal(err)
	}
	if err = client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to MongoDB!")
	collectionRecipes := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipes")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", //os.Getenv("REDIS_URI"),
		Password: "",               //os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	status := redisClient.Ping(context.Background())
	log.Println(status)

	recipesHandler = handlers.NewRecipesHandler(ctx, collectionRecipes, redisClient)

	collectionUsers := client.Database(os.Getenv("MONGO_DATABASE")).Collection("users")
	authHandler = handlers.NewAuthHandler(collectionUsers, ctx)
}

func main() {
	router := gin.Default()

	authorized := router.Group("/")
	authorized.Use(AuthJWTAuthorizationMiddleware())

	// public access
	{
		router.POST("/signin", authHandler.SignInHandler)
		router.POST("/signup", authHandler.SignUpHandler)
		router.GET("/recipes", recipesHandler.ListRecipesHandler)
	}

	// API key protected
	{
		authorized.GET("/recipes/:id", recipesHandler.GetOneRecipeHandler)
		authorized.POST("/recipes", recipesHandler.NewRecipeHandler)
		authorized.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
		authorized.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
		//TODO maybe this should not need to be authed to use
		authorized.POST("/refresh", authHandler.RefreshJWTHandler)
	}

	log.Fatal(router.Run())
}
