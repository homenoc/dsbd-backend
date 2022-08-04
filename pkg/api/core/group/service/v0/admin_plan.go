package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	dbService "github.com/homenoc/dsbd-backend/pkg/api/store/group/service/v0"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

func AddPlanByAdmin(c *gin.Context) {
	var input core.Plan

	err := c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	resultAdmin := auth.AdminAuthorization(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	input.IPID = uint(id)

	if err = dbService.JoinPlan(uint(id), input); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	noticeAddPlanByAdmin(id, input)
	c.JSON(http.StatusOK, common.Result{})
}

func DeletePlanByAdmin(c *gin.Context) {
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

	if err = dbService.DeletePlan(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	noticeDelete("Plan情報", uint(id))
	c.JSON(http.StatusOK, common.Result{})
}

func UpdatePlanByAdmin(c *gin.Context) {
	var input core.Plan

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

	before, err := dbService.GetPlan(&core.Plan{Model: gorm.Model{ID: uint(id)}})
	if err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: err.Error()})
		return
	}

	input.ID = uint(id)

	if err = dbService.UpdatePlan(input); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	noticeUpdatePlanByAdmin(before, input)
	c.JSON(http.StatusOK, common.Result{})
}
