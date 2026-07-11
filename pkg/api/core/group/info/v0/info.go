package v0

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/info"
	"github.com/homenoc/dsbd-backend/pkg/api/middleware"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/api/store/group/v0"
	dbNotice "github.com/homenoc/dsbd-backend/pkg/api/store/notice/v0"
	dbUser "github.com/homenoc/dsbd-backend/pkg/api/store/user/v0"
)

// The user-facing read surface used to be a single GET /info that returned a
// bootstrap blob mixing identity, owned resources, and a derived network view.
// It is now decomposed into one endpoint per resource; the projections below
// are shared so each endpoint keeps the exact wire shape the blob carried.

// authUser resolves the authenticated user, writing a 401 and returning ok=false
// on failure.
func authUser(c *gin.Context) (core.User, bool) {
	return middleware.CurrentUser(c), true
}

// canAccessGroup mirrors the blob's guard: a group-scoped resource is readable
// only by users at level 1..3. Returns ok=false (after writing 403) otherwise.
func canAccessGroup(c *gin.Context, u core.User) bool {
	if !core.CanViewGroup(u.Level) {
		c.JSON(http.StatusForbidden, common.Error{Error: "error: access is not permitted"})
		return false
	}
	return true
}

func projectUser(u core.User) info.User {
	var groupID uint = 0
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

// projectChat converts a ticket's chat log.
func projectChat(chats []core.Chat) []info.Chat {
	var result []info.Chat
	for _, tmpChat := range chats {
		var userID uint = 0
		if tmpChat.UserID != nil {
			userID = *tmpChat.UserID
		}
		result = append(result, info.Chat{
			CreatedAt: tmpChat.CreatedAt,
			TicketID:  tmpChat.TicketID,
			UserID:    userID,
			Admin:     tmpChat.Admin,
			Data:      tmpChat.Data,
		})
	}
	return result
}

// projectJPNICAdmin mirrors the blob's admin block, which intentionally omits
// Address/AddressEn (they stay empty).
func projectJPNICAdmin(j core.JPNICAdmin) info.JPNIC {
	return info.JPNIC{
		ID:       j.ID,
		Name:     j.Name,
		NameEn:   j.NameEn,
		Org:      j.Org,
		OrgEn:    j.OrgEn,
		PostCode: j.PostCode,
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

// serviceGraph holds the three derived projections that all iterate the group's
// services, computed once so the per-resource endpoints cannot diverge.
type serviceGraph struct {
	services    []info.Service
	connections []info.Connection
	infos       []info.Info
}

// buildServiceGraph reproduces the blob's service/connection/info derivation.
func buildServiceGraph(group core.Group) (serviceGraph, error) {
	var out serviceGraph

	for _, tmpService := range group.Services {
		resultServiceWithTemplate, err := core.GetServiceType(tmpService.ServiceType)
		if err != nil {
			return serviceGraph{}, err
		}

		if *tmpService.Enable {
			var resultServiceJPNICAdmin info.JPNIC
			var resultServiceJPNICTech []info.JPNIC
			var resultServiceIP []info.IP

			resultServiceJPNICAdmin = projectJPNICAdmin(tmpService.JPNICAdmin)

			for _, tmpJPNICTech := range tmpService.JPNICTech {
				resultServiceJPNICTech = append(resultServiceJPNICTech, projectJPNICTech(tmpJPNICTech))
			}

			for _, tmpIP := range tmpService.IP {
				if *tmpIP.Open {
					var resultIPPlan []info.Plan = nil
					if tmpIP.Plan != nil {
						for _, tmpIPPlan := range tmpIP.Plan {
							resultIPPlan = append(resultIPPlan, info.Plan{
								ID:       tmpIPPlan.ID,
								IPID:     tmpIPPlan.IPID,
								Name:     tmpIPPlan.Name,
								After:    tmpIPPlan.After,
								HalfYear: tmpIPPlan.HalfYear,
								OneYear:  tmpIPPlan.OneYear,
							})
						}
					}

					resultServiceIP = append(resultServiceIP, info.IP{
						ID:        tmpIP.ID,
						Version:   tmpIP.Version,
						Name:      tmpIP.Name,
						IP:        tmpIP.IP,
						Plan:      resultIPPlan,
						PlanJPNIC: "",
						UseCase:   tmpIP.UseCase,
					})
				}
			}

			out.services = append(out.services, info.Service{
				ID: tmpService.ID,
				ServiceID: strconv.Itoa(int(tmpService.GroupID)) + "-" + tmpService.ServiceType +
					fmt.Sprintf("%03d", tmpService.ServiceNumber),
				ServiceType:    tmpService.ServiceType,
				NeedRoute:      resultServiceWithTemplate.NeedRoute,
				NeedBGP:        resultServiceWithTemplate.NeedBGP,
				NeedJPNIC:      resultServiceWithTemplate.NeedJPNIC,
				AddAllow:       *tmpService.AddAllow,
				Pass:           *tmpService.Pass,
				Org:            tmpService.Org,
				OrgEn:          tmpService.OrgEn,
				PostCode:       tmpService.PostCode,
				Address:        tmpService.Address,
				AddressEn:      tmpService.AddressEn,
				ASN:            tmpService.ASN,
				AveUpstream:    tmpService.AveUpstream,
				MaxUpstream:    tmpService.MaxUpstream,
				AveDownstream:  tmpService.AveDownstream,
				MaxDownstream:  tmpService.MaxDownstream,
				MaxBandWidthAS: tmpService.MaxBandWidthAS,
				IP:             resultServiceIP,
				JPNICAdmin:     resultServiceJPNICAdmin,
				JPNICTech:      resultServiceJPNICTech,
			})
		}

		for _, tmpConnection := range tmpService.Connection {
			serviceID := strconv.Itoa(int(tmpService.GroupID)) + "-" + tmpService.ServiceType +
				fmt.Sprintf("%03d", tmpService.ServiceNumber) + "-" + tmpConnection.ConnectionType +
				fmt.Sprintf("%03d", tmpConnection.ConnectionNumber)

			if *tmpConnection.Enable {
				out.connections = append(out.connections, info.Connection{
					ID:        tmpConnection.ID,
					ServiceID: serviceID,
					Open:      *tmpConnection.Open,
				})
			}

			if *tmpService.Pass && *tmpService.Enable {
				var v4, v6 []string
				for _, tmpIP := range tmpService.IP {
					if *tmpIP.Open {
						if tmpIP.Version == 4 {
							v4 = append(v4, tmpIP.IP)
						} else if tmpIP.Version == 6 {
							v6 = append(v6, tmpIP.IP)
						}
					}
				}

				if *tmpConnection.Open && *tmpConnection.Enable {
					var asn uint = 0
					if tmpService.ASN != nil {
						asn = *tmpService.ASN
					}

					out.infos = append(out.infos, info.Info{
						ServiceID:  serviceID,
						Service:    resultServiceWithTemplate.Name,
						Assign:     resultServiceWithTemplate.NeedJPNIC,
						ASN:        asn,
						V4:         v4,
						V6:         v6,
						Fee:        "Free",
						NOC:        tmpConnection.BGPRouter.NOC.Name,
						NOCIP:      tmpConnection.TunnelEndPointRouterIP.IP,
						TermIP:     tmpConnection.TermIP,
						RFC8950:    tmpConnection.RFC8950,
						LinkV4Our:  tmpConnection.LinkV4Our,
						LinkV4Your: tmpConnection.LinkV4Your,
						LinkV6Our:  tmpConnection.LinkV6Our,
						LinkV6Your: tmpConnection.LinkV6Your,
					})
				}
			}
		}
	}

	return out, nil
}

// loadGroupScoped resolves the authed user, enforces the group guard, and loads
// the group's full association graph. ok=false means a response was already
// written. group is the zero value when the user has no group.
func loadGroupScoped(c *gin.Context) (user core.User, group core.Group, hasGroup bool, ok bool) {
	user, ok = authUser(c)
	if !ok {
		return
	}
	if user.GroupID == nil {
		return user, core.Group{}, false, true
	}
	if !canAccessGroup(c, user) {
		return user, core.Group{}, false, false
	}
	groupResult := dbGroup.GetByID(*user.GroupID)
	if groupResult.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: groupResult.Err.Error()})
		return user, core.Group{}, false, false
	}
	if len(groupResult.Group) == 0 {
		return user, core.Group{}, false, true
	}
	return user, groupResult.Group[0], true, true
}

// GetMe returns the authenticated user (GET /user/me).
func GetMe(c *gin.Context) {
	user, ok := authUser(c)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": projectUser(user)})
}

// GetGroup returns the user's group plus its member list (GET /group).
func GetGroup(c *gin.Context) {
	user, ok := authUser(c)
	if !ok {
		return
	}

	// No group yet: the member list is just the user themselves.
	if user.GroupID == nil {
		c.JSON(http.StatusOK, gin.H{
			"group":     info.Group{},
			"user_list": []info.User{projectUser(user)},
		})
		return
	}

	groupResult := dbGroup.GetByID(*user.GroupID)
	if groupResult.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: groupResult.Err.Error()})
		return
	}
	if len(groupResult.Group) == 0 {
		c.JSON(http.StatusOK, gin.H{"group": info.Group{}, "user_list": []info.User{}})
		return
	}
	group := groupResult.Group[0]

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

	resultGroup := info.Group{
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
	if core.CanManageServices(user.Level) {
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

	var resultUserList []info.User
	if core.CanViewGroup(user.Level) {
		for _, tmpUser := range group.Users {
			resultUserList = append(resultUserList, projectUser(tmpUser))
		}
	}

	c.JSON(http.StatusOK, gin.H{"group": resultGroup, "user_list": resultUserList})
}

// GetService returns the group's enabled services (GET /service).
func GetService(c *gin.Context) {
	_, group, hasGroup, ok := loadGroupScoped(c)
	if !ok {
		return
	}
	if !hasGroup {
		c.JSON(http.StatusOK, gin.H{"service": []info.Service{}})
		return
	}
	graph, err := buildServiceGraph(group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"service": graph.services})
}

// GetConnection returns the group's enabled connections (GET /connection).
func GetConnection(c *gin.Context) {
	_, group, hasGroup, ok := loadGroupScoped(c)
	if !ok {
		return
	}
	if !hasGroup {
		c.JSON(http.StatusOK, gin.H{"connection": []info.Connection{}})
		return
	}
	graph, err := buildServiceGraph(group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"connection": graph.connections})
}

// Get returns the derived active-connection network summary (GET /info).
func Get(c *gin.Context) {
	_, group, hasGroup, ok := loadGroupScoped(c)
	if !ok {
		return
	}
	if !hasGroup {
		c.JSON(http.StatusOK, gin.H{"info": []info.Info{}})
		return
	}
	graph, err := buildServiceGraph(group)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"info": graph.infos})
}

// GetNotice returns notices active for the user (GET /notice).
func GetNotice(c *gin.Context) {
	user, ok := authUser(c)
	if !ok {
		return
	}
	noticeResult := dbNotice.GetActiveForUser(user.ID)
	if noticeResult.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: noticeResult.Err.Error()})
		return
	}
	var resultNotice []info.Notice
	for _, tmpNotice := range noticeResult.Notice {
		resultNotice = append(resultNotice, info.Notice{
			StartTime: tmpNotice.StartTime,
			EndTime:   tmpNotice.EndTime,
			Everyone:  *tmpNotice.Everyone,
			Important: *tmpNotice.Important,
			Fault:     *tmpNotice.Fault,
			Info:      *tmpNotice.Info,
			Title:     tmpNotice.Title,
			Data:      tmpNotice.Data,
		})
	}
	c.JSON(http.StatusOK, gin.H{"notice": resultNotice})
}

// buildTickets projects the user's support tickets, splitting request tickets
// out. Reproduces the blob's group-ticket + personal-ticket merge.
func buildTickets(userDetail core.User) (tickets []info.Ticket, requests []info.Request) {
	if userDetail.GroupID != nil {
		for _, tmpTicket := range userDetail.Group.Tickets {
			resultChat := projectChat(tmpTicket.Chat)
			var groupID uint = 0
			if tmpTicket.GroupID != nil {
				groupID = *tmpTicket.GroupID
			}
			var userID uint = 0
			if tmpTicket.UserID != nil {
				userID = *tmpTicket.UserID
			}

			if !*tmpTicket.Request {
				tickets = append(tickets, info.Ticket{
					ID:        tmpTicket.ID,
					CreatedAt: tmpTicket.CreatedAt,
					GroupID:   groupID,
					UserID:    userID,
					Chat:      resultChat,
					Solved:    tmpTicket.Solved,
					Admin:     tmpTicket.Admin,
					Title:     tmpTicket.Title,
				})
			} else {
				requests = append(requests, info.Request{
					ID:        tmpTicket.ID,
					CreatedAt: tmpTicket.CreatedAt,
					GroupID:   groupID,
					UserID:    userID,
					Chat:      resultChat,
					Reject:    tmpTicket.RequestReject,
					Solved:    tmpTicket.Solved,
					Admin:     tmpTicket.Admin,
					Title:     tmpTicket.Title,
				})
			}
		}
	}

	for _, tmpTicket := range userDetail.Ticket {
		if tmpTicket.GroupID == nil {
			resultChat := projectChat(tmpTicket.Chat)
			tickets = append(tickets, info.Ticket{
				ID:        tmpTicket.ID,
				CreatedAt: tmpTicket.CreatedAt,
				GroupID:   0,
				UserID:    userDetail.ID,
				Chat:      resultChat,
				Solved:    tmpTicket.Solved,
				Title:     tmpTicket.Title,
			})
		}
	}

	sort.Slice(tickets, func(i, j int) bool {
		return tickets[i].ID < tickets[j].ID
	})
	return tickets, requests
}

// loadUserDetail resolves the authed user and loads the ticket association graph.
func loadUserDetail(c *gin.Context) (core.User, bool) {
	user, ok := authUser(c)
	if !ok {
		return core.User{}, false
	}
	userResult := dbUser.GetDetail(user.ID)
	if userResult.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: userResult.Err.Error()})
		return core.User{}, false
	}
	if len(userResult.User) == 0 {
		c.JSON(http.StatusInternalServerError, common.Error{Error: "user not found"})
		return core.User{}, false
	}
	return userResult.User[0], true
}

// GetTicket returns the user's support tickets (GET /ticket).
func GetTicket(c *gin.Context) {
	userDetail, ok := loadUserDetail(c)
	if !ok {
		return
	}
	tickets, _ := buildTickets(userDetail)
	c.JSON(http.StatusOK, gin.H{"ticket": tickets})
}

// GetRequest returns the user's requests (GET /request).
func GetRequest(c *gin.Context) {
	userDetail, ok := loadUserDetail(c)
	if !ok {
		return
	}
	_, requests := buildTickets(userDetail)
	c.JSON(http.StatusOK, gin.H{"request": requests})
}
