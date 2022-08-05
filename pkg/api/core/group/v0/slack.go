package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	"github.com/slack-go/slack"
	"strconv"
)

func noticeAddGroup(user core.User, group group.Input) {
	// 審査ステータスのSlack通知
	notification.Notification.Slack.PostMessage(config.Conf.Slack.Channels.Main, slack.MsgOptionBlocks(
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTHeader,
			Text: &slack.TextBlockObject{Type: "plain_text", Text: "新規Group登録"},
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
				{Type: "mrkdwn", Text: "*User* " + strconv.Itoa(int(user.ID)) + "-" + user.Name},
				{Type: "mrkdwn", Text: "*Group* " + group.Org + " (" + group.OrgEn + ")"},
			},
		},
		slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "*状況*", false, false), nil, nil),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: "*Question*\n" + group.Question + "\n*Country*\n" + group.Country + "\n*Contract*\n" + group.Contract,
			},
		},
		slack.NewDividerBlock(),
	))
}

func noticeByAdmin(before, after core.Group) {
	// 審査ステータスのSlack通知
	notification.Notification.Slack.PostMessage(config.Conf.Slack.Channels.Main, slack.MsgOptionBlocks(
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTHeader,
			Text: &slack.TextBlockObject{Type: "plain_text", Text: "Group情報更新"},
		},
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*Group情報更新*"},
				{Type: "mrkdwn", Text: "*申請者* 管理者"},
			},
		},
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*GroupID* " + strconv.Itoa(int(before.ID)) + ":" + before.Org},
			},
		},
		slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "*Title*　Test Title", false, false), nil, nil),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: "*更新状況*\n" + changeTextByAdmin(before, after),
			},
		},
		slack.NewDividerBlock(),
	))
}

func noticeCancelSubscriptionByAdmin(group core.Group) {
	// 審査ステータスのSlack通知
	notification.Notification.Slack.PostMessage(config.Conf.Slack.Channels.Main, slack.MsgOptionBlocks(
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTHeader,
			Text: &slack.TextBlockObject{Type: "plain_text", Text: "Cancel Subscription"},
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
				{Type: "mrkdwn", Text: "*GroupID* " + strconv.Itoa(int(group.ID)) + ":" + group.Org},
			},
		},
		slack.NewDividerBlock(),
	))
}

func changeText(before core.Group, after group.Input) string {
	data := ""

	if after.Org != "" && after.Org != before.Org {
		data += "Org: " + before.Org + "=>" + after.Org + "\n"
	}

	if after.OrgEn != "" && after.OrgEn != before.OrgEn {
		data += "Org(En): " + before.OrgEn + "=>" + after.OrgEn + "\n"
	}

	if after.PostCode != "" && after.PostCode != before.PostCode {
		data += "PostCode: " + before.PostCode + "=>" + after.PostCode + "\n"
	}

	if after.Address != "" && after.Address != before.Address {
		data += "Address: " + before.Address + "=>" + after.Address + "\n"
	}

	if after.AddressEn != "" && after.AddressEn != before.AddressEn {
		data += "Address(En): " + before.AddressEn + "=>" + after.AddressEn + "\n"
	}

	if after.Tel != "" && after.Tel != before.Tel {
		data += "Tel: " + before.Tel + "=>" + after.Tel + "\n"
	}

	if after.Country != "" && after.Country != before.Country {
		data += "Country: " + before.Country + "=>" + after.Country + "\n"
	}

	return data
}

func changeTextByAdmin(before, after core.Group) string {
	data := ""

	if after.AddAllow != nil {
		if *before.AddAllow != *after.AddAllow {
			if *after.AddAllow {
				data += "サービス新規申請: 禁止 => 許可\n"
			} else {
				data += "サービス新規申請: 許可 => 禁止\n"
			}
		}
	}

	if after.Pass != nil {
		if *before.Pass != *after.Pass {
			if *after.Pass {
				data += "審査: 未審査 => 審査合格済み\n"
			} else {
				data += "審査: 審査合格 => 未審査状態\n"
			}
		}
	}

	if after.ExpiredStatus != nil {
		if *before.ExpiredStatus != *after.ExpiredStatus {
			data += "ExpiredStatus: " + expiredStatusText(*before.ExpiredStatus) + " => " +
				expiredStatusText(*after.ExpiredStatus) + "\n"
		}
	}

	if after.Org != "" && after.Org != before.Org {
		data += "Org: " + before.Org + "=>" + after.Org + "\n"
	}

	if after.OrgEn != "" && after.OrgEn != before.OrgEn {
		data += "Org(En): " + before.OrgEn + "=>" + after.OrgEn + "\n"
	}

	if after.PostCode != "" && after.PostCode != before.PostCode {
		data += "PostCode: " + before.PostCode + "=>" + after.PostCode + "\n"
	}

	if after.Address != "" && after.Address != before.Address {
		data += "Address: " + before.Address + "=>" + after.Address + "\n"
	}

	if after.AddressEn != "" && after.AddressEn != before.AddressEn {
		data += "Address(En): " + before.AddressEn + "=>" + after.AddressEn + "\n"
	}

	if after.Tel != "" && after.Tel != before.Tel {
		data += "Tel: " + before.Tel + "=>" + after.Tel + "\n"
	}

	if after.Country != "" && after.Country != before.Country {
		data += "Country: " + before.Country + "=>" + after.Country + "\n"
	}

	return data
}

func expiredStatusText(status uint) string {
	if status == 0 {
		return "0"
	} else if status == 1 {
		return "ユーザより廃止"
	} else if status == 2 {
		return "運営委員より廃止"
	} else if status == 3 {
		return "審査落ち"
	} else {
		return "status不明"
	}
}

func statusText(status uint) string {
	if status == 0 {
		return "0"
	} else if status == 1 {
		return "ネットワーク情報　記入段階"
	} else if status == 2 {
		return "審査中"
	} else if status == 3 {
		return "接続情報　記入段階"
	} else if status == 4 {
		return "開通作業中"
	} else {
		return "status不明"
	}
}
