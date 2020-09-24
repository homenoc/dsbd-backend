package v0

import (
	"database/sql"
	"fmt"
	group "github.com/homenoc/dsbd-backend/pkg/api/core/group"
	"github.com/homenoc/dsbd-backend/pkg/store"
	"log"
	"time"
)

func Create(group group.Group) error {
	db := store.ConnectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: database connection\n", time.Now())
	}
	defer db.Close()
	writeTable, err := db.Prepare(`INSERT INTO "group" ("created_at","updated_at","org_ja","org","status","tech_id","postcode",
"address_ja","address","mail","phone") VALUES (?,?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		log.Println("write error |error: ", err)
		return fmt.Errorf("(%s)error: write error\n %s", time.Now(), err)
	}
	if _, err := writeTable.Exec(time.Now().Unix(), time.Now().Unix(), group.OrgJa, group.Org, group.Status, group.TechID, group.PostCode, group.AddressJa, group.Address, group.Mail, group.Phone); err != nil {
		log.Println("apply error |error: ", err)
		return fmt.Errorf("(%s)error: apply error\n %s", time.Now(), err)
	}
	return nil
}

func Delete(id int) error {
	db := store.ConnectDB()
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

func Update(base int, g group.Group) group.Result {
	db := store.ConnectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return group.Result{
			Status: false,
			Error:  fmt.Sprintf("(%s)error: database connection\n", time.Now()),
		}
	}
	defer db.Close()

	var err error

	if group.UpdateOrg == base {
		_, err = db.Exec("UPDATE group SET updated_at = ?,org_ja = ?,org = ? WHERE id = ?", time.Now().Unix(),
			g.OrgJa, g.Org, g.ID)
	} else if group.UpdateStatus == base {
		_, err = db.Exec("UPDATE group SET updated_at = ?,status = ? WHERE id = ?", time.Now().Unix(), g.Status, g.ID)
	} else if group.UpdateTechID == base {
		_, err = db.Exec("UPDATE group SET updated_at = ?,tech_id = ? WHERE id = ?", time.Now().Unix(), g.TechID, g.ID)
	} else if group.UpdateInfo == base {
		_, err = db.Exec("UPDATE group SET updated_at = ?,postcode = ?,address_ja = ?,address = ?,mail = ?,phone = ? WHERE id = ?",
			time.Now().Unix(), g.PostCode, g.AddressJa, g.Address, g.Mail, g.Phone, g.ID)
	} else {
		log.Println("base select error")
		return group.Result{
			Status: false,
			Error:  fmt.Sprintf("(%s)error: base select\n", time.Now()),
		}
	}
	if err != nil {
		log.Println("database update table error |", err)
		return group.Result{
			Status: false,
			Error:  fmt.Sprintf("(%s)error: delete error\n %s", time.Now(), err),
		}
	}
	return group.Result{
		Status: true,
	}
}

func Get(base int, data *group.Group) group.Result {
	db := store.ConnectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return group.Result{
			Status: false,
			Error:  fmt.Sprintf("(%s)error: database connection\n", time.Now()),
		}
	}
	defer db.Close()

	var rows *sql.Row

	if base == group.ID { //ID
		rows = db.QueryRow("SELECT * FROM data WHERE id = ?", &data.ID)
	} else if base == group.OrgJa { //OrgJa
		rows = db.QueryRow("SELECT * FROM data WHERE org_ja = ?", &data.OrgJa)
	} else if base == group.Org { //Org
		rows = db.QueryRow("SELECT * FROM data WHERE org = ?", &data.Org)
	} else if base == group.Email { //Mail
		rows = db.QueryRow("SELECT * FROM data WHERE mail = ?", &data.Mail)
	} else {
		log.Println("base select error")
		return group.Result{
			Status: false,
			Error:  fmt.Sprintf("(%s)error: base select\n", time.Now()),
		}
	}

	var g group.Group
	err := rows.Scan(&g.ID, &g.CreatedAt, &g.UpdatedAt, &g.OrgJa, &g.Org, &g.Status, &g.TechID, &g.PostCode,
		&g.AddressJa, &g.Address, &g.Mail, &g.Phone)
	if err != nil {
		log.Println("database scan error")
		return group.Result{
			Status: false,
			Error:  fmt.Sprintf("(%s)error: database scan\n", time.Now()),
		}
	}
	return group.Result{
		Status:    true,
		GroupData: []group.Group{g},
	}

}

func GetAll() *group.Result {
	db := store.ConnectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return &group.Result{Error: fmt.Sprintf("(%s)error: database connection\n", time.Now())}
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM group")
	if err != nil {
		log.Println("database query error")
		return &group.Result{Error: fmt.Sprintf("(%s)error: query\n", time.Now())}
	}
	defer rows.Close()

	var allGroup *[]group.Group
	for rows.Next() {
		var data group.Group
		err := rows.Scan(&data.ID, &data.CreatedAt, &data.UpdatedAt, &data.OrgJa, &data.Org, &data.Status,
			&data.TechID, &data.PostCode, &data.AddressJa, &data.Address, &data.Mail, &data.Phone)
		if err != nil {
			log.Println("database scan error")
			return &group.Result{Error: fmt.Sprintf("(%s)error: query\n", time.Now())}
		}
		*allGroup = append(*allGroup, data)
	}
	return &group.Result{GroupData: *allGroup, Error: nil}
}
