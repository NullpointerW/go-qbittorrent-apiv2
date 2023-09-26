package qbt_apiv2

import (
	"fmt"
	"net/url"
)

// Login
func (c *Client) Login(username, password string) error {
	opts := Optional{
		"username": username,
		"password": password,
	}
	resp, err := c.postXwwwFormUrlencoded("auth/login", opts)
	err = RespOk(resp, err)
	if err != nil {
		return err
	}
	if err = RespBodyOk(resp.Body, ErrLoginfailed); err != nil {
		return err
	}
	// add the cookie to cookie jar to authenticate later requests
	if cookies := resp.Cookies(); len(cookies) > 0 {
		u, err := url.Parse(c.URL)
		if err != nil {
			return fmt.Errorf("parse url error: %w", err)
		}
		u.Path = ""
		c.httpCli.Jar.SetCookies(u, cookies)
	}
	return nil
}
