package qbt_apiv2

import "errors"

var (
	ErrBadResponse = errors.New("bad response")

	ErrLoginfailed = errors.New("login failed")

	ErrAddTorrnetfailed = errors.New("add torrnet failed")
)
