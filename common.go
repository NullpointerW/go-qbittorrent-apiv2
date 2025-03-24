package qbt_apiv2

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"unicode"
)

const (
	ResponseBodyOK   = "Ok."
	ResponseBodyFAIL = "Fails."
)

// Optional parameters when sending HTTP requests
type Optional map[string]any

func (opt Optional) StringField() map[string]string {
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
		return fmt.Errorf("%w: %s", ErrBadResponse, resp.Status)
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

func versionInt(v string) int {
	v = strings.ToLower(v)
	l := strings.LastIndex(v, ".")
	if l == -1 {
		return 0
	}
	sfx := v[l+1:]

	for i, r := range sfx {
		if !unicode.IsNumber(r) {
			sfx = sfx[:i]
			break
		}
	}
	sfxb, vb := []byte(sfx), []byte(v)
	vb = append(vb[:l+1], sfxb...)
	var major, minor, patch int
	fmt.Sscanf(string(vb), "v%d.%d.%d", &major, &minor, &patch)
	return major*100 + minor*10 + patch
}
