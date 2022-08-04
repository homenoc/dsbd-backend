package notification

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"github.com/slack-go/slack"
)

type NotifyStruct struct {
	Slack *slack.Client
}

var Notification NotifyStruct

func NewNotification() {
	// slack
	slackToken := config.Conf.Slack.Token
	Notification.Slack = slack.New(slackToken)
}
