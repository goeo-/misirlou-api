package api

import (
	"zxq.co/ripple/misirlou-api/http"
)

// Tournaments fetches tournaments, desc-sorted by ID.
func Tournaments(c *http.Context) {
	tourns, err := c.DB.Tournaments(c.QueryInt("p"))
	if err != nil {
		c.Error(err)
		return
	}
	c.SetJSON(200, tourns)
}

// Tournament retrieves a single tournament knowing its ID.
func Tournament(c *http.Context) {
	tourn, err := c.DB.Tournament(c.ParamInt("id"))
	if err != nil {
		c.Error(err)
		return
	}
	if tourn == nil {
		c.SetJSON(404, nil)
	} else {
		c.SetJSON(200, tourn)
	}
}

// TournamentRules fetches the rules of a tournament.
func TournamentRules(c *http.Context) {
	rules, err := c.DB.TournamentRules(c.ParamInt("id"))
	if err != nil {
		c.Error(err)
		return
	}
	if rules == nil {
		c.SetJSON(404, nil)
	} else {
		c.SetJSON(200, rules)
	}
}

func init() {
	http.GET("/tournaments", Tournaments)
	http.GET("/tournaments/:id", Tournament)
	http.GET("/tournaments/:id/rules", TournamentRules)
}
