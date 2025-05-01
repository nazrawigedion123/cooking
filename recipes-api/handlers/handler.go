package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/nazrawigedion123/cooking/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RecipesHandler struct {
	Collection *mongo.Collection
	Ctx        context.Context
}

func NewRecipesHandler(ctx context.Context, collection *mongo.Collection) *RecipesHandler {
	return &RecipesHandler{
		Collection: collection,
		Ctx:        ctx,
	}
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
func (h *RecipesHandler) NewRecipeHandler(c *gin.Context) {

	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()

	_, err := h.Collection.InsertOne(h.Ctx, recipe)
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
func (h *RecipesHandler) ListRecipesHandler(c *gin.Context) {
	curr, err := h.Collection.Find(h.Ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer curr.Close(h.Ctx)
	recipes := make([]models.Recipe, 0)

	for curr.Next(h.Ctx) {
		var recipe models.Recipe
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
func (h *RecipesHandler) UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe models.Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	objectId, _ := primitive.ObjectIDFromHex(id)
	_, err := h.Collection.UpdateOne(h.Ctx, bson.M{
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
func (h *RecipesHandler) DeleteRecipeHandler(c *gin.Context) {
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
	result, err := h.Collection.DeleteOne(h.Ctx, filter)
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

func (h *RecipesHandler) SearchRecipesHandler(c *gin.Context) {
	tag := c.Query("tag")

	cursor, err := h.Collection.Find(h.Ctx, bson.M{"tags": tag})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer cursor.Close(h.Ctx)

	var recipes []models.Recipe
	for cursor.Next(h.Ctx) {
		var recipe models.Recipe
		if err := cursor.Decode(&recipe); err == nil {
			recipes = append(recipes, recipe)
		}
	}
	c.JSON(http.StatusOK, recipes)
}
