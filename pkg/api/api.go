package api

import (
	"github.com/gin-gonic/gin"
	connection "github.com/homenoc/dsbd-backend/pkg/api/core/group/connection/v0"
	info "github.com/homenoc/dsbd-backend/pkg/api/core/group/info/v0"
	network "github.com/homenoc/dsbd-backend/pkg/api/core/group/network/v0"
	group "github.com/homenoc/dsbd-backend/pkg/api/core/group/v0"
	notice "github.com/homenoc/dsbd-backend/pkg/api/core/notice/v0"
	chat "github.com/homenoc/dsbd-backend/pkg/api/core/support/chat/v0"
	ticket "github.com/homenoc/dsbd-backend/pkg/api/core/support/ticket/v0"
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

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			//
			// User
			//
			// User Create
			v1.POST("/user", user.AddAdmin)
			// User Delete
			v1.DELETE("/user", user.DeleteAdmin)
			// User Update
			v1.PUT("/user", user.UpdateAdmin)
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
			v1.PUT("/token", token.UpdateAdmin)
			v1.GET("/token", token.GetAllAdmin)
			v1.GET("/token/:id", token.GetAdmin)
			//
			// Group
			//
			v1.POST("/group", group.AddAdmin)
			// Group Delete
			v1.DELETE("/group", group.DeleteAdmin)
			// Group Update
			v1.PUT("/group", group.UpdateAdmin)
			v1.GET("/group", group.GetAllAdmin)
			v1.GET("/group/:id", group.GetAdmin)

			//
			// Support
			//
			v1.POST("/support", ticket.CreateAdmin)
			v1.GET("/support", ticket.GetAllAdmin)
			//v1.POST("/support/:id", chat.AddAdmin)
			v1.GET("/support/:id", ticket.GetAdmin)
			v1.PUT("/support/:id", ticket.UpdateAdmin)
			////
			//// Network
			////
			//v1.POST("/group/network", network.AddAdmin)
			//// Group Delete
			//v1.DELETE("/group/network", network.DeleteAdmin)
			//// Group Update
			//v1.PUT("/group/network", network.UpdateAdmin)
			//v1.GET("/group/network", network.GetAllAdmin)
			//v1.GET("/group/network/:id", network.GetAdmin)
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
			////
			//// Connection
			////
			//v1.POST("/group/connection", connection.AddAdmin)
			//// Group Delete
			//v1.DELETE("/group/connection", connection.DeleteAdmin)
			//// Group Update
			//v1.PUT("/group/connection", connection.UpdateAdmin)
			//v1.GET("/group/connection", connection.GetAllAdmin)
			//v1.GET("/group/connection/:id", connection.GetAdmin)
		}
	}
	ws := router.Group("/ws")
	{
		v1 := ws.Group("/v1")
		{
			v1.GET("/support", ticket.GetWebSocket)
		}
	}

	//go ticket.HandleMessages()
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Conf.Controller.Admin.Port), router))
}

func UserRestAPI() {
	router := gin.Default()
	router.Use(cors)

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			//
			// User
			//
			// User Create
			v1.POST("/user", user.Add)
			// User Delete
			//router.DELETE("/user", user.Delete)
			// User Get
			v1.GET("/user", user.Get)
			v1.GET("/user/all", user.GetGroup)
			// User ID Get
			// v1.GET("/user/:id",user.GetId)
			// User Update
			v1.PUT("/user/:id", user.Update)
			// User Mail MailVerify
			v1.GET("/user/verify/:token", user.MailVerify)
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
			// Connection Create
			v1.POST("/group/connection", connection.Add)
			// Network add
			v1.POST("/group/network", network.Add)
			//v1.PUT("/group/network", network.Update)
			//
			// Info
			//
			v1.GET("/group/info", info.Get)

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

			// Network JPNIC Admin
			//v1.POST("/group/network/jpnic/admin", jpnicAdmin.Add)
			//v1.DELETE("/group/network/jpnic/admin", jpnicAdmin.Delete)
			//v1.GET("/group/network/jpnic/admin", jpnicAdmin.Get)
			// Network JPNIC Tech
			//v1.POST("/group/network/jpnic/tech", jpnicTech.Add)
			//v1.DELETE("/group/network/jpnic/tech", jpnicTech.Delete)
			//v1.GET("/group/network/jpnic/tech", jpnicTech.Get)
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
