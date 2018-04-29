package api

import (
	"zxq.co/ripple/misirlou-api/http"
	"zxq.co/ripple/misirlou-api/models"
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
	http.GET("/teams", Teams)
	http.GET("/teams/:id", Team)
	http.GET("/teams/:id/members", TeamMembers)
}
