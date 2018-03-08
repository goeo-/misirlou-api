package models

import (
	"time"
)

// Team represents a team playing in Misirlou.
type Team struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Tournament int       `json:"tournament"`
	Captain    int       `json:"captain"`
	CreatedAt  time.Time `json:"created_at"`
}

// TODO: we should check that the team's tournament status is not 0.

// TeamFilters are options that can be passed to Teams for filtering teams.
type TeamFilters struct {
	Tournament      int
	ForceTournament bool
	Member          int
}

// Teams returns at most 50 teams with the specified page.
func (db *DB) Teams(filters TeamFilters, page int) ([]Team, error) {
	d := db.db
	if filters.ForceTournament || filters.Tournament != 0 {
		d = d.Where("tournament = ?", filters.Tournament)
	}
	if filters.Member != 0 {
		d = d.Joins("INNER JOIN team_users ON team_users.team = teams.id AND team_users.user = ?", filters.Member)
	}
	teams := make([]Team, 0, 50)
	d = d.Offset(positivePage(page) * 50).Limit(50).Find(&teams)
	return teams, d.Error
}

// Team returns a single team knowing its ID.
func (db *DB) Team(id int) (*Team, error) {
	var t Team
	res := db.db.First(&t, id)
	if res.Error != nil {
		if res.RecordNotFound() {
			return nil, nil
		}
		return nil, res.Error
	}
	return &t, nil
}
