package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/notice"
	"github.com/homenoc/dsbd-backend/pkg/api/core/payment"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	dbPayment "github.com/homenoc/dsbd-backend/pkg/api/store/payment/v0"
	"github.com/stripe/stripe-go/v73"
	"github.com/stripe/stripe-go/v73/refund"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func AddByAdmin(c *gin.Context) {
	//var input notice.Input
	//
	//resultAdmin := auth.AdminAuthorization(c.Request.Header.Get("ACCESS_TOKEN"))
	//if resultAdmin.Err != nil {
	//	c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
	//	return
	//}
	//err := c.BindJSON(&input)
	//if err != nil {
	//	log.Println(err)
	//	c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
	//	return
	//}
	//
	//if err = check(input); err != nil {
	//	c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
	//	return
	//}
	//
	//// 時間はJST基準
	//jst, _ := time.LoadLocation(config.Conf.Controller.TimeZone)
	//
	//// 9999年12月31日 23:59:59.59
	//var endTime = time.Date(9999, time.December, 31, 23, 59, 59, 59, jst)
	//
	//startTime, _ := time.ParseInLocation("2006-01-02 15:04:05", input.StartTime, jst)
	//if input.EndTime != nil {
	//	endTime, _ = time.ParseInLocation("2006-01-02 15:04:05", *input.EndTime, jst)
	//}
	//
	//var userIDArray []uint
	//
	//for _, tmpID := range userExtraction(input.UserID, input.GroupID, input.NOCID) {
	//	userIDArray = append(userIDArray, tmpID)
	//}
	//
	//resultUser := dbUser.GetArray(userIDArray)
	//if resultUser.Err != nil {
	//	log.Println(resultUser.Err.Error())
	//	c.JSON(http.StatusInternalServerError, common.Error{Error: resultUser.Err.Error()})
	//	return
	//}
	//
	//if _, err = dbNotice.Create(&core.Notice{
	//	User:      resultUser.User,
	//	Everyone:  input.Everyone,
	//	StartTime: startTime,
	//	EndTime:   endTime,
	//	Important: input.Important,
	//	Fault:     input.Fault,
	//	Info:      input.Info,
	//	Title:     input.Title,
	//	Data:      input.Data,
	//}); err != nil {
	//	c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
	//	return
	//}
	//noticeSlackAddByAdmin(input)
	//c.JSON(http.StatusOK, notice.Result{})
}

func DeleteByAdmin(c *gin.Context) {
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

	if err = dbPayment.Delete(&core.Payment{Model: gorm.Model{ID: uint(id)}}); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, notice.Result{})
}

func GetByAdmin(c *gin.Context) {
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

	result, err := dbPayment.Get(payment.ID, core.Payment{Model: gorm.Model{ID: uint(id)}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, payment.ResultByAdmin{Payment: result})
}

func GetAllByAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthorization(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	if result, err := dbPayment.GetAll(); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
	} else {
		c.JSON(http.StatusOK, payment.ResultByAdmin{Payment: result})
	}
}

func RefundByAdmin(c *gin.Context) {
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

	result, err := dbPayment.Get(payment.ID, core.Payment{Model: gorm.Model{ID: uint(id)}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	stripe.Key = config.Conf.Stripe.SecretKey

	_, err = refund.New(&stripe.RefundParams{
		PaymentIntent: stripe.String(result[0].PaymentIntentID),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, common.Result{})
}
