package v0

import (
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/homenoc/dsbd-backend/pkg/api/core/notice"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
)

func noticeSlackAddAdmin(input notice.Input) {
	// 審査ステータスのSlack通知
	attachment := slack.Attachment{}

	endTime := "無期限(9999年12月31日 23:59:59.59)"
	if input.EndTime != nil {
		endTime = *input.EndTime
	}

	attachment.AddField(slack.Field{Title: "Title", Value: "通知の追加"}).
		AddField(slack.Field{Title: "申請者", Value: "管理者"}).
		AddField(slack.Field{Title: "通知時期", Value: input.StartTime + " => " + endTime}).
		AddField(slack.Field{Title: "title", Value: input.Title}).
		AddField(slack.Field{Title: "data", Value: input.Data})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})
}
