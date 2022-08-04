package notification

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"github.com/slack-go/slack"
)

func NoticeUpdateStatus(groupID, info, history string) {
	// 審査ステータスのSlack通知
	Notification.Slack.PostMessage(config.Conf.Slack.Channels.Main, slack.MsgOptionBlocks(
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTHeader,
			Text: &slack.TextBlockObject{Type: "plain_text", Text: "接続情報登録"},
		},
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*申請者* System"},
			},
		},
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*GroupID* " + groupID},
			},
		},
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*現在ステータス情報* " + info},
				{Type: "mrkdwn", Text: "*ステータス履歴* " + history},
			},
		},
		slack.NewDividerBlock(),
	))
}
