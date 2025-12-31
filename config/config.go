package config

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
)

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func CORSConfig() cors.Config {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{
		"http://localhost:3000",
		"http://127.0.0.1:5500", // Local dev
		"https://ptrainsbbq.com",
		"https://codycomingsoon.netlify.app", // Netlify preview
		"https://cody.live",                  // Main domain
		"https://www.cody.live",              // WWW variant
		"https://coming-soon.cody.live",      // Subdomain
	}
	corsConfig.AllowCredentials = true
	corsConfig.AddAllowHeaders("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers", "Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization")
	corsConfig.AddAllowMethods("GET", "POST", "PUT", "DELETE", "OPTIONS")
	return corsConfig
}
