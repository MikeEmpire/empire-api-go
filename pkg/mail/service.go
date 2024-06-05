package mail

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// @Summary Authenticate Gmail Account
// @Tags Mail
// @Accept json
// @Success 200 {string} Success
// @Failure 400 {string} Error
// @Router /mail/auth [get]
func AuthenticateGmailAccount(c *gin.Context) {

	ctx := context.Background()
	b, err := os.ReadFile("creds.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailReadonlyScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	client := getClient(config)

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))

	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	usr := "me"
	r, err := srv.Users.Labels.List(usr).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve labels: %v", err)
	}

	if len(r.Labels) == 0 {
		fmt.Println("No labels found.")
		return
	}

	c.JSON(http.StatusOK, gin.H{"labels": r.Labels})
}

// @Summary Authorization callback to get token from Gmail API
// @Tags Mail
// @Accept json
// @Success 200 {string} token
// @Failure 400 {string} Error
// @Router /mail/auth/callback [get]
func AuthCallback(c *gin.Context) {
	queryParams := c.Request.URL.Query()
	var token string
	for key, values := range queryParams {
		if key == "code" {
			for _, value := range values {
				token = value
			}
		}
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success!", "token": token})
}
