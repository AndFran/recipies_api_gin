package main

import (
	"github.com/gin-gonic/gin"
	"time"
)

type Recipe struct {
	Name         string    `json:"name"`
	Tags         []string  `json:"tags"`
	Ingredients  []string  `json:"ingredients"`
	Instructions []string  `json:"instructions"`
	PublishAt    time.Time `json:"publishAt"`
}

func main() {
	router := gin.Default()

	router.Run()
}
