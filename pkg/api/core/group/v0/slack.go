package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	"github.com/homenoc/dsbd-backend/pkg/api/notify"
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
				Text: "*追加状況*\nQuestion: " + group.Question +
					"\nOrg: " + group.Org + " (" + group.OrgEn + ")" +
					"\nCountry: " + group.Country +
					"\nContract: " + group.Contract,
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

// changeTextByAdmin summarises the group fields an admin changed. Scalar fields
// are driven by the `notify:"..."` tags on core.Group; MemberType and
// ExpiredStatus translate codes to human text, so they stay custom here.
func changeTextByAdmin(before, after core.Group) string {
	data := notify.Diff(before, after)

	if before.MemberType != after.MemberType {
		beforeMemberType, _ := core.GetMembershipTypeID(before.MemberType)
		afterMemberType, _ := core.GetMembershipTypeID(after.MemberType)
		data += "MemberType: " + beforeMemberType.Name + " => " + afterMemberType.Name + "\n"
	}

	if before.ExpiredStatus != nil && after.ExpiredStatus != nil && *before.ExpiredStatus != *after.ExpiredStatus {
		data += "ExpiredStatus: " + core.ExpiredLabel(*before.ExpiredStatus) + " => " +
			core.ExpiredLabel(*after.ExpiredStatus) + "\n"
	}

	return data
}
