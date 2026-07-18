package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	dbService "github.com/homenoc/dsbd-backend/pkg/api/store/group/service/v0"
	"log"
	"net/http"
	"strconv"
)

func AddJPNICTechByAdmin(c *gin.Context) {
	var input core.JPNICTech

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
	noticeAddJPNICTechByAdmin(id, input)
	c.JSON(http.StatusOK, common.Result{})
}

func DeleteJPNICTechByAdmin(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	if err = dbService.DeleteJPNICTech(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	noticeDelete("JPNIC技術連絡担当者", uint(id))
	c.JSON(http.StatusOK, common.Result{})
}

func UpdateJPNICTechByAdmin(c *gin.Context) {
	var input core.JPNICTech

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

	before, err := dbService.GetJPNICTech(uint(id))
	if err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: err.Error()})
		return
	}

	input.ID = uint(id)

	if err = dbService.UpdateJPNICTech(input); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	noticeUpdateJPNICTechByAdmin(before, input)
	c.JSON(http.StatusOK, common.Result{})
}
