package v0

import (
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	"strconv"
)

func noticeSlackAdd(memo *core.Memo) {
	// 審査ステータスのSlack通知
	attachment := slack.Attachment{}

	attachment.Text = &[]string{"Memoの登録"}[0]
	attachment.AddField(slack.Field{Title: "申請者", Value: "管理者"}).
		AddField(slack.Field{Title: "GroupID", Value: strconv.Itoa(int(memo.GroupID))}).
		AddField(slack.Field{Title: "Type", Value: strconv.Itoa(int(memo.Type))}).
		AddField(slack.Field{Title: "Title", Value: memo.Title}).
		AddField(slack.Field{Title: "Message", Value: memo.Message})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})
}

func noticeSlackDelete(id int) {
	// 審査ステータスのSlack通知
	attachment := slack.Attachment{}

	attachment.Text = &[]string{"Memoの削除"}[0]
	attachment.AddField(slack.Field{Title: "申請者", Value: "管理者"}).
		AddField(slack.Field{Title: "ID", Value: strconv.Itoa(id)})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})
}
