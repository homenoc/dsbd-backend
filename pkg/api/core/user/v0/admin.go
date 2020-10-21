package v0

import (
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	user "github.com/homenoc/dsbd-backend/pkg/api/core/user"
	dbUser "github.com/homenoc/dsbd-backend/pkg/store/user/v0"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func AddAdmin(c *gin.Context) {
	var input user.User

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	c.BindJSON(&input)

	if err := dbUser.Create(&input); err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, user.Result{Status: true})
}

func DeleteAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: err.Error()})
		return
	}

	if err := dbUser.Delete(&user.User{Model: gorm.Model{ID: uint(id)}}); err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, user.Result{Status: true})
}

func UpdateAdmin(c *gin.Context) {
	var input user.User

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	c.BindJSON(&input)

	if !strings.Contains(input.Email, "@") {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: "wrong email address"})
		return
	}
	tmp := dbUser.Get(user.Email, &user.User{Email: input.Email})
	if tmp.Err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: tmp.Err.Error()})
		return
	}

	if tmp.User[0].ID != input.ID && len(tmp.User) != 0 {
		log.Println("error: this email is already registered: " + input.Email)
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: "error: this email is already registered"})
		return
	}

	if err := dbUser.Update(user.UpdateAll, &input); err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, user.Result{Status: true})
}

func GetAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: err.Error()})
		return
	}

	result := dbUser.Get(user.ID, &user.User{Model: gorm.Model{ID: uint(id)}})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, user.Result{Status: true, User: result.User})
}

func GetAllAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	if result := dbUser.GetAll(); result.Err != nil {
		c.JSON(http.StatusInternalServerError, user.Result{Status: false, Error: result.Err.Error()})
	} else {
		c.JSON(http.StatusOK, user.Result{Status: true, User: result.User})
	}
}
