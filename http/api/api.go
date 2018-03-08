package api

import (
	"strconv"

	"zxq.co/ripple/misirlou-api/http"
)

func Home(c *http.Context) {
	c.WriteString("Misirlou API 2.0\nhttps://zxq.co/ripple/misirlou-api\n")
	s := c.Session()
	if s != nil {
		c.WriteString("Ripple User: " + strconv.Itoa(s.UserID))
	}
}

func init() {
	http.GET("/", Home)
}
