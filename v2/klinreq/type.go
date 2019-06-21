package klinreq

import (
	"net/http"
)

type ReqQ struct {
	Url       *string
	Headers   map[string]string
	Method    *string
	NoVerify  bool
	BodyBytes []byte
	Json      *interface{}
	Client    *http.Client
}
type ReqBuilder struct {
	ReqQ
}
