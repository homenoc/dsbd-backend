package jpnic

import (
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	jpnicTransaction "github.com/homenoc/jpnic"
	"log"
	"net/http"
	"strconv"
)

func ManualRegistration(c *gin.Context) {
	var input jpnicTransaction.WebTransaction

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	conf := jpnicTransaction.Config{
		URL:        config.Conf.JPNIC.URL,
		CAFilePath: config.Conf.JPNIC.CAFilePath,
	}

	err := c.BindJSON(&input)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, common.Error{Error: err.Error()})
		return
	}
	log.Println(input)

	input.Network.KindID = "10"

	if input.Network.KindID == strconv.Itoa(jpnicTransaction.IPv4Register) ||
		input.Network.KindID == strconv.Itoa(jpnicTransaction.IPv4Edit) {
		// IPv4の場合
		log.Println(config.Conf.JPNIC.V4KeyFilePath)
		conf.KeyFilePath = config.Conf.JPNIC.V4KeyFilePath
		conf.CertFilePath = config.Conf.JPNIC.V4CertFilePath
	} else if input.Network.KindID == strconv.Itoa(jpnicTransaction.IPv6Register) ||
		input.Network.KindID == strconv.Itoa(jpnicTransaction.IPv6Edit) {
		// IPv6の場合
		conf.KeyFilePath = config.Conf.JPNIC.V6KeyFilePath
		conf.CertFilePath = config.Conf.JPNIC.V6CertFilePath
	} else {
		c.JSON(http.StatusBadRequest, common.Error{Error: "Kind ID is invalid."})
	}

	result := conf.Send(input)
	if result.Err != nil {
		bad(result)
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
	} else {
		success(result)
		c.JSON(http.StatusOK, common.Result{})
	}
}
