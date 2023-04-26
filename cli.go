package qbt_apiv2

import (
	"bytes"
	errwrp "github.com/pkg/errors"
	"golang.org/x/net/publicsuffix"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path"
)

type Client struct {
	httpCli *http.Client
	URL     string
	// API `sync/maindata`` Parameter `rid`
	rid int
}

// NewCli v2
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

	nreq := len(auth) == 0

	var (
		resp *http.Response
		err  error
	)
	if nreq {
		resp, err = client.Login("", "")
	} else {
		resp, err = client.Login(auth[0], auth[1])
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if string(b) != "Ok." {
		return nil, ErrLoginfailed
	}

	return client, nil
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
