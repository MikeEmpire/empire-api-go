package routes

import (
	"empire-api-go/pkg/cody"
	"empire-api-go/pkg/esp32"
	"empire-api-go/pkg/mail"
	"empire-api-go/pkg/ptrains"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine, esp32Handlers *esp32.Handlers) {
	// Default route
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to the Empire API!"})
	})
	// Setup group routes
	v1 := r.Group("/api/v1")
	{
		v1.GET("/auth", mail.AuthenticateGmailAccount)
		v1.GET("/auth/callback", mail.AuthCallback)
		v1.POST("/test", mail.TestSendEmail)

		// Cody.live API routes
		codyGroup := v1.Group("/cody")
		{
			codyGroup.POST("/newsletter-signup", cody.HandleNewsletterSignup)
		}

		// PTRAIN API routes
		ptrain := v1.Group("/ptrains")
		{
			ptrain.POST("/contact", ptrains.HandleContactForm)
		}

		// ESP32 API routes
		esp32Group := v1.Group("/esp32")
		{
			esp32Group.POST("/sensor-data", esp32Handlers.SaveSensorData)
			esp32Group.GET("/readings/:device_id", esp32Handlers.GetReadings)
			esp32Group.GET("/stats/:device_id", esp32Handlers.GetStats)
		}

	}
}
