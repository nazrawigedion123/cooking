// NewRecipeHandler godoc
// @Summary Add a new recipe
// @Description Create a new recipe and add it to the list
// @Tags recipes
// @Accept  json
// @Produce  json
// @Param   recipe body     Recipe true "Recipe to create"
// @Success 200    {object} Recipe
// @Failure 400    {object} gin.H
// @Router /recipes [post]
package main

import (
	"context"
	"encoding/json"

	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/nazrawigedion123/cooking/docs"

	swaggerFiles "github.com/swaggo/files"

	ginSwagger "github.com/swaggo/gin-swagger"

	// for mangodb

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	// redis
	"github.com/go-redis/redis/v8"

	//models and handlers

	"github.com/nazrawigedion123/cooking/handlers"
	"github.com/nazrawigedion123/cooking/models"
)

var recipes []models.Recipe

// data persistance
var ctx context.Context
var err error
var client *mongo.Client
var collection *mongo.Collection
var recipesHandler *handlers.RecipesHandler

func init() {
	recipes = make([]models.Recipe, 0)
	// read from a file
	file, _ := os.ReadFile("recipes.json")
	json.Unmarshal(file, &recipes)
	ctx = context.Background()
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGO_URI")))

	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collection = client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipies")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	recipesHandler = handlers.NewRecipesHandler(ctx, collection, redisClient)

}

// @title Cooking API
// @version 1.0
// @description This is a simple REST API to manage recipes using Gin and Swagger.
// @host localhost:8080
// @BasePath /
func main() {
	// constants

	router := gin.Default()

	// Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Recipe routes
	router.POST("/recipes", recipesHandler.NewRecipeHandler)
	router.GET("/recipes", recipesHandler.ListRecipesHandler)
	router.PUT("/recipes/:id", recipesHandler.UpdateRecipeHandler)
	router.DELETE("/recipes/:id", recipesHandler.DeleteRecipeHandler)
	router.GET("/recipes/search", recipesHandler.SearchRecipesHandler)

	router.Run()
}
