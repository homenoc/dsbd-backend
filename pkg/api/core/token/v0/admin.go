package v0

import (
	"github.com/gin-gonic/gin"
	"github.com/homenoc/dsbd-backend/pkg/api/core"
	authInterface "github.com/homenoc/dsbd-backend/pkg/api/core/auth"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	toolToken "github.com/homenoc/dsbd-backend/pkg/api/core/tool/token"
	dbToken "github.com/homenoc/dsbd-backend/pkg/api/store/token/v0"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
	"time"
)

func GenerateByAdmin(c *gin.Context) {
	resultAuth := auth.AdminRadiusAuthorization(authInterface.AdminStruct{
		User: c.Request.Header.Get("USER"), Pass: c.Request.Header.Get("PASS")})
	if resultAuth.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAuth.Err.Error()})
		return
	}
	accessToken, _ := toolToken.Generate(2)

	if err := dbToken.Create(&core.Token{
		UserID:      nil,
		ExpiredAt:   time.Now().Add(60 * time.Minute),
		Admin:       &[]bool{true}[0],
		AccessToken: accessToken,
		Debug:       "User: " + c.Request.Header.Get("USER"),
	}); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, token.Result{Token: []core.Token{{AccessToken: accessToken}}})
}

func AddByAdmin(c *gin.Context) {
	var input core.Token

	resultAdmin := auth.AdminAuthorization(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	err := c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	accessToken, _ := toolToken.Generate(2)

	if err = dbToken.Create(&core.Token{
		Admin:       &[]bool{true}[0],
		AccessToken: accessToken,
		Debug:       "User: " + c.Request.Header.Get("USER"),
	}); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, token.ResultTmpToken{Token: accessToken})
}

func DeleteByAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthorization(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	if err = dbToken.Delete(&core.Token{Model: gorm.Model{ID: uint(id)}}); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, token.Result{})
}

func DeleteAllByAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthorization(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	if err := dbToken.DeleteAll(); err != nil {
		//エラー時はTokenがすでに消えている状態なので、問題なし
		c.JSON(http.StatusOK, common.Result{})
		return
	}
	c.JSON(http.StatusOK, token.Result{})
}

func UpdateByAdmin(c *gin.Context) {
	var input core.Token

	resultAdmin := auth.AdminAuthorization(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	err := c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	if err := dbToken.Update(token.UpdateAll, &input); err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, token.Result{})
}

func GetByAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthorization(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}

	result := dbToken.Get(token.ID, &core.Token{Model: gorm.Model{ID: uint(id)}})
	if result.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
		return
	}
	c.JSON(http.StatusOK, token.Result{Token: result.Token})
}

func GetAllByAdmin(c *gin.Context) {
	resultAdmin := auth.AdminAuthorization(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	if result := dbToken.GetAll(); result.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
	} else {
		c.JSON(http.StatusOK, token.Result{Token: result.Token})
	}
}
