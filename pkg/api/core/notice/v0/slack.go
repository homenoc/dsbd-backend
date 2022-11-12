package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/notice"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	"github.com/slack-go/slack"
	"strconv"
	"time"
)

const layoutInput = "2006-01-02 15:04:05"

func noticeSlackAddByAdmin(input notice.Input) {
	// 審査ステータスのSlack通知
	endTime := "無期限(9999年12月31日 23:59:59.59)"
	if input.EndTime != nil {
		endTime = *input.EndTime
	}

	notification.Notification.Slack.PostMessage(config.Conf.Slack.Channels.Main, slack.MsgOptionBlocks(
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTHeader,
			Text: &slack.TextBlockObject{Type: "plain_text", Text: "通知追加"},
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
				{Type: "mrkdwn", Text: "*通知時期* " + input.StartTime + " => " + endTime},
			},
		},
		slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "*"+input.Title+"*", false, false), nil, nil),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: input.Body,
			},
		},
		slack.NewDividerBlock(),
	))
}

func noticeSlackReplaceByAdmin(before core.Notice, after notice.Input) {
	// 審査ステータスのSlack通知
	notification.Notification.Slack.PostMessage(config.Conf.Slack.Channels.Main, slack.MsgOptionBlocks(
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTHeader,
			Text: &slack.TextBlockObject{Type: "plain_text", Text: "通知変更"},
		},
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*申請者* 管理者"},
			},
		},
		slack.NewDividerBlock(),
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

func changeText(before core.Notice, after notice.Input) string {
	data := ""
	//Title
	if after.Title != "" && after.Title != before.Title {
		data += "Title: " + before.Title + " => " + after.Title + "\n"
	}

	//Data
	if after.Body != "" && after.Body != before.Data {
		data += "Contents: " + before.Data + " => " + after.Body + "\n"
	}

	//if after.UserID != before.UserID {
	//	data += "UserID: " + strconv.Itoa(int(before.UserID)) + " => " + strconv.Itoa(int(after.UserID)) + "\n"
	//}
	//
	//if after.GroupID != before.GroupID {
	//	data += "GroupID: " + strconv.Itoa(int(before.UserID)) + " => " + strconv.Itoa(int(after.UserID)) + "\n"
	//}
	//
	//if after.NOCID != before.NOCID {
	//	data += "NOCID: " + strconv.Itoa(int(before.NOCID)) + " => " + strconv.Itoa(int(after.NOCID)) + "\n"
	//}

	if after.StartTime != before.StartTime.Add(9*time.Hour).Format(layoutInput) {
		data += "Start Time: " + before.StartTime.Add(9*time.Hour).Format(layoutInput) + " => " + after.StartTime + "\n"
	}

	if *after.EndTime != before.EndTime.Add(9*time.Hour).Format(layoutInput) {
		data += "End Time: " + before.EndTime.Add(9*time.Hour).Format(layoutInput) + " => " + *after.EndTime + "\n"
	}

	if *after.Everyone != *before.Everyone {
		data += "EveryOne: " + strconv.FormatBool(*before.Everyone) + " => " + strconv.FormatBool(*after.Everyone) + "\n"
	}

	if *after.Info != *before.Info {
		data += "Info: " + strconv.FormatBool(*before.Info) + " => " + strconv.FormatBool(*after.Info) + "\n"
	}

	if *after.Fault != *before.Fault {
		data += "Fault: " + strconv.FormatBool(*before.Fault) + " => " + strconv.FormatBool(*after.Fault) + "\n"
	}

	if *after.Important != *before.Important {
		data += "Important: " + strconv.FormatBool(*before.Important) + " => " + strconv.FormatBool(*after.Important) + "\n"
	}

	return data
}
