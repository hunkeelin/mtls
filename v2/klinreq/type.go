package klinreq

import (
	"net/http"
)

type ReqQ struct {
	Url      *string
	Headers  map[string]string
	Method   *string
	NoVerify bool
	Json     *interface{}
	Client   *http.Cilent
}
type ReqBuilder struct {
	ReqQ
}
