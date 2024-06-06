package mail

import (
	"context"
	"fmt"
	"net/http"
	"net/smtp"
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
		c.JSON(http.StatusNotFound, gin.H{"error": "Unable to find credentials"})
		return
	}

	config, err := google.ConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Unable to parse client secret file to config: %v", err)})
		return
	}

	client, clientErr := getClient(config)
	if clientErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Unable to retrieve Gmail client: %v", clientErr)})
		return
	}

	srv, err := gmail.NewService(ctx, option.WithHTTPClient(client))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Unable to retrieve Gmail client: %v", err)})
		return
	}

	usr := "me"
	r, err := srv.Users.Labels.List(usr).Do()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Unable to retrieve labels: %v", err)})
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

func TestSendEmail(c *gin.Context) {
	from := os.Getenv("EMAIL_USERNAME")
	pass := os.Getenv("EMAIL_PASSWORD")
	to := "aolie1794@gmail.com"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Hello there\n\n" +
		"Helllooooo"

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Success!"})
}
