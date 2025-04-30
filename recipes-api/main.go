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
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/nazrawigedion123/cooking/docs"
	"github.com/rs/xid"

	swaggerFiles "github.com/swaggo/files"

	ginSwagger "github.com/swaggo/gin-swagger"
)

// Recipe represents a cooking recipe
type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishedAt  time.Time `json:"time"`
}

var recipes []Recipe

func init() {
	recipes = make([]Recipe, 0)
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
	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)
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
	id := c.Params.ByName("id")
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}
	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}
	recipe.ID = id
	recipes[index] = recipe
	c.JSON(http.StatusOK, recipe)
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
	id := c.Params.ByName("id")
	index := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}
	if index == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Recipe not found"})
		return
	}
	recipes = append(recipes[:index], recipes[index+1:]...)
	c.JSON(http.StatusOK, gin.H{"message": "recipe " + id + " has been deleted"})
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
