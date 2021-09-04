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
	if input.Network.KindID == strconv.Itoa(jpnicTransaction.IPv4Register) ||
		input.Network.KindID == strconv.Itoa(jpnicTransaction.IPv4Edit) {
		// IPv4の場合
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
		registrationBad(result)
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
	} else {
		registrationSuccess(result)
		c.JSON(http.StatusOK, common.Result{})
	}
}

func Get(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	url := c.Param("url")

	log.Println(url)

	if url == "" {
		c.JSON(http.StatusBadRequest, common.Error{Error: "invalid url"})
		return
	}

	conf := jpnicTransaction.Config{
		URL:        config.Conf.JPNIC.URL,
		CAFilePath: config.Conf.JPNIC.CAFilePath,
	}

	if url[0:1] == "4" {
		// IPv4の場合
		conf.KeyFilePath = config.Conf.JPNIC.V4KeyFilePath
		conf.CertFilePath = config.Conf.JPNIC.V4CertFilePath

		result, err := conf.GetIPUser("/jpnic/entryinfo_v4.do?netwrk_id=" + url[1:])
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
			return
		}
		log.Println(result)
		c.JSON(http.StatusOK, result)
	} else if url[0:1] == "6" {
		// IPv6の場合
		conf.KeyFilePath = config.Conf.JPNIC.V6KeyFilePath
		conf.CertFilePath = config.Conf.JPNIC.V6CertFilePath

		result, err := conf.GetIPUser("/jpnic/G11320.do?netwrk_id=" + url[1:])
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusBadRequest, common.Error{Error: "Kind ID is invalid."})
	}
}

func GetHandle(c *gin.Context) {
	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
	if resultAdmin.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: resultAdmin.Err.Error()})
		return
	}

	handle := c.Param("handle")

	if handle == "" {
		c.JSON(http.StatusBadRequest, common.Error{Error: "invalid url"})
		return
	}

	conf := jpnicTransaction.Config{
		URL:        config.Conf.JPNIC.URL,
		CAFilePath: config.Conf.JPNIC.CAFilePath,
	}

	if handle[0:1] == "4" {
		// IPv4の場合
		conf.KeyFilePath = config.Conf.JPNIC.V4KeyFilePath
		conf.CertFilePath = config.Conf.JPNIC.V4CertFilePath

		result, err := conf.GetJPNICHandle(handle[1:])
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
			return
		}
		log.Println(result)
		c.JSON(http.StatusOK, result)
	} else if handle[0:1] == "6" {
		// IPv6の場合
		conf.KeyFilePath = config.Conf.JPNIC.V6KeyFilePath
		conf.CertFilePath = config.Conf.JPNIC.V6CertFilePath

		result, err := conf.GetJPNICHandle(handle[1:])
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusBadRequest, common.Error{Error: "Kind ID is invalid."})
	}
}

func GetAll(c *gin.Context) {
	var input GetAllInput

	resultAdmin := auth.AdminAuthentication(c.Request.Header.Get("ACCESS_TOKEN"))
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

	conf := jpnicTransaction.Config{
		URL:        config.Conf.JPNIC.URL,
		CAFilePath: config.Conf.JPNIC.CAFilePath,
	}

	log.Println(input)

	if input.Version == 4 {
		// IPv4の場合
		conf.KeyFilePath = config.Conf.JPNIC.V4KeyFilePath
		conf.CertFilePath = config.Conf.JPNIC.V4CertFilePath

		result, err := conf.GetAllIPv4(input.Org)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	} else if input.Version == 6 {
		// IPv6の場合
		conf.KeyFilePath = config.Conf.JPNIC.V6KeyFilePath
		conf.CertFilePath = config.Conf.JPNIC.V6CertFilePath

		result, err := conf.GetAllIPv6(input.Org)
		if err != nil {
			c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	} else {
		c.JSON(http.StatusBadRequest, common.Error{Error: "Kind ID is invalid."})
	}
}

func Return(c *gin.Context) {
	var input ReturnInput

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

	var result string

	if input.Version == 4 {
		// IPv4の場合
		conf.KeyFilePath = config.Conf.JPNIC.V4KeyFilePath
		conf.CertFilePath = config.Conf.JPNIC.V4CertFilePath
		result, err = conf.ReturnIPv4(input.Address[0], input.NetworkName, input.ReturnDate, input.NotifyEMail)
	} else if input.Version == 6 {
		// IPv6の場合
		conf.KeyFilePath = config.Conf.JPNIC.V6KeyFilePath
		conf.CertFilePath = config.Conf.JPNIC.V6CertFilePath
	} else {
		c.JSON(http.StatusBadRequest, common.Error{Error: "Version is invalid."})
	}

	if err != nil {
		returnBad(input, result)
		c.JSON(http.StatusInternalServerError, common.Error{Error: err.Error()})
	} else {
		returnSuccess(input, result)
		c.JSON(http.StatusOK, common.Result{})
	}
}
