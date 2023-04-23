package qbt_apiv4

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
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
			defer resp.Body.Close()
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
	err = RespOk(resp, err)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) TorrentList(opt optional) (BasicTorrent, error) {
	values := url.Values{}
	for k, v := range opt.StringField() {
		values.Set(k, v)
	}

	req, err := http.NewRequest("POST", c.URL+"torrents/info", bytes.NewBufferString(values.Encode()))
	if err != nil {
		return BasicTorrent{}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.http.Do(req)

	err = RespOk(resp, err)

	if err != nil {
		return BasicTorrent{}, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(req.Body)
	if err != nil {
		return BasicTorrent{}, err
	}
	bt := new(BasicTorrent)
	err = json.Unmarshal(b, bt)
	if err != nil {
		return BasicTorrent{}, err
	}
	return *bt, nil
}

// writeOptions will write a map to the buffer through multipart.NewWriter
func writeOptions(writer *multipart.Writer, opts optional) {
	ws := opts.StringField()
	for key, val := range ws {
		writer.WriteField(key, val)
	}
}

// postMultipart will perform a multiple part POST request
func (client *Client) postMultipart(endpoint string, buffer bytes.Buffer, contentType string) (*http.Response, error) {
	req, err := http.NewRequest("POST", client.URL+endpoint, &buffer)
	if err != nil {
		return nil, errwrp.Wrap(err, "error creating request")
	}

	// add the content-type so qbittorrent knows what to expect
	req.Header.Set("Content-Type", contentType)

	resp, err := client.http.Do(req)
	if err != nil {
		return nil, errwrp.Wrap(err, "failed to perform request")
	}

	return resp, nil
}

// postMultipartData will perform a multiple part POST request without a file
func (client *Client) postMultipartData(endpoint string, opts optional) (*http.Response, error) {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	// write the options to the buffer
	// will contain the link string
	writeOptions(writer, opts)

	// close the writer before doing request to get closing line on multipart request
	if err := writer.Close(); err != nil {
		return nil, errwrp.Wrap(err, "failed to close writer")
	}

	resp, err := client.postMultipart(endpoint, buffer, writer.FormDataContentType())
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// postMultipartFile will perform a multiple part POST request with a file
func (client *Client) postMultipartFile(endpoint string, fileName string, opts optional) (*http.Response, error) {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	// open the file for reading
	file, err := os.Open(fileName)
	if err != nil {
		return nil, errwrp.Wrap(err, "error opening file")
	}
	// defer the closing of the file until the end of function
	// so we can still copy its contents
	defer file.Close()

	// create form for writing the file to and give it the filename
	formWriter, err := writer.CreateFormFile("torrents", path.Base(fileName))
	if err != nil {
		return nil, errwrp.Wrap(err, "error adding file")
	}

	// write the options to the buffer
	writeOptions(writer, opts)

	// copy the file contents into the form
	if _, err = io.Copy(formWriter, file); err != nil {
		return nil, errwrp.Wrap(err, "error copying file")
	}

	// close the writer before doing request to get closing line on multipart request
	if err := writer.Close(); err != nil {
		return nil, errwrp.Wrap(err, "failed to close writer")
	}

	resp, err := client.postMultipart(endpoint, buffer, writer.FormDataContentType())
	if err != nil {
		return nil, err
	}

	return resp, nil
}
