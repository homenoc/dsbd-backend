package v0

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/info"
	"github.com/homenoc/dsbd-backend/pkg/api/core/notice"
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

	authResult := auth.UserAuthentication(core.Token{UserToken: userToken, AccessToken: accessToken})
	if authResult.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: authResult.Err.Error()})
	}

	// User
	var resultUser info.User
	dbUserResult := dbUser.Get(user.IDDetail, &core.User{Model: gorm.Model{ID: authResult.User.ID}})
	if dbUserResult.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: dbUserResult.Err.Error()})
		return
	}
	resultUser = info.User{
		ID:         authResult.User.ID,
		GroupID:    authResult.User.GroupID,
		Name:       authResult.User.Name,
		NameEn:     authResult.User.NameEn,
		Email:      authResult.User.Email,
		Level:      authResult.User.Level,
		MailVerify: authResult.User.MailVerify,
	}

	//log.Println(*authResult.User.Group.PaymentCouponTemplateID)
	//log.Println(*authResult.User.Group.PaymentMembershipTemplateID)

	// Membership Info
	membershipInfo := "一般会員"
	membershipPlan := "未設定"
	paid := false
	automaticUpdate := false
	var discountRate uint = 0
	if authResult.User.Group.PaymentCouponTemplateID != nil && *authResult.User.Group.PaymentCouponTemplateID != 0 {
		membershipInfo = dbUserResult.User[0].Group.PaymentCouponTemplate.Title
		discountRate = dbUserResult.User[0].Group.PaymentCouponTemplate.DiscountRate
	}
	if authResult.User.Group.PaymentMembershipTemplateID != nil && *authResult.User.Group.PaymentMembershipTemplateID != 0 {
		membershipPlan = dbUserResult.User[0].Group.PaymentMembershipTemplate.Plan
		automaticUpdate = true
	}
	if dbUserResult.User[0].Group.MemberExpired != nil && dbUserResult.User[0].Group.MemberExpired.Unix() > time.Now().Unix() {
		paid = true
	}

	// Group and UserList
	var resultGroup info.Group
	var resultUserList []info.User

	if authResult.User.GroupID != 0 {
		resultGroup = info.Group{
			ID:                        authResult.User.Group.ID,
			Student:                   authResult.User.Group.Student,
			Fee:                       dbUserResult.User[0].Group.Fee,
			Pass:                      authResult.User.Group.Pass,
			Lock:                      authResult.User.Group.Lock,
			ExpiredStatus:             authResult.User.Group.ExpiredStatus,
			MemberInfo:                membershipInfo,
			Status:                    authResult.User.Group.Status,
			AutomaticUpdate:           automaticUpdate,
			Paid:                      &paid,
			DiscountRate:              discountRate,
			PaymentMembershipTemplate: membershipPlan,
			MemberExpired:             authResult.User.Group.MemberExpired,
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
					GroupID:    tmpUser.GroupID,
					Name:       tmpUser.Name,
					NameEn:     tmpUser.NameEn,
					Email:      tmpUser.Email,
					Level:      tmpUser.Level,
					MailVerify: tmpUser.MailVerify,
				})
			}
		}
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

	if dbUserResult.User[0].GroupID != 0 {
		for _, tmpTicket := range dbUserResult.User[0].Group.Tickets {
			var resultChat []info.Chat
			for _, tmpChat := range tmpTicket.Chat {
				resultChat = append(resultChat, info.Chat{
					CreatedAt: tmpChat.CreatedAt,
					TicketID:  tmpChat.TicketID,
					UserID:    tmpChat.UserID,
					Admin:     tmpChat.Admin,
					Data:      tmpChat.Data,
				})
			}

			if !*tmpTicket.Request {
				// Ticket
				resultTicket = append(resultTicket, info.Ticket{
					ID:        tmpTicket.ID,
					CreatedAt: tmpTicket.CreatedAt,
					GroupID:   tmpTicket.GroupID,
					UserID:    tmpTicket.UserID,
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
					GroupID:   tmpTicket.GroupID,
					UserID:    tmpTicket.UserID,
					Chat:      resultChat,
					Reject:    tmpTicket.RequestReject,
					Solved:    tmpTicket.Solved,
					Admin:     tmpTicket.Admin,
					Title:     tmpTicket.Title,
				})
			}
		}
		for _, tmpTicket := range dbUserResult.User[0].Ticket {
			var resultChat []info.Chat
			if tmpTicket.GroupID == 0 {
				for _, tmpChat := range tmpTicket.Chat {
					resultChat = append(resultChat, info.Chat{
						CreatedAt: tmpChat.CreatedAt,
						TicketID:  tmpChat.TicketID,
						UserID:    tmpChat.UserID,
						Admin:     tmpChat.Admin,
						Data:      tmpChat.Data,
					})
				}
				resultTicket = append(resultTicket, info.Ticket{
					ID:        tmpTicket.ID,
					CreatedAt: tmpTicket.CreatedAt,
					GroupID:   tmpTicket.GroupID,
					UserID:    tmpTicket.UserID,
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
	}

	// Info
	var resultInfo []info.Info
	var resultService []info.Service
	var resultConnection []info.Connection

	if authResult.User.GroupID != 0 {

		if !(0 < authResult.User.Level && authResult.User.Level <= 3) {
			c.JSON(http.StatusForbidden, common.Error{Error: "error: access is not permitted"})
			return
		}

		for _, tmpService := range dbUserResult.User[0].Group.Services {
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
					ServiceID: strconv.Itoa(int(tmpService.GroupID)) + "-" + tmpService.ServiceTemplate.Type +
						fmt.Sprintf("%03d", tmpService.ServiceNumber),
					ServiceType:    tmpService.ServiceTemplate.Type,
					NeedRoute:      *tmpService.ServiceTemplate.NeedRoute,
					NeedJPNIC:      *tmpService.ServiceTemplate.NeedJPNIC,
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
				serviceID := strconv.Itoa(int(tmpService.GroupID)) + "-" + tmpService.ServiceTemplate.Type +
					fmt.Sprintf("%03d", tmpService.ServiceNumber) + "-" + tmpConnection.ConnectionTemplate.Type +
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
					var fee string
					var v4, v6 []string
					if *tmpService.Fee == 0 {
						fee = "Free"
					} else {
						fee = strconv.Itoa(int(*tmpService.Fee)) + "円"
					}

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
						resultInfo = append(resultInfo, info.Info{
							ServiceID:  serviceID,
							Service:    tmpService.ServiceTemplate.Name,
							Assign:     *tmpService.ServiceTemplate.NeedJPNIC,
							ASN:        *tmpService.ASN,
							V4:         v4,
							V6:         v6,
							NOC:        tmpConnection.NOC.Name,
							NOCIP:      tmpConnection.TunnelEndPointRouterIP.IP,
							TermIP:     tmpConnection.TermIP,
							LinkV4Our:  tmpConnection.LinkV4Our,
							LinkV4Your: tmpConnection.LinkV4Your,
							LinkV6Our:  tmpConnection.LinkV6Our,
							LinkV6Your: tmpConnection.LinkV6Your,
							Fee:        fee,
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
