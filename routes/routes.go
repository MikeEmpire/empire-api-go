package routes

import (
	"empire-api-go/pkg/mail"

	"github.com/gin-gonic/gin"
)

func SetupRouter(r *gin.Engine) {
	// Setup group routes
	v1 := r.Group("/api/v1")
	{
		v1.GET("/auth", mail.AuthenticateGmailAccount)
		v1.GET("/auth/callback", mail.AuthCallback)
		// Add more routes here
	}
}
