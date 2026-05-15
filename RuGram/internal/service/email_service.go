package service

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
	"strconv"
)

type EmailService struct {
	host     string
	port     int
	username string
	password string
	from     string
	secure   bool
}

func NewEmailService() (*EmailService, error) {
	host := os.Getenv("SMTP_HOST")
	portStr := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USER")
	password := os.Getenv("SMTP_PASS")
	from := os.Getenv("SMTP_FROM")
	secure := os.Getenv("SMTP_SECURE") == "true"

	if host == "" || username == "" || password == "" {
		return nil, fmt.Errorf("SMTP configuration is incomplete")
	}

	port := 465
	if portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	if from == "" {
		from = username
	}

	return &EmailService{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
		secure:   secure,
	}, nil
}

func (s *EmailService) SendWelcomeEmail(to, displayName, userID string) error {
	subject := "Добро пожаловать в RuGram!"
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Здравствуйте, %s!</h2>
			<p>Спасибо за регистрацию в RuGram!</p>
			<p>Ваш аккаунт успешно создан. Теперь вы можете:</p>
			<ul>
				<li>Создавать посты</li>
				<li>Загружать фотографии</li>
				<li>Общаться с друзьями</li>
			</ul>
			<p>Ваш ID пользователя: %s</p>
			<p>С уважением,<br>Команда RuGram</p>
		</body>
		</html>
	`, displayName, userID)

	plainBody := fmt.Sprintf(`
		Здравствуйте, %s!
		Спасибо за регистрацию в RuGram!
		Ваш аккаунт успешно создан.
		Ваш ID пользователя: %s
		С уважением, Команда RuGram
	`, displayName, userID)

	msg := fmt.Sprintf("From: %s\r\n", s.from)
	msg += fmt.Sprintf("To: %s\r\n", to)
	msg += fmt.Sprintf("Subject: %s\r\n", subject)
	msg += "MIME-Version: 1.0\r\n"
	msg += "Content-Type: text/html; charset=UTF-8\r\n"
	msg += "\r\n" + body

	plainMsg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		s.from, to, subject, plainBody)

	auth := smtp.PlainAuth("", s.username, s.password, s.host)

	var err error
	if s.secure && s.port == 465 {
		err = s.sendWithTLS(msg, auth)
	} else {
		err = smtp.SendMail(fmt.Sprintf("%s:%d", s.host, s.port), auth, s.from, []string{to}, []byte(plainMsg))
	}

	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	return nil
}

func (s *EmailService) sendWithTLS(msg string, auth smtp.Auth) error {
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", s.host, s.port), &tls.Config{
		ServerName: s.host,
	})
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, s.host)
	if err != nil {
		return err
	}
	defer client.Quit()

	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return err
		}
	}

	if err = client.Mail(s.from); err != nil {
		return err
	}
	if err = client.Rcpt(s.username); err != nil {
		return err
	}

	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *EmailService) IsConfigured() bool {
	return s.host != "" && s.username != "" && s.password != ""
}