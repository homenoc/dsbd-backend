package v0

import (
	"bytes"
	"encoding/base64"
	mailStruct "github.com/homenoc/dsbd-backend/pkg/api/core/mail"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"log"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"
)

//
// https://qiita.com/yamasaki-masahide/items/a9f8b43eeeaddbfb6b44
//

func utf8Split(utf8string string, length int) []string {
	var resultString []string
	var buffer bytes.Buffer
	for k, c := range strings.Split(utf8string, "") {
		buffer.WriteString(c)
		if k%length == length-1 {
			resultString = append(resultString, buffer.String())
			buffer.Reset()
		}
	}
	if buffer.Len() > 0 {
		resultString = append(resultString, buffer.String())
	}
	return resultString
}

func encodeSubject(subject string) string {
	var buffer bytes.Buffer
	buffer.WriteString("Subject:")
	for _, line := range utf8Split(subject, 13) {
		buffer.WriteString(" =?utf-8?B?")
		buffer.WriteString(base64.StdEncoding.EncodeToString([]byte(line)))
		buffer.WriteString("?=\r\n")
	}
	return buffer.String()
}

func SendMail(d mailStruct.Mail) error {
	from := mail.Address{Name: "Home NOC Operators' Group Support", Address: config.Conf.Mail.From}
	to := mail.Address{Address: d.ToMail}
	receivers := []string{to.Address}

	msg := "" +
		"From:" + from.String() + "\r\n" +
		"To:" + to.String() + "\r\n" +
		encodeSubject(d.Subject) + "\r\n" +
		"\r\n" + d.Content + "\r\n"

	auth := smtp.PlainAuth("", config.Conf.Mail.User, config.Conf.Mail.Pass, config.Conf.Mail.Host)
	err := smtp.SendMail(config.Conf.Mail.Host+":"+strconv.Itoa(config.Conf.Mail.Port), auth,
		from.Address, receivers, []byte(msg))

	noticeSlack(err, d)

	if err != nil {
		log.Printf("Error: %v\n", err)
		return err
	}

	return nil
}
