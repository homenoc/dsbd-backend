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
	dbChat "github.com/homenoc/dsbd-backend/pkg/api/store/support/chat/v0"
	dbTicket "github.com/homenoc/dsbd-backend/pkg/api/store/support/ticket/v0"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"time"
)

const timeLayout = "2006-01-02 15:04:05 JST"

func Create(c *gin.Context) {
	var input support.FirstInput
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	err := c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	// input check
	if err = check(input); err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	resultTicket := &core.Ticket{
		Solved:        &[]bool{false}[0],
		Title:         input.Title,
		Admin:         &[]bool{false}[0],
		Request:       &[]bool{false}[0],
		RequestReject: &[]bool{false}[0],
	}
	var userValue, groupValue string

	// isn't group
	if !input.IsGroup {
		result := auth.UserAuthorization(core.Token{UserToken: userToken, AccessToken: accessToken})
		if result.Err != nil {
			c.JSON(http.StatusUnauthorized, common.Error{Error: result.Err.Error()})
			return
		}
		resultTicket.GroupID = nil
		resultTicket.UserID = &result.User.ID
		userValue = "[" + strconv.Itoa(int(result.User.ID)) + "] " + result.User.Name + " (" + result.User.NameEn + ")"
		groupValue = "個人ユーザ"
	} else {
		//is group
		// Group authentication
		result := auth.GroupAuthorization(1, core.Token{UserToken: userToken, AccessToken: accessToken})
		if result.Err != nil {
			c.JSON(http.StatusUnauthorized, common.Error{Error: result.Err.Error()})
			return
		}
		resultTicket.GroupID = result.User.GroupID
		resultTicket.UserID = &result.User.ID
		userValue = "[" + strconv.Itoa(int(result.User.ID)) + "] " + result.User.Name + "(" + result.User.NameEn + ")"
		groupValue = "[" + strconv.Itoa(int(result.User.Group.ID)) + "] " + result.User.Group.Org + "(" + result.User.Group.OrgEn + ")"
	}

	// Tickets DBに登録
	ticketResult, err := dbTicket.Create(resultTicket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	// Chat DBに登録
	_, err = dbChat.Create(&core.Chat{
		UserID:   resultTicket.UserID,
		Admin:    false,
		Data:     input.Data,
		TicketID: ticketResult.ID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	//HomeNOC Slackに送信
	noticeAdd("新規チケット作成", userValue, groupValue, input)

	c.JSON(http.StatusOK, ticket.Ticket{ID: ticketResult.ID})
}

func Request(c *gin.Context) {
	var input support.FirstInput
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	err := c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	// Group authentication
	result := auth.GroupAuthorization(1, core.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: result.Err.Error()})
		return
	}

	// input check
	if err = check(input); err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	resultTicket := &core.Ticket{
		GroupID:       result.User.GroupID,
		UserID:        &result.User.ID,
		Solved:        &[]bool{false}[0],
		Title:         input.Title,
		Admin:         &[]bool{false}[0],
		Request:       &[]bool{true}[0],
		RequestReject: &[]bool{false}[0],
	}

	// Tickets DBに登録
	ticketResult, err := dbTicket.Create(resultTicket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	// Chat DBに登録
	_, err = dbChat.Create(&core.Chat{
		UserID:   resultTicket.UserID,
		Admin:    false,
		Data:     input.Data,
		TicketID: ticketResult.ID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	userValue := "[" + strconv.Itoa(int(result.User.ID)) + "] " + result.User.Name + "(" + result.User.NameEn + ")"
	groupValue := "[" + strconv.Itoa(int(result.User.Group.ID)) + "] " + result.User.Group.Org + "(" + result.User.Group.OrgEn + ")"

	//HomeNOC Slackに送信
	noticeAdd("[新規] 追加・変更手続き", userValue, groupValue, input)

	c.JSON(http.StatusOK, ticket.Ticket{ID: ticketResult.ID})
}

func Update(c *gin.Context) {
	var input core.Ticket

	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	err := c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: fmt.Sprintf("id error")})
		return
	}

	// Tickets DBからデータを取得
	ticketResult := dbTicket.Get(ticket.ID, &core.Ticket{Model: gorm.Model{ID: uint(id)}})
	if ticketResult.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: ticketResult.Err.Error()})
		return
	}

	updateTicketData := ticketResult.Tickets[0]
	var userValue, groupValue string

	// isn't group
	if ticketResult.Tickets[0].GroupID == nil {
		result := auth.UserAuthorization(core.Token{UserToken: userToken, AccessToken: accessToken})
		if result.Err != nil {
			c.JSON(http.StatusUnauthorized, common.Error{Error: result.Err.Error()})
			return
		}
		userValue = "[" + strconv.Itoa(int(result.User.ID)) + "] " + result.User.Name + "(" + result.User.NameEn + ")"
	} else {
		//is group
		// Group authentication
		result := auth.GroupAuthorization(1, core.Token{UserToken: userToken, AccessToken: accessToken})
		if result.Err != nil {
			c.JSON(http.StatusUnauthorized, common.Error{Error: result.Err.Error()})
			return
		}
		userValue = "[" + strconv.Itoa(int(result.User.ID)) + "] " + result.User.Name + "(" + result.User.NameEn + ")"
		groupValue = "[" + strconv.Itoa(int(result.User.Group.ID)) + "] " + result.User.Group.Org + "(" + result.User.Group.OrgEn + ")"
	}

	updateTicketData.Solved = input.Solved

	// Ticketのアップデート
	err = dbTicket.Update(ticket.UpdateAll, updateTicketData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	noticeUpdate(ticketResult.Tickets[0], input, userValue, groupValue)

	c.JSON(http.StatusOK, support.Result{})
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

	result := auth.UserAuthorization(core.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		log.Println("ws:// support error:Auth error")
		conn.WriteMessage(websocket.TextMessage, []byte("error: auth error"))
		return
	}

	ticketResult := dbTicket.Get(ticket.ID, &core.Ticket{Model: gorm.Model{ID: uint(id)}})
	if ticketResult.Err != nil {
		log.Println("ws:// support error: db error")
		conn.WriteMessage(websocket.TextMessage, []byte("error: db error"))
		return
	}

	// [group ticket] check groupID
	if ticketResult.Tickets[0].GroupID != nil {
		// [group ticket] check groupID
		if *ticketResult.Tickets[0].GroupID != *result.User.GroupID {
			log.Println("groupID not match.")
			return
		}
	} else {
		// [user ticket] check userID
		if ticketResult.Tickets[0].UserID != nil && *ticketResult.Tickets[0].UserID != result.User.ID {
			log.Println("userID not match.")
			return
		}
	}

	var groupID uint = 0

	if ticketResult.Tickets[0].GroupID != nil {
		groupID = *ticketResult.Tickets[0].GroupID
	}

	// WebSocket送信
	support.Clients[&support.WebSocket{
		TicketID: uint(id),
		Admin:    false,
		UserID:   result.User.ID,
		UserName: result.User.Name,
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
				Admin:    false,
				UserID:   result.User.ID,
				UserName: result.User.Name,
				GroupID:  groupID,
				Socket:   conn,
			})
			break
		}
		// 入力されたデータをTokenにて認証
		resultAuth := auth.UserAuthorization(core.Token{
			UserToken:   msg.UserToken,
			AccessToken: msg.AccessToken,
		})
		if resultAuth.Err != nil {
			log.Println(resultAuth.Err)
			return
		}

		if result.User.ID != resultAuth.User.ID {
			log.Println("UserID is not match")
			return
		}

		if !*ticketResult.Tickets[0].Solved {
			_, err = dbChat.Create(&core.Chat{
				TicketID: ticketResult.Tickets[0].ID,
				UserID:   &result.User.ID,
				Admin:    false,
				Data:     msg.Message,
			})
			if err != nil {
				conn.WriteJSON(&support.WebSocketResult{Err: "db write error"})
			} else {
				msg.TicketID = ticketResult.Tickets[0].ID
				msg.UserID = result.User.ID
				msg.GroupID = groupID
				msg.Admin = false
				msg.UserName = result.User.Name
				// Token関連の初期化
				msg.AccessToken = ""
				msg.UserToken = ""

				//管理側に送信
				controller.SendChatUser(controllerInterface.Chat{
					TicketID:  ticketResult.Tickets[0].ID,
					CreatedAt: msg.CreatedAt,
					UserID:    result.User.ID,
					UserName:  result.User.Name,
					GroupID:   groupID,
					Admin:     msg.Admin,
					Message:   msg.Message,
				})

				//Slackに送信
				userValue := "[" + strconv.Itoa(int(result.User.ID)) + "] " + result.User.Name + "(" + result.User.NameEn + ")"
				groupValue := "[" + strconv.Itoa(int(result.User.Group.ID)) + "] " + result.User.Group.Org + "(" + result.User.Group.OrgEn + ")"
				noticeNewMessage(false, userValue, groupValue, ticketResult.Tickets[0], msg.Message)

				support.Broadcast <- msg
			}
		}
	}
}

func HandleMessages() {
	for {
		msg := <-support.Broadcast

		//登録されているクライアント宛にデータ送信する
		for client := range support.Clients {
			// ユーザのみの場合
			if client.TicketID == msg.TicketID {
				if msg.GroupID == 0 && client.GroupID == 0 && client.UserID == msg.UserID {
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
