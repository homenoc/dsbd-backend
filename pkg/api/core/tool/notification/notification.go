package notification

import (
	"github.com/slack-go/slack"
)

type NotifyStruct struct {
	Slack *slack.Client
}

var Notification NotifyStruct

func NewNotification() {
	// slack
	slackToken := "xoxb-5041561262-3888879509782-BmR612XtlEObOzKRB95bGUMM"
	Notification.Slack = slack.New(slackToken)
}
