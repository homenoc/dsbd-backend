package v0

import (
	"fmt"
	"git.bgp.ne.jp/dsbd/backend/pkg/api/core/user"
	"git.bgp.ne.jp/dsbd/backend/pkg/store"
	"log"
	"strconv"
	"time"
)

func Create(u *user.User) user.Result {
	db := store.ConnectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return user.Result{
			Status: fmt.Sprintf("(%s)error: database connection", time.Now()),
		}
	}
	defer db.Close()

	writeTable, err := db.Prepare(`INSERT INTO "user" ("created_at","gid","name","email","pass","level","status","is_verify","mail_token") VALUES (?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		log.Println("write error |error: ", err)
		return user.Result{
			Status: fmt.Sprintf("(%s)error: write error\n %s", time.Now(), err),
		}
	}
	if _, err := writeTable.Exec(time.Now().Unix(), u.GID, u.Name, u.Email, u.Pass, 1, u.Status, u.MailVerify, u.MailToken); err != nil {
		log.Println("apply error |error: ", err)
		return user.Result{
			Status: fmt.Sprintf("(%s)error: apply error\n %s", time.Now(), err),
		}
	}
	return user.Result{
		Status: "OK",
	}
}

func Delete(u *user.User) user.Result {
	db := store.ConnectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return user.Result{
			Status: fmt.Sprintf("(%s)error: database connection\n", time.Now()),
		}
	}
	defer db.Close()

	if _, err := db.Exec("DELETE FROM user WHERE name = ?", u.ID); err != nil {
		log.Println("database delete table error |", err)
		return user.Result{
			Status: fmt.Sprintf("(%s)error: delete error\n %s", time.Now(), err),
		}
	}
	return user.Result{
		Status: "OK",
	}
}

func Update(u *user.User) user.Result {
	db := store.ConnectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return user.Result{
			Status: fmt.Sprintf("(%s)error: database connection\n", time.Now()),
		}
	}
	defer db.Close()

	if _, err := db.Exec("UPDATE user SET updated_at = ?,name = ?,email = ?,pass = ?,level = ?,status = ?,is_verify = ? WHERE id = ?",
		time.Now().Unix(), u.Name, u.Email, u.Pass, u.Level, u.Status, u.MailVerify, u.ID); err != nil {
		log.Println("database update table error |", err)
		return user.Result{
			Status: fmt.Sprintf("(%s)error: database connection\n", time.Now()),
		}
	}
	return user.Result{
		Status: "OK",
	}
}

// 0:ID ,1:GID, 2:Mail, 3:Token, 4:MailVerify
func Get(base int, data *user.User) user.Result {
	db := store.ConnectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return user.Result{
			Status: fmt.Sprintf("(%s)error: database connection\n", time.Now()),
		}
	}
	defer db.Close()

	var database, baseData string

	if base == 0 { //ID
		database = "SELECT * FROM data WHERE id = ?"
		baseData = strconv.Itoa(data.ID)
	} else if base == 1 { //GID
		database = "SELECT * FROM data WHERE gid = ?"
		baseData = strconv.Itoa(data.GID)
	} else if base == 2 { //Mail
		database = "SELECT * FROM data WHERE email = ?"
		baseData = data.Email
	} else if base == 3 { //Token
		database = "SELECT * FROM data WHERE mail_token = ?"
		baseData = data.MailToken
		//} else if base == 4 { //MailVerify
		//	database = "SELECT * FROM data WHERE is_verify = ?"
		//	baseData =
	} else {
		log.Println("base select error")
		return user.Result{
			Status: fmt.Sprintf("(%s)error: base select\n", time.Now()),
		}
	}

	rows := db.QueryRow(database, baseData)
	var u user.User
	err := rows.Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt, &u.GID, &u.Name, &u.Email, &u.Pass, &u.Level, &u.Status, &u.MailVerify, &u.MailToken)
	if err != nil {
		log.Println("database scan error")
		return user.Result{
			Status: fmt.Sprintf("(%s)error: database scan\n", time.Now()),
		}
	}
	return user.Result{
		Status:   "OK",
		UserData: []user.User{u},
	}
}

func GetAll() user.Result {
	db := store.ConnectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return user.Result{
			Status: fmt.Sprintf("(%s)error: database connection\n", time.Now()),
		}
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM user")
	if err != nil {
		log.Println("database query error")
		return user.Result{
			Status: fmt.Sprintf("(%s)error: database query\n", time.Now()),
		}
	}
	defer rows.Close()

	var allUser []user.User
	for rows.Next() {
		var u user.User
		err := rows.Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt, &u.GID, &u.Name, &u.Email, &u.Pass, &u.Level, &u.Status, &u.MailVerify, &u.MailToken)
		if err != nil {
			log.Println("database scan error")
			return user.Result{
				Status: fmt.Sprintf("(%s)error: query\n", time.Now()),
			}
		}
		allUser = append(allUser, u)
	}
	return user.Result{
		Status:   "OK",
		UserData: allUser,
	}
}
