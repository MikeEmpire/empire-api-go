package main

import (
	"empire-api-go/config"
	"empire-api-go/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize the router
	r := gin.Default()

	// Initialize CORS middleware
	r.Use(cors.New(config.CORSConfig()))

	// Setup routes
	routes.SetupRouter(r)

	// Start the server
	r.Run(":8080") // Default listens and serves on 0.0.0.0:8080
}
