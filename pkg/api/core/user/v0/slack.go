package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	"github.com/slack-go/slack"
	"strconv"
)

func noticeAdd(input user.Input) {
	// 審査ステータスのSlack通知
	notification.Notification.Slack.PostMessage(config.Conf.Slack.Channels.Main, slack.MsgOptionBlocks(
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTHeader,
			Text: &slack.TextBlockObject{Type: "plain_text", Text: "新規ユーザ登録"},
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
				{Type: "mrkdwn", Text: "*Name* " + input.Name + " (" + input.NameEn + ")"},
				{Type: "mrkdwn", Text: "*E-Mail* " + input.Email},
			},
		},
		slack.NewDividerBlock(),
	))
}

func noticeAddFromGroup(user user.Input, group core.Group) {
	// 審査ステータスのSlack通知
	notification.Notification.Slack.PostMessage(config.Conf.Slack.Channels.Main, slack.MsgOptionBlocks(
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTHeader,
			Text: &slack.TextBlockObject{Type: "plain_text", Text: "追加ユーザ登録"},
		},
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*申請者* ユーザ"},
				{Type: "mrkdwn", Text: "*Group* " + "(" + strconv.Itoa(int(group.ID)) + ")" + group.Org + " (" + group.OrgEn + ")"},
			},
		},
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*Name* " + user.Name + " (" + user.NameEn + ")"},
				{Type: "mrkdwn", Text: "*E-Mail* " + user.Email},
			},
		},
		slack.NewDividerBlock(),
	))
}

func noticeRenew(loginUser, before core.User, after user.Input) {
	// 審査ステータスのSlack通知
	groupStr := "なし"
	if loginUser.GroupID != nil {
		groupStr = strconv.Itoa(int(*loginUser.GroupID)) + "-" + loginUser.Group.Org
	}

	notification.Notification.Slack.PostMessage(config.Conf.Slack.Channels.Main, slack.MsgOptionBlocks(
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTHeader,
			Text: &slack.TextBlockObject{Type: "plain_text", Text: "User情報更新"},
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
				{Type: "mrkdwn", Text: "*GroupID* " + groupStr},
				{Type: "mrkdwn", Text: "*User* " + strconv.Itoa(int(loginUser.ID)) + "-" + loginUser.Name},
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

func changeText(before core.User, after user.Input) string {
	data := ""
	//Name
	if after.Name != "" {
		data += "名前: " + before.Name + " => " + after.Name + "\n"
	}

	//Name (English)
	if after.NameEn != "" {
		data += "名前(English): " + before.NameEn + " => " + after.NameEn + "\n"
	}

	if after.Email != "" {
		data += "メールアドレスの変更: " + before.Email + " => " + after.Email + "\n"
	}

	if after.Pass != "" {
		data += "パスワードの変更処理\n"
	}

	if after.Level != 0 {
		data += "ServiceID: " + levelTemplateText(before.Level) + " => " + levelTemplateText(after.Level)
	}

	return data
}

func levelTemplateText(status uint) string {
	if status == 1 {
		return "Group内の申請、変更、閲覧(Master扱い)"
	} else if status == 2 {
		return "Group内の申請、変更、閲覧(Masterから割当)"
	} else if status == 3 {
		return "Group内の情報閲覧のみ"
	} else if status == 4 {
		return "障害情報閲覧のみ"
	}
	return "不明"
}
