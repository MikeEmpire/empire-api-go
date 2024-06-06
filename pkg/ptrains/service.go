package ptrains

import (
	"net/http"
	"net/smtp"
	"os"

	"github.com/gin-gonic/gin"
)

type ContactForm struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
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
	contactFormMessage := contactFormInput.Message

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Contact Form Submission\n\n" +
		contactFormMessage

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success!"})

}
