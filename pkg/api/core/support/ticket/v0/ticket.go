package v0

import (
	"fmt"
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	controllerInterface "github.com/homenoc/dsbd-backend/pkg/api/core/controller"
	controller "github.com/homenoc/dsbd-backend/pkg/api/core/controller/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/support"
	"github.com/homenoc/dsbd-backend/pkg/api/core/support/chat"
	"github.com/homenoc/dsbd-backend/pkg/api/core/support/ticket"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	dbChat "github.com/homenoc/dsbd-backend/pkg/api/store/support/chat/v0"
	dbTicket "github.com/homenoc/dsbd-backend/pkg/api/store/support/ticket/v0"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
)

func Create(c *gin.Context) {
	var input support.FirstInput
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	c.BindJSON(&input)

	// Group authentication
	result := auth.GroupAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusUnauthorized, support.Result{Status: false, Error: result.Err.Error()})
		return
	}

	// input check
	if err := check(input); err != nil {
		c.JSON(http.StatusBadRequest, support.Result{Status: false, Error: err.Error()})
		return
	}

	// Chat DBに登録
	chatResult, err := dbChat.Create(&chat.Chat{UserID: result.User.ID, Admin: false, Data: input.Data})
	if err != nil {
		c.JSON(http.StatusInternalServerError, support.Result{Status: false, Error: err.Error()})
		return
	}

	// Ticket DBに登録
	ticketResult, err := dbTicket.Create(&ticket.Ticket{GroupID: result.Group.ID, UserID: result.User.ID,
		ChatIDStart: chatResult.ID, ChatIDEnd: chatResult.ID, Solved: &[]bool{false}[0], Title: input.Title})
	if err != nil {
		c.JSON(http.StatusInternalServerError, support.Result{Status: false, Error: err.Error()})
		return
	}

	// Chat DBにTicketIDを登録
	err = dbChat.Update(chat.UpdateAll, chat.Chat{Admin: false, Data: chatResult.Data, TicketID: ticketResult.ID})
	if err != nil {
		c.JSON(http.StatusInternalServerError, support.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, support.Result{Status: true, Ticket: []ticket.Ticket{*ticketResult},
		Chat: []chat.Chat{*chatResult}})
}

func Get(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, support.Result{Status: false, Error: fmt.Sprintf("id error")})
		return
	}

	// Group authentication
	result := auth.GroupAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusUnauthorized, support.Result{Status: false, Error: result.Err.Error()})
		return
	}

	// IDからDBからチケットを検索
	resultTicket := dbTicket.Get(ticket.ID, &ticket.Ticket{Model: gorm.Model{ID: uint(id)}})
	if resultTicket.Err != nil {
		c.JSON(http.StatusInternalServerError, support.Result{Status: false, Error: resultTicket.Err.Error()})
		return
	}

	// GroupIDが一致しない場合はここでエラーを返す
	if resultTicket.Ticket[0].GroupID != result.Group.ID {
		c.JSON(http.StatusInternalServerError, support.Result{Status: false, Error: "Auth Error: group id failed..."})
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

func GetTitle(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	result := auth.GroupAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusUnauthorized, support.Result{Status: false, Error: result.Err.Error()})
		return
	}

	// Ticket DBからGroup IDのTicketデータを抽出
	resultTicket := dbTicket.Get(ticket.GID, &ticket.Ticket{GroupID: result.Group.ID})
	if resultTicket.Err != nil {
		c.JSON(http.StatusInternalServerError, support.Result{Status: false, Error: resultTicket.Err.Error()})
		return
	}

	log.Println(resultTicket)

	c.JSON(http.StatusOK, support.Result{Status: true, Ticket: resultTicket.Ticket})
}

func GetWebSocket(c *gin.Context) {
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
	support.Clients[&support.WebSocket{TicketID: uint(id), Admin: false,
		UserID: result.User.ID, GroupID: result.Group.ID, Socket: conn}] = true

	//WebSocket受信
	for {
		var msg support.WebSocketResult
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(support.Clients, &support.WebSocket{TicketID: uint(id), Admin: false, UserID: result.User.ID,
				GroupID: result.Group.ID, Socket: conn})
			break
		}
		// 入力されたデータをTokenにて認証
		resultGroup := auth.GroupAuthentication(token.Token{UserToken: msg.UserToken, AccessToken: msg.AccessToken})
		if resultGroup.Err != nil {
			log.Println(resultGroup.Err)
			return
		}

		_, err = dbChat.Create(&chat.Chat{TicketID: ticketResult.Ticket[0].ID, UserID: result.User.ID, Admin: false,
			Data: msg.Message})
		if err != nil {
			conn.WriteJSON(&support.WebSocketResult{Err: "db write error"})
		} else {

			msg.UserID = result.User.ID
			msg.GroupID = resultGroup.Group.ID
			msg.Admin = false
			// Token関連の初期化
			msg.AccessToken = ""
			msg.UserToken = ""

			//ユーザ側に送信
			controller.SendChatUser(controllerInterface.Chat{CreatedAt: msg.CreatedAt,
				UserID: result.User.ID, GroupID: resultGroup.Group.ID, Admin: msg.Admin, Message: msg.Message})

			//HomeNOC Slackに送信
			attachment := slack.Attachment{}
			attachment.AddField(slack.Field{Title: "Title", Value: "Supportメッセージ"}).
				AddField(slack.Field{Title: "UserID", Value: strconv.Itoa(int(result.User.ID))}).
				AddField(slack.Field{Title: "GroupID", Value: strconv.Itoa(int(resultGroup.Group.ID)) + "-" + resultGroup.Group.Org}).
				AddField(slack.Field{Title: "Message", Value: msg.Message})
			notification.SendSlack(notification.Slack{Attachment: attachment, Channel: "user", Status: true})

			support.Broadcast <- msg
		}
	}
}

func HandleMessages() {
	for {
		msg := <-support.Broadcast

		//登録されているクライアント宛にデータ送信する
		for client := range support.Clients {
			// ユーザのみの場合
			if client.GroupID == 0 {
				return
			} else if client.GroupID == msg.GroupID {
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
