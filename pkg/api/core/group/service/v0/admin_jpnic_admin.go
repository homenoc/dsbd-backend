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

func AddJPNICAdminByAdmin(c *gin.Context) {
	var input core.JPNICAdmin

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

	if err = dbService.JoinJPNICByAdmin(uint(id), input); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	noticeAddJPNICByAdmin(id, input)
	c.JSON(http.StatusOK, common.Result{})
}

func DeleteJPNICAdminByAdmin(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	if err = dbService.DeleteJPNICByAdmin(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	noticeDelete("JPNIC管理者連絡窓口", uint(id))
	c.JSON(http.StatusOK, common.Result{})
}

func UpdateJPNICAdminByAdmin(c *gin.Context) {
	var input core.JPNICAdmin

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

	before, err := dbService.GetJPNICAdmin(uint(id))
	if err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: err.Error()})
		return
	}

	input.ID = uint(id)

	if err = dbService.UpdateJPNICByAdmin(input); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	noticeUpdateJPNICByAdmin(before, input)
	c.JSON(http.StatusOK, common.Result{})
}
