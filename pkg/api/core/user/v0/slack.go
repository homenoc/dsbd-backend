package v0

import (
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	"strconv"
)

func noticeSlack(before core.User, after user.Input) {
	// 審査ステータスのSlack通知
	attachment := slack.Attachment{}

	attachment.AddField(slack.Field{Title: "Title", Value: "Service情報の更新"}).
		AddField(slack.Field{Title: "申請者", Value: "管理者"}).
		AddField(slack.Field{Title: "Group", Value: strconv.Itoa(int(before.ID)) + "-" + before.Group.Org}).
		AddField(slack.Field{Title: "更新状況", Value: changeText(before, after)})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})
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
