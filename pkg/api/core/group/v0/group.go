package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/api/store/group/v0"
	dbUser "github.com/homenoc/dsbd-backend/pkg/api/store/user/v0"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

// 参照関連のエラーが出る可能性あるかもしれない
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

	userResult := auth.UserAuthorization(core.Token{UserToken: userToken, AccessToken: accessToken})
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

	// double check
	resultDB := dbGroup.Get(group.Org, &core.Group{Org: input.Org})
	if resultDB.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: "error: can't get data from db."})
		return
	}
	if len(resultDB.Group) != 0 {
		c.JSON(http.StatusBadRequest, common.Error{Error: "error: Change to a different group name."})
		return
	}

	if err = check(input); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	memberType := core.MemberTypeStandard.ID
	var memberExpired *time.Time = nil
	if *input.Student {
		tmpMemberExpired, _ := time.Parse("2006-01-02", *input.StudentExpired)
		memberExpired = &tmpMemberExpired
		memberType = core.MemberTypeStudent.ID
	}

	groupData := core.Group{
		Agree:         &[]bool{*input.Agree}[0],
		Question:      input.Question,
		Org:           input.Org,
		OrgEn:         input.OrgEn,
		PostCode:      input.PostCode,
		Address:       input.Address,
		AddressEn:     input.AddressEn,
		Tel:           input.Tel,
		Country:       input.Country,
		ExpiredStatus: &[]uint{0}[0],
		Contract:      input.Contract,
		MemberType:    memberType,
		MemberExpired: memberExpired,
		Pass:          &[]bool{false}[0],
		AddAllow:      &[]bool{true}[0],
	}

	_, err = dbGroup.Create(&groupData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}

	noticeAddGroup(userResult.User, input)

	if err = dbUser.Update(user.UpdateGID, &core.User{Model: gorm.Model{ID: userResult.User.ID}, GroupID: &groupData.ID}); err != nil {
		log.Println(dbGroup.Delete(&core.Group{Model: gorm.Model{ID: groupData.ID}}))
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
	} else {
		c.JSON(http.StatusOK, common.Result{})
	}
}
