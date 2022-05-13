package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service/jpnicTech"
	dbJPNICTech "github.com/homenoc/dsbd-backend/pkg/api/store/group/service/jpnicTech/v0"
	dbService "github.com/homenoc/dsbd-backend/pkg/api/store/group/service/v0"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

func AddJPNICTechByAdmin(c *gin.Context) {
	var input core.JPNICTech

	resultAdmin := auth.AdminAuthorization(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

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

	input.ID = 0
	input.ServiceID = uint(id)

	if err = dbService.JoinJPNICTech(input); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	noticeSlackAddJPNICTech(id, input)
	c.JSON(http.StatusOK, common.Result{})
}

func DeleteJPNICTechByAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthorization(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	if err = dbService.DeleteJPNICTech(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	noticeSlackDelete("JPNIC技術連絡担当者", uint(id))
	c.JSON(http.StatusOK, common.Result{})
}

func UpdateJPNICTechByAdmin(c *gin.Context) {
	var input core.JPNICTech

	resultAdmin := auth.AdminAuthorization(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

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

	before := dbJPNICTech.Get(jpnicTech.ID, &core.JPNICTech{Model: gorm.Model{ID: uint(id)}})
	if before.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: before.Err.Error()})
		return
	}

	input.ID = uint(id)

	if err = dbJPNICTech.Update(jpnicTech.UpdateAll, input); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	noticeSlackUpdateJPNICTech(before.Tech[0], input)
	c.JSON(http.StatusOK, common.Result{})
}
