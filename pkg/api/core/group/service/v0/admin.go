package v0

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/service"
	dbService "github.com/homenoc/dsbd-backend/pkg/api/store/group/service/v0"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"time"
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

	// check input.ConnectionType and getting connection template
	resultServiceTemplate, err := core.GetServiceType(input.ServiceType)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	if resultServiceTemplate.NeedJPNIC {
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

		// IPトランジット以外
		if !resultServiceTemplate.NeedGlobalAS {
			grpIP, err = ipProcess(true, true, input.IP)
			if err != nil {
				c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
				return
			}
		}
	}

	if resultServiceTemplate.NeedComment && input.ServiceComment == "" {
		c.JSON(http.StatusBadRequest, common.Error{Error: "no data: comment"})
		return
	}

	startDate, err := time.Parse("2006-01-02", input.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: "invalid data: start date"})
		return
	}
	var endDate *time.Time = nil
	if input.EndDate != nil {
		tmpEndDate, err := time.Parse("2006-01-02", *input.EndDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, common.Error{Error: "invalid data: end date"})
			return
		}
		endDate = &tmpEndDate
	}

	// pattern check: IP Transit
	var bgpComment = ""
	if resultServiceTemplate.NeedGlobalAS {
		if input.ASN == 0 {
			c.JSON(http.StatusBadRequest, common.Error{Error: "no data: ASN"})
			return
		}
		bgpComment = input.BGPComment
	}

	resultNetwork := dbService.GetByGroupID(uint(id))
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
		GroupID:        uint(id),
		ServiceType:    input.ServiceType,
		ServiceComment: input.ServiceComment,
		ServiceNumber:  number,
		Org:            input.Org,
		OrgEn:          input.OrgEn,
		PostCode:       input.Postcode,
		Address:        input.Address,
		AddressEn:      input.AddressEn,
		Abuse:          input.Abuse,
		AveUpstream:    input.AveUpstream,
		MaxUpstream:    input.MaxUpstream,
		AveDownstream:  input.AveDownstream,
		MaxDownstream:  input.MaxDownstream,
		MaxBandWidthAS: input.MaxBandWidthAS,
		StartDate:      startDate,
		EndDate:        endDate,
		ASN:            &[]uint{input.ASN}[0],
		IP:             grpIP,
		JPNICAdmin:     input.JPNICAdmin,
		JPNICTech:      input.JPNICTech,
		Enable:         &[]bool{true}[0],
		Pass:           &[]bool{false}[0],
		AddAllow:       &[]bool{false}[0],
		Comment:        input.Comment,
		BGPComment:     bgpComment,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	grp := getGroupInfo(uint(id))
	groupValue := "[" + strconv.Itoa(int(grp.ID)) + "] " + grp.Org + "(" + grp.OrgEn + ")"
	serviceCodeNew := resultServiceTemplate.Type + fmt.Sprintf("%03d", number)
	serviceCodeComment := input.ServiceComment
	noticeAdd("", groupValue, serviceCodeNew, serviceCodeComment)

	c.JSON(http.StatusOK, service.Result{})
}

func DeleteByAdmin(c *gin.Context) {

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	if err = dbService.Delete(&core.Service{Model: gorm.Model{ID: uint(id)}}); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	noticeDelete("Service情報", uint(id))
	c.JSON(http.StatusOK, service.Result{})
}

func UpdateByAdmin(c *gin.Context) {
	var input core.Service

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	err = common.BindJSONTolerant(c, &input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	before := dbService.GetByID(uint(id))
	if before.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: before.Err.Error()})
		return
	}

	input.ID = uint(id)

	if err = dbService.UpdateAll(input); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	noticeUpdate(before.Service[0], input)
	c.JSON(http.StatusOK, service.Result{})
}

func GetByAdmin(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	result := dbService.GetByID(uint(id))
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
		return
	}

	c.JSON(http.StatusOK, service.Result{Service: result.Service})
}

func GetAllByAdmin(c *gin.Context) {

	if result := dbService.GetAll(); result.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
	} else {
		c.JSON(http.StatusOK, service.Result{Service: result.Service})
	}
}
