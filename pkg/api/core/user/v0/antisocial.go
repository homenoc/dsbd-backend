package v0

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/middleware"
	dbUser "github.com/homenoc/dsbd-backend/pkg/api/store/user/v0"
	"gorm.io/gorm"
)

func AgreeAntisocialCheck(c *gin.Context) {
	user := middleware.CurrentUser(c)

	now := time.Now()
	if err := dbUser.UpdateAntisocialCheck(&core.User{
		Model:             gorm.Model{ID: user.ID},
		AntisocialCheck:   &[]bool{true}[0],
		AntisocialCheckAt: &now,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, common.Result{})
}
