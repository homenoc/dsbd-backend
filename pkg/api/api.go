package api

import (
	"github.com/gin-gonic/gin"
	connection "github.com/homenoc/dsbd-backend/pkg/api/core/group/connection/v0"
	jpnicUser "github.com/homenoc/dsbd-backend/pkg/api/core/group/network/jpnic_user/v0"
	networkUser "github.com/homenoc/dsbd-backend/pkg/api/core/group/network/network_user/v0"
	network "github.com/homenoc/dsbd-backend/pkg/api/core/group/network/v0"
	group "github.com/homenoc/dsbd-backend/pkg/api/core/group/v0"
	token "github.com/homenoc/dsbd-backend/pkg/api/core/token/v0"
	user "github.com/homenoc/dsbd-backend/pkg/api/core/user/v0"
	"log"
	"net/http"
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
			v1.GET("/user/:id", user.GetAllAdmin)
			//
			// Token
			//
			v1.POST("/token", token.AddAdmin)
			// User Delete
			v1.DELETE("/token", token.DeleteAdmin)
			// User Update
			v1.PUT("/token", token.UpdateAdmin)
			v1.GET("/token", token.GetAllAdmin)
			v1.GET("/token/:id", token.GetAllAdmin)
			//
			// Group
			//
			v1.POST("/group", group.AddAdmin)
			// Group Delete
			v1.DELETE("/group", group.DeleteAdmin)
			// Group Update
			v1.PUT("/group", group.UpdateAdmin)
			v1.GET("/group", group.GetAllAdmin)
			v1.GET("/group/:id", group.GetAllAdmin)
			//
			// Network
			//
			v1.POST("/group/network", network.AddAdmin)
			// Group Delete
			v1.DELETE("/group/network", network.DeleteAdmin)
			// Group Update
			v1.PUT("/group/network", network.UpdateAdmin)
			v1.GET("/group/network", network.GetAllAdmin)
			v1.GET("/group/network/:id", network.GetAllAdmin)
			//
			// JPNIC User
			//
			v1.POST("/group/network/jpnic", jpnicUser.AddAdmin)
			// Group Delete
			v1.DELETE("/group/network/jpnic", jpnicUser.DeleteAdmin)
			// Group Update
			v1.PUT("/group/network/jpnic", jpnicUser.UpdateAdmin)
			v1.GET("/group/network/jpnic", jpnicUser.GetAllAdmin)
			v1.GET("/group/network/jpnic/:id", jpnicUser.GetAllAdmin)
			//
			// Connection
			//
			v1.POST("/group/connection", connection.AddAdmin)
			// Group Delete
			v1.DELETE("/group/connection", connection.DeleteAdmin)
			// Group Update
			v1.PUT("/group/connection", connection.UpdateAdmin)
			v1.GET("/group/connection", connection.GetAllAdmin)
			v1.GET("/group/connection/:id", connection.GetAllAdmin)
		}
	}
	log.Fatal(http.ListenAndServe(":8080", router))
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
			// User Update
			v1.PUT("/user/:id", user.Update)
			// User Mail MailVerify
			v1.POST("/user/verify/:token", user.MailVerify)
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
			// Network User
			v1.POST("/group/network/user", networkUser.Add)
			v1.DELETE("/group/network/user", networkUser.Delete)
			v1.PUT("/group/network/user", networkUser.Update)
			// Network add
			v1.POST("/group/network", network.Add)
			v1.PUT("/group/network", network.Update)
			// Network Confirm
			v1.POST("/group/network/confirm", network.Confirm)
			// Network JPNIC User
			v1.POST("/group/network/jpnic", jpnicUser.Add)
			v1.PUT("/group/network/jpnic", jpnicUser.Update)

		}
	}
	log.Fatal(http.ListenAndServe(":8080", router))
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
