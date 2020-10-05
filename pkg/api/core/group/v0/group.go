package v0

import (
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/store/group/v0"
	dbUser "github.com/homenoc/dsbd-backend/pkg/store/user/v0"
	"github.com/jinzhu/gorm"
	"net/http"
)

//参照関連のエラーが出る可能性あるかもしれない
func Add(c *gin.Context) {
	var input group.Group
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	c.BindJSON(&input)

	userResult := auth.UserAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	if userResult.Err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: userResult.Err.Error()})
		return
	}

	// check authority
	if userResult.User.Level > 1 {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: "You don't have authority this operation"})
		return
	}

	if userResult.User.GID != 0 {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: "error: You can't create new group", GroupData: nil})
		return
	}

	if err := check(input); err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: err.Error()})
		return
	}

	result, err := dbGroup.Create(&group.Group{
		Agree: true, Question: input.Question, Org: input.Org, Status: 0, Bandwidth: input.Bandwidth, Name: input.Name,
		PostCode: input.PostCode, Address: input.Address, Mail: input.Mail, Phone: input.Phone, Country: input.Country,
		Comment: input.Comment, Monitor: input.Monitor, Contract: input.Contract,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: err.Error(), GroupData: nil})
		return
	}
	if err := dbUser.Update(user.UpdateStatus, &user.User{Model: gorm.Model{ID: userResult.User.ID}, Status: 10}); err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: err.Error(), GroupData: nil})
		return
	}
	if err := dbUser.Update(user.UpdateGID, &user.User{Model: gorm.Model{ID: userResult.User.ID}, GID: result.Model.ID}); err != nil {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: err.Error()})
	} else {
		c.JSON(http.StatusOK, group.Result{Status: true})
	}
}

func Update(c *gin.Context) {

}

func groupOrgChange() {

}
