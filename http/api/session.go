package api

import "zxq.co/ripple/misirlou-api/http"

// PutFCMToken sets the current's user FCM token.
func PutFCMToken(c *http.Context) {
	var data struct {
		FCMToken string `json:"fcm_token"`
	}
	if err := c.JSON(&data); err != nil {
		return
	}
	sess := c.Session()
	if sess == nil {
		c.SetCode(403)
		c.WriteString("Missing Authentication header (can't find session)")
		return
	}
	sess.FCMToken = data.FCMToken
	err := c.DB.SetSession(sess)
	if err != nil {
		c.Error(err)
		return
	}
	c.SetJSON(data, false)
}

func init() {
	http.PUT("/fcm_token", PutFCMToken)
}
