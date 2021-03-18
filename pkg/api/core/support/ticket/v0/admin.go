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
	"github.com/homenoc/dsbd-backend/pkg/api/core/support"
	"github.com/homenoc/dsbd-backend/pkg/api/core/support/ticket"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/mail"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbChat "github.com/homenoc/dsbd-backend/pkg/api/store/support/chat/v0"
	dbTicket "github.com/homenoc/dsbd-backend/pkg/api/store/support/ticket/v0"
	dbUser "github.com/homenoc/dsbd-backend/pkg/api/store/user/v0"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
	"time"
)

func CreateAdmin(c *gin.Context) {
	var input support.FirstInput

	// Admin authentication
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	err := c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	// input check
	if err = checkAdmin(input); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	// Tickets DBに登録
	ticketResult, err := dbTicket.Create(&core.Ticket{
		GroupID: input.GroupID,
		UserID:  0,
		Solved:  &[]bool{false}[0],
		Title:   input.Title,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	// Chat DBに登録
	chatResult, err := dbChat.Create(&core.Chat{
		UserID:   0,
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

func UpdateAdmin(c *gin.Context) {
	var input core.Ticket
	// Admin authentication
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultAdmin.Err.Error()})
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

func GetAdmin(c *gin.Context) {
	// Admin authentication
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultAdmin.Err.Error()})
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

func GetAllAdmin(c *gin.Context) {
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
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	ticketResult := dbTicket.Get(ticket.ID, &core.Ticket{Model: gorm.Model{ID: uint(id)}})
	if ticketResult.Err != nil {
		log.Println("ws:// support error: db error")
		conn.WriteMessage(websocket.TextMessage, []byte("error: db error"))
		return
	}

	// WebSocket送信
	support.Clients[&support.WebSocket{
		TicketID: uint(id),
		UserID:   resultAdmin.AdminID,
		UserName: "HomeNOC",
		GroupID:  ticketResult.Tickets[0].GroupID,
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
				GroupID:  ticketResult.Tickets[0].GroupID,
				Socket:   conn,
			})
			break
		}

		_, err = dbChat.Create(&core.Chat{
			TicketID: ticketResult.Tickets[0].ID,
			UserID:   resultAdmin.AdminID,
			Admin:    true,
			Data:     msg.Message,
		})
		if err != nil {
			conn.WriteJSON(&support.WebSocketResult{Err: "db write error"})
		} else {
			msg.UserID = resultAdmin.AdminID
			msg.GroupID = ticketResult.Tickets[0].GroupID
			msg.UserName = "HomeNOC(運営)"
			msg.Admin = true
			// Token関連の初期化
			msg.AccessToken = ""
			msg.UserToken = ""

			//Admin側に送信
			controller.SendChatAdmin(controllerInterface.Chat{
				CreatedAt: msg.CreatedAt,
				Admin:     msg.Admin,
				UserID:    resultAdmin.AdminID,
				UserName:  msg.UserName,
				GroupID:   ticketResult.Tickets[0].GroupID,
				Message:   msg.Message,
			})

			resultTicket := dbTicket.Get(ticket.ID, &core.Ticket{Model: gorm.Model{ID: ticketResult.Tickets[0].ID}})
			if resultTicket.Err != nil {
				log.Println(resultTicket.Err)
			}
			if len(resultTicket.Tickets) != 0 {
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
						mail.SendMail(mail.Mail{
							ToMail:  userTmp.Email,
							Subject: "Supportより新着メッセージ",
							Content: " " + userTmp.Name + "様\n\n" + "チャットより新着メッセージがあります\n" +
								"Webシステムよりご覧いただけます。\n",
						})
					}
				}
			}

			support.Broadcast <- msg
		}
	}
}

func HandleMessagesAdmin() {
	for {
		msg := <-support.Broadcast
		//登録されているクライアント宛にデータ送信する
		for client := range support.Clients {
			// ユーザのみの場合
			if client.GroupID == 0 {
				return
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
