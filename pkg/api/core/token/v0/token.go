package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbToken "github.com/homenoc/dsbd-backend/pkg/store/token/v0"
	dbUser "github.com/homenoc/dsbd-backend/pkg/store/user/v0"
	"github.com/homenoc/dsbd-backend/pkg/tool/hash"
	toolToken "github.com/homenoc/dsbd-backend/pkg/tool/token"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

func GenerateInit(c *gin.Context) {
	ip := c.ClientIP()
	userToken := c.Param("token1")
	tmpToken, _ := toolToken.Generate(2)
	err := dbToken.Create(&token.Token{ExpiredAt: time.Now().Add(30 * time.Minute), UID: 0, Status: 0,
		UserToken: userToken, TmpToken: tmpToken, Debug: ip})
	if err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{})
	}

	c.JSON(http.StatusOK, &token.Result{Status: true, Token: []token.Token{{TmpToken: tmpToken}}})
}

func Generate(c *gin.Context) {
	userToken := c.Param("USER_TOKEN")
	hashPass := c.Param("HASH_PASS")
	mail := c.Param("Email")
	tokenResult, err := dbToken.Get(token.UserToken, &token.Token{UserToken: userToken})
	if err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: err.Error()})
	}
	userResult, err := dbUser.Get(user.Email, &user.User{Email: mail})
	if err != nil {
		c.JSON(http.StatusInternalServerError, &token.Result{Status: false, Error: err.Error()})
	}
	h, err := hash.Generate(userResult.Pass + tokenResult.TmpToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &token.Result{Status: false, Error: "error: hash process"})
	}
	if h != hashPass {
		c.JSON(http.StatusInternalServerError, &token.Result{Status: false, Error: "failed pass"})
	}
	accessToken, _ := toolToken.Generate(2)
	err = dbToken.Update(token.AddToken, &token.Token{Model: gorm.Model{ID: tokenResult.Model.ID},
		ExpiredAt: time.Now().Add(30 * time.Minute), UID: userResult.ID, Status: 1, AccessToken: accessToken})
	if err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: err.Error()})
	} else {
		tmp := []token.Token{{AccessToken: accessToken}}
		c.JSON(http.StatusOK, &token.Result{Status: true, Token: tmp})
	}
}
