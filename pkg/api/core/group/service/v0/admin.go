package v0

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service"
	serviceTemplate "github.com/homenoc/dsbd-backend/pkg/api/core/template/service"
	dbService "github.com/homenoc/dsbd-backend/pkg/api/store/group/service/v0"
	dbServiceTemplate "github.com/homenoc/dsbd-backend/pkg/api/store/template/service/v0"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

func AddByAdmin(c *gin.Context) {
	// ID取得
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	// serviceIDが0の時エラー処理
	if id == 0 {
		c.JSON(http.StatusBadRequest, common.Error{Error: fmt.Sprintf("This id is wrong... ")})
		return
	}

	var input service.Input

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}
	err = c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	// check json
	if err = check(input); err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	log.Println(input)

	var grpIP []core.IP = nil

	// get service template
	resultServiceTemplate := dbServiceTemplate.Get(serviceTemplate.ID, &core.ServiceTemplate{Model: gorm.Model{ID: input.ServiceTemplateID}})
	if resultServiceTemplate.Err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: resultServiceTemplate.Err.Error()})
		return
	}

	if *resultServiceTemplate.Services[0].NeedJPNIC {
		if err = checkJPNIC(input); err != nil {
			c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
			return
		}

		if err = checkJPNICAdminUser(input.JPNICAdmin); err != nil {
			c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
			return
		}

		// 0 < jpnic_tech count < 3
		if len(input.JPNICTech) == 0 || len(input.JPNICTech) > 2 {
			c.JSON(http.StatusBadRequest, common.Error{Error: "error: user tech count"})
			return
		}

		for _, tmp := range input.JPNICTech {
			if err = checkJPNICTechUser(tmp); err != nil {
				c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
				return
			}
		}

		grpIP, err = ipProcess(true, true, input.IP)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
			return
		}
	}

	if *resultServiceTemplate.Services[0].NeedComment && input.ServiceComment == "" {
		c.JSON(http.StatusBadRequest, common.Error{Error: "no data: comment"})
		return
	}

	if *resultServiceTemplate.Services[0].NeedGlobalAS {
		if input.ASN == 0 {
			c.JSON(http.StatusBadRequest, common.Error{Error: "no data: ASN"})
			return
		}

		grpIP, err = ipProcess(true, false, input.IP)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
			return
		}
	}

	resultNetwork := dbService.Get(service.SearchNewNumber, &core.Service{GroupID: uint(id)})
	if resultNetwork.Err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: resultNetwork.Err.Error()})
		return
	}
	var number uint = 1
	for _, tmp := range resultNetwork.Service {
		if tmp.ServiceNumber >= 1 {
			number = tmp.ServiceNumber + 1
		}
	}

	if number >= 999 {
		c.JSON(http.StatusInternalServerError, common.Error{Error: "error: over number"})
		return
	}

	//db create for network
	_, err = dbService.Create(&core.Service{
		GroupID:           uint(id),
		ServiceTemplateID: &input.ServiceTemplateID,
		ServiceComment:    input.ServiceComment,
		ServiceNumber:     number,
		Org:               input.Org,
		OrgEn:             input.OrgEn,
		PostCode:          input.Postcode,
		Address:           input.Address,
		AddressEn:         input.AddressEn,
		AveUpstream:       input.AveUpstream,
		MaxUpstream:       input.MaxUpstream,
		AveDownstream:     input.AveDownstream,
		MaxDownstream:     input.MaxDownstream,
		MaxBandWidthAS:    input.MaxBandWidthAS,
		ASN:               &[]uint{input.ASN}[0],
		IP:                grpIP,
		JPNICAdmin:        input.JPNICAdmin,
		JPNICTech:         input.JPNICTech,
		Enable:            &[]bool{true}[0],
		Pass:              &[]bool{false}[0],
		AddAllow:          &[]bool{false}[0],
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	noticeSlackAdd(id, resultServiceTemplate.Services[0].Type+fmt.Sprintf("%03d", number), input.ServiceComment)
	c.JSON(http.StatusOK, service.Result{})
}

func DeleteByAdmin(c *gin.Context) {
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

	if err = dbService.Delete(&core.Service{Model: gorm.Model{ID: uint(id)}}); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	noticeSlackDelete("Service情報", uint(id))
	c.JSON(http.StatusOK, service.Result{})
}

func UpdateByAdmin(c *gin.Context) {
	var input core.Service

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

	before := dbService.Get(service.ID, &core.Service{Model: gorm.Model{ID: uint(id)}})
	if before.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: before.Err.Error()})
		return
	}

	input.ID = uint(id)

	if err = dbService.Update(service.UpdateAll, input); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	noticeSlackUpdate(before.Service[0], input)
	c.JSON(http.StatusOK, service.Result{})
}

func GetByAdmin(c *gin.Context) {
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

	result := dbService.Get(service.ID, &core.Service{Model: gorm.Model{ID: uint(id)}})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
		return
	}

	c.JSON(http.StatusOK, service.Result{Service: result.Service})
}

func GetAllByAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	if result := dbService.GetAll(); result.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
	} else {
		c.JSON(http.StatusOK, service.Result{Service: result.Service})
	}
}
