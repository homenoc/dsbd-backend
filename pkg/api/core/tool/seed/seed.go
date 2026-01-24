package seed

import (
	"log"
	"strings"

	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/tool/hash"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	dbGroup "github.com/homenoc/dsbd-backend/pkg/api/store/group/v0"
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
	if err := createNOCData(); err != nil {
		log.Printf("[Seed] Warning: NOC creation skipped: %v", err)
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

	log.Println("[Seed] Seed data creation completed!")
	log.Println("")
	log.Println("[Seed] Test accounts (User API):")
	log.Println("  - Master: master@example.com / password (Level 1: 申請・変更・閲覧)")
	log.Println("  - Member: member@example.com / password (Level 2: 一般メンバー)")
	log.Println("  ※ 反社チェック未同意状態（PUT /api/v1/user/antisocial/agree で同意）")
	log.Println("")
	log.Println("[Seed] Admin API uses Basic Auth from config.json (default: admin/admin)")

	return nil
}

func createNOCData() error {
	enable := true

	noc := &core.NOC{
		Name:      "Tokyo NOC",
		Location:  "Tokyo",
		Bandwidth: "10Gbps",
		Enable:    &enable,
		Comment:   "Development NOC",
	}

	_, err := dbNOC.Create(noc)
	return err
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
