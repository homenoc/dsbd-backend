package seed

import (
	"log"
	"strings"
	"time"

	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/hash"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	dbConnection "github.com/homenoc/dsbd-backend/pkg/api/store/group/connection/v0"
	dbService "github.com/homenoc/dsbd-backend/pkg/api/store/group/service/v0"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/api/store/group/v0"
	dbBGPRouter "github.com/homenoc/dsbd-backend/pkg/api/store/noc/bgpRouter/v0"
	dbTunnelEndPointRouter "github.com/homenoc/dsbd-backend/pkg/api/store/noc/tunnelEndPointRouter/v0"
	dbTunnelEndPointRouterIP "github.com/homenoc/dsbd-backend/pkg/api/store/noc/tunnelEndPointRouterIP/v0"
	dbNOC "github.com/homenoc/dsbd-backend/pkg/api/store/noc/v0"
	dbUser "github.com/homenoc/dsbd-backend/pkg/api/store/user/v0"
)

func Run() error {
	log.Println("[Seed] Starting seed data creation...")

	db, err := store.ConnectDB()
	if err != nil {
		return err
	}

	// NOC
	log.Println("[Seed] Creating NOC data...")
	noc, err := createNOCData()
	if err != nil {
		log.Printf("[Seed] Warning: NOC creation skipped: %v", err)
	}

	// BGP Router
	log.Println("[Seed] Creating BGP router...")
	bgpRouter, err := createBGPRouter(noc)
	if err != nil {
		log.Printf("[Seed] Warning: BGP router creation skipped: %v", err)
	}

	// Tunnel EndPoint Router (+ IP)
	log.Println("[Seed] Creating tunnel endpoint router...")
	tunnelEndPointRouterIP, err := createTunnelEndPointRouter(noc)
	if err != nil {
		log.Printf("[Seed] Warning: Tunnel endpoint router creation skipped: %v", err)
	}

	// Group
	log.Println("[Seed] Creating test group...")
	group, err := createTestGroup()
	if err != nil {
		log.Printf("[Seed] Warning: Group creation skipped: %v", err)
	}

	// Users
	log.Println("[Seed] Creating test users...")
	if err := createTestUsers(db, group); err != nil {
		return err
	}

	// Service
	log.Println("[Seed] Creating test service...")
	service, err := createTestService(group)
	if err != nil {
		log.Printf("[Seed] Warning: Service creation skipped: %v", err)
	}

	// Connection
	log.Println("[Seed] Creating test connection...")
	if err := createTestConnection(service, bgpRouter, tunnelEndPointRouterIP); err != nil {
		log.Printf("[Seed] Warning: Connection creation skipped: %v", err)
	}

	log.Println("[Seed] Seed data creation completed!")
	log.Println("")
	log.Println("[Seed] Test accounts (User API):")
	log.Println("  - Master: master@example.com / password (Level 1: 申請・変更・閲覧)")
	log.Println("  - Member: member@example.com / password (Level 2: 一般メンバー)")
	log.Println("  ※ 反社チェック未同意状態（PUT /api/v1/user/antisocial/agree で同意）")
	log.Println("")
	log.Println("[Seed] Test NOC & routers:")
	log.Println("  - NOC: NOC01 (神奈川県横浜市)")
	log.Println("  - BGP Router: noc01er01")
	log.Println("  - Tunnel EndPoint Router: noc01er01 (IP: 2404:7a81:920:9600::1111/64)")
	log.Println("")
	log.Println("[Seed] Test service & connection:")
	log.Println("  - Service: L3 BGP (1-3B00001)")
	log.Println("  - Connection: EtherIP (1-3B00001-EIP001)")
	log.Println("")
	log.Println("[Seed] Admin API uses Basic Auth from config.json (default: admin/admin)")

	return nil
}

func createNOCData() (*core.NOC, error) {
	enable := true

	noc := &core.NOC{
		Name:      "NOC01",
		Location:  "神奈川県横浜市",
		Bandwidth: "2Gbps",
		Enable:    &enable,
	}

	return dbNOC.Create(noc)
}

func createBGPRouter(noc *core.NOC) (*core.BGPRouter, error) {
	if noc == nil {
		return nil, nil
	}

	enable := true

	bgpRouter := &core.BGPRouter{
		NOCID:    noc.ID,
		HostName: "noc01er01",
		Enable:   &enable,
	}

	return dbBGPRouter.Create(bgpRouter)
}

func createTunnelEndPointRouter(noc *core.NOC) (*core.TunnelEndPointRouterIP, error) {
	if noc == nil {
		return nil, nil
	}

	enable := true

	router := &core.TunnelEndPointRouter{
		NOCID:    &noc.ID,
		HostName: "noc01er01",
		Enable:   &enable,
	}

	router, err := dbTunnelEndPointRouter.Create(router)
	if err != nil {
		return nil, err
	}

	routerIP := &core.TunnelEndPointRouterIP{
		TunnelEndPointRouterID: &router.ID,
		IP:                     "2404:7a81:920:9600::1111/64",
		Enable:                 &enable,
	}

	return dbTunnelEndPointRouterIP.Create(routerIP)
}

func createTestGroup() (*core.Group, error) {
	agree := true
	pass := true
	addAllow := true
	expiredStatus := uint(0)

	group := &core.Group{
		Agree:         &agree,
		Question:      "開発用テストグループです",
		Org:           "テスト組織",
		OrgEn:         "Test Organization",
		PostCode:      "100-0001",
		Address:       "東京都千代田区",
		AddressEn:     "Chiyoda-ku, Tokyo",
		Tel:           "03-1234-5678",
		Country:       "Japan",
		MemberType:    1,
		Pass:          &pass,
		ExpiredStatus: &expiredStatus,
		AddAllow:      &addAllow,
	}

	return dbGroup.Create(group)
}

func createTestUsers(db interface{}, group *core.Group) error {
	mailVerify := true
	antisocialCheck := false // 実際のフローに合わせて未同意状態
	expiredStatus := uint(0)

	var groupID *uint
	if group != nil {
		groupID = &group.ID
	}

	// Master user (Level 1: グループ内の申請・変更・閲覧可能)
	masterUser := &core.User{
		GroupID:         groupID,
		Name:            "マスターユーザー",
		NameEn:          "Master User",
		Email:           "master@example.com",
		Pass:            strings.ToLower(hash.Generate("password")),
		ExpiredStatus:   &expiredStatus,
		Level:           1, // Master: グループ内の申請・変更・閲覧可能
		MailVerify:      &mailVerify,
		AntisocialCheck: &antisocialCheck,
	}

	if err := dbUser.Create(masterUser); err != nil {
		log.Printf("[Seed] Warning: Master user creation skipped: %v", err)
	}

	// Member user (Level 2: 一般メンバー)
	memberUser := &core.User{
		GroupID:         groupID,
		Name:            "一般ユーザー",
		NameEn:          "Member User",
		Email:           "member@example.com",
		Pass:            strings.ToLower(hash.Generate("password")),
		ExpiredStatus:   &expiredStatus,
		Level:           2, // Member: 一般メンバー
		MailVerify:      &mailVerify,
		AntisocialCheck: &antisocialCheck,
	}

	if err := dbUser.Create(memberUser); err != nil {
		log.Printf("[Seed] Warning: Member user creation skipped: %v", err)
	}

	return nil
}

func createTestService(group *core.Group) (*core.Service, error) {
	if group == nil {
		return nil, nil
	}

	pass := true
	enable := true
	addAllow := true

	service := &core.Service{
		GroupID:        group.ID,
		ServiceType:    "3B00", // L3 BGP
		ServiceComment: "開発用テストサービス",
		ServiceNumber:  1,
		Org:            "テスト組織",
		OrgEn:          "Test Organization",
		PostCode:       "100-0001",
		Address:        "東京都千代田区",
		AddressEn:      "Chiyoda-ku, Tokyo",
		AveUpstream:    100,
		MaxUpstream:    1000,
		AveDownstream:  100,
		MaxDownstream:  1000,
		StartDate:      time.Now(),
		Pass:           &pass,
		Enable:         &enable,
		AddAllow:       &addAllow,
	}

	return dbService.Create(service)
}

func createTestConnection(service *core.Service, bgpRouter *core.BGPRouter, tunnelEndPointRouterIP *core.TunnelEndPointRouterIP) error {
	if service == nil {
		return nil
	}

	open := false
	enable := true
	monitor := false

	connection := &core.Connection{
		ServiceID:         service.ID,
		ConnectionType:    "EIP", // EtherIP
		ConnectionComment: "開発用テスト接続",
		ConnectionNumber:  1,
		NTT:               "はい（IPoEによりIPv6インターネットへ接続可能）",
		PreferredAP:       "東日本",
		TermIP:            "192.0.2.1",
		Address:           "東京都",
		Open:              &open,
		Enable:            &enable,
		Monitor:           &monitor,
	}

	if bgpRouter != nil {
		connection.BGPRouterID = &bgpRouter.ID
	}
	if tunnelEndPointRouterIP != nil {
		connection.TunnelEndPointRouterIPID = &tunnelEndPointRouterIP.ID
	}

	_, err := dbConnection.Create(connection)
	return err
}
