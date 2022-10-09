package main

import (
	"stanza-api/src/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	// Routes
	routes.UserRoute(router)
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	router.Run("localhost:6000")
}
