package client

import "errors"

var (
	Err401StatusCode     = errors.New("unauthorized - 401")
	Err400StatusCode     = errors.New("bad request - 400")
	ErrWrongStatusCode   = errors.New("wrong status code")
	ErrCantCreateRequest = errors.New("can't create request")
)
