package v0

import (
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	serviceTemplate "github.com/homenoc/dsbd-backend/pkg/api/core/template/service"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/api/store/group/v0"
	dbServiceTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/service/v0"
	"github.com/jinzhu/gorm"
	"strconv"
)

func getGroupInfo(groupID uint) core.Group {
	result := dbGroup.Get(group.ID, &core.Group{Model: gorm.Model{ID: groupID}})
	return result.Group[0]
}

func noticeSlackAdd(groupID int, serviceCode, serviceCodeComment string) {
	grpInfo := getGroupInfo(uint(groupID))
	attachment := slack.Attachment{}

	attachment.AddField(slack.Field{Title: "Title", Value: "Service情報登録(管理者実行)"}).
		AddField(slack.Field{Title: "申請者", Value: "管理者"}).
		AddField(slack.Field{Title: "Group", Value: strconv.Itoa(groupID) + "-" + grpInfo.Org}).
		AddField(slack.Field{Title: "サービスコード（新規発番）", Value: serviceCode}).
		AddField(slack.Field{Title: "サービスコード（補足情報）", Value: serviceCodeComment})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})
}

func noticeSlackAddJPNICByAdmin(serviceID int, input core.JPNICAdmin) {
	attachment := slack.Attachment{}

	attachment.AddField(slack.Field{Title: "Title", Value: "JPNIC管理者連絡窓口の追加"}).
		AddField(slack.Field{Title: "申請者", Value: "管理者"}).
		AddField(slack.Field{Title: "Service", Value: strconv.Itoa(serviceID)}).
		AddField(slack.Field{Title: "Name", Value: input.Name + " (" + input.NameEn + ")"}).
		AddField(slack.Field{Title: "Org", Value: input.Org + " (" + input.OrgEn + ")"})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})
}

func noticeSlackAddJPNICTech(serviceID int, input core.JPNICTech) {
	attachment := slack.Attachment{}

	attachment.AddField(slack.Field{Title: "Title", Value: "JPNIC技術連絡担当者の追加"}).
		AddField(slack.Field{Title: "申請者", Value: "管理者"}).
		AddField(slack.Field{Title: "Service", Value: strconv.Itoa(serviceID)}).
		AddField(slack.Field{Title: "Name", Value: input.Name + " (" + input.NameEn + ")"}).
		AddField(slack.Field{Title: "Org", Value: input.Org + " (" + input.OrgEn + ")"})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})
}

func noticeSlackAddIP(serviceID int, inputName string) {
	// 審査ステータスのSlack通知
	attachment := slack.Attachment{}

	attachment.AddField(slack.Field{Title: "Title", Value: "IPの追加"}).
		AddField(slack.Field{Title: "申請者", Value: "管理者"}).
		AddField(slack.Field{Title: "Service", Value: strconv.Itoa(serviceID)}).
		AddField(slack.Field{Title: "Name", Value: inputName})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})
}

func noticeSlackAddPlan(ipID int, input core.Plan) {
	attachment := slack.Attachment{}

	attachment.AddField(slack.Field{Title: "Title", Value: "Planの追加"}).
		AddField(slack.Field{Title: "申請者", Value: "管理者"}).
		AddField(slack.Field{Title: "IP", Value: strconv.Itoa(ipID)}).
		AddField(slack.Field{Title: "Name", Value: input.Name}).
		AddField(slack.Field{Title: "Plan", Value: strconv.Itoa(int(input.After)) + "/" +
			strconv.Itoa(int(input.HalfYear)) + "/" + strconv.Itoa(int(input.OneYear))})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})
}

func noticeSlackDelete(str string, id uint) {
	color := "warning"
	attachment := slack.Attachment{Color: &color}

	attachment.AddField(slack.Field{Title: "Title", Value: str + "の削除"}).
		AddField(slack.Field{Title: "申請者", Value: "管理者"}).
		AddField(slack.Field{Title: "削除処理", Value: "ID: " + strconv.Itoa(int(id))})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})
}

func noticeSlackUpdate(before, after core.Service) {
	attachment := slack.Attachment{}

	attachment.AddField(slack.Field{Title: "Title", Value: "Service情報の更新"}).
		AddField(slack.Field{Title: "申請者", Value: "管理者"}).
		AddField(slack.Field{Title: "Group", Value: strconv.Itoa(int(before.ID)) + "-" + before.Group.Org}).
		AddField(slack.Field{Title: "更新状況", Value: changeText(before, after)})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})
}

func noticeSlackUpdateJPNICByAdmin(before, after core.JPNICAdmin) {
	attachment := slack.Attachment{}

	attachment.AddField(slack.Field{Title: "Title", Value: "JPNIC管理者連絡窓口の更新"}).
		AddField(slack.Field{Title: "申請者", Value: "管理者"}).
		AddField(slack.Field{Title: "JPNICAdmin", Value: strconv.Itoa(int(before.ID))}).
		AddField(slack.Field{Title: "更新状況", Value: changeTextJPNICByAdmin(before, after)})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})
}

func noticeSlackUpdateJPNICTech(before, after core.JPNICTech) {
	attachment := slack.Attachment{}

	attachment.AddField(slack.Field{Title: "Title", Value: "JPNIC技術連絡担当者の更新"}).
		AddField(slack.Field{Title: "申請者", Value: "管理者"}).
		AddField(slack.Field{Title: "JPNICTech", Value: strconv.Itoa(int(before.ID))}).
		AddField(slack.Field{Title: "更新状況", Value: changeTextJPNICTech(before, after)})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})
}

func noticeSlackUpdateIP(before, after core.IP) {
	attachment := slack.Attachment{}

	attachment.AddField(slack.Field{Title: "Title", Value: "IPの更新"}).
		AddField(slack.Field{Title: "申請者", Value: "管理者"}).
		AddField(slack.Field{Title: "IP", Value: strconv.Itoa(int(before.ID))}).
		AddField(slack.Field{Title: "更新状況", Value: changeTextIP(before, after)})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})
}

func noticeSlackUpdatePlan(before, after core.Plan) {
	attachment := slack.Attachment{}

	attachment.AddField(slack.Field{Title: "Title", Value: "Planの更新"}).
		AddField(slack.Field{Title: "申請者", Value: "管理者"}).
		AddField(slack.Field{Title: "Plan", Value: strconv.Itoa(int(before.ID))}).
		AddField(slack.Field{Title: "更新状況", Value: changeTextPlan(before, after)})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})
}

func changeText(before, after core.Service) string {
	data := ""
	if after.Pass != nil {
		if *before.Pass != *after.Pass {
			if *after.Pass {
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

	if before.RouteV4 != after.RouteV4 {
		data += "広報方法(v4): " + before.RouteV4 + "=>" + after.RouteV4 + "\n"
	}

	if before.RouteV6 != after.RouteV6 {
		data += "広報方法(v6): " + before.RouteV6 + "=>" + after.RouteV6 + "\n"
	}

	if before.V4Name != after.V4Name {
		data += "ネットワーク名(v4): " + before.V4Name + "=>" + after.V4Name + "\n"
	}

	if before.V6Name != after.V6Name {
		data += "ネットワーク名(v6): " + before.V6Name + "=>" + after.V6Name + "\n"
	}

	if before.Org != after.Org {
		data += "Org: " + before.Org + "=>" + after.Org + "\n"
	}

	if before.OrgEn != after.OrgEn {
		data += "Org(En): " + before.OrgEn + "=>" + after.OrgEn + "\n"
	}

	if before.PostCode != after.PostCode {
		data += "PostCode: " + before.PostCode + "=>" + after.PostCode + "\n"
	}

	if before.Address != after.Address {
		data += "Address: " + before.Address + "=>" + after.Address + "\n"
	}

	if before.AddressEn != after.AddressEn {
		data += "Address(En): " + before.AddressEn + "=>" + after.AddressEn + "\n"
	}

	return data
}

func changeTextJPNICByAdmin(before, after core.JPNICAdmin) string {
	data := ""

	if after.Lock != nil {
		if *before.Lock != *after.Lock {
			if !*after.Lock {
				data += "ユーザ変更: 禁止 => 許可\n"
			} else {
				data += "ユーザ変更: 許可 => 禁止\n"
			}
		}
	}

	if before.JPNICHandle != after.JPNICHandle {
		data += "JPNICHandle: " + before.JPNICHandle + "=>" + after.JPNICHandle + "\n"
	}

	if before.Name != after.Name {
		data += "Name: " + before.Name + "=>" + after.Name + "\n"
	}

	if before.NameEn != after.NameEn {
		data += "Name(En): " + before.NameEn + "=>" + after.NameEn + "\n"
	}

	if before.Mail != after.Mail {
		data += "Mail: " + before.Mail + "=>" + after.Mail + "\n"
	}

	if before.Org != after.Org {
		data += "Org: " + before.Org + "=>" + after.Org + "\n"
	}

	if before.OrgEn != after.OrgEn {
		data += "Org(En): " + before.OrgEn + "=>" + after.OrgEn + "\n"
	}

	if before.PostCode != after.PostCode {
		data += "PostCode: " + before.PostCode + "=>" + after.PostCode + "\n"
	}

	if before.Address != after.Address {
		data += "Address: " + before.Address + "=>" + after.Address + "\n"
	}

	if before.AddressEn != after.AddressEn {
		data += "Address(En): " + before.AddressEn + "=>" + after.AddressEn + "\n"
	}

	if before.Dept != after.Dept {
		data += "Dept: " + before.Dept + "=>" + after.Dept + "\n"
	}

	if before.DeptEn != after.DeptEn {
		data += "Dept(En): " + before.DeptEn + "=>" + after.DeptEn + "\n"
	}

	if before.Tel != after.Tel {
		data += "Tel: " + before.Tel + "=>" + after.Tel + "\n"
	}

	if before.Fax != after.Fax {
		data += "Fax: " + before.Fax + "=>" + after.Fax + "\n"
	}

	if before.Country != after.Country {
		data += "Country: " + before.Country + "=>" + after.Country + "\n"
	}

	return data
}

func changeTextJPNICTech(before, after core.JPNICTech) string {
	data := ""

	if after.Lock != nil {
		if *before.Lock != *after.Lock {
			if !*after.Lock {
				data += "ユーザ変更: 禁止 => 許可\n"
			} else {
				data += "ユーザ変更: 許可 => 禁止\n"
			}
		}
	}

	if before.JPNICHandle != after.JPNICHandle {
		data += "JPNICHandle: " + before.JPNICHandle + "=>" + after.JPNICHandle + "\n"
	}

	if before.Name != after.Name {
		data += "Name: " + before.Name + "=>" + after.Name + "\n"
	}

	if before.NameEn != after.NameEn {
		data += "Name(En): " + before.NameEn + "=>" + after.NameEn + "\n"
	}

	if before.Mail != after.Mail {
		data += "Mail: " + before.Mail + "=>" + after.Mail + "\n"
	}

	if before.Org != after.Org {
		data += "Org: " + before.Org + "=>" + after.Org + "\n"
	}

	if before.OrgEn != after.OrgEn {
		data += "Org(En): " + before.OrgEn + "=>" + after.OrgEn + "\n"
	}

	if before.PostCode != after.PostCode {
		data += "PostCode: " + before.PostCode + "=>" + after.PostCode + "\n"
	}

	if before.Address != after.Address {
		data += "Address: " + before.Address + "=>" + after.Address + "\n"
	}

	if before.AddressEn != after.AddressEn {
		data += "Address(En): " + before.AddressEn + "=>" + after.AddressEn + "\n"
	}

	if before.Dept != after.Dept {
		data += "Dept: " + before.Dept + "=>" + after.Dept + "\n"
	}

	if before.DeptEn != after.DeptEn {
		data += "Dept(En): " + before.DeptEn + "=>" + after.DeptEn + "\n"
	}

	if before.Tel != after.Tel {
		data += "Tel: " + before.Tel + "=>" + after.Tel + "\n"
	}

	if before.Fax != after.Fax {
		data += "Fax: " + before.Fax + "=>" + after.Fax + "\n"
	}

	if before.Country != after.Country {
		data += "Country: " + before.Country + "=>" + after.Country + "\n"
	}

	return data
}

func changeTextIP(before, after core.IP) string {
	data := ""

	if before.Name != after.Name {
		data += "Name: " + before.Name + "=>" + after.Name + "\n"
	}

	if before.IP != after.IP {
		data += "IP: " + before.IP + "=>" + after.IP + "\n"
	}

	if before.UseCase != after.UseCase {
		data += "UseCase: " + before.UseCase + "=>" + after.UseCase + "\n"
	}

	if *before.Open != *after.Open {
		if *after.Open {
			data += "Open: 未開通 => 開通\n"
		} else {
			data += "Open: 開通 => 未開通\n"
		}
	}

	return data
}

func changeTextPlan(before, after core.Plan) string {
	data := ""

	if before.Name != after.Name {
		data += "Name: " + before.Name + "=>" + after.Name + "\n"
	}

	if before.After != after.After {
		data += "直後: " + strconv.Itoa(int(before.After)) + "=>" + strconv.Itoa(int(after.After)) + "\n"
	}

	if before.HalfYear != after.HalfYear {
		data += "半年後: " + strconv.Itoa(int(before.HalfYear)) + "=>" + strconv.Itoa(int(after.HalfYear)) + "\n"
	}

	if before.OneYear != after.OneYear {
		data += "1年後: " + strconv.Itoa(int(before.OneYear)) + "=>" + strconv.Itoa(int(after.OneYear)) + "\n"
	}

	return data
}

func serviceTemplateText(status uint) string {
	result := dbServiceTemplate.Get(serviceTemplate.ID, &core.ServiceTemplate{Model: gorm.Model{ID: status}})
	return result.Services[0].Type
}
