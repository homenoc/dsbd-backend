package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/support"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	"github.com/slack-go/slack"
	"strconv"
)

func noticeAdd(title, user, group string, input support.FirstInput) {
	notification.Notification.Slack.PostMessage(config.Conf.Slack.Channels.Main, slack.MsgOptionBlocks(
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTHeader,
			Text: &slack.TextBlockObject{Type: "plain_text", Text: title},
		},
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*申請者* ユーザ"},
			},
		},
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*User* " + user},
				{Type: "mrkdwn", Text: "*Group* " + group},
			},
		},
		slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "*"+input.Title+"*", false, false), nil, nil),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: input.Data,
			},
		},
		slack.NewDividerBlock(),
	))
}
func noticeUpdateByAdmin(before, after core.Ticket) {
	title := "Ticket"
	if *before.Request {
		title = "Request"
	}

	notification.Notification.Slack.PostMessage(config.Conf.Slack.Channels.Main, slack.MsgOptionBlocks(
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTHeader,
			Text: &slack.TextBlockObject{Type: "plain_text", Text: title + " 更新"},
		},
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*申請者* 管理者"},
			},
		},
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*Title* " + strconv.Itoa(int(before.ID)) + ": " + before.Title},
			},
		},
		slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "*更新状況*", false, false), nil, nil),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: changeText(before, after),
			},
		},
		slack.NewDividerBlock(),
	))
}

func noticeUpdate(before, after core.Ticket, user, group string) {
	title := "Ticket"
	if *before.Request {
		title = "Request"
	}

	notification.Notification.Slack.PostMessage(config.Conf.Slack.Channels.Main, slack.MsgOptionBlocks(
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTHeader,
			Text: &slack.TextBlockObject{Type: "plain_text", Text: title + " 更新"},
		},
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*申請者* ユーザ"},
			},
		},
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*User* " + user},
				{Type: "mrkdwn", Text: "*Group* " + group},
				{Type: "mrkdwn", Text: "*Title* " + strconv.Itoa(int(before.ID)) + ": " + before.Title},
			},
		},
		slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "*更新状況*", false, false), nil, nil),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: changeText(before, after),
			},
		},
		slack.NewDividerBlock(),
	))
}

func noticeNewMessage(isAdmin bool, user, group string, ticket core.Ticket, message string) {
	applicant := "管理者"
	if isAdmin {
		applicant = "ユーザ"
	}
	title := "Ticket"
	if *ticket.Request {
		title = "Request"
	}

	notification.Notification.Slack.PostMessage(config.Conf.Slack.Channels.Main, slack.MsgOptionBlocks(
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTHeader,
			Text: &slack.TextBlockObject{Type: "plain_text", Text: "[" + title + "] Support(新規メッセージ)"},
		},
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*申請者* " + applicant},
			},
		},
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*User* " + user},
				{Type: "mrkdwn", Text: "*Group* " + group},
			},
		},
		slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", ticket.Title, false, false), nil, nil),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: message,
			},
		},
		slack.NewDividerBlock(),
	))
}

func changeText(before, after core.Ticket) string {
	data := ""
	// Title
	if !*before.Solved && *after.Solved {
		if *before.Request {
			data += "Solved: " + "未対応 => 変更/承諾済み\n"
		} else {
			data += "Solved: " + "未解決 => 解決済み\n"
		}
	}

	if *before.Solved && !*after.Solved {
		if *before.Request {
			data += "Solved: " + "変更/承諾済み => 未対応\n"
		} else {
			data += "Solved: " + "解決済み => 未解決\n"
		}
	}

	if *before.Request && !*before.RequestReject && *after.RequestReject {
		data += "Solved: " + "未対応or変更/承諾済み => 却下\n"
	}

	return data
}
