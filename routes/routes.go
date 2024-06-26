package routes

import (
	"empire-api-go/pkg/mail"
	"empire-api-go/pkg/ptrains"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
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

		// PTRAIN API routes
		ptrain := v1.Group("/ptrains")
		{
			ptrain.POST("/contact", ptrains.HandleContactForm)
		}
	}
}
