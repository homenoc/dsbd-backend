package jpnic

import (
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	jpnicTransaction "github.com/homenoc/jpnic"
)

func success(data jpnicTransaction.Result) {
	attachment := slack.Attachment{}

	attachment.Text = &[]string{"JPNIC登録(手動実行)"}[0]
	attachment.AddField(slack.Field{Title: "申請者", Value: "管理者"}).
		AddField(slack.Field{Title: "RecepNo", Value: data.RecepNo}).
		AddField(slack.Field{Title: "AdmJPNICHdl（新規発番）", Value: data.AdmJPNICHdl}).
		AddField(slack.Field{Title: "Tech1JPNICHdl（補足情報）", Value: data.Tech1JPNICHdl}).
		AddField(slack.Field{Title: "Tech2JPNICHdl（補足情報）", Value: data.Tech2JPNICHdl})

	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})
}

func bad(data jpnicTransaction.Result) {
	attachment := slack.Attachment{}

	var errStr string
	for _, tmp := range data.ResultErr {
		errStr += tmp.Error() + "\n"
	}

	attachment.Text = &[]string{"[失敗]JPNIC登録(手動実行)"}[0]
	attachment.AddField(slack.Field{Title: "申請者", Value: "管理者"}).
		AddField(slack.Field{Title: "Error", Value: data.Err.Error()}).
		AddField(slack.Field{Title: "Error(詳細)", Value: errStr})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: false})
}
