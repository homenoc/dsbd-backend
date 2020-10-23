package v0

import (
	"github.com/gin-gonic/gin"
	authInterface "github.com/homenoc/dsbd-backend/pkg/api/core/auth"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	toolToken "github.com/homenoc/dsbd-backend/pkg/api/core/tool/token"
	dbToken "github.com/homenoc/dsbd-backend/pkg/api/store/token/v0"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"time"
)

func GenerateAdmin(c *gin.Context) {
	resultAuth := auth.AdminRadiusAuthentication(authInterface.AdminStruct{
		User: c.Request.Header.Get("USER"), Pass: c.Request.Header.Get("PASS")})
	if resultAuth.Err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: resultAuth.Err.Error()})
		return
	}
	accessToken, _ := toolToken.Generate(2)

	if err := dbToken.Create(&token.Token{AdminID: resultAuth.AdminID, UID: 0, ExpiredAt: time.Now().Add(60 * time.Minute),
		Admin: true, AccessToken: accessToken, Debug: "User: " + c.Request.Header.Get("USER")}); err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, token.Result{Status: true, Token: []token.Token{{AccessToken: accessToken}}})
}

func AddAdmin(c *gin.Context) {
	var input token.Token

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	c.BindJSON(&input)

	accessToken, _ := toolToken.Generate(2)

	if err := dbToken.Create(&token.Token{
		Admin: true, AccessToken: accessToken, Debug: "User: " + c.Request.Header.Get("USER")}); err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, token.ResultTmpToken{Status: true, Token: accessToken})
}

func DeleteAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: err.Error()})
		return
	}

	if err := dbToken.Delete(&token.Token{Model: gorm.Model{ID: uint(id)}}); err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, token.Result{Status: true})
}

func DeleteAllAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	if err := dbToken.DeleteAll(); err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, token.Result{Status: true})
}

func UpdateAdmin(c *gin.Context) {
	var input token.Token

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	c.BindJSON(&input)

	if err := dbToken.Update(token.UpdateAll, &input); err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, token.Result{Status: true})
}

func GetAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: err.Error()})
		return
	}

	result := dbToken.Get(token.ID, &token.Token{Model: gorm.Model{ID: uint(id)}})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, token.Result{Status: true, Token: result.Token})
}

func GetAllAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: resultAdmin.Err.Error()})
		return
	}

	if result := dbToken.GetAll(); result.Err != nil {
		c.JSON(http.StatusInternalServerError, token.Result{Status: false, Error: result.Err.Error()})
	} else {
		c.JSON(http.StatusOK, token.Result{Status: true, Token: result.Token})
	}
}
