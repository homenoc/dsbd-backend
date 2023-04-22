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
	slackToken := config.Conf.Slack.BotToken
	appToken := slack.OptionAppLevelToken(config.Conf.Slack.AppToken)
	Notification.Slack = slack.New(slackToken, appToken, slack.OptionDebug(config.IsDebug))
}
