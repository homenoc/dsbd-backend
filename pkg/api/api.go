package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	catalog "github.com/homenoc/dsbd-backend/pkg/api/core/catalog/v0"
	controller "github.com/homenoc/dsbd-backend/pkg/api/core/controller/v0"
	connection "github.com/homenoc/dsbd-backend/pkg/api/core/group/connection/v0"
	memo "github.com/homenoc/dsbd-backend/pkg/api/core/group/memo/v0"
	service "github.com/homenoc/dsbd-backend/pkg/api/core/group/service/v0"
	group "github.com/homenoc/dsbd-backend/pkg/api/core/group/v0"
	mail "github.com/homenoc/dsbd-backend/pkg/api/core/mail/v0"
	bgpRouter "github.com/homenoc/dsbd-backend/pkg/api/core/noc/bgpRouter/v0"
	tunnelEndPointRouter "github.com/homenoc/dsbd-backend/pkg/api/core/noc/tunnelEndPointRouter/v0"
	tunnelEndPointRouterIP "github.com/homenoc/dsbd-backend/pkg/api/core/noc/tunnelEndPointRouterIP/v0"
	noc "github.com/homenoc/dsbd-backend/pkg/api/core/noc/v0"
	notice "github.com/homenoc/dsbd-backend/pkg/api/core/notice/v0"
	payment "github.com/homenoc/dsbd-backend/pkg/api/core/payment/v0"
	ticket "github.com/homenoc/dsbd-backend/pkg/api/core/support/ticket/v0"
	token "github.com/homenoc/dsbd-backend/pkg/api/core/token/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/config"
	user "github.com/homenoc/dsbd-backend/pkg/api/core/user/v0"
	"github.com/homenoc/dsbd-backend/pkg/api/middleware"
)

// NewAdminRouter builds the admin API router (routes only, no background
// goroutines and no listener) so tests can drive it via httptest.
func NewAdminRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// Controller
			//noc
			v1.POST("/controller/chat", controller.ReceiveChatByAdmin)

			// Notice
			//
			v1.POST("/notice", notice.AddByAdmin)
			v1.DELETE("/notice/:id", notice.DeleteByAdmin)
			v1.GET("/notice", notice.GetAllByAdmin)
			v1.GET("/notice/:id", notice.GetByAdmin)
			v1.PUT("/notice/:id", notice.UpdateByAdmin)

			//
			// User
			//
			// User Create
			v1.POST("/user", user.AddByAdmin)
			// User Delete
			v1.DELETE("/user", user.DeleteByAdmin)
			// User Update
			v1.PUT("/user/:id", user.UpdateByAdmin)
			v1.GET("/user", user.GetAllByAdmin)
			v1.GET("/user/:id", user.GetByAdmin)
			//
			// Login / Logout
			//
			v1.POST("/login", token.GenerateByAdmin)
			v1.POST("/logout", token.DeleteAdminUser)

			//
			// BotToken
			//
			v1.POST("/token/generate", token.GenerateByAdmin)

			v1.POST("/token", token.AddByAdmin)
			// BotToken Delete
			v1.DELETE("/token", token.DeleteAllByAdmin)
			v1.DELETE("/token/:id", token.DeleteByAdmin)
			// BotToken Update
			v1.PUT("/token/:id", token.UpdateByAdmin)
			v1.GET("/token", token.GetAllByAdmin)
			v1.GET("/token/:id", token.GetByAdmin)
			//
			// Group
			//
			v1.POST("/group", group.AddByAdmin)
			// Group Delete
			v1.DELETE("/group", group.DeleteByAdmin)
			// Group Update
			v1.PUT("/group/:id", group.UpdateByAdmin)
			v1.GET("/group", group.GetAllByAdmin)
			v1.GET("/group/:id", group.GetByAdmin)

			//
			// Memo
			//
			v1.POST("/memo", memo.AddByAdmin)
			v1.DELETE("/memo/:id", memo.DeleteByAdmin)

			// Template
			v1.GET("/catalog", catalog.GetByAdmin)

			//
			// NOC
			//
			v1.POST("/noc", noc.AddByAdmin)
			v1.GET("/noc", noc.GetAllByAdmin)
			v1.DELETE("/noc/:id", noc.DeleteByAdmin)
			v1.GET("/noc/:id", noc.GetByAdmin)
			v1.PUT("/noc/:id", noc.UpdateByAdmin)

			//
			// NOC Router
			//
			v1.POST("/router", bgpRouter.AddByAdmin)
			v1.GET("/router", bgpRouter.GetAllByAdmin)
			v1.DELETE("/router/:id", bgpRouter.DeleteByAdmin)
			v1.GET("/router/:id", bgpRouter.GetByAdmin)
			v1.PUT("/router/:id", bgpRouter.UpdateByAdmin)

			//
			// NOC Gateway
			//
			v1.POST("/gateway", tunnelEndPointRouter.AddByAdmin)
			v1.GET("/gateway", tunnelEndPointRouter.GetAllByAdmin)
			v1.DELETE("/gateway/:id", tunnelEndPointRouter.DeleteByAdmin)
			v1.GET("/gateway/:id", tunnelEndPointRouter.GetByAdmin)
			v1.PUT("/gateway/:id", tunnelEndPointRouter.UpdateByAdmin)

			//
			// NOC Gateway IP
			//
			v1.POST("/gateway_ip", tunnelEndPointRouterIP.AddByAdmin)
			v1.GET("/gateway_ip", tunnelEndPointRouterIP.GetAllByAdmin)
			v1.DELETE("/gateway_ip/:id", tunnelEndPointRouterIP.DeleteByAdmin)
			v1.GET("/gateway_ip/:id", tunnelEndPointRouterIP.GetByAdmin)
			v1.PUT("/gateway_ip/:id", tunnelEndPointRouterIP.UpdateByAdmin)

			//
			// Support
			//
			v1.POST("/support", ticket.CreateByAdmin)
			v1.GET("/support", ticket.GetAllByAdmin)
			//v1.POST("/support/:id", chat.AddByAdmin)
			v1.GET("/support/:id", ticket.GetByAdmin)
			v1.PUT("/support/:id", ticket.UpdateByAdmin)

			////
			//// Connection
			////
			v1.POST("/service/:id/connection", connection.AddByAdmin)
			// Group Delete
			v1.DELETE("/connection/:id", connection.DeleteByAdmin)
			// Group Update
			v1.PUT("/connection/:id", connection.UpdateByAdmin)
			v1.GET("/connection", connection.GetAllByAdmin)
			v1.GET("/connection/:id", connection.GetByAdmin)

			//
			// Service
			//
			v1.POST("/group/:id/service", service.AddByAdmin)
			// Service Delete
			v1.DELETE("/service/:id", service.DeleteByAdmin)
			// Service Update
			v1.PUT("/service/:id", service.UpdateByAdmin)
			v1.GET("/service", service.GetAllByAdmin)
			v1.GET("/service/:id", service.GetByAdmin)

			//
			// Payment
			//
			v1.POST("/group/:id/payment/subscribe", payment.PostAdminSubscribeGettingURL)
			v1.GET("/group/:id/payment/subscribe", payment.GetAdminDashboardSubscribeURL)
			v1.GET("/group/:id/payment", payment.GetAdminBillingPortalURL)
			v1.GET("/group/:id/payment/customer", payment.GetAdminDashboardCustomerURL)

			//
			// JPNIC ByAdmin
			//
			v1.POST("/service/:id/jpnic_admin", service.AddJPNICAdminByAdmin)
			v1.DELETE("/jpnic_admin/:id", service.DeleteJPNICAdminByAdmin)
			v1.PUT("/jpnic_admin/:id", service.UpdateJPNICAdminByAdmin)

			//
			// JPNIC Tech
			//
			v1.POST("/service/:id/jpnic_tech", service.AddJPNICTechByAdmin)
			v1.DELETE("/jpnic_tech/:id", service.DeleteJPNICTechByAdmin)
			v1.PUT("/jpnic_tech/:id", service.UpdateJPNICTechByAdmin)

			//
			// IP
			//
			v1.POST("/service/:id/ip", service.AddIPByAdmin)
			v1.DELETE("/ip/:id", service.DeleteIPByAdmin)
			v1.PUT("/ip/:id", service.UpdateIPByAdmin)

			//
			// Plan
			//
			v1.POST("/ip/:id/plan", service.AddPlanByAdmin)
			v1.DELETE("/plan/:id", service.DeletePlanByAdmin)
			v1.PUT("/plan/:id", service.UpdatePlanByAdmin)

			//
			// Mail
			//
			v1.POST("/mail", mail.SendByAdmin)
		}
	}
	ws := router.Group("/ws")
	{
		v1 := ws.Group("/v1")
		{
			v1.GET("/support", ticket.GetAdminWebSocket)
		}
	}

	return router
}

func AdminRestAPI() {
	if !config.IsDebug {
		gin.SetMode(gin.ReleaseMode)
	}
	router := NewAdminRouter()

	go token.TokenRemove()
	go ticket.HandleMessagesByAdmin()
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Conf.Controller.Admin.Port), router))
}

// NewUserRouter builds the user API router (routes only, no background
// goroutines and no listener) so tests can drive it via httptest.
func NewUserRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors)

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	api := router.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// Controller
			//
			v1.POST("/controller/chat", controller.ReceiveChatUser)

			// Stripe
			//
			//v1.POST("/stripe", payment.GetStripeWebHook)
			v1.Any("/stripe", payment.GetStripeWebHook)

			// User Mail MailVerify
			v1.GET("/verify/:token", user.MailVerify)

			//
			// Login / Logout
			//
			v1.POST("/login", token.Generate)
			v1.GET("/login", token.GenerateInit)
			v1.POST("/logout", token.Delete)

			//
			// User
			//
			// User Create (public: registration)
			v1.POST("/user", user.Add)
			// User Create(Group)
			v1.POST("/group/:id/user", middleware.GroupAuth(0), user.AddGroup)
			// Antisocial Check
			v1.PUT("/user/antisocial/agree", middleware.UserAuth, user.AgreeAntisocialCheck)
			// User Update
			v1.PUT("/user/:id", middleware.UserAuth, user.Update)
			// User Delete
			v1.DELETE("/user/:id", middleware.GroupAuth(0), user.Delete)

			//
			// Per-resource user reads (formerly the /info bootstrap blob)
			//
			v1.GET("/user/me", middleware.UserAuth, user.GetMe)
			v1.GET("/group", middleware.UserAuth, group.Get)
			v1.GET("/service", middleware.UserAuth, service.Get)
			v1.GET("/connection", middleware.UserAuth, connection.Get)
			v1.GET("/notice", middleware.UserAuth, notice.GetActive)
			v1.GET("/ticket", middleware.UserAuth, ticket.GetUserTickets)
			v1.GET("/request", middleware.UserAuth, ticket.GetUserRequests)
			// Derived network summary of the user's opened connections
			v1.GET("/info", middleware.UserAuth, connection.GetNetworkInfo)

			//
			// Group
			//
			// Group Create
			v1.POST("/group", middleware.UserAuth, group.Add)

			// Template
			v1.GET("/catalog", middleware.UserAuth, catalog.Get)

			// Service add
			v1.POST("/service", middleware.GroupAuth(0), service.Add)
			v1.GET("/service/add_allow", middleware.GroupAuth(0), service.GetAddAllow)
			//v1.PUT("/group/network", network.Update)

			// Connection Create
			v1.POST("/service/:id/connection", middleware.GroupAuth(0), connection.Add)

			//
			// Payment
			//
			v1.POST("/payment/subscribe", middleware.GroupAuth(0), payment.PostSubscribeGettingURL)
			v1.GET("/payment", middleware.GroupAuth(0), payment.GetBillingPortalURL)

			//
			// Support/Request
			//
			v1.POST("/support", ticket.Create)
			v1.POST("/request", ticket.Request)
			v1.PUT("/support/:id", ticket.Update)

			// Group Delete
			//v1.DELETE("/group", group.Delete)

			//v1.POST("/support/:id", chat.Add)
		}
	}

	//
	// Stripe
	//
	router.POST("/stripe", payment.GetStripeWebHook)

	ws := router.Group("/ws")
	{
		v1 := ws.Group("/v1")
		{
			v1.GET("/support", ticket.GetWebSocket)
		}
	}

	return router
}

func UserRestAPI() {
	if !config.IsDebug {
		gin.SetMode(gin.ReleaseMode)
	}
	router := NewUserRouter()

	go ticket.HandleMessages()

	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Conf.Controller.User.Port), router))
}

func cors(c *gin.Context) {
	// A wildcard origin ("*") together with Allow-Credentials:true is rejected by
	// browsers for credentialed requests, so we reflect a specific allowed origin
	// instead. config.cors.origins is an exact-match allowlist; when empty the
	// request Origin is echoed back (valid with credentials).
	if origin := c.Request.Header.Get("Origin"); origin != "" && originAllowed(origin) {
		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Vary", "Origin")
		c.Header("Access-Control-Allow-Credentials", "true")
	}
	c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers",
		"Content-Type, USER_TOKEN, ACCESS_TOKEN, HASH_PASS, Email, USER, PASS")

	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.AbortWithStatus(http.StatusOK)
	}
}

// originAllowed reports whether origin is in the config.cors.origins
// allowlist. Fail-safe: an empty/missing allowlist allows NOTHING cross-origin
// — a deployment that forgets to configure it must not silently accept
// credentialed requests from anywhere. Dev origins are listed in
// configs/config.json.
func originAllowed(origin string) bool {
	for _, o := range config.Conf.CORS.Origins {
		if o == origin {
			return true
		}
	}
	return false
}
