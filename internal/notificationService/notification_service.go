package notificationService

type NotificationService interface {
	Send(target, subject, text string) error
}
