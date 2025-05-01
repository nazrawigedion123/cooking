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
	"fmt"

	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/nazrawigedion123/cooking/docs"

	swaggerFiles "github.com/swaggo/files"

	ginSwagger "github.com/swaggo/gin-swagger"

	// for mangodb
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Recipe represents a cooking recipe
//
//	type Recipe struct {
//		ID           string    `json:"id"`
//		Name         string    `json:"name"`
//		Tags         []string  `json:"tags"`
//		Ingredients  []string  `json:"ingredients"`
//		Instructions []string  `json:"instructions"`
//		PublishedAt  time.Time `json:"time"`
//	}
type Recipe struct {
	//swagger:ignore
	ID           primitive.ObjectID `json:"id" bson:"_id"`
	Name         string             `json:"name" bson:"name"`
	Tags         []string           `json:"tags" bson:"tags"`
	Ingredients  []string           `json:"ingredients" bson:"ingredients"`
	Instructions []string           `json:"instructions" bson:"instructions"`
	PublishedAt  time.Time          `json:"publishedAt" bson:"publishedAt"`
}

var recipes []Recipe

// data persistance
var ctx context.Context
var err error
var client *mongo.Client
var collection *mongo.Collection

func init() {
	recipes = make([]Recipe, 0)
	// read from a file
	file, _ := os.ReadFile("recipes.json")
	json.Unmarshal(file, &recipes)
	ctx = context.Background()
	client, err = mongo.Connect(ctx,
		options.Client().ApplyURI(os.Getenv("MONGO_URI")))

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB")
	collection = client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipies")

	// used on time to add data to our database
	// var listOfRecipes []interface{}

	// for _, recipe := range recipes {
	// 	listOfRecipes = append(listOfRecipes, recipe)
	// }
	// collection := client.Database(os.Getenv("MONGO_DATABASE")).Collection("recipies")
	// insertManyResult, err := collection.InsertMany(ctx, listOfRecipes)
	// if err != nil {
	// 	log.Fatal(err)

	// }
	// log.Println("Inserted recipies: ", len(insertManyResult.InsertedIDs))

}

// NewRecipeHandler godoc
// @Summary      Create a new recipe
// @Description  Add a new recipe to the list
// @Tags         recipes
// @Accept       json
// @Produce      json
// @Param        recipe  body      Recipe  true  "Recipe to create"
// @Success      200     {object}  Recipe
// @Failure      400     {object}  ErrorResponse
// @Router       /recipes [post]
func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()

	_, err = collection.InsertOne(ctx, recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error while inserting the recipe"})
		return
	}

	c.JSON(http.StatusOK, recipe)
}

// ListRecipesHandler godoc
// @Summary Get all recipes
// @Description Retrieve list of all recipes
// @Tags recipes
// @Produce json
// @Success 200 {array} Recipe
// @Router /recipes [get]
func ListRecipesHandler(c *gin.Context) {
	curr, err := collection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer curr.Close(ctx)
	recipes = make([]Recipe, 0)

	for curr.Next(ctx) {
		var recipe Recipe
		curr.Decode(&recipe)
		recipes = append(recipes, recipe)

	}

	c.JSON(http.StatusOK, recipes)
}

// UpdateRecipeHandler godoc
// @Summary Update a recipe
// @Description Update recipe by ID
// @Tags recipes
// @Accept json
// @Produce json
// @Param id path string true "Recipe ID"
// @Param recipe body Recipe true "Updated recipe data"
// @Success 200 {object} Recipe
// @Failure 400 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /recipes/{id} [put]
func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err = collection.UpdateOne(ctx, bson.M{
		"_id": objectId},
		bson.D{{Key: "$set", Value: bson.D{
			{Key: "name", Value: recipe.Name},
			{Key: "instructions", Value: recipe.Instructions},
			{Key: "ingredients", Value: recipe.Ingredients},
			{Key: "tags", Value: recipe.Tags}}}})

	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Recipe has been updated"})

}

// DeleteRecipeHandler godoc
// @Summary Delete a recipe
// @Description Delete a recipe by ID
// @Tags recipes
// @Produce json
// @Param id path string true "Recipe ID"
// @Success 200 {object} gin.H
// @Failure 404 {object} gin.H
// @Router /recipes/{id} [delete]
func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")

	// Convert string to MongoDB ObjectID
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Filter by _id to find the specific document
	filter := bson.M{"_id": objectId}

	// Attempt to delete the document
	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if result.DeletedCount == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No recipe found with ID: " + id})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Recipe " + id + " has been deleted"})
}

// SearchRecipesHandler godoc
// @Summary Search recipes by tag
// @Description Filter recipes that contain a specific tag
// @Tags recipes
// @Produce json
// @Param tag query string true "Tag to filter by"
// @Success 200 {array} Recipe
// @Router /recipes/search [get]
func SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")
	listOfRecipes := make([]Recipe, 0)
	for i := 0; i < len(recipes); i++ {
		found := false
		for _, t := range recipes[i].Tags {
			if strings.EqualFold(t, tag) {
				found = true
			}
		}
		if found {
			listOfRecipes = append(listOfRecipes, recipes[i])
		}
	}
	c.JSON(http.StatusOK, listOfRecipes)
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
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.GET("/recipes/search", SearchRecipesHandler)

	router.Run()
}
