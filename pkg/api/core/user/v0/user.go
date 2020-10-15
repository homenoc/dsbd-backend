package v0

import (
	"fmt"
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbUser "github.com/homenoc/dsbd-backend/pkg/store/user/v0"
	toolToken "github.com/homenoc/dsbd-backend/pkg/tool/token"
	"github.com/jinzhu/gorm"
	"github.com/vmmgr/controller/etc"
	"net/http"
	"strconv"
	"strings"
)

func Add(c *gin.Context) {
	var input, data user.User
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	c.BindJSON(&input)

	if !strings.Contains(input.Email, "@") {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: fmt.Sprintf("wrong email address")})
		return
	}
	if input.Pass == "" || input.Name == "" {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: fmt.Sprintf("wrong name or pass")})
		return
	}

	if err := check(input); err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: err.Error()})
		return
	}

	mailToken, _ := toolToken.Generate(4)

	if input.GID == 0 { //new user
		data = user.User{GID: 0, Name: input.Name, Email: input.Email, Pass: input.Pass, Status: 0, Level: 1,
			MailVerify: false, MailToken: mailToken}
	} else { //new users for group
		if input.Level == 0 || input.Level > 5 {
			c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: fmt.Sprintf("wrong user level")})
			return
		}
		authResult := auth.UserAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
		if authResult.Err != nil {
			c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: authResult.Err.Error()})
			return
		}
		if authResult.User.GID != input.GID && authResult.User.GID > 0 {
			c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: "gid mismatch"})
			return
		}

		data = user.User{GID: input.GID, Name: input.Name, Email: input.Email, Pass: etc.GenerateUUID(), Status: 0,
			Tech: input.Tech, Level: input.Level, MailVerify: false, MailToken: mailToken}
	}

	//check exist for database
	if err := dbUser.Create(&data); err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: err.Error()})
	} else {
		c.JSON(http.StatusOK, user.Result{Status: true})
	}
}

func MailVerify(c *gin.Context) {
	token := c.Param("token")

	result := dbUser.Get(user.MailToken, &user.User{MailToken: token})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: result.Err.Error() + "| we can't find token data"})
		return
	}

	if result.User[0].MailVerify {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: fmt.Sprintf("This email has already been checked")})
		return
	}
	if result.User[0].Status >= 100 {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: fmt.Sprintf("error: user status")})
		return
	}

	if err := dbUser.Update(user.UpdateVerifyMail, &user.User{MailVerify: true}); err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: err.Error()})
	} else {
		c.JSON(http.StatusOK, &user.Result{Status: true})
	}
}

func Update(c *gin.Context) {
	var input user.User

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, user.Result{Status: false, Error: fmt.Sprintf("id error")})
		return
	}
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	c.BindJSON(&input)

	authResult := auth.UserAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	if authResult.Err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: authResult.Err.Error()})
		return
	}

	if !authResult.User.MailVerify {
		c.JSON(http.StatusBadRequest, user.Result{Status: false, Error: "not verify for user mail"})
		return
	}

	var u, serverData user.User

	if authResult.User.ID == uint(id) || id == 0 {
		serverData = authResult.User
		u.Model.ID = authResult.User.ID
		u.Status = authResult.User.Status
	} else {
		if authResult.User.GID == 0 {
			c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: "error: Group ID = 0"})
			return
		}
		if authResult.User.Level > 1 {
			c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: "error: failed user level"})
			return
		}
		userResult := dbUser.Get(user.ID, &user.User{Model: gorm.Model{ID: uint(id)}})
		if userResult.Err != nil {
			c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: userResult.Err.Error()})
			return
		}
		if userResult.User[0].GID != authResult.User.GID {
			c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: fmt.Sprintf("failed group authentication")})
			return
		}
		serverData = userResult.User[0]
		u.Model.ID = uint(id)
		u.Status = userResult.User[0].Status
	}

	u, err = replaceUser(serverData, input, u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: err.Error()})
		return
	}

	if err := dbUser.Update(user.UpdateInfo, &u); err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: err.Error()})
	} else {
		c.JSON(http.StatusOK, user.Result{Status: true})
	}
}

func Get(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	authResult := auth.UserAuthentication(token.Token{UserToken: userToken, AccessToken: accessToken})
	authResult.User.Pass = ""
	authResult.User.MailToken = ""
	if authResult.Err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: authResult.Err.Error()})
	} else {
		c.JSON(http.StatusOK, user.ResultOne{Status: true, User: authResult.User})
	}
}
