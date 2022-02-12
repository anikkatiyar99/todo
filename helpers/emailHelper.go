package helper

import (
	"fmt"
	"log"
	"os"

	"github.com/anikkatiyar99/todo/models"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// SendEmail sends email to the given email address using sendgrid
func SendEmail(task *models.Task) error {
	sendgridKey := os.Getenv("SENDGRID_API_KEY")

	from := mail.NewEmail("Todo app", "home.katiyar@gmail.com")
	to := mail.NewEmail("User", task.AlertEmail)
	subject := "Your task ||" + *task.Title + "|| due soon"
	plainTextContent := "Your task" + *task.Title + " is due at " + task.DueDate.String()
	htmlContent := "<strong>Your task" + *task.Title + " is due at " + task.DueDate.String() + "</strong>"

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	alertSet := int(task.AlertAt.Unix())
	log.Println(alertSet)
	message.SetSendAt(alertSet)

	client := sendgrid.NewSendClient(sendgridKey)
	sendgridResp, err := client.Send(message)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("SendgridID",
		//sendgridResp.Headers["X-Message-Id"][0],
		sendgridResp.StatusCode,
		task.AlertEmail,
		sendgridResp.Body,
	)

	if sendgridResp.StatusCode != 202 {
		return fmt.Errorf("%d: %s ", sendgridResp.StatusCode, sendgridResp.Body)
	}

	return nil
}
