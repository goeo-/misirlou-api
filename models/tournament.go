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
	case nil:
		s = "null"
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
	if m == "" {
		return []byte("null"), nil
	}
	return []byte(m), nil
}

// TournamentStatus represents one of the statuses of the tournament, as
// described in the constants
type TournamentStatus int

// A tournament will always be in one of these states.
const (
	// StatusOrganising means the tournament is currently private as the owner
	// needs to create the rules, staff, etc.
	StatusOrganising TournamentStatus = iota
	// StatusOpen means the tournament currently accepts new teams.
	StatusOpen
	// StatusRegsClosed means that registrations for the tournament have been
	// closed, and we're waiting for a bracket.
	StatusRegsClosed
	// StatusAwaitRound means we're waiting for the next round of the tournament.
	StatusAwaitRound
	// StatusPlaying means the tournament game is currently being played (so, for
	// instance, we need to allow inputs of results of games from referees).
	StatusPlaying
	// StatusClosed means the tournament has been terminated.
	StatusClosed
)

// Tournament represents a tournament managed by Misirlou.
type Tournament struct {
	ID                ID               `json:"id"`
	Name              string           `json:"name"`
	Description       string           `json:"description"`
	Mode              int              `json:"mode"`
	Status            TournamentStatus `json:"status"`
	StatusData        rawMessageSQL    `json:"status_data"`
	TeamSize          int              `json:"team_size"`
	MinTeamSize       int              `json:"min_team_size"`
	ExclusivityStarts time.Time        `json:"exclusivity_starts"`
	ExclusivityEnds   time.Time        `json:"exclusivity_ends"`
	UpdatedAt         time.Time        `json:"updated_at"`
}

// Tournaments returns the tournaments sorted by their ID.
func (db *DB) Tournaments(page int) ([]Tournament, error) {
	tourns := make([]Tournament, 0, 50)
	err := db.db.Order("id desc").Offset(positivePage(page)).Limit(50).
		Find(&tourns, "status != 0").Error
	return tourns, err
}

// Tournament returns a single tournament knowing its ID.
func (db *DB) Tournament(id ID) (*Tournament, error) {
	var t Tournament
	res := db.db.Where("status != 0").First(&t, "id = ?", id)
	if res.Error != nil {
		return nil, ignoreNotFound(res)
	}
	return &t, nil
}

// UserInTournament checks whether an user is taking part in a tournament.
func (db *DB) UserInTournament(tournID ID, userID int) (bool, error) {
	var i []int
	err := db.db.Raw(`SELECT COUNT(*) as c FROM teams
	INNER JOIN team_users ON teams.id = team_users.team
	WHERE teams.tournament = ? AND team_users.user = ? AND team_users.attributes > 0
	LIMIT 1`, tournID, userID).Pluck("c", &i).Error
	return len(i) > 0 && i[0] > 0, err
}

// UserIsBusy checks whether the user is 'busy' with another tournament during
// the period in which they would play in the tournament.
func (db *DB) UserIsBusy(tourn *Tournament, userID int) (bool, error) {
	var i []int
	// We check if another tournament overlaps with ours.
	// In plain words: we check if the given user (team_users.user = ?)
	// takes part (team_users.attributes > 0) in another (tournaments.id != ?)
	// tournament that begins while our tournament is in session OR
	// our tournament begins while another tournament is in session.
	err := db.db.Raw(
		`SELECT COUNT(*) as c FROM team_users
	INNER JOIN teams ON teams.id = team_users.team
	INNER JOIN tournaments ON teams.tournament = tournaments.id
	WHERE
		team_users.user = ? AND team_users.attributes > 0 AND
		tournaments.id != ? AND
		((tournaments.exclusivity_starts >= ? AND tournaments.exclusivity_starts <= ?) OR
		 (tournaments.exclusivity_ends >= ? AND tournaments.exclusivity_starts <= ?))
	LIMIT 1`,
		userID, tourn.ID, tourn.ExclusivityStarts, tourn.ExclusivityEnds,
		tourn.ExclusivityStarts, tourn.ExclusivityStarts,
	).Pluck("c", &i).Error
	return len(i) > 0 && i[0] > 0, err
}

// TournamentRules represents a collection rules set out for a given tournament,
// which is represented by the ID field in the struct.
type TournamentRules struct {
	ID    ID     `json:"id"`
	Rules string `json:"rules"`
}

// TournamentRules returns the tournament rules for the given tournament.
func (db *DB) TournamentRules(id ID) (*TournamentRules, error) {
	// make sure the status of the tournament is not 0
	var status []int
	err := db.db.Table("tournaments").Where("id = ?", id).
		Pluck("status", &status).Error
	if err != nil {
		return nil, err
	}
	if len(status) == 0 || status[0] == 0 {
		return nil, nil
	}
	// fetch tourn rules
	var rules TournamentRules
	res := db.db.First(&rules)
	if res.Error != nil {
		return nil, ignoreNotFound(res)
	}
	return &rules, nil
}

// TournamentStaff represents what privileges, if any, a user has in a
// tournament.
type TournamentStaff struct {
	ID         int `json:"id"`
	Tournament ID  `json:"tournament"`
	Privileges int `json:"privileges"`
}
