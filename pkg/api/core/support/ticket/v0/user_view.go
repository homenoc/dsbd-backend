package v0

import (
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/support/ticket"
	"github.com/homenoc/dsbd-backend/pkg/api/middleware"
	dbUser "github.com/homenoc/dsbd-backend/pkg/api/store/user/v0"
)

// GetUserTickets returns the caller's support tickets (GET /ticket).
func GetUserTickets(c *gin.Context) {
	userDetail, ok := loadUserDetail(c)
	if !ok {
		return
	}
	tickets, _ := buildUserTickets(userDetail)
	c.JSON(http.StatusOK, gin.H{"ticket": tickets})
}

// GetUserRequests returns the caller's request tickets (GET /request).
func GetUserRequests(c *gin.Context) {
	userDetail, ok := loadUserDetail(c)
	if !ok {
		return
	}
	_, requests := buildUserTickets(userDetail)
	c.JSON(http.StatusOK, gin.H{"request": requests})
}

// loadUserDetail loads the caller with the ticket association graph.
func loadUserDetail(c *gin.Context) (core.User, bool) {
	u := middleware.CurrentUser(c)
	userResult := dbUser.GetDetail(u.ID)
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

// buildUserTickets projects the user's support tickets, splitting request
// tickets out. Group tickets and the user's personal (group-less) tickets are
// merged, sorted by ID.
func buildUserTickets(userDetail core.User) (tickets []ticket.UserView, requests []ticket.RequestView) {
	if userDetail.GroupID != nil {
		for _, t := range userDetail.Group.Tickets {
			chat := chatViews(t.Chat)
			var groupID uint = 0
			if t.GroupID != nil {
				groupID = *t.GroupID
			}
			var userID uint = 0
			if t.UserID != nil {
				userID = *t.UserID
			}

			if !*t.Request {
				tickets = append(tickets, ticket.UserView{
					ID:        t.ID,
					CreatedAt: t.CreatedAt,
					GroupID:   groupID,
					UserID:    userID,
					Chat:      chat,
					Solved:    t.Solved,
					Admin:     t.Admin,
					Title:     t.Title,
				})
			} else {
				requests = append(requests, ticket.RequestView{
					ID:        t.ID,
					CreatedAt: t.CreatedAt,
					GroupID:   groupID,
					UserID:    userID,
					Chat:      chat,
					Reject:    t.RequestReject,
					Solved:    t.Solved,
					Admin:     t.Admin,
					Title:     t.Title,
				})
			}
		}
	}

	for _, t := range userDetail.Ticket {
		if t.GroupID == nil {
			tickets = append(tickets, ticket.UserView{
				ID:        t.ID,
				CreatedAt: t.CreatedAt,
				GroupID:   0,
				UserID:    userDetail.ID,
				Chat:      chatViews(t.Chat),
				Solved:    t.Solved,
				Title:     t.Title,
			})
		}
	}

	sort.Slice(tickets, func(i, j int) bool {
		return tickets[i].ID < tickets[j].ID
	})
	return tickets, requests
}

func chatViews(chats []core.Chat) []ticket.ChatView {
	var out []ticket.ChatView
	for _, ch := range chats {
		var userID uint = 0
		if ch.UserID != nil {
			userID = *ch.UserID
		}
		out = append(out, ticket.ChatView{
			CreatedAt: ch.CreatedAt,
			TicketID:  ch.TicketID,
			UserID:    userID,
			Admin:     ch.Admin,
			Data:      ch.Data,
		})
	}
	return out
}
