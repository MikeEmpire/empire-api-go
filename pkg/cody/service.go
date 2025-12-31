package cody

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type NewsletterSignup struct {
	Email string `json:"email" binding:"required,email"`
}

// Subscriber model - add this to your models file if you have one
type Subscriber struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"uniqueIndex;not null"`
	CreatedAt time.Time
	IsActive  bool `gorm:"default:true"`
}

func HandleNewsletterSignup(c *gin.Context) {
	var input NewsletterSignup
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	email := strings.TrimSpace(strings.ToLower(input.Email))

	// 1. Save to database
	// subscriber := Subscriber{
	// 	Email:    email,
	// 	IsActive: true,
	// }

	// if err := db.Create(&subscriber).Error; err != nil {
	// 	if strings.Contains(err.Error(), "duplicate") {
	// 		c.JSON(http.StatusOK, gin.H{"message": "Already subscribed!"})
	// 		return
	// 	}
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to subscribe"})
	// 	return
	// }

	// 2. Send confirmation email to user (async)
	go sendConfirmationEmail(email)

	// 3. Notify yourself (async)
	go sendAdminNotification(email)

	c.JSON(http.StatusOK, gin.H{"message": "Success!"})
}

// Send confirmation to the user
func sendConfirmationEmail(userEmail string) {
	from := os.Getenv("TITAN_EMAIL") // info@cody.live
	pass := os.Getenv("TITAN_PASSWORD")
	smtpHost := "smtp.titan.email"
	smtpAddr := smtpHost + ":587"

	subject := "Thanks for signing up for Cody.live!"
	body := `
		<html>
		<body>
			<h2>Thanks for your interest in Cody.live!</h2>
			<p>We'll keep you updated on our launch and new features.</p>
			<p>Stay tuned!</p>
		</body>
		</html>
	`

	msg := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"MIME-Version: 1.0\r\n"+
			"Content-Type: text/html; charset=UTF-8\r\n"+
			"\r\n"+
			"%s",
		from, userEmail, subject, body,
	)

	auth := smtp.PlainAuth("", from, pass, smtpHost)
	err := smtp.SendMail(smtpAddr, auth, from, []string{userEmail}, []byte(msg))
	if err != nil {
		log.Printf("Failed to send confirmation to %s: %v", userEmail, err)
	}
}

// Notify yourself about new signup
func sendAdminNotification(userEmail string) {
	from := os.Getenv("TITAN_EMAIL") // info@cody.live
	pass := os.Getenv("TITAN_PASSWORD")
	smtpHost := "smtp.titan.email"
	smtpAddr := smtpHost + ":587"
	log.Printf("TITAN_EMAIL: %s", from)
	log.Printf("TITAN_PASSWORD: %s (length: %d)", "***REDACTED***", len(pass))

	subject := "New Cody.live Signup!"
	body := fmt.Sprintf("New subscriber: %s\n", userEmail)

	msg := fmt.Sprintf(
		"From: %s\r\n"+
			"To: %s\r\n"+
			"Subject: %s\r\n"+
			"\r\n"+
			"%s",
		from, from, subject, body,
	)

	auth := smtp.PlainAuth("", from, pass, smtpHost)
	err := smtp.SendMail(smtpAddr, auth, from, []string{from}, []byte(msg))
	if err != nil {
		log.Printf("Failed to send admin notification: %v", err)
	}
}
