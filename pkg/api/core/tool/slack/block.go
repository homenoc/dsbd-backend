package slack

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbIP "github.com/homenoc/dsbd-backend/pkg/api/store/group/service/ip/v0"
	dbService "github.com/homenoc/dsbd-backend/pkg/api/store/group/service/v0"
	dbUser "github.com/homenoc/dsbd-backend/pkg/api/store/user/v0"
	"github.com/slack-go/slack"
	"gorm.io/gorm"
	"net"
	"strconv"
)

func errorProcess(blocks []slack.Block, err1, err2 string) slack.MsgOption {
	blocks = append(blocks, []slack.Block{
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: "*Error*\n" +
					err1 + "\n" +
					err2,
			},
		}}...)
	return slack.MsgOptionBlocks(blocks...)
}

func invalidCommand() slack.MsgOption {
	return slack.MsgOptionText(":scream: Invalid Command", false)
}

func getHelpInfo() slack.MsgOption {
	return slack.MsgOptionBlocks(
		&slack.SectionBlock{
			Type: slack.MBTHeader,
			Text: &slack.TextBlockObject{
				Type: "plain_text",
				Text: ":question:  HELP",
			},
		},
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: "*April 14, 2023*",
			},
		},
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: "*Common*\n" +
					"`/dsbd help` help page",
			},
		},
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: "*Search*\n" +
					"`/dsbd user {userID}` search UserID\n" +
					//"`/dsbd group {groupID}` search Group\n" +
					//"`/dsbd service {id or serviceID}` search Service\n" +
					//"`/dsbd connection {id or connectionID}` search Connection\n" +
					"`/dsbd asn {AS Number}` search AS Number\n" +
					"`/dsbd addr {ipaddr}` search IP Address",
			},
		},
		slack.NewDividerBlock(),
	)
}

func getUserInfo(userId uint) slack.MsgOption {
	blocks := []slack.Block{
		&slack.SectionBlock{
			Type: slack.MBTHeader,
			Text: &slack.TextBlockObject{
				Type: "plain_text",
				Text: ":bust_in_silhouette:  Search User",
			},
		},
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: "取得時間: *" + tool.GetNowStr() + "*",
			},
		},
	}
	userDetail := dbUser.Get(user.ID, &core.User{Model: gorm.Model{ID: userId}})
	if userDetail.Err != nil {
		return errorProcess(blocks, "データ取得エラー", userDetail.Err.Error())
	}

	groupIDStr := "None"
	if userDetail.User[0].GroupID != nil {
		groupIDStr = strconv.Itoa(int(*userDetail.User[0].GroupID))
	}
	blocks = append(blocks, []slack.Block{
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				&slack.TextBlockObject{
					Type: "mrkdwn",
					Text: "*Created by*\n" + tool.DateToStr(userDetail.User[0].CreatedAt),
				},
				&slack.TextBlockObject{
					Type: "mrkdwn",
					Text: "*Updated by*\n" + tool.DateToStr(userDetail.User[0].UpdatedAt),
				},
			},
		},
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				&slack.TextBlockObject{
					Type: "mrkdwn",
					Text: "*UserName:*\n" + userDetail.User[0].Name,
				},
				&slack.TextBlockObject{
					Type: "mrkdwn",
					Text: "*UserName(En):*\n" + userDetail.User[0].NameEn,
				},
			},
		},
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				&slack.TextBlockObject{
					Type: "mrkdwn",
					Text: "*E-Mail*\n" + userDetail.User[0].Email,
				},
				&slack.TextBlockObject{
					Type: "mrkdwn",
					Text: "*GroupID*\n" + groupIDStr,
				},
			},
		},
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: "plain_text",
				Text: "Dashboard Link",
			},
			Accessory: &slack.Accessory{
				ButtonElement: &slack.ButtonBlockElement{
					Type: "button",
					Text: &slack.TextBlockObject{
						Type: "plain_text",
						Text: "Go Page",
					},
					URL: config.Conf.Controller.Admin.ReturnURL + "/dashboard/user/" + strconv.Itoa(int(userDetail.User[0].ID)),
				},
			},
		},
		slack.NewDividerBlock(),
	}...)
	return slack.MsgOptionBlocks(blocks...)
}

func getASNInfo(asn int) slack.MsgOption {
	blocks := []slack.Block{
		&slack.SectionBlock{
			Type: slack.MBTHeader,
			Text: &slack.TextBlockObject{
				Type: "plain_text",
				Text: ":office:  Search ASN [" + strconv.Itoa(asn) + "]",
			},
		},
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: "取得時間: *" + tool.GetNowStr() + "*",
			},
		},
	}
	resultService := dbService.Get(service.ASN, &core.Service{ASN: tool.ToUintP(uint(asn))})
	if resultService.Err != nil {
		return errorProcess(blocks, "データ取得エラー", resultService.Err.Error())
	}

	for _, detailSer := range resultService.Service {
		serviceCode := strconv.Itoa(int(detailSer.GroupID)) + "-" + detailSer.ServiceType + fmt.Sprintf("%03d", detailSer.ServiceNumber)
		connection_blocks := []*slack.TextBlockObject{}
		ip_blocks := []*slack.TextBlockObject{}
		for _, detailConnection := range detailSer.Connection {
			connectionCode := serviceCode + "-" + detailConnection.ConnectionType + fmt.Sprintf("%03d", detailConnection.ConnectionNumber)
			connection_blocks = append(connection_blocks, []*slack.TextBlockObject{
				&slack.TextBlockObject{
					Type: "mrkdwn",
					Text: connectionCode,
				},
			}...)
		}
		for _, ip := range detailSer.IP {
			ip_blocks = append(ip_blocks, []*slack.TextBlockObject{
				&slack.TextBlockObject{
					Type: "mrkdwn",
					Text: ip.IP,
				},
			}...)
		}

		blocks = append(blocks, []slack.Block{
			&slack.SectionBlock{
				Type: slack.MBTSection,
				Text: &slack.TextBlockObject{
					Type: "mrkdwn",
					Text: "ServiceID: *" + serviceCode + "*",
				},
			},
			&slack.SectionBlock{
				Type: slack.MBTSection,
				Fields: []*slack.TextBlockObject{
					&slack.TextBlockObject{
						Type: "mrkdwn",
						Text: "*Created by*\n" + tool.DateToStr(detailSer.CreatedAt),
					},
					&slack.TextBlockObject{
						Type: "mrkdwn",
						Text: "*Updated by*\n" + tool.DateToStr(detailSer.UpdatedAt),
					},
				},
			},
			&slack.SectionBlock{
				Type: slack.MBTSection,
				Fields: []*slack.TextBlockObject{
					&slack.TextBlockObject{
						Type: "mrkdwn",
						Text: "*ID*\n" + strconv.Itoa(int(detailSer.ID)),
					},
					&slack.TextBlockObject{
						Type: "mrkdwn",
						Text: "*GroupID/Name*\n[" + strconv.Itoa(int(detailSer.GroupID)) + "] " + detailSer.Group.Org,
					},
				},
			},
			&slack.SectionBlock{
				Type: slack.MBTSection,
				Text: &slack.TextBlockObject{
					Type: "plain_text",
					Text: "　",
				},
				Accessory: &slack.Accessory{
					ButtonElement: &slack.ButtonBlockElement{
						Type: "button",
						Text: &slack.TextBlockObject{
							Type: "plain_text",
							Text: "Group Page",
						},
						URL: config.Conf.Controller.Admin.ReturnURL + "/dashboard/group/" + strconv.Itoa(int(detailSer.GroupID)),
					},
				},
			},
		}...)

		if len(ip_blocks) != 0 {
			blocks = append(blocks, []slack.Block{
				&slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: "mrkdwn",
						Text: "*IP一覧*",
					},
				},
				&slack.SectionBlock{
					Type:   slack.MBTSection,
					Fields: ip_blocks,
				}}...)
		}
		if len(connection_blocks) != 0 {
			blocks = append(blocks, []slack.Block{
				&slack.SectionBlock{
					Type: slack.MBTSection,
					Text: &slack.TextBlockObject{
						Type: "mrkdwn",
						Text: "*接続情報一覧*",
					},
				},
				&slack.SectionBlock{
					Type:   slack.MBTSection,
					Fields: connection_blocks,
				}}...)
		}

		blocks = append(blocks, []slack.Block{slack.NewDividerBlock()}...)
	}

	return slack.MsgOptionBlocks(blocks...)
}

func getAddrInfo(addr string) slack.MsgOption {
	blocks := []slack.Block{
		&slack.SectionBlock{
			Type: slack.MBTHeader,
			Text: &slack.TextBlockObject{
				Type: "plain_text",
				Text: ":office:  Search IP Address [" + addr + "]",
			},
		},
		slack.NewDividerBlock(),
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: "取得時間: *" + tool.GetNowStr() + "*",
			},
		},
	}
	resultIP := dbIP.GetAll()
	if resultIP.Err != nil {
		return errorProcess(blocks, "データ取得エラー", resultIP.Err.Error())
	}

	// check input value
	inputIP := net.ParseIP(addr)
	if inputIP == nil {
		return errorProcess(blocks, "入力エラー", "IPアドレスを入力してください(例: 1.1.1.1)")
	}
	// search IP
	var serviceID uint = 0
	for _, detailIP := range resultIP.IP {
		_, ipNet, err := net.ParseCIDR(detailIP.IP)
		if err != nil {
			continue
		}
		if ipNet.Contains(inputIP) {
			serviceID = detailIP.ServiceID
			break
		}
	}
	if serviceID == 0 {
		return errorProcess(blocks, "Not Found...", "一致するアドレスがありませんでした。("+addr+")")
	}
	resultService := dbService.Get(service.ID, &core.Service{Model: gorm.Model{ID: serviceID}})
	if resultService.Err != nil {
		return errorProcess(blocks, "データ取得エラー", resultService.Err.Error())
	}
	oneService := resultService.Service[0]
	serviceCode := strconv.Itoa(int(oneService.GroupID)) + "-" + oneService.ServiceType + fmt.Sprintf("%03d", oneService.ServiceNumber)
	connection_blocks := []*slack.TextBlockObject{}
	ip_blocks := []*slack.TextBlockObject{}
	for _, oneConnection := range oneService.Connection {
		connectionCode := serviceCode + "-" + oneConnection.ConnectionType + fmt.Sprintf("%03d", oneConnection.ConnectionNumber)
		connection_blocks = append(connection_blocks, []*slack.TextBlockObject{
			&slack.TextBlockObject{
				Type: "mrkdwn",
				Text: connectionCode,
			},
		}...)
	}
	for _, ip := range oneService.IP {
		ip_blocks = append(ip_blocks, []*slack.TextBlockObject{
			&slack.TextBlockObject{
				Type: "mrkdwn",
				Text: ip.IP,
			},
		}...)
	}

	blocks = append(blocks, []slack.Block{
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: "mrkdwn",
				Text: "ServiceID: *" + serviceCode + "*",
			},
		},
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				&slack.TextBlockObject{
					Type: "mrkdwn",
					Text: "*Created by*\n" + tool.DateToStr(oneService.CreatedAt),
				},
				&slack.TextBlockObject{
					Type: "mrkdwn",
					Text: "*Updated by*\n" + tool.DateToStr(oneService.UpdatedAt),
				},
			},
		},
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Fields: []*slack.TextBlockObject{
				&slack.TextBlockObject{
					Type: "mrkdwn",
					Text: "*ID*\n" + strconv.Itoa(int(oneService.ID)),
				},
				&slack.TextBlockObject{
					Type: "mrkdwn",
					Text: "*GroupID/Name*\n[" + strconv.Itoa(int(oneService.GroupID)) + "] " + oneService.Group.Org,
				},
			},
		},
		&slack.SectionBlock{
			Type: slack.MBTSection,
			Text: &slack.TextBlockObject{
				Type: "plain_text",
				Text: "　",
			},
			Accessory: &slack.Accessory{
				ButtonElement: &slack.ButtonBlockElement{
					Type: "button",
					Text: &slack.TextBlockObject{
						Type: "plain_text",
						Text: "Group Page",
					},
					URL: config.Conf.Controller.Admin.ReturnURL + "/dashboard/group/" + strconv.Itoa(int(oneService.GroupID)),
				},
			},
		},
	}...)

	if len(ip_blocks) != 0 {
		blocks = append(blocks, []slack.Block{
			&slack.SectionBlock{
				Type: slack.MBTSection,
				Text: &slack.TextBlockObject{
					Type: "mrkdwn",
					Text: "*IP一覧*",
				},
			},
			&slack.SectionBlock{
				Type:   slack.MBTSection,
				Fields: ip_blocks,
			}}...)
	}
	if len(connection_blocks) != 0 {
		blocks = append(blocks, []slack.Block{
			&slack.SectionBlock{
				Type: slack.MBTSection,
				Text: &slack.TextBlockObject{
					Type: "mrkdwn",
					Text: "*接続情報一覧*",
				},
			},
			&slack.SectionBlock{
				Type:   slack.MBTSection,
				Fields: connection_blocks,
			}}...)
	}

	blocks = append(blocks, []slack.Block{slack.NewDividerBlock()}...)

	return slack.MsgOptionBlocks(blocks...)
}
