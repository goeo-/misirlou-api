package api

import (
	"zxq.co/ripple/misirlou-api/http"
)

// Teams fetches the teams that are playing in the given tournament ID.
func Teams(c *http.Context) {
	teams, err := c.DB.Teams(c.QueryInt("tournament"), c.QueryInt("p"))
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
	http.GET("/teams", Teams)
	http.GET("/teams/:id", Team)
}
