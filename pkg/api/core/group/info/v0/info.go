package v0

import (
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/info"
	"github.com/homenoc/dsbd-backend/pkg/api/middleware"
	dbNotice "github.com/homenoc/dsbd-backend/pkg/api/store/notice/v0"
	dbUser "github.com/homenoc/dsbd-backend/pkg/api/store/user/v0"
)

// Get returns the user-facing bootstrap payload (GET /info): the caller's
// profile, group, member list, services, connections, notices, tickets,
// requests, and the derived network summary of opened connections.
func Get(c *gin.Context) {
	u := middleware.CurrentUser(c)

	// One GetDetail load carries the whole association graph this payload
	// needs (group -> services -> IP/connections/routers, plus tickets).
	userResult := dbUser.GetDetail(u.ID)
	if userResult.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: userResult.Err.Error()})
		return
	}
	if len(userResult.User) == 0 {
		c.JSON(http.StatusInternalServerError, common.Error{Error: "user not found"})
		return
	}
	userDetail := userResult.User[0]

	resultUser := projectUser(u)

	// Group and member list
	var resultGroup info.Group
	var resultUserList []info.User
	var group core.Group
	hasGroup := false

	if userDetail.GroupID != nil && userDetail.Group != nil {
		group = *userDetail.Group
		hasGroup = true
	}

	if hasGroup {
		membership, err := core.GetMembershipTypeID(group.MemberType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
			return
		}

		// isExpired (課金確認)
		isExpired := false
		if core.IsPaidMemberType(group.MemberType) && group.MemberExpired != nil {
			jst, err := time.LoadLocation("Asia/Tokyo")
			if err != nil {
				panic(err)
			}
			nowJST := time.Now().In(jst)
			if nowJST.Unix() > group.MemberExpired.Add(time.Hour*24).Unix() {
				isExpired = true
			}
		} else if core.IsPaidMemberType(group.MemberType) && group.MemberExpired == nil {
			isExpired = true
		}

		// isStripeID
		isStripeID := true
		if group.StripeCustomerID == nil || *group.StripeCustomerID == "" ||
			group.StripeSubscriptionID == nil || *group.StripeSubscriptionID == "" {
			isStripeID = false
		}

		couponID := ""
		if group.CouponID != nil {
			couponID = *group.CouponID
		}

		resultGroup = info.Group{
			ID:            group.ID,
			Pass:          group.Pass,
			ExpiredStatus: group.ExpiredStatus,
			IsExpired:     isExpired,
			IsStripeID:    isStripeID,
			MemberTypeID:  membership.ID,
			MemberType:    membership.Name,
			MemberExpired: group.MemberExpired,
			CouponID:      couponID,
		}
		if core.CanManageServices(u.Level) {
			resultGroup.Agree = group.Agree
			resultGroup.Question = group.Question
			resultGroup.Org = group.Org
			resultGroup.OrgEn = group.OrgEn
			resultGroup.PostCode = group.PostCode
			resultGroup.Address = group.Address
			resultGroup.AddressEn = group.AddressEn
			resultGroup.Tel = group.Tel
			resultGroup.Country = group.Country
			resultGroup.Contract = group.Contract
			resultGroup.AddAllow = group.AddAllow
		}

		if core.CanViewGroup(u.Level) {
			for _, member := range group.Users {
				resultUserList = append(resultUserList, projectUser(member))
			}
		}
	} else {
		// No group yet: the member list is just the user themselves.
		resultUserList = append(resultUserList, resultUser)
	}

	// Notice
	noticeResult := dbNotice.GetActiveForUser(u.ID)
	if noticeResult.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: noticeResult.Err.Error()})
		return
	}
	var resultNotice []info.Notice
	for _, n := range noticeResult.Notice {
		resultNotice = append(resultNotice, info.Notice{
			StartTime: n.StartTime,
			EndTime:   n.EndTime,
			Everyone:  *n.Everyone,
			Important: *n.Important,
			Fault:     *n.Fault,
			Info:      *n.Info,
			Title:     n.Title,
			Data:      n.Data,
		})
	}

	// Ticket / Request
	resultTicket, resultRequest := buildTickets(userDetail)

	// Service / Connection / derived network info (group-scoped; level-gated)
	var resultService []info.Service
	var resultConnection []info.Connection
	var resultInfo []info.Info

	if u.GroupID != nil {
		if !core.CanViewGroup(u.Level) {
			c.JSON(http.StatusForbidden, common.Error{Error: "error: access is not permitted"})
			return
		}
		if hasGroup {
			var err error
			resultService, resultConnection, resultInfo, err = buildServiceGraph(group)
			if err != nil {
				c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
				return
			}
		}
	}

	c.JSON(http.StatusOK, info.Result{
		User:       resultUser,
		Group:      resultGroup,
		UserList:   resultUserList,
		Service:    resultService,
		Connection: resultConnection,
		Notice:     resultNotice,
		Ticket:     resultTicket,
		Request:    resultRequest,
		Info:       resultInfo,
	})
}

func projectUser(u core.User) info.User {
	var groupID uint
	if u.GroupID != nil {
		groupID = *u.GroupID
	}
	return info.User{
		ID:                u.ID,
		GroupID:           groupID,
		Name:              u.Name,
		NameEn:            u.NameEn,
		Email:             u.Email,
		Level:             u.Level,
		MailVerify:        u.MailVerify,
		AntisocialCheck:   u.AntisocialCheck,
		AntisocialCheckAt: u.AntisocialCheckAt,
	}
}

// The admin contact intentionally omits its address fields on the wire.
func projectJPNICAdmin(j core.JPNICAdmin) info.JPNIC {
	return info.JPNIC{
		ID:       j.ID,
		Org:      j.Org,
		OrgEn:    j.OrgEn,
		PostCode: j.PostCode,
		Name:     j.Name,
		NameEn:   j.NameEn,
		Dept:     j.Dept,
		DeptEn:   j.DeptEn,
		Tel:      j.Tel,
		Fax:      j.Fax,
		Mail:     j.Mail,
		Country:  j.Country,
	}
}

func projectJPNICTech(j core.JPNICTech) info.JPNIC {
	return info.JPNIC{
		ID:        j.ID,
		Name:      j.Name,
		NameEn:    j.NameEn,
		Org:       j.Org,
		OrgEn:     j.OrgEn,
		PostCode:  j.PostCode,
		Address:   j.Address,
		AddressEn: j.AddressEn,
		Dept:      j.Dept,
		DeptEn:    j.DeptEn,
		Tel:       j.Tel,
		Fax:       j.Fax,
		Mail:      j.Mail,
		Country:   j.Country,
	}
}

// buildServiceGraph derives the three service-scoped projections in one pass
// so they cannot diverge.
func buildServiceGraph(g core.Group) ([]info.Service, []info.Connection, []info.Info, error) {
	var services []info.Service
	var connections []info.Connection
	var infos []info.Info

	for _, s := range g.Services {
		serviceType, err := core.GetServiceType(s.ServiceType)
		if err != nil {
			return nil, nil, nil, err
		}
		serviceCode := core.ServiceCode(s.GroupID, s.ServiceType, s.ServiceNumber)

		if *s.Enable {
			admin := projectJPNICAdmin(s.JPNICAdmin)

			var tech []info.JPNIC
			for _, t := range s.JPNICTech {
				tech = append(tech, projectJPNICTech(t))
			}

			var ips []info.IP
			for _, ip := range s.IP {
				if !*ip.Open {
					continue
				}
				var plans []info.Plan
				for _, plan := range ip.Plan {
					plans = append(plans, info.Plan{
						ID:       plan.ID,
						IPID:     plan.IPID,
						Name:     plan.Name,
						After:    plan.After,
						HalfYear: plan.HalfYear,
						OneYear:  plan.OneYear,
					})
				}
				ips = append(ips, info.IP{
					ID:      ip.ID,
					Version: ip.Version,
					Name:    ip.Name,
					IP:      ip.IP,
					Plan:    plans,
					UseCase: ip.UseCase,
				})
			}

			services = append(services, info.Service{
				ID:             s.ID,
				ServiceID:      serviceCode,
				ServiceType:    s.ServiceType,
				NeedRoute:      serviceType.NeedRoute,
				NeedBGP:        serviceType.NeedBGP,
				NeedJPNIC:      serviceType.NeedJPNIC,
				AddAllow:       *s.AddAllow,
				Pass:           *s.Pass,
				Org:            s.Org,
				OrgEn:          s.OrgEn,
				PostCode:       s.PostCode,
				Address:        s.Address,
				AddressEn:      s.AddressEn,
				ASN:            s.ASN,
				AveUpstream:    s.AveUpstream,
				MaxUpstream:    s.MaxUpstream,
				AveDownstream:  s.AveDownstream,
				MaxDownstream:  s.MaxDownstream,
				MaxBandWidthAS: s.MaxBandWidthAS,
				IP:             ips,
				JPNICAdmin:     admin,
				JPNICTech:      tech,
			})
		}

		var v4, v6 []string
		for _, ip := range s.IP {
			if !*ip.Open {
				continue
			}
			switch ip.Version {
			case 4:
				v4 = append(v4, ip.IP)
			case 6:
				v6 = append(v6, ip.IP)
			}
		}

		for _, conn := range s.Connection {
			connectionCode := core.ConnectionCode(serviceCode, conn.ConnectionType, conn.ConnectionNumber)

			if *conn.Enable {
				connections = append(connections, info.Connection{
					ID:        conn.ID,
					ServiceID: connectionCode,
					Open:      *conn.Open,
				})
			}

			if *s.Pass && *s.Enable && *conn.Open && *conn.Enable {
				var asn uint
				if s.ASN != nil {
					asn = *s.ASN
				}
				infos = append(infos, info.Info{
					ServiceID:  connectionCode,
					Service:    serviceType.Name,
					Assign:     serviceType.NeedJPNIC,
					ASN:        asn,
					V4:         v4,
					V6:         v6,
					Fee:        "Free",
					NOC:        conn.BGPRouter.NOC.Name,
					NOCIP:      conn.TunnelEndPointRouterIP.IP,
					TermIP:     conn.TermIP,
					RFC8950:    conn.RFC8950,
					LinkV4Our:  conn.LinkV4Our,
					LinkV4Your: conn.LinkV4Your,
					LinkV6Our:  conn.LinkV6Our,
					LinkV6Your: conn.LinkV6Your,
				})
			}
		}
	}
	return services, connections, infos, nil
}

// buildTickets projects the user's support tickets, splitting request tickets
// out. Group tickets and the user's personal (group-less) tickets are merged,
// sorted by ID.
func buildTickets(userDetail core.User) (tickets []info.Ticket, requests []info.Request) {
	if userDetail.GroupID != nil {
		for _, t := range userDetail.Group.Tickets {
			chat := projectChats(t.Chat)
			var groupID uint
			if t.GroupID != nil {
				groupID = *t.GroupID
			}
			var userID uint
			if t.UserID != nil {
				userID = *t.UserID
			}

			if !*t.Request {
				tickets = append(tickets, info.Ticket{
					ID:        t.ID,
					CreatedAt: t.CreatedAt,
					GroupID:   groupID,
					UserID:    userID,
					Chat:      chat,
					Solved:    t.Solved,
					Admin:     t.Admin,
					Title:     t.Title,
				})
			} else {
				requests = append(requests, info.Request{
					ID:        t.ID,
					CreatedAt: t.CreatedAt,
					GroupID:   groupID,
					UserID:    userID,
					Chat:      chat,
					Reject:    t.RequestReject,
					Solved:    t.Solved,
					Admin:     t.Admin,
					Title:     t.Title,
				})
			}
		}
	}

	for _, t := range userDetail.Ticket {
		if t.GroupID == nil {
			tickets = append(tickets, info.Ticket{
				ID:        t.ID,
				CreatedAt: t.CreatedAt,
				GroupID:   0,
				UserID:    userDetail.ID,
				Chat:      projectChats(t.Chat),
				Solved:    t.Solved,
				Title:     t.Title,
			})
		}
	}

	sort.Slice(tickets, func(i, j int) bool {
		return tickets[i].ID < tickets[j].ID
	})
	return tickets, requests
}

func projectChats(chats []core.Chat) []info.Chat {
	var out []info.Chat
	for _, ch := range chats {
		var userID uint
		if ch.UserID != nil {
			userID = *ch.UserID
		}
		out = append(out, info.Chat{
			CreatedAt: ch.CreatedAt,
			TicketID:  ch.TicketID,
			UserID:    userID,
			Admin:     ch.Admin,
			Data:      ch.Data,
		})
	}
	return out
}
