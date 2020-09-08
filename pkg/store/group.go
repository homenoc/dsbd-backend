package store

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/auth"
	"log"
	"time"
)

func CreateGroupToDB(group auth.Group) error {
	db := connectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: database connection\n", time.Now())
	}
	defer db.Close()
	writeTable, err := db.Prepare(`INSERT INTO "group" ("created_at","org_ja","org","status","tech_id","postcode",
"address_ja","address","mail","phone") VALUES (?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		log.Println("write error |error: ", err)
		return fmt.Errorf("(%s)error: write error\n %s", time.Now(), err)
	}
	if _, err := writeTable.Exec(time.Now().Unix(), group.OrgJa, group.Org, group.Status, group.TechID, group.PostCode, group.AddressJa, group.Address, group.Mail, group.Phone); err != nil {
		log.Println("apply error |error: ", err)
		return fmt.Errorf("(%s)error: apply error\n %s", time.Now(), err)
	}
	return nil
}

func DeleteGroupToDB(id int) error {
	db := connectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: database connection\n", time.Now())
	}
	defer db.Close()

	if _, err := db.Exec("DELETE FROM group WHERE id = ?", id); err != nil {
		log.Println("database delete table error |", err)
		return fmt.Errorf("(%s)error: delete table\n", time.Now())
	}
	return nil
}

func UpdateGroupToDB(group auth.Group) error {
	db := connectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: database connection\n", time.Now())
	}
	defer db.Close()

	if _, err := db.Exec("UPDATE group SET updated_at = ?,org_ja = ?,org = ?,status = ?,tech_id = ?,"+
		"postcode = ?,address_ja = ?,address = ?,mail = ?,phone = ? WHERE id = ?",
		time.Now().Unix(), group.OrgJa, group.Org, group.Status, group.TechID, group.PostCode, group.AddressJa,
		group.Address, group.Mail, group.Phone); err != nil {
		log.Println("database update table error |", err)
		return fmt.Errorf("(%s)error: update table\n", time.Now())
	}
	return nil
}

func GetGroupOrgNameFromDB(orgName string) *groupResult {
	db := connectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return &groupResult{err: fmt.Errorf("(%s)error: database connection\n", time.Now())}
	}
	defer db.Close()

	rows := db.QueryRow("SELECT * FROM group WHERE org_ja = ?", orgName)
	var group auth.Group
	err := rows.Scan(&group.ID, &group.CreatedAt, &group.UpdatedAt, &group.OrgJa, &group.Org, &group.Status,
		&group.TechID, &group.PostCode, &group.AddressJa, &group.Address, &group.Mail, &group.Phone)
	if err != nil {
		log.Println("database scan error")
		return &groupResult{err: fmt.Errorf("(%s)error: query\n", time.Now())}
	}
	return &groupResult{group: group, err: nil}
}

func GetGroupMailFromDB(mail string) *groupResult {
	db := connectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return &groupResult{err: fmt.Errorf("(%s)error: database connection\n", time.Now())}
	}
	defer db.Close()

	rows := db.QueryRow("SELECT * FROM group WHERE email = ?", mail)
	var group auth.Group
	err := rows.Scan(&group.ID, &group.CreatedAt, &group.UpdatedAt, &group.OrgJa, &group.Org, &group.Status,
		&group.TechID, &group.PostCode, &group.AddressJa, &group.Address, &group.Mail, &group.Phone)
	if err != nil {
		log.Println("database scan error")
		return &groupResult{err: fmt.Errorf("(%s)error: query\n", time.Now())}
	}
	return &groupResult{group: group, err: nil}
}

func GetGroupIDFromDB(id int) *groupResult {
	db := connectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return &groupResult{err: fmt.Errorf("(%s)error: database connection\n", time.Now())}
	}
	defer db.Close()

	rows := db.QueryRow("SELECT * FROM group WHERE id = ?", id)
	var group auth.Group
	err := rows.Scan(&group.ID, &group.CreatedAt, &group.UpdatedAt, &group.OrgJa, &group.Org, &group.Status,
		&group.TechID, &group.PostCode, &group.AddressJa, &group.Address, &group.Mail, &group.Phone)
	if err != nil {
		log.Println("database scan error")
		return &groupResult{err: fmt.Errorf("(%s)error: query\n", time.Now())}
	}
	return &groupResult{group: group, err: nil}
}

func GetAllGroupFromDB() *allGroupResult {
	db := connectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return &allGroupResult{err: fmt.Errorf("(%s)error: database connection\n", time.Now())}
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM group")
	if err != nil {
		log.Println("database query error")
		return &allGroupResult{err: fmt.Errorf("(%s)error: query\n", time.Now())}
	}
	defer rows.Close()

	var allGroup *[]auth.Group
	for rows.Next() {
		var group auth.Group
		err := rows.Scan(&group.ID, &group.CreatedAt, &group.UpdatedAt, &group.OrgJa, &group.Org, &group.Status,
			&group.TechID, &group.PostCode, &group.AddressJa, &group.Address, &group.Mail, &group.Phone)
		if err != nil {
			log.Println("database scan error")
			return &allGroupResult{err: fmt.Errorf("(%s)error: query\n", time.Now())}
		}
		*allGroup = append(*allGroup, group)
	}
	return &allGroupResult{group: *allGroup, err: nil}
}
