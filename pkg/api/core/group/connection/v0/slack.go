package v0

import (
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/noc"
	"github.com/homenoc/dsbd-backend/pkg/api/core/noc/bgpRouter"
	"github.com/homenoc/dsbd-backend/pkg/api/core/noc/tunnelEndPointRouterIP"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	dbBGPRouter "github.com/homenoc/dsbd-backend/pkg/api/store/noc/bgpRouter/v0"
	dbTunnelEndPointRouterIP "github.com/homenoc/dsbd-backend/pkg/api/store/noc/tunnelEndPointRouterIP/v0"
	dbNOC "github.com/homenoc/dsbd-backend/pkg/api/store/noc/v0"
	"github.com/slack-go/slack"
	"gorm.io/gorm"
	"strconv"
)

func noticeAdd(applicant, groupID, serviceCode, connectionCodeNew, connectionCodeComment string) {
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
				{Type: "mrkdwn", Text: "*サービスコード* " + serviceCode},
				{Type: "mrkdwn", Text: "*接続コード（新規発番）* " + connectionCodeNew},
				{Type: "mrkdwn", Text: "*接続コード（補足情報）* " + connectionCodeComment},
			},
		},
		slack.NewDividerBlock(),
	))
}

func noticeUpdateByAdmin(before, after core.Connection) {
	// 審査ステータスのSlack通知
	notification.Notification.Slack.PostMessage(config.Conf.Slack.Channels.Main, slack.MsgOptionBlocks(
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTHeader,
			Text: &slack.TextBlockObject{Type: "plain_text", Text: "接続情報の更新"},
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
				{Type: "mrkdwn", Text: "*Group* [" + strconv.Itoa(int(before.ID)) + "] " + before.Service.Group.Org + "(" + before.Service.Group.OrgEn + ")"},
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

func changeText(before, after core.Connection) string {
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

	if after.ConnectionType != "" {
		if before.ConnectionType != after.ConnectionType {
			data += "接続ID: " + before.ConnectionType + " => " + after.ConnectionType + "\n"
		}
	}

	if after.BGPRouterID != nil {
		if before.BGPRouterID == nil || *before.BGPRouterID != *after.BGPRouterID {
			data += "BGPルータ: " + before.BGPRouter.HostName + " => " + bgpRouterText(*after.BGPRouterID) + "\n"
		}
	}

	if after.TunnelEndPointRouterIPID != nil {
		if before.TunnelEndPointRouterIPID == nil || *before.TunnelEndPointRouterIPID != *after.TunnelEndPointRouterIPID {
			data += "トンネルエンドポイントルータ: " + before.TunnelEndPointRouterIP.TunnelEndPointRouter.HostName + " " +
				before.TunnelEndPointRouterIP.IP + " => " +
				tunnelEndPointRouterIPText(*after.TunnelEndPointRouterIPID) + "\n"
		}
	}

	if after.NTT != "" {
		data += "インターネット接続: " + before.NTT + " => " + after.NTT + "\n"

	}

	if after.TermIP != "" && after.TermIP != before.TermIP {
		data += "終端アドレス: " + before.TermIP + "=>" + after.TermIP + "\n"
	}

	if after.LinkV4Our != "" && after.LinkV4Our != before.LinkV4Our {
		data += "v4(HomeNOC側): " + before.LinkV4Our + "=>" + after.LinkV4Our + "\n"
	}

	if after.LinkV4Your != "" && after.LinkV4Your != before.LinkV4Your {
		data += "v4(相手団体側): " + before.LinkV4Your + "=>" + after.LinkV4Your + "\n"
	}

	if after.LinkV6Our != "" && after.LinkV6Our != before.LinkV6Our {
		data += "v6(HomeNOC側): " + before.LinkV6Our + "=>" + after.LinkV6Our + "\n"
	}

	if after.LinkV6Your != "" && after.LinkV6Your != before.LinkV6Your {
		data += "v6(相手団体側): " + before.LinkV6Your + "=>" + after.LinkV6Your + "\n"
	}

	return data
}

func bgpRouterText(status uint) string {
	if status != 0 {
		result := dbBGPRouter.Get(bgpRouter.ID, &core.BGPRouter{Model: gorm.Model{ID: status}})
		return result.BGPRouter[0].HostName
	} else {
		return "なし"
	}
}

func tunnelEndPointRouterIPText(status uint) string {
	if status != 0 {
		result := dbTunnelEndPointRouterIP.Get(tunnelEndPointRouterIP.ID,
			&core.TunnelEndPointRouterIP{Model: gorm.Model{ID: status}})
		return result.TunnelEndPointRouterIP[0].TunnelEndPointRouter.HostName + " " +
			result.TunnelEndPointRouterIP[0].IP
	} else {
		return "なし"
	}
}

func nocText(status uint) string {
	result := dbNOC.Get(noc.ID, &core.NOC{Model: gorm.Model{ID: status}})
	return result.NOC[0].Name
}
