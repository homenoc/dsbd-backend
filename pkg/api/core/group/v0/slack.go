package v0

import (
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	"strconv"
)

func noticeSlack(loginUser core.User, before core.Group, after group.Input) {
	// 審査ステータスのSlack通知
	attachment := slack.Attachment{}

	attachment.AddField(slack.Field{Title: "Title", Value: "Group情報の更新"}).
		AddField(slack.Field{Title: "申請者", Value: strconv.Itoa(int(loginUser.ID)) + "-" + loginUser.Name}).
		AddField(slack.Field{Title: "Group", Value: strconv.Itoa(int(before.ID)) + ":" + before.Org}).
		AddField(slack.Field{Title: "更新状況", Value: changeText(before, after)})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})
}

func noticeSlackByAdmin(before, after core.Group) {
	// 審査ステータスのSlack通知
	attachment := slack.Attachment{}

	attachment.AddField(slack.Field{Title: "Title", Value: "Group情報の更新"}).
		AddField(slack.Field{Title: "申請者", Value: "管理者"}).
		AddField(slack.Field{Title: "Group", Value: strconv.Itoa(int(before.ID)) + ":" + before.Org}).
		AddField(slack.Field{Title: "更新状況", Value: changeTextByAdmin(before, after)})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})
}

func noticeSlackCancelSubscriptionByAdmin(group core.Group) {
	// 審査ステータスのSlack通知
	attachment := slack.Attachment{}

	attachment.AddField(slack.Field{Title: "Title", Value: "Cancel Subscription"}).
		AddField(slack.Field{Title: "申請者", Value: "管理者"}).
		AddField(slack.Field{Title: "Group", Value: strconv.Itoa(int(group.ID)) + ":" + group.Org})

	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: false})
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
	if after.Open != nil {
		if *before.Open != *after.Open {
			if *after.Open {
				data += "開通: 未開通 => 開通済み\n"
			} else {
				data += "開通: 開通 => 未開通\n"
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
