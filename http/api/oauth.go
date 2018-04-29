package api

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"golang.org/x/oauth2"
	"zxq.co/ripple/misirlou-api/http"
	"zxq.co/ripple/misirlou-api/models"
	"zxq.co/x/ripple"
)

var oauthEndpoint = oauth2.Endpoint{
	AuthURL:  "https://ripple.moe/oauth/authorize",
	TokenURL: "https://ripple.moe/oauth/token",
}

func getConfig(c *http.Context) oauth2.Config {
	return oauth2.Config{
		ClientID:     c.OAuth2ClientID,
		ClientSecret: c.OAuth2ClientSecret,
		Endpoint:     oauthEndpoint,
		RedirectURL:  c.BaseURL + "/oauth/finish",
	}
}

// OAuthStart starts the oauth flow.
func OAuthStart(c *http.Context) {
	conf := getConfig(c)

	id := randomStr(9)
	val := randomStr(15)

	err := c.Redis.Set("misirlou:oauth_state:"+id, val, time.Minute*40).Err()
	if err != nil {
		c.Error(err)
		return
	}

	c.SetCookie("oauth_state", id, time.Minute*40)

	c.Redirect(302, conf.AuthCodeURL(val))
}

func randomStr(l int) string {
	b := make([]byte, l)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

// OAuthFinish starts the oauth flow.
func OAuthFinish(c *http.Context) {
	// Check that the state is valid
	sec, err := c.Redis.Get("misirlou:oauth_state:" + c.Cookie("oauth_state")).
		Result()
	if err != nil {
		c.Error(err)
		return
	}
	defer c.DeleteCookie("oauth_state")
	if c.Query("state") != sec {
		c.SetCode(401)
		c.WriteString("State is invalid; please try logging again!")
		return
	}

	// Exchange token with Ripple OAuth
	conf := getConfig(c)
	ctx, _ := context.WithTimeout(context.Background(), 5000*time.Millisecond)
	code, err := conf.Exchange(ctx, c.Query("code"))
	if err != nil {
		c.Error(err)
		return
	}
	// Create temporary Ripple Client to get info about self
	rc := &ripple.Client{
		IsBearer: true,
		Token:    code.AccessToken,
	}
	u, err := rc.User(ripple.Self)
	if err != nil {
		c.SetCode(403)
		c.WriteString("Error from Ripple API: " + err.Error())
		return
	}
	if u == nil {
		c.SetCode(404)
		c.WriteString("User not found?")
		return
	}

	// Save session
	sess := &models.Session{
		ID:          randomStr(15),
		UserID:      u.ID,
		AccessToken: code.AccessToken,
	}
	err = c.DB.SetSession(sess)
	if err != nil {
		c.Error(err)
		return
	}

	c.Redirect(302, c.StoreTokensURL+"?session="+sess.ID+"&access="+sess.AccessToken)
}

func init() {
	http.GET("/oauth/start", OAuthStart)
	http.GET("/oauth/finish", OAuthFinish)
}
