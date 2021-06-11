package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service/ip"
	dbIP "github.com/homenoc/dsbd-backend/pkg/api/store/group/service/ip/v0"
	"log"
	"net/http"
	"strconv"
)

//func DeleteByAdmin(c *gin.Context) {
//	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
//	if resultAdmin.Err != nil {
//		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
//		return
//	}
//
//	id, err := strconv.Atoi(c.Param("id"))
//	if err != nil {
//		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
//		return
//	}
//
//	if err := dbService.Delete(&core.Service{Model: gorm.Model{ID: uint(id)}}); err != nil {
//		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
//		return
//	}
//	c.JSON(http.StatusOK, service.Result{})
//}

func UpdateByAdmin(c *gin.Context) {
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

	input.ID = uint(id)

	if err = dbIP.Update(ip.UpdateAll, input); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, service.Result{})
}
