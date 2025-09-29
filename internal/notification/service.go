package notification

import (
	"easy-monitor/internal/config"
	"easy-monitor/internal/monitor"
	"fmt"
	"log"
	"net/smtp"
	"os"
)

func Notify(monitorResults []monitor.MonitorResult) {
	fmt.Printf("Monitor failed:, %v", monitorResults)

	if os.Getenv("SMTP_ENABLED") != "true" {
		return
	}

	config := config.GetConfig()
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_FROM") // must be a verified sender/domain in Mailjet
	to := config.Notify
	subject := fmt.Sprintf("[Easy Monitor] %v monitors have failed", len(monitorResults))

	lines := []string{
		fmt.Sprintf("The following %v monitors have failed. See the details below", len(monitorResults)),
		"",
		"Check the Details below",
	}

	for _, mr := range monitorResults {
		lines = append(lines, "")
		lines = append(lines, "-----------------------")
		lines = append(lines, "")
		lines = append(lines, fmt.Sprintf("Name: %v", mr.Name))
		lines = append(lines, fmt.Sprintf("Endpoint: %v", mr.Endpoint))
		lines = append(lines, fmt.Sprintf("Method: %v", mr.Method))
		lines = append(lines, fmt.Sprintf("Body: %v", mr.Body))
		lines = append(lines, fmt.Sprintf("Status: %v", mr.Status))
		lines = append(lines, "Expected:")
		lines = append(lines, fmt.Sprintf("   Status: %v", mr.Expected.Status))
		lines = append(lines, "   Headers:")
		for k, v := range mr.Expected.Headers {
			lines = append(lines, fmt.Sprintf("      %v: %v", k, v))
		}
		lines = append(lines, fmt.Sprintf("   Body: %v", mr.Expected.Body))
		lines = append(lines, "Received:")
		lines = append(lines, fmt.Sprintf("   Status: %v", mr.Received.Status))
		lines = append(lines, "   Headers:")
		for k, v := range mr.Expected.Headers {
			lines = append(lines, fmt.Sprintf("      %v: %v", k, v))
		}
		lines = append(lines, fmt.Sprintf("   Body: %v", mr.Received.Body))
		if mr.Error != "" {
			lines = append(lines, fmt.Sprintf("Error: %v", mr.Error))
		}
	}

	lines = append(lines, "")
	lines = append(lines, "-----------------------")
	lines = append(lines, "")
	lines = append(lines, "Best Regards,")
	lines = append(lines, "Your Easy Monitor")

	body := ""
	for _, l := range lines {
		body += fmt.Sprintln(l)
	}

	auth := smtp.PlainAuth("", username, password, host)

	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s\r\n", to, subject, body))

	addr := host + ":" + port
	if err := smtp.SendMail(addr, auth, from, to, msg); err != nil {
		log.Println("Could not send mail due to the following error. Make sure you have set host, port, user, pass and sender adress.")
		log.Println("For some smtp providers the sender address must be verified.")
		log.Println(err)
		return
	}

	fmt.Println("Mail sent successfully 🚀")
}
