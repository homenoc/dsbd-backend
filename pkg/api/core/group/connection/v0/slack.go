package v0

import (
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/noc"
	"github.com/homenoc/dsbd-backend/pkg/api/core/noc/bgpRouter"
	"github.com/homenoc/dsbd-backend/pkg/api/core/noc/tunnelEndPointRouterIP"
	connectionTemplate "github.com/homenoc/dsbd-backend/pkg/api/core/template/connection"
	ntt "github.com/homenoc/dsbd-backend/pkg/api/core/template/ntt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	dbBGPRouter "github.com/homenoc/dsbd-backend/pkg/api/store/noc/bgpRouter/v0"
	dbTunnelEndPointRouterIP "github.com/homenoc/dsbd-backend/pkg/api/store/noc/tunnelEndPointRouterIP/v0"
	dbNOC "github.com/homenoc/dsbd-backend/pkg/api/store/noc/v0"
	dbConnectionTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/connection/v0"
	dbNTTTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/ntt/v0"
	"github.com/jinzhu/gorm"
	"strconv"
)

func noticeSlackAdmin(before, after core.Connection) {
	// 審査ステータスのSlack通知
	attachment := slack.Attachment{}

	attachment.AddField(slack.Field{Title: "Title", Value: "Connection情報の更新"}).
		AddField(slack.Field{Title: "申請者", Value: "管理者"}).
		AddField(slack.Field{Title: "Group", Value: strconv.Itoa(int(before.ID)) + "-" + before.Service.Group.Org}).
		AddField(slack.Field{Title: "更新状況", Value: changeText(before, after)})
	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})
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

	if after.Lock != nil {
		if *before.Lock != *after.Lock {
			if !*after.Lock {
				data += "ユーザ変更: 禁止 => 許可\n"
			} else {
				data += "ユーザ変更: 許可 => 禁止\n"
			}
		}
	}

	if after.ConnectionTemplateID != nil {
		if *before.ConnectionTemplateID != *after.ConnectionTemplateID {
			data += "接続ID: " + before.ConnectionTemplate.Type + " => " +
				connectionTemplateText(*after.ConnectionTemplateID) + "\n"
		}
	}

	if after.BGPRouterID != nil {
		if *before.BGPRouterID != *after.BGPRouterID {
			data += "BGPルータ: " + before.BGPRouter.HostName + " => " + bgpRouterText(*after.BGPRouterID) + "\n"
		}
	}

	if after.TunnelEndPointRouterIPID != nil {
		if *before.TunnelEndPointRouterIPID != *after.TunnelEndPointRouterIPID {
			data += "トンネルエンドポイントルータ: " + before.TunnelEndPointRouterIP.TunnelEndPointRouter.HostName + " " +
				before.TunnelEndPointRouterIP.IP + " => " +
				tunnelEndPointRouterIPText(*after.TunnelEndPointRouterIPID) + "\n"
		}
	}

	if after.NTTTemplateID != nil {
		if *before.NTTTemplateID != *after.NTTTemplateID {
			data += "インターネット接続: " + before.NTTTemplate.Name + " => " + nttTemplateText(*after.NTTTemplateID) + "\n"
		}
	}

	if after.NOCID != nil {
		if *before.NOCID != *after.NOCID {
			data += "希望NOC: " + before.NOC.Name + " => " + nocText(*after.NOCID) + "\n"
		}
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

func connectionTemplateText(status uint) string {
	result := dbConnectionTemplate.Get(connectionTemplate.ID, &core.ConnectionTemplate{Model: gorm.Model{ID: status}})
	return result.Connections[0].Type
}

func bgpRouterText(status uint) string {
	if status != 0 {
		result := dbBGPRouter.Get(bgpRouter.ID, &core.BGPRouter{Model: gorm.Model{ID: status}})
		return result.BGPRouter[0].HostName
	} else {
		return "なし"
	}
}

func nttTemplateText(status uint) string {
	result := dbNTTTemplate.Get(ntt.ID, &core.NTTTemplate{Model: gorm.Model{ID: status}})
	return result.NTTs[0].Name
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
