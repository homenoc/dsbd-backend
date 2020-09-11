package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbInitToken "github.com/homenoc/dsbd-backend/pkg/store/initToken/v0"
	dbUser "github.com/homenoc/dsbd-backend/pkg/store/user/v0"
	"github.com/homenoc/dsbd-backend/pkg/tool/hash"
	"github.com/homenoc/dsbd-backend/pkg/tool/token"
	"net/http"
)

//func SendEmailVerify(user auth.User) auth.ResultAuth {
//
//}

func Add(c *gin.Context) {
	var input user.User
	c.BindJSON(&input)
	//check exist for database
	//Status struct
	//1: Group Admin
	if result := dbUser.Create(&input); !result.Status {
		c.JSON(http.StatusInternalServerError, result)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func Delete(c *gin.Context) {
	var input user.User
	c.BindJSON(&input)

	if result := dbUser.Delete(&input); !result.Status {
		c.JSON(http.StatusInternalServerError, result)
	} else {
		c.JSON(http.StatusOK, result)
	}
}

func Verify(c *gin.Context) {
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

	result.UserData[0].MailToken = ""
	result.UserData[0].MailVerify = 0

	result = dbUser.Update(&result.UserData[0])
	c.JSON(http.StatusOK, result)
}

func UserPassChange() {

}

func UserNameChange() {

}

func GetInitToken(c *gin.Context) {
	ip := c.ClientIP()
	token1 := c.Param("token1")
	token2, _ := token.Generate(2)
	result := dbInitToken.Create(&user.InitToken{Token1: token1, Token2: token2, IP: ip})
	if result.Status {
		c.JSON(http.StatusOK, result)
	}
	c.JSON(http.StatusInternalServerError, result)
}

func GetToken(c *gin.Context) {
	token1 := c.Param("token1")
	token2 := c.Param("token2")
	mail := c.Param("mail")
	result := dbInitToken.Get(token1)
	if !result.Status {
		c.JSON(http.StatusInternalServerError, result)
	}
	userResult := dbUser.Get(user.Email, &user.User{Email: mail})
	h, err := hash.Generate(userResult.UserData[0].Pass + result.TokenData[0].Token2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, &user.InitTokenResult{Status: false, Error: "error: hash process"})
	}
	if h == token2 {

	} else {
		c.JSON(http.StatusInternalServerError, &user.InitTokenResult{Status: false, Error: "failed pass"})
	}

}
