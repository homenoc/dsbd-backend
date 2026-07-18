package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service"
	dbService "github.com/homenoc/dsbd-backend/pkg/api/store/group/service/v0"
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

	resultIP, err := ipProcess(true, false, []service.IPInput{input})
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	resultIP[0].ServiceID = uint(id)

	if err = dbService.CreateIP(resultIP[0]); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	noticeAddIPByAdmin(id, input)
	c.JSON(http.StatusOK, service.Result{})
}

func DeleteIPByAdmin(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	if err = dbService.DeleteIP(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	noticeDelete("IP情報", uint(id))

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

	before, err := dbService.GetIP(uint(id))
	if err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: err.Error()})
		return
	}

	input.ID = uint(id)

	if err = dbService.UpdateIP(input); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	noticeUpdateIPByAdmin(before, input)
	c.JSON(http.StatusOK, service.Result{})
}
