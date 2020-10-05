package v0

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbToken "github.com/homenoc/dsbd-backend/pkg/store/token/v0"
	dbUser "github.com/homenoc/dsbd-backend/pkg/store/user/v0"
	"github.com/homenoc/dsbd-backend/pkg/tool/hash"
	toolToken "github.com/homenoc/dsbd-backend/pkg/tool/token"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"time"
)

func GenerateInit(c *gin.Context) {
	ip := c.ClientIP()
	userToken := c.Request.Header.Get("USER_TOKEN")
	log.Println("userToken: " + userToken)
	tmpToken, _ := toolToken.Generate(2)
	err := dbToken.Create(&token.Token{ExpiredAt: time.Now().Add(30 * time.Minute), UID: 0, Status: 0,
		UserToken: userToken, TmpToken: tmpToken, Debug: ip})
	if err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: err.Error()})
	} else {
		c.JSON(http.StatusOK, &token.ResultTmpToken{Status: true, Token: tmpToken})
	}
}

func Generate(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	hashPass := c.Request.Header.Get("HASH_PASS")
	mail := c.Request.Header.Get("Email")
	tokenResult := dbToken.Get(token.UserToken, &token.Token{UserToken: userToken})
	if tokenResult.Err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: tokenResult.Err.Error()})
		return
	}
	userResult := dbUser.Get(user.Email, &user.User{Email: mail})
	if userResult.Err != nil {
		c.JSON(http.StatusInternalServerError, &token.Result{Status: false, Error: userResult.Err.Error()})
		return
	}

	if !userResult.User[0].MailVerify {
		c.JSON(http.StatusInternalServerError, &token.Result{Status: false, Error: fmt.Sprintf("You don't have email verification.")})
		return
	}

	if userResult.User[0].Status >= 100 {
		c.JSON(http.StatusInternalServerError, &token.Result{Status: false, Error: fmt.Sprintf("status error")})
		return
	}

	if hash.Generate(userResult.User[0].Pass+tokenResult.Token[0].TmpToken) != hashPass {
		c.JSON(http.StatusInternalServerError, &token.Result{Status: false, Error: "not match"})
		return
	}
	accessToken, _ := toolToken.Generate(2)
	err := dbToken.Update(token.AddToken, &token.Token{Model: gorm.Model{ID: tokenResult.Token[0].Model.ID},
		ExpiredAt: time.Now().Add(30 * time.Minute), UID: userResult.User[0].ID, Status: 1, AccessToken: accessToken})
	if err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: err.Error()})
	} else {
		tmp := []token.Token{{AccessToken: accessToken}}
		c.JSON(http.StatusOK, &token.Result{Status: true, Token: tmp})
	}
}
