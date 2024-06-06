package main

import (
	"empire-api-go/config"
	"empire-api-go/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize the router
	r := gin.Default()

	// Setup routes
	routes.SetupRouter(r)

	// Start the server
	r.Run(":8080") // Default listens and serves on 0.0.0.0:8080
}
