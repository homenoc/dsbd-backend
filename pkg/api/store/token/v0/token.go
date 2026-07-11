package v0

import (
	"time"

	"github.com/homenoc/dsbd-backend/pkg/api/core"
	"github.com/homenoc/dsbd-backend/pkg/api/core/token"
	"github.com/homenoc/dsbd-backend/pkg/api/store"
	"gorm.io/gorm"
)

// Create inserts a token. tx may be nil to use the shared pool.
func Create(t *core.Token) error {
	return store.DB().Create(t).Error
}

func Delete(t *core.Token) error {
	return store.DB().Delete(t).Error
}

func DeleteAll() error {
	return store.DB().Exec("DELETE FROM tokens").Error
}

// UpdateSession sets the fields written on login (was Update(AddToken, ...)).
func UpdateSession(t *core.Token) error {
	return store.DB().Model(&core.Token{Model: gorm.Model{ID: t.ID}}).Updates(core.Token{
		ExpiredAt: t.ExpiredAt, UserID: t.UserID, Status: t.Status, AccessToken: t.AccessToken}).Error
}

// Renew extends only the expiry (was Update(UpdateToken, ...)).
func Renew(t *core.Token) error {
	return store.DB().Model(&core.Token{Model: gorm.Model{ID: t.ID}}).
		Updates(core.Token{ExpiredAt: t.ExpiredAt}).Error
}

// UpdateAll updates all mutable token fields (was Update(UpdateAll, ...)).
func UpdateAll(t *core.Token) error {
	return store.DB().Model(&core.Token{Model: gorm.Model{ID: t.ID}}).Updates(core.Token{
		ExpiredAt:   t.ExpiredAt,
		UserID:      t.UserID,
		Status:      t.Status,
		UserToken:   t.UserToken,
		TmpToken:    t.TmpToken,
		AccessToken: t.AccessToken,
		Debug:       t.Debug,
	}).Error
}

// GetByID looks up a single token by primary key.
func GetByID(id uint) ([]core.Token, error) {
	var tokens []core.Token
	err := store.DB().Where(&core.Token{Model: gorm.Model{ID: id}}).Find(&tokens).Error
	return tokens, err
}

// GetByUserToken returns valid, non-admin tokens matching a user token.
func GetByUserToken(userToken string) ([]core.Token, error) {
	var tokens []core.Token
	err := store.DB().Where("user_token = ? AND admin = ? AND expired_at > ?",
		userToken, false, time.Now()).Find(&tokens).Error
	return tokens, err
}

// GetSession resolves a user session (user token + access token) and preloads
// the owning user and group.
func GetSession(userToken, accessToken string) ([]core.Token, error) {
	var tokens []core.Token
	err := store.DB().Where("user_token = ? AND access_token = ? AND admin = ? AND expired_at > ?",
		userToken, accessToken, false, time.Now()).
		Preload("User").
		Preload("User.Group").
		Find(&tokens).Error
	return tokens, err
}

// GetAdminBotToken returns admin tokens matching an access token (no expiry filter).
func GetAdminBotToken(accessToken string) ([]core.Token, error) {
	var tokens []core.Token
	err := store.DB().Where("access_token = ? AND admin = ?", accessToken, true).Find(&tokens).Error
	return tokens, err
}

// GetValidAdminToken returns valid (non-expired) admin tokens matching an access token.
func GetValidAdminToken(accessToken string) ([]core.Token, error) {
	var tokens []core.Token
	err := store.DB().Where("access_token = ? AND admin = ? AND expired_at > ?",
		accessToken, true, time.Now()).Find(&tokens).Error
	return tokens, err
}

// GetExpired returns tokens past their expiry (for the cleanup goroutine).
func GetExpired() ([]core.Token, error) {
	var tokens []core.Token
	err := store.DB().Where("expired_at < ?", time.Now()).Find(&tokens).Error
	return tokens, err
}

func GetAll() token.ResultDatabase {
	var tokens []core.Token
	err := store.DB().Find(&tokens).Error
	return token.ResultDatabase{Token: tokens, Err: err}
}
