package v0

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/info"
	"github.com/homenoc/dsbd-backend/pkg/api/core/notice"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbNotice "github.com/homenoc/dsbd-backend/pkg/api/store/notice/v0"
	dbUser "github.com/homenoc/dsbd-backend/pkg/api/store/user/v0"
	"gorm.io/gorm"
	"net/http"
	"sort"
	"strconv"
	"time"
)

func Get(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	authResult := auth.UserAuthorization(core.Token{UserToken: userToken, AccessToken: accessToken})
	if authResult.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: authResult.Err.Error()})
		return
	}

	// User
	var resultUser info.User
	dbUserResult := dbUser.Get(user.IDDetail, &core.User{Model: gorm.Model{ID: authResult.User.ID}})
	if dbUserResult.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: dbUserResult.Err.Error()})
		return
	}

	var groupID uint = 0
	if authResult.User.GroupID != nil {
		groupID = *authResult.User.GroupID
	}

	resultUser = info.User{
		ID:         authResult.User.ID,
		GroupID:    groupID,
		Name:       authResult.User.Name,
		NameEn:     authResult.User.NameEn,
		Email:      authResult.User.Email,
		Level:      authResult.User.Level,
		MailVerify: authResult.User.MailVerify,
	}

	//log.Println(*authResult.User.Group.PaymentCouponTemplateID)
	//log.Println(*authResult.User.Group.PaymentMembershipTemplateID)

	// Group and UserList
	var resultGroup info.Group
	var resultUserList []info.User

	if authResult.User.GroupID != nil {
		// Membership Info
		membership, err := core.GetMembershipTypeID(authResult.User.Group.MemberType)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
			return
		}

		// isExpired(課金確認)
		isExpired := false
		if authResult.User.Group.MemberType < 50 && authResult.User.Group.MemberExpired != nil {
			jst, err := time.LoadLocation("Asia/Tokyo")
			if err != nil {
				panic(err)
			}
			nowJST := time.Now().In(jst)
			if nowJST.Unix() > authResult.User.Group.MemberExpired.Add(time.Hour*24).Unix() {
				isExpired = true
			}
		} else if authResult.User.Group.MemberType < 50 && authResult.User.Group.MemberExpired == nil {
			isExpired = true
		}

		// isStripeID
		isStripeID := true
		if authResult.User.Group.StripeCustomerID == nil || *authResult.User.Group.StripeCustomerID == "" ||
			authResult.User.Group.StripeSubscriptionID == nil || *authResult.User.Group.StripeSubscriptionID == "" {
			isStripeID = false
		}

		// coupon
		couponID := ""
		if authResult.User.Group.CouponID != nil {
			couponID = *authResult.User.Group.CouponID
		}

		resultGroup = info.Group{
			ID:            authResult.User.Group.ID,
			Pass:          authResult.User.Group.Pass,
			ExpiredStatus: authResult.User.Group.ExpiredStatus,
			IsExpired:     isExpired,
			IsStripeID:    isStripeID,
			MemberTypeID:  membership.ID,
			MemberType:    membership.Name,
			MemberExpired: authResult.User.Group.MemberExpired,
			CouponID:      couponID,
		}
		if authResult.User.Level < 3 {
			resultGroup.Agree = dbUserResult.User[0].Group.Agree
			resultGroup.Question = dbUserResult.User[0].Group.Question
			resultGroup.Org = dbUserResult.User[0].Group.Org
			resultGroup.OrgEn = dbUserResult.User[0].Group.OrgEn
			resultGroup.PostCode = dbUserResult.User[0].Group.PostCode
			resultGroup.Address = dbUserResult.User[0].Group.Address
			resultGroup.AddressEn = dbUserResult.User[0].Group.AddressEn
			resultGroup.Tel = dbUserResult.User[0].Group.Tel
			resultGroup.Country = dbUserResult.User[0].Group.Country
			resultGroup.Contract = dbUserResult.User[0].Group.Contract
			resultGroup.AddAllow = dbUserResult.User[0].Group.AddAllow
		}

		if 0 < authResult.User.Level && authResult.User.Level <= 3 {
			for _, tmpUser := range dbUserResult.User[0].Group.Users {
				resultUserList = append(resultUserList, info.User{
					ID:         tmpUser.ID,
					GroupID:    *tmpUser.GroupID,
					Name:       tmpUser.Name,
					NameEn:     tmpUser.NameEn,
					Email:      tmpUser.Email,
					Level:      tmpUser.Level,
					MailVerify: tmpUser.MailVerify,
				})
			}
		}
	} else {
		// GroupID == 0
		resultUserList = append(resultUserList, resultUser)
	}

	// Notice
	var resultNotice []info.Notice
	noticeResult := dbNotice.Get(notice.UIDOrAll, &core.Notice{User: []core.User{{Model: gorm.Model{ID: authResult.User.ID}}}})
	if noticeResult.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: noticeResult.Err.Error()})
		return
	}
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

	// Ticket/Request
	var resultTicket []info.Ticket
	var resultRequest []info.Request

	if dbUserResult.User[0].GroupID != nil {
		for _, tmpTicket := range dbUserResult.User[0].Group.Tickets {
			var resultChat []info.Chat
			for _, tmpChat := range tmpTicket.Chat {
				var userID uint = 0
				if tmpChat.UserID != nil {
					userID = *tmpChat.UserID
				}

				resultChat = append(resultChat, info.Chat{
					CreatedAt: tmpChat.CreatedAt,
					TicketID:  tmpChat.TicketID,
					UserID:    userID,
					Admin:     tmpChat.Admin,
					Data:      tmpChat.Data,
				})
			}
			var groupIDTicketAndRequest uint = 0
			if tmpTicket.GroupID != nil {
				groupIDTicketAndRequest = *tmpTicket.GroupID
			}
			var userIDTicketAndRequest uint = 0
			if tmpTicket.UserID != nil {
				userIDTicketAndRequest = *tmpTicket.UserID
			}

			if !*tmpTicket.Request {
				// Ticket
				resultTicket = append(resultTicket, info.Ticket{
					ID:        tmpTicket.ID,
					CreatedAt: tmpTicket.CreatedAt,
					GroupID:   groupIDTicketAndRequest,
					UserID:    userIDTicketAndRequest,
					Chat:      resultChat,
					Solved:    tmpTicket.Solved,
					Admin:     tmpTicket.Admin,
					Title:     tmpTicket.Title,
				})
			} else {
				// Request
				resultRequest = append(resultRequest, info.Request{
					ID:        tmpTicket.ID,
					CreatedAt: tmpTicket.CreatedAt,
					GroupID:   groupIDTicketAndRequest,
					UserID:    userIDTicketAndRequest,
					Chat:      resultChat,
					Reject:    tmpTicket.RequestReject,
					Solved:    tmpTicket.Solved,
					Admin:     tmpTicket.Admin,
					Title:     tmpTicket.Title,
				})
			}
		}
	}

	for _, tmpTicket := range dbUserResult.User[0].Ticket {
		var resultChat []info.Chat
		if tmpTicket.GroupID == nil {
			for _, tmpChat := range tmpTicket.Chat {

				var userID uint = 0
				if tmpChat.UserID != nil {
					userID = *tmpChat.UserID
				}

				resultChat = append(resultChat, info.Chat{
					CreatedAt: tmpChat.CreatedAt,
					TicketID:  tmpChat.TicketID,
					UserID:    userID,
					Admin:     tmpChat.Admin,
					Data:      tmpChat.Data,
				})
			}

			resultTicket = append(resultTicket, info.Ticket{
				ID:        tmpTicket.ID,
				CreatedAt: tmpTicket.CreatedAt,
				GroupID:   0,
				UserID:    authResult.User.ID,
				Chat:      resultChat,
				Solved:    tmpTicket.Solved,
				Title:     tmpTicket.Title,
			})
		}
	}
	sort.Slice(resultTicket, func(i, j int) bool {
		if resultTicket[i].ID < resultTicket[j].ID {
			return true
		}
		return false
	})

	// Info
	var resultInfo []info.Info
	var resultService []info.Service
	var resultConnection []info.Connection

	if authResult.User.GroupID != nil {

		if !(0 < authResult.User.Level && authResult.User.Level <= 3) {
			c.JSON(http.StatusForbidden, common.Error{Error: "error: access is not permitted"})
			return
		}

		for _, tmpService := range dbUserResult.User[0].Group.Services {
			// getting service detail info
			resultServiceWithTemplate, err := config.GetServiceTemplate(tmpService.ServiceType)
			if err != nil {
				c.JSON(http.StatusInternalServerError, common.Error{Error: dbUserResult.Err.Error()})
				return
			}

			// Service
			if *tmpService.Enable {
				var resultServiceJPNICAdmin info.JPNIC
				var resultServiceJPNICTech []info.JPNIC
				var resultServiceIP []info.IP

				// JPNIC Admin
				resultServiceJPNICAdmin.ID = tmpService.JPNICAdmin.ID
				resultServiceJPNICAdmin.Org = tmpService.JPNICAdmin.Org
				resultServiceJPNICAdmin.OrgEn = tmpService.JPNICAdmin.OrgEn
				resultServiceJPNICAdmin.PostCode = tmpService.JPNICAdmin.PostCode
				resultServiceJPNICAdmin.Name = tmpService.JPNICAdmin.Name
				resultServiceJPNICAdmin.NameEn = tmpService.JPNICAdmin.NameEn
				resultServiceJPNICAdmin.Dept = tmpService.JPNICAdmin.Dept
				resultServiceJPNICAdmin.DeptEn = tmpService.JPNICAdmin.DeptEn
				resultServiceJPNICAdmin.Tel = tmpService.JPNICAdmin.Tel
				resultServiceJPNICAdmin.Fax = tmpService.JPNICAdmin.Fax
				resultServiceJPNICAdmin.Mail = tmpService.JPNICAdmin.Mail
				resultServiceJPNICAdmin.Country = tmpService.JPNICAdmin.Country

				// JPNIC Tech
				for _, tmpJPNICTech := range tmpService.JPNICTech {
					resultServiceJPNICTech = append(resultServiceJPNICTech, info.JPNIC{
						ID:        tmpJPNICTech.ID,
						Name:      tmpJPNICTech.Name,
						NameEn:    tmpJPNICTech.NameEn,
						Org:       tmpJPNICTech.Org,
						OrgEn:     tmpJPNICTech.OrgEn,
						PostCode:  tmpJPNICTech.PostCode,
						Address:   tmpJPNICTech.Address,
						AddressEn: tmpJPNICTech.AddressEn,
						Dept:      tmpJPNICTech.Dept,
						DeptEn:    tmpJPNICTech.DeptEn,
						Tel:       tmpJPNICTech.Tel,
						Fax:       tmpJPNICTech.Fax,
						Mail:      tmpJPNICTech.Mail,
						Country:   tmpJPNICTech.Country,
					})
				}

				// IP
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

				resultService = append(resultService, info.Service{
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

				// Connection
				if *tmpConnection.Enable {
					resultConnection = append(resultConnection, info.Connection{
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

						resultInfo = append(resultInfo, info.Info{
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
							LinkV4Our:  tmpConnection.LinkV4Our,
							LinkV4Your: tmpConnection.LinkV4Your,
							LinkV6Our:  tmpConnection.LinkV6Our,
							LinkV6Your: tmpConnection.LinkV6Your,
						})
					}
				}
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
