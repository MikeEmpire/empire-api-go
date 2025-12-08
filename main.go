package main

import (
	"empire-api-go/config"
	"empire-api-go/pkg/esp32"
	"empire-api-go/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	config.LoadConfig()
	esp32Service, err := esp32.NewService("./data/sensors.db")
	if err != nil {
		panic("Failed to initialize ESP32 service: " + err.Error())
	}
	esp32Handlers := esp32.NewHandlers(esp32Service)

	// Initialize the router
	r := gin.Default()

	// Initialize CORS middleware
	r.Use(cors.New(config.CORSConfig()))

	// Setup routes
	routes.SetupRouter(r, esp32Handlers)

	// Start the server
	r.Run(":8080") // Default listens and serves on 0.0.0.0:8080
}
