package notification

import (
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"log"
)

type Slack struct {
	ID         string
	Status     bool
	Attachment slack.Attachment
}

func SendSlack(data Slack) {
	for _, tmp := range config.Conf.Slack {
		if tmp.ID == data.ID {
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
				Markdown:    true,
			}
			err := slack.Send(tmp.WebHookUrl, "", payload)
			if err != nil {
				log.Println(err)
			}
			return
		}
	}
}
