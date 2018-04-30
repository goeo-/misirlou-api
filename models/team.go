package models

// Team represents a team playing in Misirlou.
type Team struct {
	ID         ID     `json:"id"`
	Name       string `json:"name"`
	Tournament ID     `json:"tournament"`
	Captain    int    `json:"captain"`
}

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

// CreateTeam creates a new team, which is to say, registers in a tournament.
func (db *DB) CreateTeam(t *Team) error {
	return db.db.Create(t).Error
}

// TeamAttributes represents the attributes a user may have inside of a team,
// such as being invited, a member, or a captain.
type TeamAttributes int

// Various TeamAttributes a team member might have.
const (
	TeamAttributeInvited TeamAttributes = iota
	TeamAttributeMember
	TeamAttributeCaptain
)

// TeamMember represents a member of a team and information about their
// relationship to the team.
type TeamMember struct {
	Team       ID             `json:"team"`
	User       int            `json:"user"`
	Attributes TeamAttributes `json:"attributes"`
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

// AddTeamMembers adds the given team members to the database.
func (db *DB) AddTeamMembers(ms []TeamMember) error {
	// not using for-range because this way we can keep the reference
	// straight into the slice
	for i := 0; i < len(ms); i++ {
		if err := db.db.Create(&ms[i]).Error; err != nil {
			return err
		}
	}
	return nil
}
