package v0

import (
	"fmt"
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/support"
	"github.com/homenoc/dsbd-backend/pkg/api/core/support/chat"
	"github.com/homenoc/dsbd-backend/pkg/api/core/support/ticket"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	dbChat "github.com/homenoc/dsbd-backend/pkg/api/store/support/chat/v0"
	dbTicket "github.com/homenoc/dsbd-backend/pkg/api/store/support/ticket/v0"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

func Add(c *gin.Context) {
	var input support.FirstInput
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	c.BindJSON(&input)

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

	// input check
	if err := check(input); err != nil {
		c.JSON(http.StatusBadRequest, support.Result{Status: false, Error: err.Error()})
		return
	}
	if uint(id) == 0 {
		c.JSON(http.StatusBadRequest, support.Result{Status: false, Error: "valid id"})
		return
	}

	// IDからDBからチケットを検索
	resultTicket := dbTicket.Get(ticket.ID, &ticket.Ticket{Model: gorm.Model{ID: uint(id)}})
	if resultTicket.Err != nil {
		c.JSON(http.StatusInternalServerError, support.Result{Status: false, Error: resultTicket.Err.Error()})
		return
	}
	// 問題解決時はここでエラーを返す
	if *resultTicket.Ticket[0].Solved {
		c.JSON(http.StatusInternalServerError, support.Result{Status: false, Error: "This problem is closed..."})
		return
	}
	// GroupIDが一致しない場合はここでエラーを返す
	if resultTicket.Ticket[0].GroupID != result.Group.ID {
		c.JSON(http.StatusInternalServerError, support.Result{Status: false, Error: "Auth Error: group id failed..."})
		return
	}

	// Chat DBに登録
	resultChat, err := dbChat.Create(&chat.Chat{TicketID: resultTicket.Ticket[0].ID, UserID: result.User.ID, Admin: false, Data: input.Data})
	if err != nil {
		c.JSON(http.StatusInternalServerError, support.Result{Status: false, Error: err.Error()})
		return
	}

	support.Broadcast <- support.WebSocketResult{
		ID:        uint(id),
		CreatedAt: resultChat.CreatedAt,
		UserID:    resultChat.UserID,
		Message:   resultChat.Data,
	}
	c.JSON(http.StatusOK, support.Result{Status: true})
}
