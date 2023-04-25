package qbt_apiv2

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"

	errwrp "github.com/pkg/errors"
	"golang.org/x/net/publicsuffix"
)

type Client struct {
	httpCli *http.Client
	URL     string
}

// NewCli v4
func NewCli(url string, auth ...string) (*Client, error) {
	client := new(Client)

	// ensure url ends with "/"
	if url[len(url)-1:] != "/" {
		url += "/"
	}

	client.URL = url + "api/v2/"

	// create cookie jar
	cliJar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	client.httpCli = &http.Client{
		Jar: cliJar,
	}

	nrequired := len(auth) == 0

	var (
		resp *http.Response
		err  error
	)
	if nrequired {
		resp, err = client.Login("", "")
	} else {
		resp, err = client.Login(auth[0], auth[1])
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(resp.Body)
	if string(b) != "Ok." {
		return nil, errors.New("login failed")
	}

	return client, nil
}

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

// Torrent management
func (c *Client) AddNewTorrent(urlLink, path string) (*http.Response, error) {
	ap, err := filepath.Abs(path)
	if err != nil {
		return nil, errwrp.Wrapf(err, "cannot conv abs_path %s", path)
	}
	opt := optional{
		"urls":     urlLink,
		"savepath": ap,
	}
	resp, err := c.postMultipartData("torrents/add", opt)
	err = RespOk(resp, err)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) TorrentList(opt optional) ([]BasicTorrent, error) {
	resp, err := c.postXwwwFormUrlencoded("torrents/info", opt)

	err = RespOk(resp, err)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	bt := new([]BasicTorrent)
	err = json.Unmarshal(b, bt)
	if err != nil {
		return nil, err
	}
	return *bt, nil
}

func (c *Client) GetTorrentProperties(hash string) (Torrent, error) {
	resp, err := c.postXwwwFormUrlencoded("torrents/properties", optional{
		"hash": hash,
	})
	err = RespOk(resp, err)
	if err != nil {
		return Torrent{}, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Torrent{}, err
	}
	t := new(Torrent)
	err = json.Unmarshal(b, t)
	if err != nil {
		return Torrent{}, err
	}
	return *t, nil
}

func (c *Client) GetTorrentContents(hash string, indexes ...int) ([]TorrentFile, error) {
	opt := optional{
		"hash": hash,
	}
	if len(indexes) > 0 {
		var idxes string
		for _, idx := range indexes {
			idxes += strconv.Itoa(idx) + "|"
		}
		idxes = idxes[:len(idxes)-1]
		opt["indexes"] = idxes
	}

	resp, err := c.postXwwwFormUrlencoded("torrents/files", opt)
	err = RespOk(resp, err)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	tf := new([]TorrentFile)
	err = json.Unmarshal(b, tf)
	if err != nil {
		return nil, err
	}
	return *tf, nil
}

// Sync
// Sync API implements requests for obtaining changes since the last request. All Sync API methods are under "sync"
func (c *Client) GetMainData(rid int) (Sync, error) {
	resp, err := c.postXwwwFormUrlencoded("sync/maindata", optional{
		"rid": rid,
	})
	err = RespOk(resp, err)
	if err != nil {
		return Sync{}, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Sync{}, err
	}
	s := new(Sync)
	err = json.Unmarshal(b, s)
	if err != nil {
		return Sync{}, err
	}
	return *s, nil
}

// RSS All RSS API methods are under "rss", e.g.: /api/v2/rss/methodName.
func (c *Client) AddFolder(path string) error {
	resp, err := c.postXwwwFormUrlencoded("rss/addFolder", optional{
		"path": path,
	})
	err = RespOk(resp, err)
	if err != nil {
		return err
	}
	ignrBody(resp.Body)
	return nil
}

func (c *Client) AddFeed(url, path string) error {
	opt := optional{
		"url": url,
	}
	if path != "" {
		opt["path"] = path
	}
	resp, err := c.postXwwwFormUrlencoded("rss/addFeed", opt)
	err = RespOk(resp, err)
	if err != nil {
		return err
	}
	ignrBody(resp.Body)
	return nil
}

func (c *Client) RemoveItem(path string) error {
	resp, err := c.postXwwwFormUrlencoded("rss/removeItem", optional{
		"path": path,
	})
	err = RespOk(resp, err)
	if err != nil {
		return err
	}
	ignrBody(resp.Body)
	return nil
}

func (c *Client) MoveItem(dst, src string) error {
	resp, err := c.postXwwwFormUrlencoded("rss/moveItem", optional{
		"itemPath": src,
		"destPath": dst,
	})
	err = RespOk(resp, err)
	if err != nil {
		return err
	}
	ignrBody(resp.Body)
	return nil
}

func (c *Client) GetAllItems(withData bool) (RssItem, error) {
	opt := optional{}
	if withData {
		opt["withData"] = true
	}
	resp, err := c.postXwwwFormUrlencoded("rss/items", opt)
	err = RespOk(resp, err)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	ri := new(RssItem)
	json.Unmarshal(b, ri)
	return *ri, nil
}

func (c *Client) MarkAsRead(itemPath, articleId string) error {
	opt := optional{
		"itemPath": itemPath,
	}
	if articleId != "" {
		opt["articleId"] = articleId
	}
	resp, err := c.postXwwwFormUrlencoded("rss/markAsRead", opt)
	err = RespOk(resp, err)
	if err != nil {
		return err
	}
	ignrBody(resp.Body)
	return nil
}

// Common Methods for HTTP Requests
// Use POST request to send x-www-form-urlencoded encoding.
func (c *Client) postXwwwFormUrlencoded(endpoint string, opts optional) (*http.Response, error) {
	values := url.Values{}
	for k, v := range opts.StringField() {
		values.Set(k, v)
	}

	req, err := http.NewRequest("POST", c.URL+endpoint, bytes.NewBufferString(values.Encode()))
	if err != nil {
		return nil, errwrp.Wrap(err, "error creating request")
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpCli.Do(req)
	if err != nil {
		return nil, errwrp.Wrap(err, "failed to perform request")
	}
	return resp, nil
}

// writeOptions will write a map to the buffer through multipart.NewWriter
func writeOptions(writer *multipart.Writer, opts optional) {
	ws := opts.StringField()
	for key, val := range ws {
		writer.WriteField(key, val)
	}
}

// postMultipart will perform a multiple part POST request
func (c *Client) postMultipart(endpoint string, buffer bytes.Buffer, contentType string) (*http.Response, error) {
	req, err := http.NewRequest("POST", c.URL+endpoint, &buffer)
	if err != nil {
		return nil, errwrp.Wrap(err, "error creating request")
	}

	// add the content-type so qbittorrent knows what to expect
	req.Header.Set("Content-Type", contentType)

	resp, err := c.httpCli.Do(req)
	if err != nil {
		return nil, errwrp.Wrap(err, "failed to perform request")
	}

	return resp, nil
}

// postMultipartData will perform a multiple part POST request without a file
func (c *Client) postMultipartData(endpoint string, opts optional) (*http.Response, error) {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)

	// write the options to the buffer
	// will contain the link string
	writeOptions(writer, opts)

	// close the writer before doing request to get closing line on multipart request
	if err := writer.Close(); err != nil {
		return nil, errwrp.Wrap(err, "failed to close writer")
	}

	resp, err := c.postMultipart(endpoint, buffer, writer.FormDataContentType())
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// postMultipartFile will perform a multiple part POST request with a file
func (c *Client) postMultipartFile(endpoint string, fileName string, opts optional) (*http.Response, error) {
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

	resp, err := c.postMultipart(endpoint, buffer, writer.FormDataContentType())
	if err != nil {
		return nil, err
	}

	return resp, nil
}
