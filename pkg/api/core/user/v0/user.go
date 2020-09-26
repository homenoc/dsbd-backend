package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbUser "github.com/homenoc/dsbd-backend/pkg/store/user/v0"
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
		authResult := authentication(token.Token{UserToken: userToken, AccessToken: accessToken})
		if authResult.Status == false {
			c.JSON(http.StatusInternalServerError, authResult)
			return
		}
		if authResult.UserData[0].GID != input.GID && authResult.UserData[0].GID > 0 {
			c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: "gid mismatch"})
			return
		}
		data = user.User{GID: input.GID, Name: input.Name, Email: input.Email, Pass: input.Pass, Status: 0, Level: input.Level}
	}
	//check exist for database
	if result := dbUser.Create(&data); !result.Status {
		c.JSON(http.StatusInternalServerError, result)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func MailVerify(c *gin.Context) {
	token := c.Param("token")

	result := dbUser.Get(user.MailToken, &user.User{MailToken: token})
	if !result.Status {
		result.Error = "Mail Token failed"
		c.JSON(http.StatusInternalServerError, result)
	}

	if result.UserData[0].Status > 100 || 0 > result.UserData[0].Status {
		result.Status = false
		result.Error = "I have already checked."
		c.JSON(http.StatusInternalServerError, result)
	}

	result = dbUser.Update(user.UpdateVerifyMail, &user.User{MailVerify: 1})
	c.JSON(http.StatusOK, &user.Result{Status: true})
}

func Update(c *gin.Context) {
	var input user.User
	userToken := c.Param("user_token")
	accessToken := c.Param("access_token")
	targetString := c.Param("target")

	c.BindJSON(&input)

	authResult := authentication(token.Token{UserToken: userToken, AccessToken: accessToken})

	if !authResult.Status {
		c.JSON(http.StatusInternalServerError, authResult)
		return
	}

	target, _ := strconv.Atoi(targetString)
	if target == user.UpdateName {
		if result := dbUser.Update(user.UpdatePass, &user.User{ID: authResult.UserData[0].ID, Name: input.Name}); !result.Status {
			c.JSON(http.StatusInternalServerError, result)
		} else {
			c.JSON(http.StatusOK, result)
		}
	} else if target == user.UpdatePass {
		if result := dbUser.Update(user.UpdatePass, &user.User{ID: authResult.UserData[0].ID, Pass: input.Pass}); !result.Status {
			c.JSON(http.StatusInternalServerError, result)
		} else {
			c.JSON(http.StatusOK, result)
		}
	} else if target == user.UpdateMail {
		if result := dbUser.Update(user.UpdatePass, &user.User{ID: authResult.UserData[0].ID, Email: input.Email}); !result.Status {
			c.JSON(http.StatusInternalServerError, result)
		} else {
			c.JSON(http.StatusOK, result)
		}
	} else {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: "error target"})
	}
}
