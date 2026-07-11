package v0

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	"github.com/homenoc/dsbd-backend/pkg/api/middleware"
)

// GetMe returns the authenticated user's profile (GET /user/me).
func GetMe(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"user": user.NewUser(middleware.CurrentUser(c))})
}
