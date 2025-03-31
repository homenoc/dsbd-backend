package v0

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	mailStruct "github.com/homenoc/dsbd-backend/pkg/api/core/mail"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"github.com/wneessen/go-mail"
)

func generateMessageID() string {
	host := config.Conf.Mail.Domain
	timestamp := time.Now().UnixNano()
	random := rand.Int63()
	return fmt.Sprintf("%d.%d@%s", timestamp, random, host)
}

func SendMail(d mailStruct.Mail) error {
	from := "" + config.Conf.Mail.FromName + "<" + config.Conf.Mail.From + ">"
	to := d.ToMail
	messageID := generateMessageID()

	message := mail.NewMsg()
	message.SetMIMEVersion(mail.MIME10)
	if err := message.From(from); err != nil {
		log.Fatalf("failed to set FROM address: %s", err)
	}
	if err := message.To(to); err != nil {
		log.Fatalf("failed to set TO address: %s", err)
	}
	message.SetMessageIDWithValue(messageID)

	message.Subject(d.Subject)
	message.SetBodyString(mail.TypeTextPlain, d.Content)

	client, err := mail.NewClient(config.Conf.Mail.Host, mail.WithPort(config.Conf.Mail.Port),
		mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithUsername(config.Conf.Mail.User), mail.WithPassword(config.Conf.Mail.Pass),
	)
	if err != nil {
		log.Fatalf("failed to create new mail delivery client: %s", err)
	}
	err = client.DialAndSend(message)
	if err != nil {
		log.Fatalf("failed to deliver mail: %s", err)
	}
	noticeSlack(err, d)

	return nil
}
