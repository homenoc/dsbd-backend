package v0

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	controllerInterface "github.com/homenoc/dsbd-backend/pkg/api/core/controller"
	controller "github.com/homenoc/dsbd-backend/pkg/api/core/controller/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/mail"
	"github.com/homenoc/dsbd-backend/pkg/api/core/mail/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/support"
	"github.com/homenoc/dsbd-backend/pkg/api/core/support/ticket"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbChat "github.com/homenoc/dsbd-backend/pkg/api/store/support/chat/v0"
	dbTicket "github.com/homenoc/dsbd-backend/pkg/api/store/support/ticket/v0"
	dbMailTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/mail/v0"
	dbUser "github.com/homenoc/dsbd-backend/pkg/api/store/user/v0"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"time"
)

func CreateByAdmin(c *gin.Context) {
	var input support.FirstInput

	// Admin authentication
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	err := c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	// input check
	if err = checkByAdmin(input); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	resultTicket := &core.Ticket{
		Solved:  &[]bool{false}[0],
		Title:   input.Title,
		Admin:   &[]bool{true}[0],
		Request: &[]bool{false}[0],
	}

	// isn't group
	if !input.IsGroup {
		if input.UserID == 0 {
			c.JSON(http.StatusBadRequest, common.Error{Error: "UserID is wrong"})
			return
		}

		resultTicket.GroupID = nil
		resultTicket.UserID = &input.UserID
	} else {
		//is group
		if input.UserID == 0 {
			c.JSON(http.StatusBadRequest, common.Error{Error: "GroupID is wrong"})
			return
		}

		resultTicket.GroupID = &input.GroupID
		resultTicket.UserID = nil
	}

	// Tickets DBに登録
	ticketResult, err := dbTicket.Create(resultTicket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	// Chat DBに登録
	chatResult, err := dbChat.Create(&core.Chat{
		UserID:   nil,
		Admin:    true,
		Data:     input.Data,
		TicketID: ticketResult.ID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, support.Result{
		Ticket: []core.Ticket{*ticketResult},
		Chat:   []core.Chat{*chatResult},
	})
}

func UpdateByAdmin(c *gin.Context) {
	var input core.Ticket
	// Admin authentication
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: fmt.Sprintf("id error")})
		return
	}

	err = c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	// Tickets DBからデータを取得
	ticketResult := dbTicket.Get(ticket.ID, &core.Ticket{Model: gorm.Model{ID: uint(id)}})
	if ticketResult.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: ticketResult.Err.Error()})
		return
	}

	// input check
	replace, err := updateAdminTicket(input, ticketResult.Tickets[0])
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	// Ticketのアップデート
	err = dbTicket.Update(ticket.UpdateAll, replace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, support.Result{})
}

func GetByAdmin(c *gin.Context) {
	// Admin authentication
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: fmt.Sprintf("id error")})
		return
	}

	// IDからDBからチケットを検索
	resultTicket := dbTicket.Get(ticket.ID, &core.Ticket{Model: gorm.Model{ID: uint(id)}})
	if resultTicket.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultTicket.Err.Error()})
		return
	}
	c.JSON(http.StatusOK, support.Result{Ticket: resultTicket.Tickets})
}

func GetAllByAdmin(c *gin.Context) {
	// Admin authentication
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	// Tickets DBからGroup IDのTicketデータを抽出
	resultTicket := dbTicket.GetAll()
	if resultTicket.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultTicket.Err.Error()})
		return
	}

	c.JSON(http.StatusOK, ticket.ResultAdminAll{Tickets: resultTicket.Tickets})
}

func GetAdminWebSocket(c *gin.Context) {
	//
	// /support?id=0?user_token=accessID?access_token=token
	// id = ticketID, access_token = AccessToken

	accessToken := c.Query("access_token")

	id, err := strconv.Atoi(c.Query("id"))
	if err != nil {
		log.Println("id wrong: ", err)
		return
	}
	//wsHandle(c.Writer, c.Request)
	conn, err := ticket.WsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v\n", err)
		return
	}

	defer conn.Close()

	// Admin authentication
	resultAdmin := auth.AdminAuthentication(accessToken)
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	ticketResult := dbTicket.Get(ticket.ID, &core.Ticket{Model: gorm.Model{ID: uint(id)}})
	if ticketResult.Err != nil {
		log.Println("ws:// support error: db error")
		conn.WriteMessage(websocket.TextMessage, []byte("error: db error"))
		return
	}

	var groupID uint = 0

	if ticketResult.Tickets[0].GroupID != nil {
		groupID = *ticketResult.Tickets[0].GroupID
	}

	// WebSocket送信
	support.Clients[&support.WebSocket{
		TicketID: uint(id),
		UserID:   resultAdmin.AdminID,
		UserName: "HomeNOC",
		GroupID:  groupID,
		Socket:   conn,
	}] = true

	//WebSocket受信
	for {
		var msg support.WebSocketResult
		err = conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(support.Clients, &support.WebSocket{
				TicketID: uint(id),
				UserID:   resultAdmin.AdminID,
				UserName: "HomeNOC(運営)",
				GroupID:  groupID,
				Socket:   conn,
			})
			break
		}

		_, err = dbChat.Create(&core.Chat{
			TicketID: ticketResult.Tickets[0].ID,
			UserID:   nil,
			Admin:    true,
			Data:     msg.Message,
		})
		if err != nil {
			conn.WriteJSON(&support.WebSocketResult{Err: "db write error"})
		} else {
			msg.TicketID = uint(id)
			msg.UserID = resultAdmin.AdminID
			msg.GroupID = groupID
			msg.UserName = "HomeNOC(運営)"
			msg.Admin = true
			// Token関連の初期化
			msg.AccessToken = ""
			msg.UserToken = ""

			//Admin側に送信
			controller.SendChatByAdmin(controllerInterface.Chat{
				TicketID:  uint(id),
				CreatedAt: msg.CreatedAt,
				Admin:     msg.Admin,
				UserID:    resultAdmin.AdminID,
				UserName:  msg.UserName,
				GroupID:   groupID,
				Message:   msg.Message,
			})

			resultTicket := dbTicket.Get(ticket.ID, &core.Ticket{Model: gorm.Model{ID: ticketResult.Tickets[0].ID}})
			if resultTicket.Err != nil {
				log.Println(resultTicket.Err)
			}
			mailTemplate := core.MailTemplate{ProcessID: "signature"}
			err = dbMailTemplate.Get(&mailTemplate)
			if err != nil {
				log.Println(err)
			}

			if len(resultTicket.Tickets) != 0 {
				if groupID != 0 {

					resultUser := dbUser.Get(user.GIDAndLevel, &core.User{
						GroupID: resultTicket.Tickets[0].GroupID,
						Level:   1,
					})
					if resultUser.Err != nil {
						log.Println(resultUser.Err)
					}

					if len(resultUser.User) != 0 {
						for _, userTmp := range resultUser.User {
							//グループ側にメール送信
							v0.SendMail(mail.Mail{
								ToMail:  userTmp.Email,
								Subject: "Supportより新着メッセージ",
								Content: " " + userTmp.Name + "様\n\n" + "チャットより新着メッセージがあります\n" +
									"Webシステムよりご覧いただけます。" + mailTemplate.Message,
							})
						}
					}
				} else {
					resultUser := dbUser.Get(user.ID, &core.User{
						Model: gorm.Model{ID: *resultTicket.Tickets[0].UserID},
					})
					if resultUser.Err != nil {
						log.Println(resultUser.Err)
					}

					if len(resultUser.User) != 0 {
						//グループ側にメール送信
						v0.SendMail(mail.Mail{
							ToMail:  resultUser.User[0].Email,
							Subject: "Supportより新着メッセージ",
							Content: " " + resultUser.User[0].Name + "様\n\n" + "チャットより新着メッセージがあります\n" +
								"Webシステムよりご覧いただけます。" + mailTemplate.Message,
						})
					}
				}

			}

			support.Broadcast <- msg
		}
	}
}

func HandleMessagesByAdmin() {
	for {
		msg := <-support.Broadcast

		//登録されているクライアント宛にデータ送信する
		for client := range support.Clients {
			// ユーザのみの場合
			log.Println(msg)
			if client.TicketID == msg.TicketID {
				if msg.GroupID == 0 {
					err := client.Socket.WriteJSON(support.WebSocketChatResponse{
						Time:     time.Now().UTC().Add(9 * time.Hour).Format(timeLayout),
						UserID:   msg.UserID,
						UserName: msg.UserName,
						GroupID:  0,
						Admin:    msg.Admin,
						Message:  msg.Message,
					})
					if err != nil {
						log.Printf("error: %v", err)
						client.Socket.Close()
						delete(support.Clients, client)
					}
				} else if client.GroupID == msg.GroupID {
					err := client.Socket.WriteJSON(support.WebSocketChatResponse{
						Time:     time.Now().UTC().Add(9 * time.Hour).Format(timeLayout),
						UserID:   msg.UserID,
						UserName: msg.UserName,
						GroupID:  msg.GroupID,
						Admin:    msg.Admin,
						Message:  msg.Message,
					})
					if err != nil {
						log.Printf("error: %v", err)
						client.Socket.Close()
						delete(support.Clients, client)
					}
				} else {
					// 認証失敗時の処理
				}
			}
		}
	}
}
