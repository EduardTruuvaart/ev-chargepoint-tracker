package service

type EmailService struct {
	send func(to []string, msg string)
}
