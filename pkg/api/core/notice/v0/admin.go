package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/notice"
	dbNotice "github.com/homenoc/dsbd-backend/pkg/api/store/notice/v0"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
	"time"
)

func AddAdmin(c *gin.Context) {
	var input notice.Input

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}
	err := c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	if err = check(input); err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	// 時間はJST基準
	jst, _ := time.LoadLocation("Asia/Tokyo")

	// 9999年12月31日 23:59:59.59
	var endTime = time.Date(9999, time.December, 31, 23, 59, 59, 59, jst)

	startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", input.StartTime, jst)
	if input.EndTime != nil {
		endTime, _ = time.ParseInLocation("2006-01-02 15:04:05", *input.EndTime, jst)
	}

	noticeSlackAddAdmin(input)

	var userArray []core.User

	for _, tmpID := range userExtraction(input.UserID, input.GroupID, input.NOCID) {
		userArray = append(userArray, core.User{Model: gorm.Model{ID: tmpID}})
	}

	if _, err = dbNotice.Create(&core.Notice{
		User:      userArray,
		Everyone:  input.Everyone,
		StartTime: startTime,
		EndTime:   endTime,
		Important: input.Important,
		Fault:     input.Fault,
		Info:      input.Info,
		Title:     input.Title,
		Data:      input.Data,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, notice.Result{})
}

func DeleteAdmin(c *gin.Context) {
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

	if err := dbNotice.Delete(&core.Notice{Model: gorm.Model{ID: uint(id)}}); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, notice.Result{})
}

func UpdateAdmin(c *gin.Context) {
	var input notice.Input

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

	if err = c.BindJSON(&input); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	// 時間はJST基準
	jst, _ := time.LoadLocation("Asia/Tokyo")

	startTime, _ := time.ParseInLocation(layoutInput, input.StartTime, jst)
	endTime, _ := time.ParseInLocation(layoutInput, *input.EndTime, jst)

	tmp := dbNotice.Get(notice.ID, &core.Notice{Model: gorm.Model{ID: uint(id)}})
	if tmp.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: tmp.Err.Error()})
		return
	}

	log.Println(startTime)
	log.Println(endTime)

	noticeSlackReplaceAdmin(tmp.Notice[0], input)

	if err = dbNotice.Update(notice.UpdateAll, core.Notice{
		Model:     gorm.Model{ID: uint(id)},
		StartTime: startTime,
		EndTime:   endTime,
		Everyone:  input.Everyone,
		Important: input.Important,
		Fault:     input.Fault,
		Info:      input.Info,
		Title:     input.Title,
		Data:      input.Data,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, notice.ResultAdmin{})
}

func GetAdmin(c *gin.Context) {
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

	result := dbNotice.Get(notice.ID, &core.Notice{Model: gorm.Model{ID: uint(id)}})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
		return
	}
	c.JSON(http.StatusOK, notice.ResultAdmin{Notice: result.Notice})
}

func GetAllAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	if result := dbNotice.GetAll(); result.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
	} else {
		c.JSON(http.StatusOK, notice.ResultAdmin{Notice: result.Notice})
	}
}
