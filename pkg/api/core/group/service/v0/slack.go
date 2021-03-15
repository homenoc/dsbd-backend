package v0

import (
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	serviceTemplate "github.com/homenoc/dsbd-backend/pkg/api/core/template/service"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	dbServiceTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/service/v0"
	"github.com/jinzhu/gorm"
	"strconv"
)

func noticeSlackAdmin(before, after core.Service) {
	// 審査ステータスのSlack通知
	attachment := slack.Attachment{}

	attachment.AddField(slack.Field{Title: "Title", Value: "Service情報の更新"}).
		AddField(slack.Field{Title: "申請者", Value: "管理者"}).
		AddField(slack.Field{Title: "Group", Value: strconv.Itoa(int(before.ID)) + "-" + before.Group.Org}).
		AddField(slack.Field{Title: "更新状況", Value: changeText(before, after)})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})
}

func changeText(before, after core.Service) string {
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

	if after.Lock != nil {
		if *before.Lock != *after.Lock {
			if !*after.Lock {
				data += "ユーザ変更: 禁止 => 許可\n"
			} else {
				data += "ユーザ変更: 許可 => 禁止\n"
			}
		}
	}

	if after.AddAllow != nil {
		if *before.AddAllow != *after.AddAllow {
			if *after.AddAllow {
				data += "ユーザ側にて接続追加の許可: 禁止 => 許可\n"
			} else {
				data += "ユーザ側にて接続追加の許可: 許可 => 禁止\n"
			}
		}
	}

	if after.ServiceTemplateID != nil {
		if *before.ServiceTemplateID != *after.ServiceTemplateID {
			data += "ServiceID: " + before.ServiceTemplate.Type + " => " +
				serviceTemplateText(*after.ServiceTemplateID) + "\n"
		}
	}

	if after.Fee != nil {
		if *before.Fee != *after.Fee {
			data += "Fee: " + strconv.Itoa(int(*before.Fee)) + " => " + strconv.Itoa(int(*after.Fee)) + "\n"
		}
	}

	if before.AveDownstream != after.AveDownstream {
		data += "平均ダウンロード帯域: " + strconv.Itoa(int(before.AveDownstream)) + "Kbps => " +
			strconv.Itoa(int(after.AveDownstream)) + "Kbps\n"
	}

	if before.MaxDownstream != after.MaxDownstream {
		data += "最大ダウンロード帯域: " + strconv.Itoa(int(before.MaxDownstream)) + "Kbps => " +
			strconv.Itoa(int(after.MaxDownstream)) + "Kbps\n"
	}

	if before.AveUpstream != after.AveUpstream {
		data += "平均アップロード帯域: " + strconv.Itoa(int(before.AveUpstream)) + "Kbps => " +
			strconv.Itoa(int(after.AveUpstream)) + "Kbps\n"
	}

	if before.MaxUpstream != after.MaxUpstream {
		data += "最大アップロード帯域: " + strconv.Itoa(int(before.MaxUpstream)) + "Kbps => " +
			strconv.Itoa(int(after.MaxUpstream)) + "Kbps\n"
	}

	if after.ASN != nil {
		if *before.ASN != *after.ASN {
			data += "ASN: " + strconv.Itoa(int(*before.ASN)) + " => " + strconv.Itoa(int(*after.ASN)) + "\n"
		}
	}

	if after.RouteV4 != "" {
		data += "広報方法(v4): " + before.RouteV4 + "=>" + after.RouteV4 + "\n"
	}

	if after.RouteV6 != "" {
		data += "広報方法(v6): " + before.RouteV6 + "=>" + after.RouteV6 + "\n"
	}

	if after.V4Name != "" {
		data += "ネットワーク名(v4): " + before.V4Name + "=>" + after.V4Name + "\n"
	}

	if after.V6Name != "" {
		data += "ネットワーク名(v6): " + before.V6Name + "=>" + after.V6Name + "\n"
	}

	if after.Org != "" {
		data += "Org: " + before.Org + "=>" + after.Org + "\n"
	}

	if after.OrgEn != "" {
		data += "Org(En): " + before.OrgEn + "=>" + after.OrgEn + "\n"
	}

	if after.PostCode != "" {
		data += "PostCode: " + before.PostCode + "=>" + after.PostCode + "\n"
	}

	if after.Address != "" {
		data += "Address: " + before.Address + "=>" + after.Address + "\n"
	}

	if after.AddressEn != "" {
		data += "Address(En): " + before.AddressEn + "=>" + after.AddressEn + "\n"
	}

	return data
}

func serviceTemplateText(status uint) string {
	result := dbServiceTemplate.Get(serviceTemplate.ID, &core.ServiceTemplate{Model: gorm.Model{ID: status}})
	return result.Services[0].Type
}
