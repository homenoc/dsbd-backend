package v0

import (
	"fmt"
	"github.com/homenoc/dsbd-backend/pkg/api/core/user"
	"github.com/homenoc/dsbd-backend/pkg/store"
	"log"
	"time"
)

func Create(t *user.InitToken) user.InitTokenResult {
	db := store.ConnectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return user.InitTokenResult{
			Status: false,
			Error:  fmt.Sprintf("(%s)error: database connection", time.Now()),
		}
	}
	defer db.Close()

	writeTable, err := db.Prepare(`INSERT INTO "init_token" ("created_at","expired_at","delete_at","ip","token1","token2") VALUES (?,?,?,?,?,?)`)
	if err != nil {
		log.Println("write error |error: ", err)
		return user.InitTokenResult{
			Status: false,
			Error:  fmt.Sprintf("(%s)error: write error\n %s", time.Now(), err),
		}
	}
	if _, err := writeTable.Exec(time.Now().Unix(), time.Now().Unix()+10000, time.Now().Unix()+20000, t.IP, t.Token1, t.Token2); err != nil {
		log.Println("apply error |error: ", err)
		return user.InitTokenResult{
			Status: false,
			Error:  fmt.Sprintf("(%s)error: apply error\n %s", time.Now(), err),
		}
	}
	return user.InitTokenResult{
		Status: true,
	}
}

func Delete(t *user.InitToken) user.InitTokenResult {
	db := store.ConnectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return user.InitTokenResult{
			Status: false,
			Error:  fmt.Sprintf("(%s)error: database connection\n", time.Now()),
		}
	}
	defer db.Close()

	if _, err := db.Exec("DELETE FROM user WHERE name = ?", t.ID); err != nil {
		log.Println("database delete table error |", err)
		return user.InitTokenResult{
			Status: false,
			Error:  fmt.Sprintf("(%s)error: delete error\n %s", time.Now(), err),
		}
	}
	return user.InitTokenResult{
		Status: true,
	}
}

// value of base can reference from api/core/user/interface.go
func Get(token string) user.InitTokenResult {
	db := store.ConnectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return user.InitTokenResult{
			Status: false,
			Error:  fmt.Sprintf("(%s)error: database connection\n", time.Now()),
		}
	}
	defer db.Close()

	rows := db.QueryRow("SELECT * FROM init_token WHERE token1 = ?", token)
	var t user.InitToken
	err := rows.Scan(&t.ID, &t.CreatedAt, &t.ExpiredAt, &t.DeletedAt, &t.IP, &t.Token1, &t.Token2)
	if err != nil {
		log.Println("database scan error")
		return user.InitTokenResult{
			Status: false,
			Error:  fmt.Sprintf("(%s)error: database scan\n", time.Now()),
		}
	}
	return user.InitTokenResult{
		Status:    true,
		TokenData: []user.InitToken{t},
	}
}

func GetAll() user.InitTokenResult {
	db := store.ConnectDB()
	//error check
	if db == nil {
		log.Println("database connection error")
		return user.InitTokenResult{
			Status: false,
			Error:  fmt.Sprintf("(%s)error: database connection\n", time.Now()),
		}
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM user")
	if err != nil {
		log.Println("database query error")
		return user.InitTokenResult{
			Status: false,
			Error:  fmt.Sprintf("(%s)error: database query\n", time.Now()),
		}
	}
	defer rows.Close()

	var allInitToken []user.InitToken
	for rows.Next() {
		var t user.InitToken
		err := rows.Scan(&t.ID, &t.CreatedAt, &t.ExpiredAt, &t.DeletedAt, &t.IP, &t.Token1, &t.Token2)
		if err != nil {
			log.Println("database scan error")
			return user.InitTokenResult{
				Status: false,
				Error:  fmt.Sprintf("(%s)error: query\n", time.Now()),
			}
		}
		allInitToken = append(allInitToken, t)
	}
	return user.InitTokenResult{
		Status:    true,
		TokenData: allInitToken,
	}
}
