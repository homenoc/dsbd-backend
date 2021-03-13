package api

import (
	"github.com/gin-gonic/gin"
	controller "github.com/homenoc/dsbd-backend/pkg/api/core/controller/v0"
	connection "github.com/homenoc/dsbd-backend/pkg/api/core/group/connection/v0"
	info "github.com/homenoc/dsbd-backend/pkg/api/core/group/info/v0"
	jpnicAdmin "github.com/homenoc/dsbd-backend/pkg/api/core/group/service/jpnicAdmin/v0"
	jpnicTech "github.com/homenoc/dsbd-backend/pkg/api/core/group/service/jpnicTech/v0"
	service "github.com/homenoc/dsbd-backend/pkg/api/core/group/service/v0"
	group "github.com/homenoc/dsbd-backend/pkg/api/core/group/v0"
	bgpRouter "github.com/homenoc/dsbd-backend/pkg/api/core/noc/bgpRouter/v0"
	tunnelEndPointRouter "github.com/homenoc/dsbd-backend/pkg/api/core/noc/tunnelEndPointRouter/v0"
	tunnelEndPointRouterIP "github.com/homenoc/dsbd-backend/pkg/api/core/noc/tunnelEndPointRouterIP/v0"
	noc "github.com/homenoc/dsbd-backend/pkg/api/core/noc/v0"
	notice "github.com/homenoc/dsbd-backend/pkg/api/core/notice/v0"
	chat "github.com/homenoc/dsbd-backend/pkg/api/core/support/chat/v0"
	ticket "github.com/homenoc/dsbd-backend/pkg/api/core/support/ticket/v0"
	template "github.com/homenoc/dsbd-backend/pkg/api/core/template/v0"
	token "github.com/homenoc/dsbd-backend/pkg/api/core/token/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	user "github.com/homenoc/dsbd-backend/pkg/api/core/user/v0"
	"log"
	"net/http"
	"strconv"
)

func AdminRestAPI() {
	router := gin.Default()
	router.Use(cors)

	go token.TokenRemove()

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// Controller
			//noc
			v1.POST("/controller/chat", controller.ReceiveChatAdmin)

			// Notice
			//
			v1.POST("/notice", notice.AddAdmin)
			v1.DELETE("/notice/:id", notice.DeleteAdmin)
			v1.GET("/notice", notice.GetAllAdmin)
			v1.GET("/notice/:id", notice.GetAdmin)
			v1.PUT("/notice/:id", notice.UpdateAdmin)

			//
			// User
			//
			// User Create
			v1.POST("/user", user.AddAdmin)
			// User Delete
			v1.DELETE("/user", user.DeleteAdmin)
			// User Update
			v1.PUT("/user/:id", user.UpdateAdmin)
			v1.GET("/user", user.GetAllAdmin)
			v1.GET("/user/:id", user.GetAdmin)
			//
			// Token
			//
			v1.POST("/token/generate", token.GenerateAdmin)

			v1.POST("/token", token.AddAdmin)
			// Token Delete
			v1.DELETE("/token", token.DeleteAllAdmin)
			v1.DELETE("/token/:id", token.DeleteAdmin)
			// Token Update
			v1.PUT("/token/:id", token.UpdateAdmin)
			v1.GET("/token", token.GetAllAdmin)
			v1.GET("/token/:id", token.GetAdmin)
			//
			// Group
			//
			v1.POST("/group", group.AddAdmin)
			// Group Delete
			v1.DELETE("/group", group.DeleteAdmin)
			// Group Update
			v1.PUT("/group/:id", group.UpdateAdmin)
			v1.GET("/group", group.GetAllAdmin)
			v1.GET("/group/:id", group.GetAdmin)

			//
			// NOC
			//
			v1.POST("/noc", noc.AddAdmin)
			v1.GET("/noc", noc.GetAllAdmin)
			v1.DELETE("/noc/:id", noc.DeleteAdmin)
			v1.GET("/noc/:id", noc.GetAdmin)
			v1.PUT("/noc/:id", noc.UpdateAdmin)

			//
			// NOC Router
			//
			v1.POST("/router", bgpRouter.AddAdmin)
			v1.GET("/router", bgpRouter.GetAllAdmin)
			v1.DELETE("/router/:id", bgpRouter.DeleteAdmin)
			v1.GET("/router/:id", bgpRouter.GetAdmin)
			v1.PUT("/router/:id", bgpRouter.UpdateAdmin)

			//
			// NOC Gateway
			//
			v1.POST("/gateway", tunnelEndPointRouter.AddAdmin)
			v1.GET("/gateway", tunnelEndPointRouter.GetAllAdmin)
			v1.DELETE("/gateway/:id", tunnelEndPointRouter.DeleteAdmin)
			v1.GET("/gateway/:id", tunnelEndPointRouter.GetAdmin)
			v1.PUT("/gateway/:id", tunnelEndPointRouter.UpdateAdmin)

			//
			// NOC Gateway IP
			//
			v1.POST("/gateway_ip", tunnelEndPointRouterIP.AddAdmin)
			v1.GET("/gateway_ip", tunnelEndPointRouterIP.GetAllAdmin)
			v1.DELETE("/gateway_ip/:id", tunnelEndPointRouterIP.DeleteAdmin)
			v1.GET("/gateway_ip/:id", tunnelEndPointRouterIP.GetAdmin)
			v1.PUT("/gateway_ip/:id", tunnelEndPointRouterIP.UpdateAdmin)

			//
			// Support
			//
			v1.POST("/support", ticket.CreateAdmin)
			v1.GET("/support", ticket.GetAllAdmin)
			//v1.POST("/support/:id", chat.AddAdmin)
			v1.GET("/support/:id", ticket.GetAdmin)
			v1.PUT("/support/:id", ticket.UpdateAdmin)

			////
			//// Connection
			////
			v1.POST("/group/:id/connection", connection.AddAdmin)
			// Group Delete
			v1.DELETE("/connection/:id", connection.DeleteAdmin)
			// Group Update
			v1.PUT("/connection/:id", connection.UpdateAdmin)
			v1.GET("/connection", connection.GetAllAdmin)
			v1.GET("/connection/:id", connection.GetAdmin)

			//
			// Service
			//
			v1.POST("/service/:id", service.AddAdmin)
			// Service Delete
			v1.DELETE("/service/:id", service.DeleteAdmin)
			// Service Update
			v1.PUT("/service/:id", service.UpdateAdmin)
			v1.GET("/service", service.GetAllAdmin)
			v1.GET("/service/:id", service.GetAdmin)
			////
			//// JPNIC Admin
			////
			//v1.POST("/group/network/jpnic", jpnicAdmin.AddAdmin)
			//v1.DELETE("/group/network/jpnic", jpnicAdmin.DeleteAdmin)
			//v1.GET("/group/network/jpnic", jpnicAdmin.GetAdmin)
			////
			//// JPNIC Admin
			////
			//v1.POST("/group/network/jpnic", jpnicTech.AddAdmin)
			//v1.DELETE("/group/network/jpnic", jpnicTech.DeleteAdmin)
			//v1.GET("/group/network/jpnic", jpnicTech.GetAdmin)
		}
	}
	ws := router.Group("/ws")
	{
		v1 := ws.Group("/v1")
		{
			v1.GET("/support", ticket.GetAdminWebSocket)
		}
	}

	go ticket.HandleMessagesAdmin()
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Conf.Controller.Admin.Port), router))
}

func UserRestAPI() {
	router := gin.Default()
	router.Use(cors)

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// Controller
			//
			v1.POST("/controller/chat", controller.ReceiveChatUser)

			//
			// User
			//
			// User Create
			v1.POST("/user", user.Add)
			// User Delete
			//router.DELETE("/user", user.Delete)
			// User Get
			v1.GET("/login", user.GetOwn)
			v1.GET("/user/:id", user.Get)
			v1.GET("/user", user.GetGroup)
			// User ID Get
			// User Update
			v1.PUT("/user/:id", user.Update)

			// User Mail MailVerify
			v1.GET("/verify/:token", user.MailVerify)
			//
			// Token
			//
			// get token for CHAP authentication
			v1.GET("/token/init", token.GenerateInit)
			// get token for user
			v1.GET("/token", token.Generate)
			// delete
			v1.DELETE("/token", token.Delete)
			//
			// Group
			//
			// Group Create
			v1.POST("/group", group.Add)
			v1.GET("/group", group.Get)
			v1.PUT("/group", group.Update)
			v1.GET("/group/all", group.GetAll)
			// Group Delete
			//v1.DELETE("/group", group.Delete)
			// Service add
			v1.POST("/group/service", service.Add)
			v1.GET("/group/service/add_allow", service.GetAddAllow)
			//v1.PUT("/group/network", network.Update)

			// Connection Create
			v1.POST("/group/service/:id/connection", connection.Add)

			//
			// Info
			//
			v1.GET("/group/info", info.Get)

			// Template
			v1.GET("/template", template.Get)

			// NOC
			v1.GET("/noc", noc.GetAll)

			//
			// Support
			//
			v1.POST("/support", ticket.Create)
			v1.GET("/support", ticket.GetTitle)
			v1.POST("/support/:id", chat.Add)
			v1.GET("/support/:id", ticket.Get)
			//
			// Notice
			//
			v1.GET("/notice", notice.Get)

			// 現在検討中

			// Service JPNIC Admin
			v1.POST("/group/jpnic/admin", jpnicAdmin.Add)
			v1.DELETE("/group/jpnic/admin", jpnicAdmin.Delete)
			//v1.GET("/group/network/jpnic/admin", jpnicAdmin.GetAll)
			// Service JPNIC Tech
			v1.POST("/group/jpnic/tech", jpnicTech.Add)
			v1.DELETE("/group/jpnic/tech", jpnicTech.Delete)
			//v1.GET("/group/service/jpnic/tech", jpnicTech.GetAll)
		}
	}

	ws := router.Group("/ws")
	{
		v1 := ws.Group("/v1")
		{
			v1.GET("/support", ticket.GetWebSocket)
		}
	}

	go ticket.HandleMessages()

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Conf.Controller.User.Port), router))
}

func cors(c *gin.Context) {

	//c.Header("Access-Control-Allow-Headers", "Accept, Content-ID, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "*")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Content-ID", "application/json")
	c.Header("Access-Control-Allow-Credentials", "true")
	//c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}
