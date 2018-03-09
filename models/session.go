package models

import "errors"

// Session represents a single user session of an user who has authenticated
// on Misirlou.
type Session struct {
	ID          string `json:"-" gorm:"type:char(64);primary_key"`
	UserID      int    `json:"user_id"`
	AccessToken string `json:"access_token"`
	FCMToken    string `json:"-"`
}

// Session retrieves a session knowing its hashed ID
func (db *DB) Session(id string) (*Session, error) {
	var s Session
	res := db.db.First(&s, "id = ?", id)
	if res.Error != nil {
		return nil, ignoreNotFound(res)
	}
	return &s, nil
}

// SetSession creates or updates a session.
func (db *DB) SetSession(sess *Session) error {
	if sess == nil {
		return errors.New("sess is nil")
	}
	return db.db.Save(sess).Error
}
