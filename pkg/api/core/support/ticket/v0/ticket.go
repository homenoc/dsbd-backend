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
	dbChat "github.com/homenoc/dsbd-backend/pkg/store/support/chat/v0"
	dbTicket "github.com/homenoc/dsbd-backend/pkg/store/support/ticket/v0"
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
		c.JSON(http.StatusInternalServerError, support.Result{Status: false, Error: result.Err.Error()})
		return
	}

	// input check
	if err := check(input); err != nil {
		c.JSON(http.StatusInternalServerError, support.Result{Status: false, Error: err.Error()})
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
		ChatIDStart: chatResult.ID, ChatIDEnd: chatResult.ID, Solved: false, Title: input.Title})
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
		c.JSON(http.StatusInternalServerError, support.Result{Status: false, Error: result.Err.Error()})
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
		c.JSON(http.StatusInternalServerError, support.Result{Status: false, Error: result.Err.Error()})
		return
	}

	support.Broadcast <- support.Data{Message: "a" + result.Group.Org}

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

	// channel定義
	//messageRev := make(chan support.Data)
	//stopCh := make(chan struct{})
	//doneCh := make(chan struct{})

	//WebSocket受信
	//go receiveData(conn, stopCh, doneCh)
	//go sendData(conn, support.Broadcast, stopCh, doneCh)
	support.Clients[conn] = true

	//WebSocket送信
	//var tmpMyself string
	for {
		// メッセージの入力
		//_, msg, err := conn.ReadMessage()
		//if err != nil {
		//	close(doneCh)
		//	break
		//}
		//support.Broadcast <- support.Data{
		//	ID: uint(1),
		//	//CreatedAt: "0",
		//	UserID:  0,
		//	Message: string(msg),
		//}
		var msg support.Data
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(support.Clients, conn)
			break
		}
		support.Broadcast <- msg
		//log.Println("--receive--")
		log.Println(msg)

		//select {
		//// doneCh経由でクローズ信号を受けった場合、停止
		//case <-doneCh:
		//	println("stop request received.")
		//	return
		//default:
		//}
	}
	// WebSocket送信が完了すれば、stopCh経由で受信処理側にクローズを知らせる
	//close(stopCh)
}

//func sendData(conn *websocket.Conn, messageRev chan support.Data, stopCh, doneCh chan struct{}) {
//	var tmpOther support.Data
//
//	for {
//		select {
//		case <-stopCh:
//			println("stop request received.")
//			return
//		//case tmpMyself = <-messageRev:
//		//	conn.WriteMessage(1, []byte(tmpMyself))
//		case tmpOther = <-messageRev:
//			log.Println("--send--")
//			log.Println(tmpOther)
//			//if tmpOther.ID == uint(1) && tmpOther.UserID != 38 {
//				conn.WriteMessage(1, []byte(tmpOther.Message))
//			//}
//		}
//	}
//	// WebSocket送信が完了すれば、stopCh経由で受信処理側にクローズを知らせる
//	close(doneCh)
//}

func HandleMessages() {
	for {
		msg := <-support.Broadcast
		log.Println(msg)
		for client := range support.Clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(support.Clients, client)
			}
		}
	}
}
