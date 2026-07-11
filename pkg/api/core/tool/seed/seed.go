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

	db := store.DB()

	// NOC
	log.Println("[Seed] Creating NOC data...")
	nocs, err := createNOCData()
	if err != nil {
		log.Printf("[Seed] Warning: NOC creation skipped: %v", err)
	}

	// BGP Router
	log.Println("[Seed] Creating BGP router...")
	bgpRouter, err := createBGPRouter(nocs["noc01"])
	if err != nil {
		log.Printf("[Seed] Warning: BGP router creation skipped: %v", err)
	}

	// Tunnel EndPoint Routers (+ IPs)
	log.Println("[Seed] Creating tunnel endpoint routers...")
	tunnelEndPointRouterIP, err := createTunnelEndPointRouters(nocs)
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
	log.Println("  - NOC: 10拠点 (NOC01〜, ダミーロケーション)")
	log.Println("  - BGP Router: noc01er01 (NOC01)")
	log.Println("  - Tunnel EndPoint Routers: 26台 / 25 IP (ダミー: 2001:db8::/32, 203.0.113.0/24)")
	log.Println("")
	log.Println("[Seed] Test service & connection:")
	log.Println("  - Service: L3 BGP (1-3B00001)")
	log.Println("  - Connection: EtherIP (1-3B00001-EIP001)")
	log.Println("")
	log.Println("[Seed] Admin API uses Basic Auth from config.json (default: admin/admin)")

	return nil
}

type nocSeed struct {
	key       string
	name      string
	location  string
	bandwidth string
	enable    bool
}

// createNOCData は本番構成を模した複数拠点の NOC を作成し、
// 拠点コードをキーにしたマップで返す（ルータ作成時の親参照に使う）。
func createNOCData() (map[string]*core.NOC, error) {
	seeds := []nocSeed{
		{key: "noc01", name: "NOC01", location: "神奈川県横浜市", bandwidth: "10Gbps", enable: true},
		{key: "noc02", name: "NOC02", location: "東京都千代田区", bandwidth: "1Gbps", enable: false},
		{key: "noc05", name: "NOC05", location: "東京都新宿区", bandwidth: "10Gbps", enable: true},
		{key: "pop03", name: "POP03", location: "東京都品川区", bandwidth: "10Gbps", enable: true},
		{key: "pop52", name: "POP52", location: "大阪府大阪市", bandwidth: "10Gbps", enable: true},
		{key: "noc51", name: "NOC51", location: "大阪府大阪市", bandwidth: "10Gbps", enable: true},
		{key: "noc03", name: "NOC03", location: "東京都港区", bandwidth: "1Gbps", enable: false},
		{key: "pop53", name: "POP53", location: "愛知県名古屋市", bandwidth: "10Gbps", enable: true},
		{key: "noc52", name: "NOC52", location: "福岡県福岡市", bandwidth: "10Gbps", enable: true},
		{key: "noc06", name: "NOC06", location: "東京都大田区", bandwidth: "10Gbps", enable: true},
	}

	nocs := make(map[string]*core.NOC, len(seeds))
	for _, s := range seeds {
		enable := s.enable
		noc, err := dbNOC.Create(&core.NOC{
			Name:      s.name,
			Location:  s.location,
			Bandwidth: s.bandwidth,
			Enable:    &enable,
		})
		if err != nil {
			log.Printf("[Seed] Warning: NOC %s creation skipped: %v", s.name, err)
			continue
		}
		nocs[s.key] = noc
	}

	if len(nocs) == 0 {
		return nil, nil
	}
	return nocs, nil
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

type tunnelRouterIPSeed struct {
	ip     string // ダミーIP（本番値は伏せ、RFC 3849/5737 のドキュメント用レンジを使用）
	enable bool
}

type tunnelRouterSeed struct {
	nocKey   string // createNOCData のキーで親 NOC を参照
	hostName string
	enable   bool
	ips      []tunnelRouterIPSeed // 0個のルータ（IP未割当）もある
}

// createTunnelEndPointRouters は本番の tunnel_end_point_routers /
// tunnel_end_point_router_ips 構成を模したルータ群とその IP を作成する。
// IP はすべてダミー値。本番のホスト部の慣習（er=::1111, ctep=::c179 など）と
// プレフィックス長の有無、IPv4/IPv6 の別、enable フラグは踏襲している。
// 戻り値はテスト接続で使う先頭 IP（noc01er01）。
func createTunnelEndPointRouters(nocs map[string]*core.NOC) (*core.TunnelEndPointRouterIP, error) {
	if len(nocs) == 0 {
		return nil, nil
	}

	routers := []tunnelRouterSeed{
		{nocKey: "noc01", hostName: "noc01er01", enable: true, ips: []tunnelRouterIPSeed{{ip: "2001:db8:1:1::1111/64", enable: true}}},
		{nocKey: "noc01", hostName: "noc01ctep02", enable: true, ips: []tunnelRouterIPSeed{{ip: "2001:db8:1:1::c179/64", enable: true}}},
		{nocKey: "noc02", hostName: "noc02er01", enable: false},
		{nocKey: "noc02", hostName: "noc02ctep01", enable: false, ips: []tunnelRouterIPSeed{{ip: "2001:db8:2:1::c179/64", enable: false}}},
		{nocKey: "noc05", hostName: "noc05er01", enable: true, ips: []tunnelRouterIPSeed{{ip: "2001:db8:3:1::1111", enable: true}}},
		{nocKey: "noc05", hostName: "noc05ctep01", enable: true, ips: []tunnelRouterIPSeed{{ip: "2001:db8:3:1::c179", enable: true}}},
		{nocKey: "pop03", hostName: "pop03er01", enable: false, ips: []tunnelRouterIPSeed{{ip: "2001:db8:4:1::1113/64", enable: false}}},
		{nocKey: "pop03", hostName: "pop03er02", enable: false, ips: []tunnelRouterIPSeed{{ip: "2001:db8:4:2::1111", enable: false}}},
		{nocKey: "pop03", hostName: "pop03er03", enable: true, ips: []tunnelRouterIPSeed{{ip: "2001:db8:4:3::1113/64", enable: true}}},
		{nocKey: "pop03", hostName: "pop03ctep01", enable: true, ips: []tunnelRouterIPSeed{{ip: "2001:db8:4:4::c179/64", enable: true}}},
		{nocKey: "pop03", hostName: "pop03ctep02", enable: true, ips: []tunnelRouterIPSeed{{ip: "2001:db8:4:5::c179/64", enable: true}}},
		{nocKey: "pop52", hostName: "pop52er01", enable: true, ips: []tunnelRouterIPSeed{{ip: "2001:db8:5:1::1111", enable: true}}},
		{nocKey: "pop52", hostName: "pop52ctep02", enable: true, ips: []tunnelRouterIPSeed{{ip: "2001:db8:5:1::c179", enable: true}}},
		{nocKey: "noc51", hostName: "noc51er01", enable: true, ips: []tunnelRouterIPSeed{{ip: "2001:db8:6:1::1111/64", enable: true}}},
		{nocKey: "noc51", hostName: "noc51ctep01", enable: true, ips: []tunnelRouterIPSeed{{ip: "2001:db8:6:1::c179/64", enable: true}}},
		{nocKey: "noc03", hostName: "noc03er01", enable: false, ips: []tunnelRouterIPSeed{{ip: "2001:db8:7:1::5:9105:179/64", enable: false}}},
		{nocKey: "noc03", hostName: "noc03ctep01", enable: false, ips: []tunnelRouterIPSeed{{ip: "2001:db8:7:1::c179/64", enable: false}}},
		{nocKey: "pop52", hostName: "pop52er02", enable: false, ips: []tunnelRouterIPSeed{{ip: "2001:db8:5:2::1112/64", enable: false}}},
		{nocKey: "pop53", hostName: "pop53er01", enable: true, ips: []tunnelRouterIPSeed{{ip: "2001:db8:8:1::1111/64", enable: true}}},
		{nocKey: "pop53", hostName: "pop53ctep01", enable: true, ips: []tunnelRouterIPSeed{{ip: "2001:db8:8:1::c179/64", enable: true}}},
		{nocKey: "noc52", hostName: "noc52er01", enable: true, ips: []tunnelRouterIPSeed{{ip: "2001:db8:9:1::1111/64", enable: true}}},
		{nocKey: "noc52", hostName: "noc52ctep01", enable: true, ips: []tunnelRouterIPSeed{{ip: "2001:db8:9:1::6/64", enable: true}}},
		{nocKey: "noc06", hostName: "noc06er01", enable: true, ips: []tunnelRouterIPSeed{{ip: "2001:db8:a:1::1111", enable: true}}},
		{nocKey: "noc06", hostName: "noc06ctep01", enable: true, ips: []tunnelRouterIPSeed{{ip: "2001:db8:a:1::c179", enable: true}}},
		{nocKey: "pop03", hostName: "pop03er04", enable: true, ips: []tunnelRouterIPSeed{{ip: "2001:db8:4:6::face/64", enable: true}}},
	}

	var firstIP *core.TunnelEndPointRouterIP
	for _, rs := range routers {
		noc := nocs[rs.nocKey]
		if noc == nil {
			continue
		}

		routerEnable := rs.enable
		router, err := dbTunnelEndPointRouter.Create(&core.TunnelEndPointRouter{
			NOCID:    &noc.ID,
			HostName: rs.hostName,
			Enable:   &routerEnable,
		})
		if err != nil {
			log.Printf("[Seed] Warning: tunnel router %s creation skipped: %v", rs.hostName, err)
			continue
		}

		for _, ipSeed := range rs.ips {
			ipEnable := ipSeed.enable
			createdIP, err := dbTunnelEndPointRouterIP.Create(&core.TunnelEndPointRouterIP{
				TunnelEndPointRouterID: &router.ID,
				IP:                     ipSeed.ip,
				Enable:                 &ipEnable,
			})
			if err != nil {
				log.Printf("[Seed] Warning: tunnel router IP %s creation skipped: %v", ipSeed.ip, err)
				continue
			}
			if firstIP == nil {
				firstIP = createdIP
			}
		}
	}

	return firstIP, nil
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
