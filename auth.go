package qbt_apiv2

import (
	"net/http"
	"net/url"
)

// Login
func (c *Client) Login(username, password string) (*http.Response, error) {
	opts := optional{
		"username": username,
		"password": password,
	}
	resp, err := c.postXwwwFormUrlencoded("auth/login", opts)
	err = RespOk(resp, err)
	if err != nil {
		return nil, err
	}
	// add the cookie to cookie jar to authenticate later requests
	if cookies := resp.Cookies(); len(cookies) > 0 {
		u, err := url.Parse(c.URL)
		if err != nil {
			return nil, err
		}
		u.Path = ""
		c.httpCli.Jar.SetCookies(u, cookies)
	}
	return resp, nil
}
