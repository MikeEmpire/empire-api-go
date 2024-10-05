package ptrains

import (
	"net/http"
	"net/smtp"
	"os"

	"github.com/gin-gonic/gin"
)

type ContactForm struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Message     string `json:"message"`
	PhoneNumber string `json:"phoneNumber"`
}

func HandleContactForm(c *gin.Context) {
	var contactFormInput ContactForm
	if err := c.BindJSON(&contactFormInput); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	from := os.Getenv("EMAIL_USERNAME")
	pass := os.Getenv("EMAIL_PASSWORD")

	to := "cairambelu@gmail.com"
	// Construct the email body
	emailBody := "Name: " + contactFormInput.Name + "\n" +
		"Email: " + contactFormInput.Email + "\n" +
		"Message: " + contactFormInput.Message + "\n" +
		"Phone Number: " + contactFormInput.PhoneNumber + "\n"

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Contact Form Submission\n\n" +
		emailBody

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success!"})

}
