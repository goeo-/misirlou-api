// Package ripple provides a Go library to access the Ripple API.
package ripple

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"
)

const defaultInstance = "https://api.ripple.moe"

// StatusError is returned by requests made from Client when the code is not
// in the 2XX series.
type StatusError struct {
	Code     int
	Response *http.Response
}

// Error will attempt to decode the message from the Response, and will
// otherwise return a generic error stating the status code.
func (s *StatusError) Error() string {
	var m struct {
		Message string `json:"message"`
	}
	defer s.Response.Body.Close()
	err := json.NewDecoder(s.Response.Body).Decode(&m)
	if err != nil || m.Message == "" {
		return "ripple: received error with code " + strconv.Itoa(s.Code)
	}
	return "ripple: received error with code " + strconv.Itoa(s.Code) + ": " + m.Message
}

// Self is a special ID used to refer to the user making the request.
// The Ripple API has a common pattern throughout the requests related to users:
// you can refer to an user either by their username, by passing the querystring
// parameter `name`, their user ID, by passing `id`, or reference the user
// making to the request, case in which you would set `id` to `self`. Instead of
// having three methods for each request, this package will understand all
// function calls having id of the user == 0 to reference "self".
const Self = 0

// Client contains the options for the API which will be sent with each request.
// By default, the instance is located at https://api.ripple.moe and requests
// are made without a token (so function as read-only). If Client.Client is not
// set, http.DefaultClient will be used.
type Client struct {
	Instance string
	Token    string
	IsBearer bool
	Client   *http.Client
}

func (c *Client) req(endpoint string, isPost bool, body io.Reader, v interface{}) error {
	// Get base URL (instance)
	url := defaultInstance
	if c.Instance != "" {
		url = c.Instance
	}

	// Build URL and method
	url += "/api/v1/" + endpoint
	method := "GET"
	if isPost {
		method = "POST"
	}

	// Create HTTP request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	// Check if we have a token; if we do, make sure to use the right header.
	if c.Token != "" {
		if c.IsBearer {
			req.Header.Set("Authorization", "Bearer "+c.Token)
		} else {
			req.Header.Set("X-Ripple-Token", c.Token)
		}
	}

	// Pick client to use
	hc := http.DefaultClient
	if c.Client != nil {
		hc = c.Client
	}

	resp, err := hc.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return &StatusError{Code: resp.StatusCode, Response: resp}
	}

	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(v)
}

func getID(id int) string {
	if id == 0 {
		return "?id=self"
	}
	return "?id=" + strconv.Itoa(id)
}

type User struct {
	ID             int       `json:"id"`
	Username       string    `json:"username"`
	UsernameAKA    string    `json:"username_aka"`
	RegisteredOn   time.Time `json:"registered_on"`
	Privileges     uint64    `json:"privileges"`
	LatestActivity time.Time `json:"latest_activity"`
	Country        string    `json:"country"`
}

// User returns an user knowing their ID. May return nil for the user if the
// user could not be found.
func (c *Client) User(id int) (*User, error) {
	var u User
	err := c.req("users"+getID(id), false, nil, &u)
	if err == nil {
		return &u, nil
	}
	if err, ok := err.(*StatusError); ok && err.Code == 404 {
		return nil, nil
	}
	return nil, err
}
