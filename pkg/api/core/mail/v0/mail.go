package v0

import (
	"fmt"
	"time"
	"crypto/rand"
	"math/big"

	mailStruct "github.com/homenoc/dsbd-backend/pkg/api/core/mail"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"github.com/wneessen/go-mail"
)

func generateMessageID() string {
    host := config.Conf.Mail.Domain
    timestamp := time.Now().UnixNano()
    random, _ := rand.Int(rand.Reader, big.NewInt(1<<63-1))
    return fmt.Sprintf("%d.%d@%s", timestamp, random.Int64(), host)
}

func SendMail(d mailStruct.Mail) error {
	from := "" + config.Conf.Mail.FromName + " <" + config.Conf.Mail.From + ">"
	to := d.ToMail
	messageID := generateMessageID()

	message := mail.NewMsg()
	message.SetMIMEVersion(mail.MIME10)
	if err := message.From(from); err != nil {
		return fmt.Errorf("failed to set FROM address: %s", err)
	}
	if err := message.To(to); err != nil {
		return fmt.Errorf("failed to set TO address: %s", err)
	}
	message.SetMessageIDWithValue(messageID)

	message.Subject(d.Subject)
	message.SetBodyString(mail.TypeTextPlain, d.Content)

	client, err := mail.NewClient(config.Conf.Mail.Host, mail.WithPort(config.Conf.Mail.Port),
		mail.WithSMTPAuth(mail.SMTPAuthPlain), mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithUsername(config.Conf.Mail.User), mail.WithPassword(config.Conf.Mail.Pass),
	)
	if err != nil {
		return fmt.Errorf("failed to create new mail delivery client: %s", err)
	}
	err = client.DialAndSend(message)
	if err != nil {
		return fmt.Errorf("failed to deliver mail: %s", err)
	}
	noticeSlack(err, d)

	return nil
}
