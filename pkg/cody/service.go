package cody

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
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

func sendConfirmationEmail(userEmail string) {
	from := mail.NewEmail("Cody", "info@cody.live")
	subject := "Thanks for signing up for Cody.live!"
	to := mail.NewEmail("", userEmail)
	htmlContent := `
		<html>
		<body>
			<h2>Thanks for your interest in Cody.live!</h2>
			<p>We'll keep you updated on our launch and new features.</p>
		</body>
		</html>
	`
	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Printf("Failed to send email: %v", err)
	} else {
		log.Printf("Email sent! Status: %d", response.StatusCode)
	}
}

// Send confirmation to the user
// func sendConfirmationEmail(userEmail string) {
// 	from := os.Getenv("TITAN_EMAIL") // info@cody.live
// 	pass := os.Getenv("TITAN_PASSWORD")
// 	smtpHost := "smtp.titan.email"
// 	smtpAddr := smtpHost + ":587"

// 	subject := "Thanks for signing up for Cody.live!"
// 	body := `
// 		<html>
// 		<body>
// 			<h2>Thanks for your interest in Cody.live!</h2>
// 			<p>We'll keep you updated on our launch and new features.</p>
// 			<p>Stay tuned!</p>
// 		</body>
// 		</html>
// 	`

// 	msg := fmt.Sprintf(
// 		"From: %s\r\n"+
// 			"To: %s\r\n"+
// 			"Subject: %s\r\n"+
// 			"MIME-Version: 1.0\r\n"+
// 			"Content-Type: text/html; charset=UTF-8\r\n"+
// 			"\r\n"+
// 			"%s",
// 		from, userEmail, subject, body,
// 	)

// 	auth := smtp.PlainAuth("", from, pass, smtpHost)
// 	err := smtp.SendMail(smtpAddr, auth, from, []string{userEmail}, []byte(msg))
// 	if err != nil {
// 		log.Printf("Failed to send confirmation to %s: %v", userEmail, err)
// 	}
// }

// Notify yourself about new signup
func sendAdminNotification(userEmail string) {
	from := mail.NewEmail("Cody", "info@cody.live")
	subject := "New Cody.live Signup!"
	to := mail.NewEmail("", "info@cody.live") // Send to yourself

	htmlContent := `
		<html>
		<body>
			<h2>New Newsletter Signup!</h2>
			<p><strong>Email:</strong> ` + userEmail + `</p>
		</body>
		</html>
	`

	message := mail.NewSingleEmail(from, subject, to, "", htmlContent)

	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		log.Printf("Failed to send admin notification: %v", err)
	} else {
		log.Printf("Admin notification sent! Status: %d", response.StatusCode)
	}
}

// Notify yourself about new signup
// func sendAdminNotification(userEmail string) {
// 	from := os.Getenv("TITAN_EMAIL") // info@cody.live
// 	pass := os.Getenv("TITAN_PASSWORD")
// 	smtpHost := "smtp.titan.email"
// 	smtpAddr := smtpHost + ":587"
// 	log.Printf("TITAN_EMAIL: %s", from)
// 	log.Printf("TITAN_PASSWORD: %s (length: %d)", "***REDACTED***", len(pass))

// 	subject := "New Cody.live Signup!"
// 	body := fmt.Sprintf("New subscriber: %s\n", userEmail)

// 	msg := fmt.Sprintf(
// 		"From: %s\r\n"+
// 			"To: %s\r\n"+
// 			"Subject: %s\r\n"+
// 			"\r\n"+
// 			"%s",
// 		from, from, subject, body,
// 	)

// 	auth := smtp.PlainAuth("", from, pass, smtpHost)
// 	err := smtp.SendMail(smtpAddr, auth, from, []string{from}, []byte(msg))
// 	if err != nil {
// 		log.Printf("Failed to send admin notification: %v", err)
// 	}
// }
