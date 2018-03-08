package models

import "time"

// Team represents a team playing in Misirlou.
type Team struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Tournament int       `json:"tournament"`
	Captain    int       `json:"captain"`
	CreatedAt  time.Time `json:"created_at"`
}

// Teams returns at most 50 teams, given the where conditions and the page.
func (db *DB) Teams(t *Team, page int) ([]Team, error) {
	if page < 0 {
		page = 0
	}
	teams := make([]Team, 0, 50)
	res := db.db.Where(t).Offset(page * 50).Limit(50).Find(&teams)
	return teams, res.Error
}
