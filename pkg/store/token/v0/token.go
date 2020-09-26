package v0

import (
	"database/sql"
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	"github.com/homenoc/dsbd-backend/pkg/store"
	"log"
	"time"
)

func Create(t *token.Token) error {
	db := store.ConnectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: database connection", time.Now())
	}
	defer db.Close()

	writeTable, err := db.Prepare(`INSERT INTO "token" ("created_at","update_at","expired_at","delete_at","uid",
"status","user_token","tmp_token","access_token","debug") VALUES (?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		log.Println("write error |error: ", err)
		return fmt.Errorf("(%s)error: write error", time.Now())
	}
	if _, err := writeTable.Exec(time.Now().Unix(), time.Now().Unix(), time.Now().Unix()+int64(t.ExpiredAt),
		time.Now().Unix()+int64(t.ExpiredAt)+int64(t.DeletedAt), t.UID, t.Status, t.UserToken, t.TmpToken,
		t.AccessToken, t.Debug); err != nil {
		log.Println("apply error |error: ", err)
		return fmt.Errorf("(%s)error: apply error", time.Now())
	}
	return nil
}

func Delete(t *token.Token) error {
	db := store.ConnectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: database connection", time.Now())
	}
	defer db.Close()

	if _, err := db.Exec("DELETE FROM token WHERE id = ?", t.ID); err != nil {
		log.Println("database delete table error |", err)
		return fmt.Errorf("(%s)error: delete error", time.Now())
	}
	return nil
}

func Update(base int, t *token.Token) error {
	db := store.ConnectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return fmt.Errorf("(%s)error: database connection", time.Now())
	}
	defer db.Close()

	var err error

	if token.AddToken == base {
		_, err = db.Exec("UPDATE token SET updated_at = ?,expired_at = ?,delete_at = ?,uid = ?,status = ?,access_token = ? WHERE id = ?",
			time.Now().Unix(), time.Now().Unix()+int64(t.ExpiredAt), time.Now().Unix()+int64(t.ExpiredAt)+int64(t.DeletedAt),
			t.UID, t.Status, t.AccessToken, t.ID)
	} else if token.UpdateToken == base {
		_, err = db.Exec("UPDATE token SET updated_at = ?,expired_at = ?,delete_at = ? WHERE id = ?",
			time.Now().Unix(), time.Now().Unix()+int64(t.ExpiredAt), time.Now().Unix()+int64(t.ExpiredAt)+int64(t.DeletedAt), t.ID)
	} else {
		log.Println("base select error")
		return fmt.Errorf("(%s)error: base select\n %s", time.Now(), err)
	}
	if err != nil {
		log.Println("database update table error |", err)
		return fmt.Errorf("(%s)error: delete error\n %s", time.Now(), err)
	}
	return nil
}

// value of base can reference from api/core/user/interface.go
func Get(base int, input *token.Token) (token.Token, error) {
	db := store.ConnectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return token.Token{}, fmt.Errorf("(%s)error: database connection\n", time.Now())
	}
	defer db.Close()

	var rows *sql.Row

	if base == token.UserToken {
		rows = db.QueryRow("SELECT * FROM token WHERE user_token = ?", &input.UserToken)
	} else if base == token.UserTokenAndAccessToken {
		rows = db.QueryRow("SELECT * FROM token WHERE user_token = ? AND access_token = ?", &input.UserToken, &input.AccessToken)
	} else {
		log.Println("base select error")
		return token.Token{}, fmt.Errorf("(%s)error: base select\n", time.Now())
	}

	var t token.Token
	err := rows.Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt, &t.ExpiredAt, &t.DeletedAt, &t.UID, &t.Status, &t.UserToken,
		&t.TmpToken, &t.AccessToken, &t.Debug)
	if err != nil {
		log.Println("database scan error")
		return token.Token{}, fmt.Errorf("(%s)error: database scan\n", time.Now())
	}
	return t, nil

}

func GetAll() ([]token.Token, error) {
	db := store.ConnectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return []token.Token{}, fmt.Errorf("(%s)error: database connection\n", time.Now())
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM token")
	if err != nil {
		log.Println("database query error")
		return []token.Token{}, fmt.Errorf("(%s)error: database query\n", time.Now())

	}
	defer rows.Close()

	var allToken []token.Token
	for rows.Next() {
		var t token.Token
		err := rows.Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt, &t.ExpiredAt, &t.DeletedAt, &t.UID, &t.Status, &t.UserToken,
			&t.TmpToken, &t.AccessToken, &t.Debug)
		if err != nil {
			log.Println("database scan error")
			return []token.Token{}, fmt.Errorf("(%s)error: database scan\n", time.Now())
		}
		allToken = append(allToken, t)
	}
	return allToken, nil
}
