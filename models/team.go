package models

// Team represents a team playing in Misirlou.
type Team struct {
	ID         ID     `json:"id"`
	Name       string `json:"name"`
	Tournament int    `json:"tournament"`
	Captain    int    `json:"captain"`
}

// TODO: we should check that the team's tournament status is not 0.

// TeamFilters are options that can be passed to Teams for filtering teams.
type TeamFilters struct {
	Tournament      ID
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
func (db *DB) Team(id ID) (*Team, error) {
	var t Team
	res := db.db.First(&t, "id = ?", id)
	if res.Error != nil {
		return nil, ignoreNotFound(res)
	}
	return &t, nil
}

// TeamMember represents a member of a team and information about their
// relationship to the team.
type TeamMember struct {
	Team       int `json:"team"`
	User       int `json:"user"`
	Attributes int `json:"attributes"`
}

// TableName returns the correct table name so that it can correctly be used
// by gorm.
func (TeamMember) TableName() string {
	return "team_users"
}

// TeamMembers retrieves all the members of a team.
func (db *DB) TeamMembers(teamID ID, page int) ([]TeamMember, error) {
	members := make([]TeamMember, 0, 50)
	err := db.db.Offset(positivePage(page)*50).Limit(50).
		Find(&members, "team = ?", teamID).Error
	return members, err
}
