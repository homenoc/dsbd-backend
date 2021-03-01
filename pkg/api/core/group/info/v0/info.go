package v0

import (
	"fmt"
	"github.com/gin-gonic/gin"
	auth "github.com/homenoc/dsbd-backend/pkg/api/core/auth/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/common"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/info"
	"github.com/homenoc/dsbd-backend/pkg/api/core/group/network"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	dbNetwork "github.com/homenoc/dsbd-backend/pkg/api/store/group/network/v0"
	dbGatewayIP "github.com/homenoc/dsbd-backend/pkg/api/store/noc/gatewayIP/v0"
	"net/http"
	"strconv"
)

func Get(c *gin.Context) {
	userToken := c.Request.Header.Get("USER_TOKEN")
	accessToken := c.Request.Header.Get("ACCESS_TOKEN")

	result := auth.GroupAuthentication(0, token.Token{UserToken: userToken, AccessToken: accessToken})
	if result.Err != nil {
		c.JSON(http.StatusUnauthorized, common.Error{Error: result.Err.Error()})
		return
	}

	resultNetwork := dbNetwork.Get(network.Open, &network.Network{GroupID: result.Group.ID})
	if resultNetwork.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
		return
	}

	resultGatewayIP := dbGatewayIP.GetAll()
	if resultGatewayIP.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: result.Err.Error()})
		return
	}

	var information []info.Info
	var v4 []string
	var v6 []string
	var asn string

	for _, tmpNetwork := range resultNetwork.Network {
		v4 = []string{}
		v6 = []string{}
		asn = tmpNetwork.ASN

		if len(tmpNetwork.IP) > 0 {
			for _, tmpIP := range tmpNetwork.IP {
				if tmpIP.Version == 4 {
					v4 = append(v4, tmpIP.IP)
				} else if tmpIP.Version == 6 {
					v6 = append(v6, tmpIP.IP)
				}
			}
			if len(tmpNetwork.Connection) > 0 {
				for _, tmpConnection := range tmpNetwork.Connection {
					if *tmpConnection.Open {
						serviceID := strconv.Itoa(int(tmpConnection.GroupID)) + "-" +
							tmpNetwork.NetworkType + fmt.Sprintf("%03d", tmpNetwork.NetworkNumber) + "-" +
							tmpConnection.ConnectionType + fmt.Sprintf("%03d", tmpConnection.ConnectionNumber)

						var serviceName string
						var existsGatewayIP bool = false
						var nocIP string

						// Todo: 良くない実装
						// サービス名の検索(networkコンフィグから検索)
						for _, tmpNetworkCode := range config.Conf.Network {
							if tmpNetworkCode.ID == tmpNetwork.NetworkType {
								serviceName = tmpNetworkCode.Name
								break
							}
						}

						// 当団体側の終端アドレスの検索(gatewayIPから検索)
						for _, tmpGatewayIP := range resultGatewayIP.GatewayIP {
							if tmpGatewayIP.ID == *tmpConnection.GatewayIPID {
								nocIP = tmpGatewayIP.IP
								existsGatewayIP = true
								break
							}
						}

						if existsGatewayIP {
							information = append(information, info.Info{
								ServiceID:  serviceID,
								Service:    serviceName,
								UserID:     tmpConnection.UserID,
								NOC:        tmpConnection.NOC,
								V4:         v4,
								V6:         v6,
								ASN:        asn,
								Assign:     nocIP,
								TermIP:     tmpConnection.TermIP,
								NOCIP:      nocIP,
								LinkV4Our:  tmpConnection.LinkV4Our,
								LinkV4Your: tmpConnection.LinkV4Your,
								LinkV6Our:  tmpConnection.LinkV6Our,
								LinkV6Your: tmpConnection.LinkV6Your,
								Fee:        tmpConnection.Fee,
							})
						}
					}
				}
			}
		}
	}

	if len(information) == 0 {
		c.JSON(http.StatusInternalServerError, common.Error{Error: "not opening"})
		return
	}

	c.JSON(http.StatusOK, info.Result{Info: information})
}
