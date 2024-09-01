package main

import (
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/rs/xid"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Recipe struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishAt    time.Time `json:"publishAt"`
}

var recipes []Recipe

func init() {

	recipeData, err := os.ReadFile("recipes.json")
	if err != nil {
		log.Fatal(err)
	}
	err = json.Unmarshal(recipeData, &recipes)
	if err != nil {
		log.Fatal(err)
	}
}

func NewRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	recipe.ID = xid.New().String()
	recipe.PublishAt = time.Now()
	recipes = append(recipes, recipe) // append to DB
	c.JSON(http.StatusOK, recipe)
}

func ListRecipesHandler(c *gin.Context) {
	c.JSON(http.StatusOK, recipes)
}

func UpdateRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	var recipe Recipe
	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	foundIndex := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			foundIndex = i
		}
	}
	if foundIndex == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
		return
	}
	recipes[foundIndex] = recipe
	c.JSON(http.StatusOK, recipe)
}

func DeleteRecipeHandler(c *gin.Context) {
	id := c.Param("id")
	foundIndex := -1
	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			foundIndex = i
		}
	}
	if foundIndex == -1 {
		c.JSON(http.StatusNotFound, gin.H{"error": "recipe not found"})
		return
	}
	recipes = append(recipes[:foundIndex], recipes[foundIndex+1:]...)
	c.JSON(http.StatusOK, gin.H{"success": "recipe deleted"})
}

func SearchRecipeHandler(c *gin.Context) {
	tag := c.Query("tag")
	results := make([]Recipe, 0)
	for i := 0; i < len(recipes); i++ {
		found := false
		for _, t := range recipes[i].Tags {

			if strings.EqualFold(t, tag) {
				found = true
			}
		}
		if found {
			results = append(results, recipes[i])
		}
	}
	c.JSON(http.StatusOK, results)
}

func main() {
	router := gin.Default()
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	router.PUT("/recipes/:id", UpdateRecipeHandler)
	router.DELETE("/recipes/:id", DeleteRecipeHandler)
	router.GET("/recipes/search", SearchRecipeHandler)
	log.Fatal(router.Run())
}
