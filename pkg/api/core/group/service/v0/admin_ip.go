package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service/ip"
	dbIP "github.com/homenoc/dsbd-backend/pkg/api/store/group/service/ip/v0"
	dbService "github.com/homenoc/dsbd-backend/pkg/api/store/group/service/v0"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

func AddIPByAdmin(c *gin.Context) {
	var input service.IPInput

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	err = c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	resultIP, err := ipProcess(true, false, []service.IPInput{input})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	resultIP[0].ServiceID = uint(id)

	if err = dbService.JoinIP(resultIP[0]); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	noticeSlackAddIP(id, input)
	c.JSON(http.StatusOK, service.Result{})
}

func DeleteIPByAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	if err = dbService.DeleteIP(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	noticeSlackDelete("IP情報", uint(id))

	c.JSON(http.StatusOK, common.Result{})
}

func UpdateIPByAdmin(c *gin.Context) {
	var input core.IP

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	err = c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	before := dbIP.Get(ip.ID, &core.IP{Model: gorm.Model{ID: uint(id)}})
	if before.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: before.Err.Error()})
		return
	}

	input.ID = uint(id)

	if err = dbService.UpdateIP(input); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	noticeSlackUpdateIP(before.IP[0], input)
	c.JSON(http.StatusOK, service.Result{})
}
