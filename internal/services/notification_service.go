package services

type NotificationService interface {
	Send(target, subject, text string) error
}
