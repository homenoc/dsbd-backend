package v0

import (
	"git.bgp.ne.jp/dsbd/backend/pkg/tool"
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbUser "github.com/homenoc/dsbd-backend/pkg/store/user/v0"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
)

//func SendEmailVerify(user auth.User) auth.ResultAuth {
//
//}

func Add(c *gin.Context) {
	var input, data user.User
	userToken := c.Param("user_token")
	accessToken := c.Param("access_token")

	c.BindJSON(&input)

	if input.GID == 0 { //new user
		data = user.User{GID: 0, Name: input.Name, Email: input.Email, Pass: input.Pass, Status: 0, Level: 0}
	} else { //new users for group
		authResult, err := authentication(token.Token{UserToken: userToken, AccessToken: accessToken})
		if err != nil {
			c.JSON(http.StatusInternalServerError, authResult)
			return
		}
		if authResult.GID != input.GID && authResult.GID > 0 {
			c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: "gid mismatch"})
			return
		}
		data = user.User{GID: input.GID, Name: input.Name, Email: input.Email, Pass: input.Pass, Status: 0, Level: input.Level}
	}
	//check exist for database
	if err := dbUser.Create(&data); err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false})
	} else {
		c.JSON(http.StatusOK, user.Result{Status: true, Error: err.Error()})
	}
}

func MailVerify(c *gin.Context) {
	token := c.Param("token")

	result, err := dbUser.Get(user.MailToken, &user.User{MailToken: token})
	if err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: err.Error() + "| can not find token data"})
	}

	if result.Status > 100 || 0 > result.Status {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: "error: user status"})
	}

	if err := dbUser.Update(user.UpdateVerifyMail, &user.User{MailVerify: true}); err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: err.Error()})
	} else {
		c.JSON(http.StatusOK, &user.Result{Status: true})
	}
}

func Update(c *gin.Context) {
	var input user.User
	userToken := c.Param("user_token")
	accessToken := c.Param("access_token")
	targetString := c.Param("target")

	c.BindJSON(&input)

	authResult, err := authentication(token.Token{UserToken: userToken, AccessToken: accessToken})

	if err != nil {
		c.JSON(http.StatusInternalServerError, authResult)
		return
	}

	target, _ := strconv.Atoi(targetString)
	if target == user.UpdateName {
		if err := dbUser.Update(user.UpdatePass, &user.User{Model: gorm.Model{ID: authResult.ID}, Name: input.Name}); err != nil {
			c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: err.Error()})
		} else {
			c.JSON(http.StatusOK, user.Result{Status: true})
		}
	} else if target == user.UpdatePass {
		if err := dbUser.Update(user.UpdatePass, &user.User{Model: gorm.Model{ID: authResult.ID}, Pass: input.Pass}); err != nil {
			c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: err.Error()})
		} else {
			c.JSON(http.StatusOK, user.Result{Status: true})
		}
	} else if target == user.UpdateMail {
		token, _ := tool.GenerateToken(4)

		if err := dbUser.Update(user.UpdatePass, &user.User{Model: gorm.Model{ID: authResult.ID},
			Email: input.Email, MailToken: token}); err != nil {
			c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: err.Error()})
		} else {
			c.JSON(http.StatusOK, user.Result{Status: true})
		}
	} else {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: "error target"})
	}
}
