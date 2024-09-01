package main

import (
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/rs/xid"
	"log"
	"net/http"
	"os"
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

func main() {
	router := gin.Default()
	router.POST("/recipes", NewRecipeHandler)
	router.GET("/recipes", ListRecipesHandler)
	log.Fatal(router.Run())
}
