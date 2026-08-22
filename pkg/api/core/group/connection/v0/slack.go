package v0

import (
	"strconv"

	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	"github.com/homenoc/dsbd-backend/pkg/api/notify"
	dbBGPRouter "github.com/homenoc/dsbd-backend/pkg/api/store/noc/bgpRouter/v0"
	dbTunnelEndPointRouterIP "github.com/homenoc/dsbd-backend/pkg/api/store/noc/tunnelEndPointRouterIP/v0"
	dbNOC "github.com/homenoc/dsbd-backend/pkg/api/store/noc/v0"
	"github.com/slack-go/slack"
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
				Text: notify.BlockText(changeText(before, after)),
			},
		},
		slack.NewDividerBlock(),
	))
}

// changeText summarises the connection fields that changed. Scalar fields are
// driven by the `notify:"..."` tags on core.Connection; the BGP router and
// tunnel-endpoint IP need related-entity hostname resolution, so they stay as
// custom lines here.
func changeText(before, after core.Connection) string {
	data := notify.Diff(before, after)

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

	return data
}

func bgpRouterText(status uint) string {
	if status != 0 {
		result := dbBGPRouter.GetByID(status)
		return result.BGPRouter[0].HostName
	} else {
		return "なし"
	}
}

func tunnelEndPointRouterIPText(status uint) string {
	if status != 0 {
		result := dbTunnelEndPointRouterIP.GetByID(status)
		return result.TunnelEndPointRouterIP[0].TunnelEndPointRouter.HostName + " " +
			result.TunnelEndPointRouterIP[0].IP
	} else {
		return "なし"
	}
}

func nocText(status uint) string {
	result := dbNOC.GetByID(status)
	return result.NOC[0].Name
}
