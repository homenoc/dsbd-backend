// Package middleware holds gin middleware shared across the user/admin routers.
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	authv0 "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
)

const userContextKey = "authUser"

// UserAuth resolves the caller's token to a user and stashes it in the request
// context, aborting 401 on failure. It replaces the per-handler
// auth.UserAuthorization prologue; handlers read the result via CurrentUser.
func UserAuth(c *gin.Context) {
	result := authv0.UserAuthorization(tokenFromHeader(c))
	if result.Err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, common.Error{Error: result.Err.Error()})
		return
	}
	c.Set(userContextKey, result.User)
	c.Next()
}

// GroupAuth is UserAuth plus group membership/status enforcement. errorType
// mirrors auth.GroupAuthorization: 0 rejects unexamined groups, 1 permits them.
func GroupAuth(errorType uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		result := authv0.GroupAuthorization(errorType, tokenFromHeader(c))
		if result.Err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, common.Error{Error: result.Err.Error()})
			return
		}
		c.Set(userContextKey, result.User)
		c.Next()
	}
}

// CurrentUser returns the user set by UserAuth/GroupAuth. It panics if no auth
// middleware ran on the route — a wiring error, not a runtime condition.
func CurrentUser(c *gin.Context) core.User {
	return c.MustGet(userContextKey).(core.User)
}

func tokenFromHeader(c *gin.Context) core.Token {
	return core.Token{
		UserToken:   c.Request.Header.Get("USER_TOKEN"),
		AccessToken: c.Request.Header.Get("ACCESS_TOKEN"),
	}
}
