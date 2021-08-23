package v0

import (
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	"strconv"
)

func noticeSlack(before, after core.Ticket) {
	attachment := slack.Attachment{}

	title := "Ticket"
	if *before.Request {
		title = "Request"
	}

	attachment.Text = &[]string{title + " 更新"}[0]
	attachment.AddField(slack.Field{Title: "申請者", Value: "管理者"}).
		AddField(slack.Field{Title: title, Value: strconv.Itoa(int(before.ID)) + ": " + before.Title}).
		AddField(slack.Field{Title: "更新状況", Value: changeText(before, after)})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})
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
