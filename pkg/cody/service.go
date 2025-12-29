package cody

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type NewsletterSignup struct {
	Email string `json:"email" binding:"required,email"`
}

func HandleNewsletterSignup(c *gin.Context) {
	var input NewsletterSignup
	if err := c.ShouldBindJSON(&input); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	from := os.Getenv("EMAIL_USERNAME")     // e.g. newsletter@cody.live
	pass := os.Getenv("EMAIL_PASSWORD")     // titan mailbox password
	to := os.Getenv("NEWSLETTER_NOTIFY_TO") // where you want to receive signups (e.g. you@cody.live)

	smtpHost := "smtp.titan.email"
	smtpAddr := smtpHost + ":587" // STARTTLS recommended

	subject := "New Cody.live Coming Soon Signup"
	body := fmt.Sprintf("New subscriber: %s\n", strings.TrimSpace(input.Email))

	msg := "" +
		"From: " + from + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		body

	// STARTTLS flow
	client, err := smtp.Dial(smtpAddr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer client.Close()

	if err := client.Hello("cody.live"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tlsConfig := &tls.Config{ServerName: smtpHost}
	if ok, _ := client.Extension("STARTTLS"); ok {
		if err := client.StartTLS(tlsConfig); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	auth := smtp.PlainAuth("", from, pass, smtpHost)
	if err := client.Auth(auth); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := client.Mail(from); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := client.Rcpt(to); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	w, err := client.Data()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if _, err := w.Write([]byte(msg)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := w.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success!"})
}
