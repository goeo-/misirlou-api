package api

import (
	"zxq.co/ripple/misirlou-api/http"
	"zxq.co/ripple/misirlou-api/models"
	"zxq.co/x/ripple"
)

// Teams fetches the teams that are playing in the given tournament ID.
func Teams(c *http.Context) {
	teams(c, models.TeamFilters{
		Tournament: c.QueryID("tournament"),
	})
}

// TeamsInTournament returns all the teams in a tournament, which basically is
// sugar for calling Teams giving a tournament.
func TeamsInTournament(c *http.Context) {
	teams(c, models.TeamFilters{
		Tournament:      c.ParamID("id"),
		ForceTournament: true,
	})
}

func teams(c *http.Context, filters models.TeamFilters) {
	filters.Member = c.QueryInt("member")
	teams, err := c.DB.Teams(filters, c.QueryInt("p"))
	if err != nil {
		c.Error(err)
		return
	}
	c.SetJSON(teams, false)
}

// Team retrieves a single team.
func Team(c *http.Context) {
	team, err := c.DB.Team(c.ParamID("id"))
	if err != nil {
		c.Error(err)
		return
	}
	c.SetJSON(team, team == nil)
}

// CreateTeamData is the JSON data that is passed to CreateTeam.
type CreateTeamData struct {
	Tournament models.ID `json:"tournament"`
	Members    []int     `json:"members"`
	Name       string    `json:"name"`
}

// CreateTeam creates a new team and verifies that the team can actually
// take part in the tournament.
func CreateTeam(c *http.Context) {
	// Get session data
	sess := c.Session()
	if sess == nil {
		c.SetCode(401)
		c.WriteString("Missing or invalid access token.")
		return
	}

	// Get client; fetch current user.
	cl := sess.RippleClient()
	user, err := validateUser(cl)
	if err != nil {
		c.Error(err)
		return
	}

	// Get the data submitted by the user.
	var d CreateTeamData
	err = c.JSON(&d)
	if err != nil {
		c.Error(err)
		return
	}
	tourn, err := userCanRegister(c.DB, d.Tournament, user.ID)
	if err != nil {
		c.Error(err)
		return
	}

	// Single-user tournament; we force the team name and members.
	if tourn.TeamSize == 1 {
		t := &models.Team{
			Name:       user.Username,
			Tournament: tourn.ID,
			Captain:    user.ID,
		}
		err = c.DB.CreateTeam(t)
		if err != nil {
			c.Error(err)
			return
		}
		err = c.DB.AddTeamMembers([]models.TeamMember{{
			Team:       t.ID,
			User:       user.ID,
			Attributes: models.TeamAttributeCaptain,
		}})
		if err != nil {
			c.Error(err)
			return
		}
		c.SetHeader("Location", c.BaseURL+"/teams/"+t.ID.String())
		c.SetJSONWithCode(t, 201)
		return
	}

	c.WriteString("Not implemented")
}

// validateUser checks that the user has a valid access token and that the user
// has privileges UserNormal and UserPublic. (is not banned/locked/restricted).
func validateUser(cl *ripple.Client) (*ripple.User, error) {
	user, err := cl.User(ripple.Self)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, &http.ResponseError{
			Code:    401,
			Message: "Access token is invalid - did you revoke the token?",
		}
	}
	if user.Privileges&3 != 3 {
		return nil, &http.ResponseError{
			Code:    403,
			Message: "You don't have the required privileges to register for a tournament.",
		}
	}

	return user, nil
}

func userCanRegister(db *models.DB, tournID models.ID, uid int) (*models.Tournament, error) {
	tourn, err := db.Tournament(tournID)
	if err != nil {
		return nil, err
	}
	if tourn == nil || tourn.Status == models.StatusOrganising {
		return nil, &http.ResponseError{
			Code:    404,
			Message: "That tournament does not exist.",
		}
	}

	// Check whether user in another team of this tournament.
	in, err := db.UserInTournament(tourn.ID, uid)
	if err != nil {
		return nil, err
	}
	if in {
		return nil, &http.ResponseError{
			Code:    409,
			Message: "You are already in this tournament",
		}
	}

	// Check whether we're busy with another tournament
	busy, err := db.UserIsBusy(tourn, uid)
	if err != nil {
		return nil, err
	}
	if busy {
		return nil, &http.ResponseError{
			Code:    409,
			Message: "You can't join another tournament that overlaps with a tournament you're already in!",
		}
	}

	return tourn, nil
}

// TeamMembers retrieves all the members of a team.
func TeamMembers(c *http.Context) {
	members, err := c.DB.TeamMembers(c.ParamID("id"), c.QueryInt("p"))
	if err != nil {
		c.Error(err)
		return
	}
	c.SetJSON(members, false)
}

func init() {
	http.GET("/tournaments/:id/teams", TeamsInTournament)
	http.POST("/tournaments/:id/teams", CreateTeam)
	http.GET("/teams", Teams)
	http.GET("/teams/:id", Team)
	http.GET("/teams/:id/members", TeamMembers)
}
