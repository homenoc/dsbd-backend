package v0

import (
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/mail"
	"net/http"
)

func SendByAdmin(c *gin.Context) {
	var input mail.Mail

	resultAdmin := auth.AdminAuthorization(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	err := c.BindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	if err = check(input); err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	if err = SendMail(input); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
	} else {
		c.JSON(http.StatusOK, common.Result{})
	}
}
