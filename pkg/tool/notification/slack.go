package notification

import (
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/homenoc/dsbd-backend/pkg/tool/config"
	"os"
)

var (
	// change token
	IncomingUrl string = "https://hooks.slack.com/services/xxxxxxxxxxxx/xxxxxxxxxxxx/xxxxxxxxx"
)

type Slack struct {
	Status     bool
	Attachment slack.Attachment
	Channel    string
}

func SendSlack(data Slack) {
	for _, tmp := range config.Conf.Slack {
		if tmp.Channel == data.Channel {
			var color string
			if data.Status {
				color = "good"
			} else {
				color = "danger"
			}
			data.Attachment.Color = &color
			payload := slack.Payload{
				Username:    tmp.Name,
				Channel:     tmp.Channel,
				Attachments: []slack.Attachment{data.Attachment},
			}
			err := slack.Send(tmp.WebHookUrl, "", payload)
			if err != nil {
				os.Exit(1)
			}
		}
	}
}
