package v0

import (
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/notice"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	"strconv"
	"time"
)

const layoutInput = "2006-01-02 15:04:05"

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

func noticeSlackReplaceAdmin(before core.Notice, after notice.Input) {
	// 審査ステータスのSlack通知
	attachment := slack.Attachment{}

	attachment.AddField(slack.Field{Title: "Title", Value: "通知の変更"}).
		AddField(slack.Field{Title: "申請者", Value: "管理者"}).
		AddField(slack.Field{Title: "更新状況", Value: changeText(before, after)})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})
}

func changeText(before core.Notice, after notice.Input) string {
	data := ""
	//Title
	if after.Title != "" && after.Title != before.Title {
		data += "Title: " + before.Title + " => " + after.Title + "\n"
	}

	//Data
	if after.Data != "" && after.Data != before.Data {
		data += "Contents: " + before.Data + " => " + after.Data + "\n"
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
