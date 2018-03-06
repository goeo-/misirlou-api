package api

import "zxq.co/ripple/misirlou-api/http"

func Home(c *http.Context) {
	c.WriteString("Misirlou API 2.0\nhttps://zxq.co/ripple/misirlou-api")
}

func init() {
	http.GET("/", Home)
}
