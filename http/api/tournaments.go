package api

import (
	"zxq.co/ripple/misirlou-api/http"
)

// Tournaments fetches tournaments, desc-sorted by ID.
func Tournaments(c *http.Context) {
	tourns, err := c.DB.Tournaments(c.QueryInt("id"), c.QueryInt("p"))
	if err != nil {
		c.Error(err)
		return
	}
	c.SetJSON(200, tourns)
}

// TournamentRules fetches the rules of a tournament.
func TournamentRules(c *http.Context) {
	rules, err := c.DB.TournamentRules(c.QueryInt("id"))
	if err != nil {
		c.Error(err)
		return
	}
	c.SetJSON(200, rules)
}

func init() {
	http.GET("/tournaments", Tournaments)
	http.GET("/tournaments/rules", TournamentRules)
}
