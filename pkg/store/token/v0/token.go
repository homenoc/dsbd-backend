package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	"github.com/homenoc/dsbd-backend/pkg/store"
	"log"
	"time"
)

func Create(t *token.Token) token.Result {
	db := store.ConnectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return token.Result{
			Status: false,
			Error:  fmt.Sprintf("(%s)error: database connection", time.Now()),
		}
	}
	defer db.Close()

	writeTable, err := db.Prepare(`INSERT INTO "token" ("created_at","update_at","expired_at","delete_at","uid",
"status","user_token","tmp_token","access_token","debug") VALUES (?,?,?,?,?,?,?,?,?,?)`)
	if err != nil {
		log.Println("write error |error: ", err)
		return token.Result{
			Status: false,
			Error:  fmt.Sprintf("(%s)error: write error\n %s", time.Now(), err),
		}
	}
	if _, err := writeTable.Exec(time.Now().Unix(), time.Now().Unix(), time.Now().Unix()+int64(t.ExpiredAt),
		time.Now().Unix()+int64(t.ExpiredAt)+int64(t.DeletedAt), t.UID, t.Status, t.UserToken, t.TmpToken,
		t.AccessToken, t.Debug); err != nil {
		log.Println("apply error |error: ", err)
		return token.Result{
			Status: false,
			Error:  fmt.Sprintf("(%s)error: apply error\n %s", time.Now(), err),
		}
	}
	return token.Result{
		Status: true,
	}
}

func Delete(t *token.Token) token.Result {
	db := store.ConnectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return token.Result{
			Status: false,
			Error:  fmt.Sprintf("(%s)error: database connection\n", time.Now()),
		}
	}
	defer db.Close()

	if _, err := db.Exec("DELETE FROM token WHERE id = ?", t.ID); err != nil {
		log.Println("database delete table error |", err)
		return token.Result{
			Status: false,
			Error:  fmt.Sprintf("(%s)error: delete error\n %s", time.Now(), err),
		}
	}
	return token.Result{
		Status: true,
	}
}

// value of base can reference from api/core/user/interface.go
func Get(userToken, accessToken string) token.Result {
	db := store.ConnectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return token.Result{
			Status: false,
			Error:  fmt.Sprintf("(%s)error: database connection\n", time.Now()),
		}
	}
	defer db.Close()

	rows := db.QueryRow("SELECT * FROM token WHERE user_token = ? AND access_token = ?", userToken, accessToken)
	var t token.Token
	err := rows.Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt, &t.ExpiredAt, &t.DeletedAt, &t.UID, &t.Status, &t.UserToken,
		&t.TmpToken, &t.AccessToken, &t.Debug)
	if err != nil {
		log.Println("database scan error")
		return token.Result{
			Status: false,
			Error:  fmt.Sprintf("(%s)error: database scan\n", time.Now()),
		}
	}
	return token.Result{
		Status: true,
		Token:  []token.Token{t},
	}
}

func GetAll() token.Result {
	db := store.ConnectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return token.Result{
			Status: false,
			Error:  fmt.Sprintf("(%s)error: database connection\n", time.Now()),
		}
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM user")
	if err != nil {
		log.Println("database query error")
		return token.Result{
			Status: false,
			Error:  fmt.Sprintf("(%s)error: database query\n", time.Now()),
		}
	}
	defer rows.Close()

	var allToken []token.Token
	for rows.Next() {
		var t token.Token
		err := rows.Scan(&t.ID, &t.CreatedAt, &t.UpdatedAt, &t.ExpiredAt, &t.DeletedAt, &t.UID, &t.Status, &t.UserToken,
			&t.TmpToken, &t.AccessToken, &t.Debug)
		if err != nil {
			log.Println("database scan error")
			return token.Result{
				Status: false,
				Error:  fmt.Sprintf("(%s)error: query\n", time.Now()),
			}
		}
		allToken = append(allToken, t)
	}
	return token.Result{
		Status: true,
		Token:  allToken,
	}
}
