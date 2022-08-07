package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/api/store/group/v0"
	"github.com/slack-go/slack"
	"gorm.io/gorm"
	"strconv"
)

func getGroupInfo(groupID uint) core.Group {
	result := dbGroup.Get(group.ID, &core.Group{Model: gorm.Model{ID: groupID}})
	return result.Group[0]
}

func noticeAdd(applicant, groupID, serviceCodeNew, serviceCodeComment string) {
	if applicant == "" {
		applicant = "管理者"
	}
	// 審査ステータスのSlack通知
	notification.Notification.Slack.PostMessage(config.Conf.Slack.Channels.Main, slack.MsgOptionBlocks(
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTHeader,
			Text: &slack.TextBlockObject{Type: "plain_text", Text: "接続情報登録"},
		},
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*申請者* " + applicant},
			},
		},
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*GroupID* " + groupID},
				{Type: "mrkdwn", Text: "*サービスコード（新規発番）* " + serviceCodeNew},
				{Type: "mrkdwn", Text: "*サービスコード（補足情報）* " + serviceCodeComment},
			},
		},
		slack.NewDividerBlock(),
	))
}

func noticeDelete(title string, id uint) {
	notification.Notification.Slack.PostMessage(config.Conf.Slack.Channels.Main, slack.MsgOptionBlocks(
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTHeader,
			Text: &slack.TextBlockObject{Type: "plain_text", Text: title + "削除"},
		},
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*申請者* 管理者"},
			},
		},
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*削除処理 ID* " + strconv.Itoa(int(id))},
			},
		},
		slack.NewDividerBlock(),
	))
}

func noticeUpdate(before, after core.Service) {
	notification.Notification.Slack.PostMessage(config.Conf.Slack.Channels.Main, slack.MsgOptionBlocks(
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTHeader,
			Text: &slack.TextBlockObject{Type: "plain_text", Text: "Service情報更新"},
		},
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*申請者* 管理者"},
			},
		},
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*Group* " + "[" + strconv.Itoa(int(before.ID)) + "] " + before.Group.Org + "(" + before.Group.OrgEn + ")"},
			},
		},
		slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "*更新状況*", false, false), nil, nil),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: changeText(before, after),
			},
		},
		slack.NewDividerBlock(),
	))
}

func noticeAddJPNICByAdmin(serviceID int, input core.JPNICAdmin) {
	notification.Notification.Slack.PostMessage(config.Conf.Slack.Channels.Main, slack.MsgOptionBlocks(
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTHeader,
			Text: &slack.TextBlockObject{Type: "plain_text", Text: "JPNIC管理者連絡窓口の追加"},
		},
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*申請者* 管理者"},
			},
		},
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*Service* " + strconv.Itoa(serviceID)},
				{Type: "mrkdwn", Text: "*Name* " + input.Name + " (" + input.NameEn + ")"},
			},
		},
		slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "*追加状況*", false, false), nil, nil),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: changeTextJPNICByAdmin(core.JPNICAdmin{}, input),
			},
		},
		slack.NewDividerBlock(),
	))
}

func noticeAddJPNICTechByAdmin(serviceID int, input core.JPNICTech) {
	notification.Notification.Slack.PostMessage(config.Conf.Slack.Channels.Main, slack.MsgOptionBlocks(
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTHeader,
			Text: &slack.TextBlockObject{Type: "plain_text", Text: "JPNIC技術連絡担当者の追加"},
		},
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*申請者* 管理者"},
			},
		},
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*Service* " + strconv.Itoa(serviceID)},
				{Type: "mrkdwn", Text: "*Name* " + input.Name + " (" + input.NameEn + ")"},
			},
		},
		slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "*追加状況*", false, false), nil, nil),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: changeTextJPNICTech(core.JPNICTech{}, input),
			},
		},
		slack.NewDividerBlock(),
	))
}

func noticeAddIPByAdmin(serviceID int, input service.IPInput) {
	notification.Notification.Slack.PostMessage(config.Conf.Slack.Channels.Main, slack.MsgOptionBlocks(
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTHeader,
			Text: &slack.TextBlockObject{Type: "plain_text", Text: "IPの追加"},
		},
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*申請者* 管理者"},
			},
		},
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*Service* " + strconv.Itoa(serviceID)},
				{Type: "mrkdwn", Text: "*Name* " + input.Name},
				{Type: "mrkdwn", Text: "*IP* " + input.IP},
			},
		},
		slack.NewDividerBlock(),
	))
}

func noticeAddPlanByAdmin(ipID int, input core.Plan) {

	notification.Notification.Slack.PostMessage(config.Conf.Slack.Channels.Main, slack.MsgOptionBlocks(
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTHeader,
			Text: &slack.TextBlockObject{Type: "plain_text", Text: "Planの追加"},
		},
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*申請者* 管理者"},
			},
		},
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*IP* " + strconv.Itoa(ipID)},
				{Type: "mrkdwn", Text: "*Name* " + input.Name},
				{Type: "mrkdwn", Text: "*Plan* " + strconv.Itoa(int(input.After)) + "/" +
					strconv.Itoa(int(input.HalfYear)) + "/" + strconv.Itoa(int(input.OneYear))},
			},
		},
		slack.NewDividerBlock(),
	))
}

func noticeUpdateJPNICByAdmin(before, after core.JPNICAdmin) {
	notification.Notification.Slack.PostMessage(config.Conf.Slack.Channels.Main, slack.MsgOptionBlocks(
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTHeader,
			Text: &slack.TextBlockObject{Type: "plain_text", Text: "JPNIC管理者連絡窓口の更新"},
		},
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*申請者* 管理者"},
			},
		},
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*JPNICAdmin* " + strconv.Itoa(int(before.ID))},
			},
		},
		slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "*更新状況*", false, false), nil, nil),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: changeTextJPNICByAdmin(before, after),
			},
		},
		slack.NewDividerBlock(),
	))
}

func noticeUpdateJPNICTechByAdmin(before, after core.JPNICTech) {
	notification.Notification.Slack.PostMessage(config.Conf.Slack.Channels.Main, slack.MsgOptionBlocks(
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTHeader,
			Text: &slack.TextBlockObject{Type: "plain_text", Text: "JPNIC技術連絡担当者の更新"},
		},
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*申請者* 管理者"},
			},
		},
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*JPNICTech* " + strconv.Itoa(int(before.ID))},
			},
		},
		slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "*更新状況*", false, false), nil, nil),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: changeTextJPNICTech(before, after),
			},
		},
		slack.NewDividerBlock(),
	))
}

func noticeUpdateIPByAdmin(before, after core.IP) {
	notification.Notification.Slack.PostMessage(config.Conf.Slack.Channels.Main, slack.MsgOptionBlocks(
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTHeader,
			Text: &slack.TextBlockObject{Type: "plain_text", Text: "IPの更新"},
		},
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*申請者* 管理者"},
			},
		},
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*IP* " + strconv.Itoa(int(before.ID))},
			},
		},
		slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "*更新状況*", false, false), nil, nil),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: changeTextIP(before, after),
			},
		},
		slack.NewDividerBlock(),
	))
}

func noticeUpdatePlanByAdmin(before, after core.Plan) {
	notification.Notification.Slack.PostMessage(config.Conf.Slack.Channels.Main, slack.MsgOptionBlocks(
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTHeader,
			Text: &slack.TextBlockObject{Type: "plain_text", Text: "Planの更新"},
		},
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*申請者* 管理者"},
			},
		},
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				{Type: "mrkdwn", Text: "*Plan* " + strconv.Itoa(int(before.ID))},
			},
		},
		slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", "*更新状況*", false, false), nil, nil),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: changeTextPlan(before, after),
			},
		},
		slack.NewDividerBlock(),
	))
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

	if after.AddAllow != nil {
		if *before.AddAllow != *after.AddAllow {
			if *after.AddAllow {
				data += "ユーザ側にて接続追加の許可: 禁止 => 許可\n"
			} else {
				data += "ユーザ側にて接続追加の許可: 許可 => 禁止\n"
			}
		}
	}

	if after.ServiceType != "" {
		if before.ServiceType != after.ServiceType {
			data += "ServiceID: " + before.ServiceType + " => " + after.ServiceType + "\n"
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

	if before.V4JPNICHandle != after.V4JPNICHandle {
		data += "JPNICHandle(IPv4): " + before.V4JPNICHandle + "=>" + after.V4JPNICHandle + "\n"
	}

	if before.V6JPNICHandle != after.V6JPNICHandle {
		data += "JPNICHandle(IPv6): " + before.V6JPNICHandle + "=>" + after.V6JPNICHandle + "\n"
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

	if before.V4JPNICHandle != after.V4JPNICHandle {
		data += "JPNICHandle(IPv4): " + before.V4JPNICHandle + "=>" + after.V4JPNICHandle + "\n"
	}

	if before.V6JPNICHandle != after.V6JPNICHandle {
		data += "JPNICHandle(IPv6): " + before.V6JPNICHandle + "=>" + after.V6JPNICHandle + "\n"
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
