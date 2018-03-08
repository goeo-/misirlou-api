package models

import (
	"database/sql/driver"
	"errors"
	"time"
)

type rawMessageSQL string

var (
	errNotAString = errors.New("models: value is not a string")
)

func (m *rawMessageSQL) Scan(value interface{}) error {
	var s string
	switch v := value.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return errNotAString
	}
	if s == "" {
		*m = "null"
		return nil
	}
	*m = rawMessageSQL(s)
	return nil
}

func (m rawMessageSQL) Value() (driver.Value, error) {
	return string(m), nil
}

func (m *rawMessageSQL) UnmarshalJSON(data []byte) error {
	*m = rawMessageSQL(data)
	return nil
}

func (m rawMessageSQL) MarshalJSON() ([]byte, error) {
	return []byte(m), nil
}

// Tournament represents a tournament managed by Misirlou.
type Tournament struct {
	ID                int           `json:"id"`
	Name              string        `json:"name"`
	Description       string        `json:"description"`
	Mode              int           `json:"mode"`
	Status            int           `json:"status"`
	StatusData        rawMessageSQL `json:"status_data"`
	TeamSize          int           `json:"team_size"`
	MinTeamSize       int           `json:"min_team_size"`
	ExclusivityStarts time.Time     `json:"exclusivity_starts"`
	ExclusivityEnds   time.Time     `json:"exclusivity_ends"`
	UpdatedAt         time.Time     `json:"updated_at"`
	CreatedAt         time.Time     `json:"created_at"`
}

// Tournaments returns the tournaments sorted by their ID. If the id is given,
// at most one tournament will be returned, of the given ID. Tournaments with
// status = 0 will not be shown.
func (db *DB) Tournaments(id int, page int) ([]Tournament, error) {
	q := db.db.Where("status != ?", id)
	if id != 0 {
		q = q.Where("id = ?", id)
	}
	tourns := make([]Tournament, 0, 50)
	err := q.Order("id desc").Offset(positivePage(page)).Limit(50).
		Find(&tourns).Error
	return tourns, err
}

// TournamentRules represents a collection rules set out for a given tournament,
// which is represented by the ID field in the struct.
type TournamentRules struct {
	ID    int    `json:"id"`
	Rules string `json:"rules"`
}

// TournamentRules returns the tournament rules for the given tournament.
func (db *DB) TournamentRules(id int) (*TournamentRules, error) {
	var status []int
	err := db.db.Table("tournaments").Where("id = ?", id).
		Pluck("status", &status).Error
	if err != nil {
		return nil, err
	}
	if len(status) == 0 || status[0] == 0 {
		return nil, nil
	}
	var rules TournamentRules
	res := db.db.First(&rules)
	if res.RecordNotFound() {
		return nil, nil
	}
	if res.Error != nil {
		return nil, res.Error
	}
	return &rules, nil
}

// TournamentStaff represents what privileges, if any, a user has in a
// tournament.
type TournamentStaff struct {
	ID         int `json:"id"`
	Tournament int `json:"tournament"`
	Privileges int `json:"privileges"`
}
