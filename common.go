package qbt_apiv2

import (
	"fmt"
	errwrp "github.com/pkg/errors"
	"io"
	"net/http"
)

const (
	ResponseBodyOK   = "Ok."
	ResponseBodyFAIL = "Fails."
)

// optional parameters when sending HTTP requests
type optional map[string]any

func (opt optional) StringField() map[string]string {
	m := make(map[string]string)
	for k, v := range opt {
		m[k] = fmt.Sprintf("%v", v)
	}
	return m
}

func RespOk(resp *http.Response, err error) error {
	if err != nil {
		return err
	} else if resp.Status != "200 OK" { // check for correct status code
		return errwrp.Wrapf(ErrBadResponse, "status code:%s", resp.Status)
	} else {
		return nil
	}
}

func RespBodyOk(body io.ReadCloser, bizErr error) error {
	defer body.Close()
	b, err := io.ReadAll(body)
	if err != nil {
		return err
	}
	if string(b) != ResponseBodyOK {
		return bizErr
	}
	return nil
}

func ignrBody(body io.ReadCloser) error {
	_, err := io.Copy(io.Discard, body)
	return err
}
