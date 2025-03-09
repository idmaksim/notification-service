package email

import (
	"context"
	"fmt"
	"github.com/idmaksim/notification-service/internal/config"
	"github.com/mailersend/mailersend-go"
	"time"
)

type MailService struct {
	cfg *config.Config
}

func NewMailService() *MailService {
	return &MailService{cfg: config.GetConfig()}
}

func (s *MailService) Send(target, subject, text string) error {
	ms := mailersend.NewMailersend(s.cfg.MailApiKey)
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*10)

	defer cancel()

	from := mailersend.From{
		Name:  s.cfg.EmailFromName,
		Email: s.cfg.EmailFrom,
	}

	recipients := []mailersend.Recipient{{
		Email: target,
	}}

	message := ms.Email.NewMessage()
	message.SetSubject(subject)
	message.SetRecipients(recipients)
	message.SetFrom(from)
	message.SetText(text)

	if res, err := ms.Email.Send(ctx, message); err != nil {
		return err
	} else {
		fmt.Println(res.Status)
	}

	return nil
}
