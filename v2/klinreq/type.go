package klinreq

type ReqQ struct {
	Url     *string
	Headers map[string]string
	Method  *string
	Json    *interface{}
}
type ReqBuilder struct {
	ReqQ
}
