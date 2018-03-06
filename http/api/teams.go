package api

import (
	"zxq.co/ripple/misirlou-api/http"
	"zxq.co/ripple/misirlou-api/models"
)

// Teams fetches the teams that are playing in the given tournament ID.
func Teams(c *http.Context) {
	teams, err := c.DB.GetTeams(&models.Team{
		Tournament: c.QueryInt("tourn_id"),
	}, c.QueryInt("p"))
	if err != nil {
		c.Error(err)
		return
	}
	c.SetJSON(200, teams)
}

func init() {
	http.GET("/teams", Teams)
}
