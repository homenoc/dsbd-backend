package v0

import (
	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/notification"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/api/store/group/v0"
	dbUser "github.com/homenoc/dsbd-backend/pkg/api/store/user/v0"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/customer"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"time"
)

//参照関連のエラーが出る可能性あるかもしれない
func Add(c *gin.Context) {
	var input group.Input
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	err := c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	userResult := auth.UserAuthentication(core.Token{UserToken: userToken, AccessToken: accessToken})
	if userResult.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: userResult.Err.Error()})
		return
	}

	// check authority
	if userResult.User.Level > 2 {
		c.JSON(http.StatusUnauthorized, common.Error{Error: "You don't have authority this operation"})
		return
	}

	if userResult.User.GroupID != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: "error: You can't create new group"})
		return
	}

	if err = check(input); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	var studentExpired *time.Time = nil
	if *input.Student {
		tmpStudentExpired, _ := time.Parse("2006-01-02", *input.StudentExpired)
		studentExpired = &tmpStudentExpired
	}

	// added customer (stripe)
	stripe.Key = config.Conf.Stripe.SecretKey

	params := &stripe.CustomerParams{
		Description: stripe.String("Org: " + input.Org + "(" + input.OrgEn + ")"),
	}
	cus, err := customer.New(params)
	if err != nil {
		log.Println("Error: " + err.Error())
	}

	result, err := dbGroup.Create(&core.Group{
		Agree:            &[]bool{*input.Agree}[0],
		StripeCustomerID: &cus.ID,
		Question:         input.Question,
		Org:              input.Org,
		OrgEn:            input.OrgEn,
		PostCode:         input.PostCode,
		Address:          input.Address,
		AddressEn:        input.AddressEn,
		Tel:              input.Tel,
		Country:          input.Country,
		Status:           &[]uint{1}[0],
		ExpiredStatus:    &[]uint{0}[0],
		Contract:         input.Contract,
		Student:          input.Student,
		MemberExpired:    studentExpired,
		Open:             &[]bool{false}[0],
		Pass:             &[]bool{false}[0],
		Lock:             &[]bool{true}[0],
		AddAllow:         &[]bool{true}[0],
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	attachment := slack.Attachment{}
	attachment.AddField(slack.Field{Title: "Title", Value: "グループ登録"}).
		AddField(slack.Field{Title: "Question", Value: input.Question}).
		AddField(slack.Field{Title: "Org", Value: input.Org + "(" + input.OrgEn + ")"}).
		AddField(slack.Field{Title: "Country", Value: input.Country}).
		AddField(slack.Field{Title: "Student", Value: strconv.FormatBool(*input.Student)}).
		AddField(slack.Field{Title: "Contract", Value: input.Contract})

	notification.SendSlack(notification.Slack{Attachment: attachment, ID: "main", Status: true})

	if err = dbUser.Update(user.UpdateGID, &core.User{Model: gorm.Model{ID: userResult.User.ID}, GroupID: &result.Model.ID}); err != nil {
		log.Println(dbGroup.Delete(&core.Group{Model: gorm.Model{ID: result.ID}}))
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
	} else {
		c.JSON(http.StatusOK, common.Result{})
	}
}
