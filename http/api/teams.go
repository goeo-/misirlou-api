package api

import (
	"zxq.co/ripple/misirlou-api/http"
	"zxq.co/ripple/misirlou-api/models"
)

// Teams fetches the teams that are playing in the given tournament ID.
func Teams(c *http.Context) {
	teams(c, models.TeamFilters{
		Tournament: c.QueryInt("tournament"),
	})
}

// TeamsInTournament returns all the teams in a tournament, which basically is
// sugar for calling Teams giving a tournament.
func TeamsInTournament(c *http.Context) {
	teams(c, models.TeamFilters{
		Tournament:      c.ParamInt("id"),
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
	c.SetJSON(200, teams)
}

func Team(c *http.Context) {
	team, err := c.DB.Team(c.ParamInt("id"))
	if err != nil {
		c.Error(err)
		return
	}
	if team == nil {
		c.SetJSON(404, nil)
	} else {
		c.SetJSON(200, team)
	}
}

func init() {
	http.GET("/tournaments/:id/teams", TeamsInTournament)
	http.GET("/teams", Teams)
	http.GET("/teams/:id", Team)
}
