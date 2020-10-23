package v0

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/support"
	"github.com/homenoc/dsbd-backend/pkg/api/core/support/chat"
	"github.com/homenoc/dsbd-backend/pkg/api/core/support/ticket"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	dbChat "github.com/homenoc/dsbd-backend/pkg/api/store/support/chat/v0"
	dbTicket "github.com/homenoc/dsbd-backend/pkg/api/store/support/ticket/v0"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
)

func CreateAdmin(c *gin.Context) {
	var input support.FirstInput

	// Admin authentication
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	c.BindJSON(&input)

	// input check
	if err := checkAdmin(input); err != nil {
		c.JSON(http.StatusInternalServerError, support.Result{Status: false, Error: err.Error()})
		return
	}

	// Chat DBに登録
	chatResult, err := dbChat.Create(&chat.Chat{UserID: 0, Admin: true, Data: input.Data})
	if err != nil {
		c.JSON(http.StatusInternalServerError, support.Result{Status: false, Error: err.Error()})
		return
	}

	// Ticket DBに登録
	ticketResult, err := dbTicket.Create(&ticket.Ticket{GroupID: input.GroupID, UserID: 0,
		ChatIDStart: chatResult.ID, ChatIDEnd: chatResult.ID, Solved: &[]bool{false}[0], Title: input.Title})
	if err != nil {
		c.JSON(http.StatusInternalServerError, support.Result{Status: false, Error: err.Error()})
		return
	}

	// Chat DBにTicketIDを登録
	err = dbChat.Update(chat.UpdateAll, chat.Chat{Admin: true, Data: chatResult.Data, TicketID: ticketResult.ID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, support.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, support.Result{Status: true, Ticket: []ticket.Ticket{*ticketResult},
		Chat: []chat.Chat{*chatResult}})
}

func UpdateAdmin(c *gin.Context) {
	var input ticket.Ticket
	// Admin authentication
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, support.Result{Status: false, Error: fmt.Sprintf("id error")})
		return
	}

	c.BindJSON(&input)

	// Ticket DBからデータを取得
	ticketResult := dbTicket.Get(ticket.ID, &ticket.Ticket{Model: gorm.Model{ID: uint(id)}})
	if ticketResult.Err != nil {
		c.JSON(http.StatusInternalServerError, support.Result{Status: false, Error: err.Error()})
		return
	}

	// input check
	replace, err := updateAdminTicket(input, ticketResult.Ticket[0])
	if err != nil {
		c.JSON(http.StatusInternalServerError, support.Result{Status: false, Error: err.Error()})
		return
	}

	// Ticketのアップデート
	err = dbTicket.Update(ticket.UpdateAll, replace)
	if err != nil {
		c.JSON(http.StatusInternalServerError, support.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, support.Result{Status: true})
}

func GetAdmin(c *gin.Context) {
	// Admin authentication
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, support.Result{Status: false, Error: fmt.Sprintf("id error")})
		return
	}

	// IDからDBからチケットを検索
	resultTicket := dbTicket.Get(ticket.ID, &ticket.Ticket{Model: gorm.Model{ID: uint(id)}})
	if resultTicket.Err != nil {
		c.JSON(http.StatusInternalServerError, support.Result{Status: false, Error: resultTicket.Err.Error()})
		return
	}

	// Ticket DBからTicket IDのTicketデータを抽出
	// このとき、データはIDの昇順で出力
	resultChat := dbChat.Get(chat.TicketID, &chat.Chat{TicketID: resultTicket.Ticket[0].ID})
	if resultChat.Err != nil {
		c.JSON(http.StatusInternalServerError, support.Result{Status: false, Error: resultTicket.Err.Error()})
		return
	}
	c.JSON(http.StatusOK, support.Result{Status: true, Ticket: resultTicket.Ticket, Chat: resultChat.Chat})
}

func GetAllAdmin(c *gin.Context) {
	// Admin authentication
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	// Ticket DBからGroup IDのTicketデータを抽出
	resultTicket := dbTicket.GetAll()
	if resultTicket.Err != nil {
		c.JSON(http.StatusInternalServerError, support.Result{Status: false, Error: resultTicket.Err.Error()})
		return
	}

	log.Println(resultTicket)

	c.JSON(http.StatusOK, support.Result{Status: true, Ticket: resultTicket.Ticket})
}

func GetAdminWebSocket(c *gin.Context) {
	//
	// /support?id=0?user_token=accessID?access_token=token
	// id = ticketID, user_token = UserToken, access_token = AccessToken

	userToken := c.Query("user_token")
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

	result := auth.GroupAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		log.Println("ws:// support error:Auth error")
		conn.WriteMessage(websocket.TextMessage, []byte("error: auth error"))
		return
	}

	ticketResult := dbTicket.Get(ticket.ID, &ticket.Ticket{Model: gorm.Model{ID: uint(id)}})
	if ticketResult.Err != nil {
		log.Println("ws:// support error: db error")
		conn.WriteMessage(websocket.TextMessage, []byte("error: db error"))
		return
	}

	if ticketResult.Ticket[0].ID != uint(id) {
		log.Println("ticketID not match.")
	}

	// WebSocket送信
	support.Clients[&support.WebSocket{TicketID: uint(id), UserID: result.User.ID, GroupID: result.Group.ID, Socket: conn}] = true

	//WebSocket受信
	for {
		var msg support.WebSocketResult
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(support.Clients, &support.WebSocket{TicketID: uint(id), UserID: result.User.ID,
				GroupID: result.Group.ID, Socket: conn})
			break
		}

		_, err = dbChat.Create(&chat.Chat{TicketID: ticketResult.Ticket[0].ID, UserID: result.User.ID, Data: msg.Message})
		if err != nil {
			conn.WriteJSON(&support.WebSocketResult{Err: "db write error"})
		} else {
			msg.UserID = result.User.ID
			support.Broadcast <- msg
		}
	}
}

func HandleMessagesAdmin() {
	for {
		msg := <-support.Broadcast
		// 入力されたデータをTokenにて認証
		resultGroup := auth.GroupAuthentication(token.Token{UserToken: msg.UserToken, AccessToken: msg.AccessToken})
		if resultGroup.Err != nil {
			log.Println(resultGroup.Err)
			return
		}
		// Token関連の初期化
		msg.AccessToken = ""
		msg.UserToken = ""
		//登録されているクライアント宛にデータ送信する
		for client := range support.Clients {
			// ユーザのみの場合
			if client.GroupID == 0 {
				return
			} else if client.GroupID == resultGroup.Group.ID {
				err := client.Socket.WriteJSON(msg)
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
