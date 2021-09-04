package jpnic

import (
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	jpnicTransaction "github.com/homenoc/jpnic"
)

func registrationSuccess(data jpnicTransaction.Result) {
	attachment := slack.Attachment{}

	attachment.Text = &[]string{"JPNIC登録(手動実行)"}[0]
	attachment.AddField(slack.Field{Title: "申請者", Value: "管理者"}).
		AddField(slack.Field{Title: "受付番号", Value: data.RecepNo}).
		AddField(slack.Field{Title: "管理者連絡窓口(JPNICHandle)", Value: data.AdmJPNICHdl}).
		AddField(slack.Field{Title: "技術連絡窓口1(JPNICHandle)", Value: data.Tech1JPNICHdl}).
		AddField(slack.Field{Title: "技術連絡窓口2(JPNICHandle)", Value: data.Tech2JPNICHdl})

	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})
}

func registrationBad(data jpnicTransaction.Result) {
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

func returnSuccess(input ReturnInput, data string) {
	attachment := slack.Attachment{}

	var addresses string
	for _, tmp := range input.Address {
		addresses = addresses + ", " + tmp
	}
	returnDate := "[未指定]2ヶ月後"
	if input.ReturnDate != "" {
		returnDate = input.ReturnDate
	}

	attachment.Text = &[]string{"JPNICアドレス返却"}[0]
	attachment.AddField(slack.Field{Title: "申請者", Value: "管理者"}).
		AddField(slack.Field{Title: "申請アドレス", Value: addresses}).
		AddField(slack.Field{Title: "NetworkName", Value: input.NetworkName}).
		AddField(slack.Field{Title: "返却日", Value: returnDate}).
		AddField(slack.Field{Title: "Result", Value: data})

	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})
}

func returnBad(input ReturnInput, data string) {
	attachment := slack.Attachment{}

	var addresses string
	for _, tmp := range input.Address {
		addresses = addresses + ", " + tmp
	}
	returnDate := "[未指定]2ヶ月後"
	if input.ReturnDate != "" {
		returnDate = input.ReturnDate
	}

	attachment.Text = &[]string{"JPNICアドレス返却"}[0]
	attachment.AddField(slack.Field{Title: "申請者", Value: "管理者"}).
		AddField(slack.Field{Title: "申請アドレス", Value: addresses}).
		AddField(slack.Field{Title: "NetworkName", Value: input.NetworkName}).
		AddField(slack.Field{Title: "返却日", Value: returnDate}).
		AddField(slack.Field{Title: "Result", Value: data})

	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: false})
}
