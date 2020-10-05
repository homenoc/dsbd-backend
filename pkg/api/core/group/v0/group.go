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

	if userResult.User.GID != 0 {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: "error: You can't create new group", GroupData: nil})
		return
	}

	// check
	if !input.Agree {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: "error: Agreement is false", GroupData: nil})
		return
	}
	if input.Question == "" {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: "no data: question", GroupData: nil})
		return
	}
	if input.Org == "" {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: "no data: org", GroupData: nil})
		return
	}
	if input.Bandwidth == "" {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: "no data: bandwidth", GroupData: nil})
		return
	}
	if input.Name == "" {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: "no data: name", GroupData: nil})
		return
	}
	if input.PostCode == "" {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: "no data: postcode", GroupData: nil})
		return
	}
	if input.Address == "" {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: "no data: address", GroupData: nil})
		return
	}
	if input.Mail == "" {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: "no data: mail", GroupData: nil})
		return
	}
	if input.Phone == "" {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: "no data: phone", GroupData: nil})
		return
	}
	if input.Country == "" {
		c.JSON(http.StatusInternalServerError, group.Result{Status: false, Error: "no data: country", GroupData: nil})
		return
	}

	result, err := dbGroup.Create(&group.Group{
		Agree: true, Question: input.Question, Org: input.Org, Status: 0, Bandwidth: input.Bandwidth, Name: input.Name,
		PostCode: input.PostCode, Address: input.Address, Mail: input.Mail, Phone: input.Phone, Country: input.Country,
		Comment: input.Comment,
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
