package qbt_apiv4

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"path/filepath"

	errwrp "github.com/pkg/errors"
	"golang.org/x/net/publicsuffix"
)

type Client struct {
	http          *http.Client
	URL           string
	Authenticated bool
	Jar           http.CookieJar
}

// NewCli v4
func NewCli(url string, auth ...string) (*Client, error) {
	client := &Client{}

	// ensure url ends with "/"
	if url[len(url)-1:] != "/" {
		url += "/"
	}

	client.URL = url + "api/v2/"

	// create cookie jar
	client.Jar, _ = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	client.http = &http.Client{
		Jar: client.Jar,
	}

	nrequired := len(auth) == 0
	if !nrequired {
		var (
			resp *http.Response
			err  error
		)
		if resp, err = client.Login(auth[0], auth[1]); err != nil {
			return nil, err
		} else if resp.Status != "200 OK" {
			return nil, errors.New("login failed")
		} else {
			b, _ := io.ReadAll(resp.Body)
			fmt.Println(string(b))
			if string(b) != "Ok." {
				return nil, errors.New("login failed")
			}
		}
	}

	client.Authenticated = true

	return client, nil
}

// Login 
func (c *Client) Login(username, password string) (*http.Response, error) {
	v := url.Values{}
	v.Set("username", username)
	v.Set("password", password)

	req, err := http.NewRequest("POST", c.URL+"auth/login", bytes.NewBufferString(v.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.http.Do(req)

	if err != nil {
		return nil, err
	} else if resp.Status != "200 OK" { // check for correct status code
		return nil, errwrp.Wrap(qbt.ErrBadResponse, "couldnt log in")
	}

	// add the cookie to cookie jar to authenticate later requests
	if cookies := resp.Cookies(); len(cookies) > 0 {
		u, err := url.Parse(c.URL)
		if err != nil {
			return nil, err
		}
		u.Path = ""
		fmt.Println(u.String())
		c.Jar.SetCookies(u, cookies)
	}
	return resp, nil
}


// Torrent management
func (c *Client) AddNewTorrent(magnetLink, path string) (*http.Response, error) {
	ap, err := filepath.Abs(path)
	if err != nil {
		return nil, errwrp.Wrapf(err, "cannot conv abs_path %s", path)
	}
	opt := optional{
		"urls":     magnetLink,
		"savepath": ap,
	}
	resp, err := c.postMultipartData("torrents/add", opt)
	if err != nil {
		return nil, err
	}
	return resp, nil

}
