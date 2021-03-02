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
	dbRouter "github.com/homenoc/dsbd-backend/pkg/api/store/noc/router/v0"
	dbNOC "github.com/homenoc/dsbd-backend/pkg/api/store/noc/v0"
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
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultNetwork.Err.Error()})
		return
	}

	//TODO:これより下の実装はマジでよくない
	//DBを３つからすべて抽出しているため、無駄な処理が多く今後改善必要がある。

	resultNOC := dbNOC.GetAll()
	if resultNOC.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultNOC.Err.Error()})
		return
	}

	resultRouter := dbRouter.GetAll()
	if resultRouter.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultRouter.Err.Error()})
		return
	}

	resultGatewayIP := dbGatewayIP.GetAll()
	if resultGatewayIP.Err != nil {
		c.JSON(http.StatusInternalServerError, common.Error{Error: resultGatewayIP.Err.Error()})
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
						var existsData bool = false
						var nocIP, noc string

						// Todo: 良くない実装
						// サービス名の検索(networkコンフィグから検索)
						for _, tmpNetworkCode := range config.Conf.Network {
							if tmpNetworkCode.ID == tmpNetwork.NetworkType {
								serviceName = tmpNetworkCode.Name
								break
							}
						}

						// CC0　構内接続の場合を除く
						if tmpConnection.ConnectionType != "CC0" {
							// 当団体側の終端アドレスの検索(gatewayIPから検索)
							for _, tmpGatewayIP := range resultGatewayIP.GatewayIP {
								if tmpGatewayIP.ID == *tmpConnection.GatewayIPID {
									nocIP = tmpGatewayIP.IP
									existsData = true
									break
								}
							}
						} else {
							nocIP = "構内接続のため必要なし"
							existsData = true
						}

						// Todo: 読みにくい上に処理的にも問題あり
						if existsData {
							existsData = false
							// NOCの検索(router=>nocの順番に検索)
							for _, tmpRouter := range resultRouter.Router {
								if tmpRouter.ID == *tmpConnection.RouterID {
									for _, tmpNOC := range resultNOC.NOC {
										if tmpNOC.ID == tmpRouter.NOC {
											noc = tmpNOC.Name
											existsData = true
											break
										}
									}
								}
							}
						}

						if existsData {
							information = append(information, info.Info{
								ServiceID:  serviceID,
								Service:    serviceName,
								UserID:     tmpConnection.UserID,
								NOC:        noc,
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
