package store

import (
	"fmt"
	"git.bgp.ne.jp/dsbd/backend/pkg/auth"
	"log"
	"time"
)

func CreateUserToDB(user auth.User) error {
	db := connectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: database connection\n", time.Now())
	}
	defer db.Close()

	writeTable, err := db.Prepare(`INSERT INTO "user" ("created_at","gid","name","email","pass","level","status","is_verify") VALUES (?,?,?,?,?,?,?,?)`)
	if err != nil {
		log.Println("write error |error: ", err)
		return fmt.Errorf("(%s)error: write error\n %s", time.Now(), err)
	}
	if _, err := writeTable.Exec(time.Now().Unix(), user.GID, user.Name, user.Mail, user.Pass, user.Level, user.Status, user.IsVerify); err != nil {
		log.Println("apply error |error: ", err)
		return fmt.Errorf("(%s)error: apply error\n %s", time.Now(), err)
	}
	return nil
}

func DeleteUserToDB(id int) error {
	db := connectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: database connection\n", time.Now())
	}
	defer db.Close()

	if _, err := db.Exec("DELETE FROM user WHERE name = ?", id); err != nil {
		log.Println("database delete table error |", err)
		return fmt.Errorf("(%s)error: delete table\n", time.Now())
	}
	return nil
}

func UpdateUserToDB(user auth.User) error {
	db := connectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: database connection\n", time.Now())
	}
	defer db.Close()

	if _, err := db.Exec("UPDATE user SET updated_at = ?,name = ?,email = ?,pass = ?,level = ?,status = ?,is_verify = ? WHERE id = ?",
		time.Now().Unix(), user.Name, user.Mail, user.Pass, user.Level, user.Status, user.IsVerify, user.ID); err != nil {
		log.Println("database update table error |", err)
		return fmt.Errorf("(%s)error: update table\n", time.Now())
	}
	return nil
}

func GetUserMailFromDB(mail string) *userResult {
	db := connectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return &userResult{err: fmt.Errorf("(%s)error: database connection\n", time.Now())}
	}
	defer db.Close()

	rows := db.QueryRow("SELECT * FROM user WHERE email = ?", mail)
	var user auth.User
	err := rows.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.GID, &user.Name, &user.Mail, &user.Pass, &user.Level, &user.Status, &user.IsVerify)
	if err != nil {
		log.Println("database scan error")
		return &userResult{err: fmt.Errorf("(%s)error: query\n", time.Now())}
	}
	return &userResult{user: user, err: nil}
}

func GetUserIDFromDB(id int) *userResult {
	db := connectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return &userResult{err: fmt.Errorf("(%s)error: database connection\n", time.Now())}
	}
	defer db.Close()

	rows := db.QueryRow("SELECT * FROM user WHERE id = ?", id)
	var user auth.User
	err := rows.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.GID, &user.Name, &user.Mail, &user.Pass, &user.Level, &user.Status, &user.IsVerify)
	if err != nil {
		log.Println("database scan error")
		return &userResult{err: fmt.Errorf("(%s)error: query\n", time.Now())}
	}
	return &userResult{user: user, err: nil}
}

func GetAllUserFromDB() *allUserResult {
	db := connectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return &allUserResult{err: fmt.Errorf("(%s)error: database connection\n", time.Now())}
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM user")
	if err != nil {
		log.Println("database query error")
		return &allUserResult{err: fmt.Errorf("(%s)error: query\n", time.Now())}
	}
	defer rows.Close()

	var allUser *[]auth.User
	for rows.Next() {
		var user auth.User
		err := rows.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.GID, &user.Name, &user.Mail, &user.Pass, &user.Level, &user.Status, &user.IsVerify)
		if err != nil {
			log.Println("database scan error")
			return &allUserResult{err: fmt.Errorf("(%s)error: query\n", time.Now())}
		}
		*allUser = append(*allUser, user)
	}
	return &allUserResult{user: *allUser, err: nil}
}
