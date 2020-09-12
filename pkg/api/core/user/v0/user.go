package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbUser "github.com/homenoc/dsbd-backend/pkg/store/user/v0"
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
